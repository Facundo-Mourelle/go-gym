package api

import (
	"net/http"

	"github.com/Facundo-Mourelle/go-gym/internal/api/handler"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
)

func NewRouter(
	authService *service.AuthService,
	exerciseService *service.ExerciseService,
	workoutService *service.WorkoutService,
	sessionService *service.SessionService,
	metricsService *service.MetricsService,
	routineService *service.RoutineService,
	equipmentService *service.EquipmentService,
) http.Handler {

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"gym-tracker-api"}`))
	})

	// Auth endpoints
	authHandler := handler.NewAuthHandler(authService)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/auth/me", Protected(authService, authHandler.GetCurrentUser))

	// Exercise endpoints
	exerciseHandler := handler.NewExerciseHandler(exerciseService)
	mux.HandleFunc("GET /api/v1/exercises/templates", exerciseHandler.ListTemplates)
	mux.HandleFunc("GET /api/v1/exercises", Protected(authService, exerciseHandler.ListExercises))
	mux.HandleFunc("GET /api/v1/exercises/{id}", exerciseHandler.GetExercise)
	mux.HandleFunc("POST /api/v1/exercises/custom", Protected(authService, exerciseHandler.CreateCustomExercise))
	mux.HandleFunc("POST /api/v1/exercises/from-template", Protected(authService, exerciseHandler.CreateFromTemplate))
	mux.HandleFunc("PUT /api/v1/exercises/{id}", Protected(authService, exerciseHandler.UpdateExercise))
	mux.HandleFunc("DELETE /api/v1/exercises/{id}", Protected(authService, exerciseHandler.DeleteExercise))
	mux.HandleFunc("GET /api/v1/patterns", exerciseHandler.ListPatterns)

	// Workout endpoints
	workoutHandler := handler.NewWorkoutHandler(workoutService)
	mux.HandleFunc("POST /api/v1/workouts", Protected(authService, workoutHandler.CreateWorkout))
	mux.HandleFunc("GET /api/v1/workouts", Protected(authService, workoutHandler.ListWorkouts))

	// Session endpoints
	sessionHandler := handler.NewSessionHandler(sessionService)
	mux.HandleFunc("POST /api/v1/sessions", Protected(authService, sessionHandler.StartSession))
	mux.HandleFunc("GET /api/v1/sessions", Protected(authService, sessionHandler.ListSessions))
	mux.HandleFunc("GET /api/v1/sessions/{id}", Protected(authService, sessionHandler.GetSession))
	mux.HandleFunc("POST /api/v1/sessions/{id}/sets", Protected(authService, sessionHandler.RecordSet))
	mux.HandleFunc("PUT /api/v1/sessions/{id}/sets/{setId}", Protected(authService, sessionHandler.UpdateSet))
	mux.HandleFunc("DELETE /api/v1/sessions/{id}/sets/{setId}", Protected(authService, sessionHandler.DeleteSet))
	mux.HandleFunc("PUT /api/v1/sessions/{id}/complete", Protected(authService, sessionHandler.CompleteSession))

	// Metrics endpoints
	metricsHandler := handler.NewMetricsHandler(metricsService)
	mux.HandleFunc("GET /api/v1/metrics/progress/{exerciseId}", Protected(authService, metricsHandler.GetExerciseProgress))

	routineHandler := handler.NewRoutineHandler(routineService)
	mux.HandleFunc("POST /api/v1/routines", Protected(authService, routineHandler.CreateRoutine))
	mux.HandleFunc("GET /api/v1/routines", Protected(authService, routineHandler.ListRoutines))
	mux.HandleFunc("GET /api/v1/routines/{id}", Protected(authService, routineHandler.GetRoutine))
	mux.HandleFunc("DELETE /api/v1/routines/{id}", Protected(authService, routineHandler.DeleteRoutine))

	// Equipment endpoints
	equipmentHandler := handler.NewEquipmentHandler(equipmentService)
	mux.HandleFunc("GET /api/v1/equipment", Protected(authService, equipmentHandler.ListEquipment))
	mux.HandleFunc("POST /api/v1/equipment", Protected(authService, equipmentHandler.CreateEquipment))
	mux.HandleFunc("PUT /api/v1/equipment/{id}", Protected(authService, equipmentHandler.UpdateEquipment))
	mux.HandleFunc("DELETE /api/v1/equipment/{id}", Protected(authService, equipmentHandler.DeleteEquipment))

	return mux
}
