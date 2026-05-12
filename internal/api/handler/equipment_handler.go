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
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Type        domain.EquipmentType `json:"type"`
	Description string             `json:"description,omitempty"`
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
			ID:   string(eq.ID),
			Name: eq.Name,
			Type: eq.Type,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"equipment": responses,
		"count":     len(responses),
	})
}
