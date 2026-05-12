package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
)

type MetricsHandler struct {
	metricsService *service.MetricsService
}

func NewMetricsHandler(svc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		metricsService: svc,
	}
}

// GET /api/v1/metrics/progress/{exerciseId}?start_date=2025-01-01&end_date=2026-01-16
func (h *MetricsHandler) GetExerciseProgress(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	exerciseID := domain.ExerciseID(r.PathValue("exerciseId"))

	if exerciseID == "" {
		http.Error(w, "Exercise ID cannot be empty", http.StatusBadRequest)
		return
	}

	// Parse date range
	var startDate, endDate *time.Time

	if sd := r.URL.Query().Get("start_date"); sd != "" {
		parsed, err := time.Parse("2006-01-02", sd)
		if err != nil {
			http.Error(w, "Invalid start_date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		startDate = &parsed
	}

	if ed := r.URL.Query().Get("end_date"); ed != "" {
		parsed, err := time.Parse("2006-01-02", ed)
		if err != nil {
			http.Error(w, "Invalid end_date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		endDate = &parsed
	}

	// Get progress data
	response, err := h.metricsService.GetExerciseProgress(
		userID,
		exerciseID,
		startDate,
		endDate,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
