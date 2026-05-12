package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
)

type RoutineHandler struct {
	routineService *service.RoutineService
}

func NewRoutineHandler(svc *service.RoutineService) *RoutineHandler {
	return &RoutineHandler{
		routineService: svc,
	}
}

func (h *RoutineHandler) CreateRoutine(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	var req service.CreateRoutineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Routine name is required", http.StatusBadRequest)
		return
	}

	routine, err := h.routineService.CreateRoutine(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(routine)
}

func (h *RoutineHandler) GetRoutine(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	routineID := domain.RoutineID(r.PathValue("id"))

	if routineID == "" {
		http.Error(w, "Routine ID is required", http.StatusBadRequest)
		return
	}

	routine, err := h.routineService.GetRoutine(routineID, userID)
	if err != nil {
		if err.Error() == "routine not found: record not found" {
			http.Error(w, "Routine not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: routine belongs to different user" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routine)
}

func (h *RoutineHandler) ListRoutines(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	routines, err := h.routineService.ListRoutines(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"routines": routines,
		"count":   len(routines),
	})
}

func (h *RoutineHandler) DeleteRoutine(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	routineID := domain.RoutineID(r.PathValue("id"))

	if routineID == "" {
		http.Error(w, "Routine ID is required", http.StatusBadRequest)
		return
	}

	err := h.routineService.DeleteRoutine(routineID, userID)
	if err != nil {
		if err.Error() == "routine not found: record not found" {
			http.Error(w, "Routine not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: routine belongs to different user" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}