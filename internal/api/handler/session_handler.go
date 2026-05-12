package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: svc,
	}
}

// POST /api/v1/sessions
func (h *SessionHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	var req struct {
		WorkoutPlanID *string `json:"workout_plan_id"`
		Notes         string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var workoutPlanID *domain.WorkoutPlanID
	if req.WorkoutPlanID != nil {
		id := domain.WorkoutPlanID(*req.WorkoutPlanID)
		workoutPlanID = &id
	}

	response, err := h.sessionService.StartSession(userID, service.StartSessionRequest{
		WorkoutPlanID: workoutPlanID,
		Notes:         req.Notes,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /api/v1/sessions/{id}/sets
func (h *SessionHandler) RecordSet(w http.ResponseWriter, r *http.Request) {
	sessionID := domain.SessionID(r.PathValue("id"))

	var req struct {
		WorkoutExerciseID *string `json:"workout_exercise_id"`
		ExerciseID        string  `json:"exercise_id"`
		SetNumber         int     `json:"set_number"`
		Reps              int     `json:"reps"`
		RepsInReserve     int     `json:"reps_in_reserve"`
		RawLoad           float64 `json:"raw_load"`
		EquipmentID       string  `json:"equipment_id"`
		RestSeconds       int     `json:"rest_seconds"`
		Notes             string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var workoutExerciseID *domain.WorkoutExerciseID
	if req.WorkoutExerciseID != nil {
		id := domain.WorkoutExerciseID(*req.WorkoutExerciseID)
		workoutExerciseID = &id
	}

	response, err := h.sessionService.RecordSet(sessionID, service.RecordSetRequest{
		WorkoutExerciseID: workoutExerciseID,
		ExerciseID:        domain.ExerciseID(req.ExerciseID),
		SetNumber:         req.SetNumber,
		Reps:              req.Reps,
		RepsInReserve:     req.RepsInReserve,
		RawLoad:           req.RawLoad,
		EquipmentID:       domain.EquipmentID(req.EquipmentID),
		RestSeconds:       req.RestSeconds,
		Notes:             req.Notes,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// PUT /api/v1/sessions/{id}/sets/{setId}
func (h *SessionHandler) UpdateSet(w http.ResponseWriter, r *http.Request) {
	sessionID := domain.SessionID(r.PathValue("id"))
	setID := domain.PerformedSetID(r.PathValue("setId"))

	var req struct {
		Reps          *int     `json:"reps"`
		RawLoad       *float64 `json:"raw_load"`
		EquipmentID   *string  `json:"equipment_id"`
		RestSeconds   *int     `json:"rest_seconds"`
		RepsInReserve *int     `json:"reps_in_reserve"`
		Notes         *string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var equipmentID *domain.EquipmentID
	if req.EquipmentID != nil {
		id := domain.EquipmentID(*req.EquipmentID)
		equipmentID = &id
	}

	response, err := h.sessionService.UpdateSet(sessionID, setID, service.UpdateSetRequest{
		Reps:          req.Reps,
		RawLoad:       req.RawLoad,
		EquipmentID:   equipmentID,
		RestSeconds:   req.RestSeconds,
		RepsInReserve: req.RepsInReserve,
		Notes:         req.Notes,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DELETE /api/v1/sessions/{id}/sets/{setId}
func (h *SessionHandler) DeleteSet(w http.ResponseWriter, r *http.Request) {
	sessionID := domain.SessionID(r.PathValue("id"))
	setID := domain.PerformedSetID(r.PathValue("setId"))

	if err := h.sessionService.DeleteSet(sessionID, setID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PUT /api/v1/sessions/{id}/complete
func (h *SessionHandler) CompleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := domain.SessionID(r.PathValue("id"))

	var req struct {
		Notes string `json:"notes"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	response, err := h.sessionService.CompleteSession(sessionID, req.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/sessions/{id}
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := domain.SessionID(r.PathValue("id"))

	response, err := h.sessionService.GetSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/sessions
func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	sessions, err := h.sessionService.ListUserSessions(userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}
