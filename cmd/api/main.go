package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"database/sql"

	"github.com/Facundo-Mourelle/go-gym/internal/api"
	"github.com/Facundo-Mourelle/go-gym/internal/api/middleware"
	"github.com/Facundo-Mourelle/go-gym/internal/config"
	"github.com/Facundo-Mourelle/go-gym/internal/domain/calculator"
	"github.com/Facundo-Mourelle/go-gym/internal/domain/resistance"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
	"github.com/Facundo-Mourelle/go-gym/internal/repository/memory"
	"github.com/Facundo-Mourelle/go-gym/internal/repository/postgres"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	var db *sql.DB
	if cfg.DatabaseURL != "" {
		var err error
		db, err = sql.Open("postgres", cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err := runMigrations(db, cfg.DatabaseURL); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	var userRepo repository.UserRepository
	var exerciseRepo repository.ExerciseRepository
	var workoutRepo repository.WorkoutRepository
	var sessionRepo repository.SessionRepository
	var equipmentRepo repository.EquipmentRepository
	var routineRepo repository.RoutineRepository

	if db != nil {
		userRepo = postgres.NewUserPostgresRepository(db)
		exerciseRepo = postgres.NewExercisePostgresRepository(db)
		workoutRepo = postgres.NewWorkoutPostgresRepository(db)
		sessionRepo = postgres.NewSessionPostgresRepository(db)
		equipmentRepo = postgres.NewEquipmentPostgresRepository(db)
		routineRepo = postgres.NewRoutinePostgresRepository(db)
		log.Println("Using PostgreSQL repositories")
	} else {
		userRepo = memory.NewUserMemoryRepository()
		exerciseRepo = memory.NewExerciseMemoryRepository()
		workoutRepo = memory.NewWorkoutMemoryRepository()
		sessionRepo = memory.NewSessionMemoryRepository()
		equipmentRepo = memory.NewEquipmentMemoryRepository()
		routineRepo = memory.NewRoutineMemoryRepository()
		log.Println("Using In-Memory repositories")
	}

	// Initialize resistance profile registry
	profileRegistry := resistance.NewRegistry()
	profileRegistry.Register(resistance.NewFreeWeightProfile())

	// Register some machine profiles
	profileRegistry.Register(resistance.NewMachineProfile(
		"machine_2to1",
		resistance.Pulley2to1,
		0.5,  // 2:1 mechanical advantage
		0.05, // 5% friction loss
	))

	// Initialize services
	routineService := service.NewRoutineService(routineRepo)
	authService := service.NewAuthService(userRepo, routineService, cfg)
	exerciseService := service.NewExerciseService(exerciseRepo)
	equipmentService := service.NewEquipmentService(equipmentRepo)
	workoutService := service.NewWorkoutService(workoutRepo, exerciseRepo)
	sessionService := service.NewSessionService(
		sessionRepo,
		exerciseRepo,
		equipmentRepo,
		workoutRepo,
		profileRegistry,
	)

	// Initialize calculators
	volumeCalc := calculator.NewVolumeCalculator(profileRegistry)
	progressCalc := calculator.NewProgressCalculator(calculator.Epley1RM{})

	metricsService := service.NewMetricsService(
		sessionRepo,
		exerciseRepo,
		equipmentRepo,
		volumeCalc,
		progressCalc,
	)

	router := api.NewRouter(
		authService,
		exerciseService,
		workoutService,
		sessionService,
		metricsService,
		routineService,
		equipmentService,
	)

	var handler http.Handler = router

	handler = middleware.CORSMiddleware(cfg.CORSAllowedOrigins)(handler)
	handler = middleware.ContentTypeJSON(handler)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s (environment: %s)", addr, cfg.Environment)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		var err error
		if cfg.TLSEnabled {
			if cfg.TLSCertFile == "" || cfg.TLSKeyFile == "" {
				log.Fatal("TLS_CERT_FILE and TLS_KEY_FILE must be set when TLS_ENABLED=true")
			}
			log.Printf("Starting HTTPS server on %s", addr)
			err = server.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
		} else {
			log.Printf("Starting HTTP server on %s", addr)
			err = server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func runMigrations(db *sql.DB, databaseURL string) error {
	migrationsPath, err := filepath.Abs("internal/repository/postgres/migrations")
	if err != nil {
		return fmt.Errorf("failed to resolve migrations path: %w", err)
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("Database schema is up to date")
	} else {
		log.Println("Database migrations applied successfully")
	}

	return nil
}
