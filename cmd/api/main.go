package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"database/sql"

	"github.com/Facundo-Mourelle/go-gym/internal/api"
	"github.com/Facundo-Mourelle/go-gym/internal/config"
	"github.com/Facundo-Mourelle/go-gym/internal/domain/calculator"
	"github.com/Facundo-Mourelle/go-gym/internal/domain/resistance"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
	"github.com/Facundo-Mourelle/go-gym/internal/repository/memory"
	"github.com/Facundo-Mourelle/go-gym/internal/repository/postgres"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
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

	handler := corsMiddleware(router)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)

	log.Fatal(http.ListenAndServe(addr, handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
