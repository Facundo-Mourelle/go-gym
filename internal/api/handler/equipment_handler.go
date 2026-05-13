package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
)

type EquipmentHandler struct {
	equipmentService *service.EquipmentService
}

func NewEquipmentHandler(svc *service.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{equipmentService: svc}
}

type EquipmentResponse struct {
	ID                    string               `json:"id"`
	Name                  string               `json:"name"`
	Type                  domain.EquipmentType `json:"type"`
	Manufacturer          string               `json:"manufacturer,omitempty"`
	UserID                string               `json:"user_id,omitempty"`
	ActualWeight          float64              `json:"actual_weight,omitempty"`
	PulleyType            string               `json:"pulley_type,omitempty"`
	StackWeights         []float64            `json:"stack_weights,omitempty"`
	WeightIncrement       float64              `json:"weight_increment,omitempty"`
	ResistanceProfileID  string               `json:"resistance_profile_id,omitempty"`
	ResistanceProfileName string              `json:"resistance_profile_name,omitempty"`
	MovementPattern       string               `json:"movement_pattern,omitempty"`
}

func toEquipmentResponse(eq domain.Equipment) EquipmentResponse {
	return EquipmentResponse{
		ID:                    string(eq.ID),
		Name:                  eq.Name,
		Type:                  eq.Type,
		Manufacturer:          eq.Manufacturer,
		UserID:                eq.UserID,
		ActualWeight:          eq.ActualWeight,
		PulleyType:            eq.PulleyType,
		StackWeights:          eq.StackWeights,
		WeightIncrement:       eq.WeightIncrement,
		ResistanceProfileID:  eq.ResistanceProfileID,
		ResistanceProfileName: eq.ResistanceProfileName,
		MovementPattern:       string(eq.MovementPattern),
	}
}

func (h *EquipmentHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
	all, err := h.equipmentService.ListAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]EquipmentResponse, len(all))
	for i, eq := range all {
		responses[i] = EquipmentResponse{
			ID:                    string(eq.ID),
			Name:                  eq.Name,
			Type:                  eq.Type,
			Manufacturer:          eq.Manufacturer,
			UserID:                eq.UserID,
			ActualWeight:          eq.ActualWeight,
			PulleyType:            eq.PulleyType,
			StackWeights:          eq.StackWeights,
			WeightIncrement:       eq.WeightIncrement,
			ResistanceProfileID:  eq.ResistanceProfileID,
			ResistanceProfileName: eq.ResistanceProfileName,
			MovementPattern:       string(eq.MovementPattern),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"equipment": responses,
		"count":     len(responses),
	})
}

func (h *EquipmentHandler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	var req struct {
		Name                  string               `json:"name"`
		Type                  domain.EquipmentType `json:"type"`
		Manufacturer          string               `json:"manufacturer"`
		ActualWeight          float64              `json:"actual_weight"`
		PulleyType            string               `json:"pulley_type"`
		StackWeights         []float64            `json:"stack_weights"`
		WeightIncrement       float64              `json:"weight_increment"`
		ResistanceProfileID  string               `json:"resistance_profile_id"`
		ResistanceProfileName string              `json:"resistance_profile_name"`
		MovementPattern       string               `json:"movement_pattern"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Equipment name is required", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		http.Error(w, "Equipment type is required", http.StatusBadRequest)
		return
	}

	equipment := domain.Equipment{
		Name:                  req.Name,
		Type:                  req.Type,
		Manufacturer:          req.Manufacturer,
		UserID:                userID,
		ActualWeight:          req.ActualWeight,
		PulleyType:            req.PulleyType,
		StackWeights:         req.StackWeights,
		WeightIncrement:       req.WeightIncrement,
		ResistanceProfileID:  req.ResistanceProfileID,
		ResistanceProfileName: req.ResistanceProfileName,
		MovementPattern:       domain.MovementPattern(req.MovementPattern),
	}

	created, err := h.equipmentService.Create(equipment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toEquipmentResponse(created))
}

func (h *EquipmentHandler) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Equipment ID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name                  string               `json:"name"`
		Type                  domain.EquipmentType `json:"type"`
		Manufacturer          string               `json:"manufacturer"`
		ActualWeight          float64              `json:"actual_weight"`
		PulleyType            string               `json:"pulley_type"`
		StackWeights         []float64            `json:"stack_weights"`
		WeightIncrement       float64              `json:"weight_increment"`
		ResistanceProfileID  string               `json:"resistance_profile_id"`
		ResistanceProfileName string              `json:"resistance_profile_name"`
		MovementPattern       string               `json:"movement_pattern"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	equipment := domain.Equipment{
		ID:                    domain.EquipmentID(id),
		Name:                  req.Name,
		Type:                  req.Type,
		Manufacturer:          req.Manufacturer,
		UserID:                userID,
		ActualWeight:          req.ActualWeight,
		PulleyType:            req.PulleyType,
		StackWeights:         req.StackWeights,
		WeightIncrement:       req.WeightIncrement,
		ResistanceProfileID:  req.ResistanceProfileID,
		ResistanceProfileName: req.ResistanceProfileName,
		MovementPattern:       domain.MovementPattern(req.MovementPattern),
	}

	if err := h.equipmentService.Update(equipment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updated, err := h.equipmentService.GetEquipment(domain.EquipmentID(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toEquipmentResponse(updated))
}

func (h *EquipmentHandler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Equipment ID is required", http.StatusBadRequest)
		return
	}

	if err := h.equipmentService.Delete(domain.EquipmentID(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
