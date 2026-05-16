package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
)

type ExerciseHandler struct {
	exerciseService *service.ExerciseService
}

func NewExerciseHandler(svc *service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: svc,
	}
}

// POST /api/v1/exercises/custom
func (h *ExerciseHandler) CreateCustomExercise(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	var req struct {
		Name              string                       `json:"name"`
		Description       string                       `json:"description"`
		PrimaryPatterns   []PatternContributionRequest `json:"primary_patterns"`
		SecondaryPatterns []PatternContributionRequest `json:"secondary_patterns"`
		Equipment         []string                     `json:"equipment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Exercise name is required", http.StatusBadRequest)
		return
	}

	if len(req.PrimaryPatterns) == 0 {
		http.Error(w, "At least one primary pattern is required", http.StatusBadRequest)
		return
	}

	// Convert to service request
	primaryPatterns := make([]domain.MovementPatternContribution, len(req.PrimaryPatterns))
	for i, p := range req.PrimaryPatterns {
		primaryPatterns[i] = domain.MovementPatternContribution{
			Pattern:       domain.MovementPattern(p.Pattern),
			Contribution:  p.Contribution,
			RangeOfMotion: p.RangeOfMotion,
			Notes:         p.Notes,
		}
	}

	secondaryPatterns := make([]domain.MovementPatternContribution, len(req.SecondaryPatterns))
	for i, p := range req.SecondaryPatterns {
		secondaryPatterns[i] = domain.MovementPatternContribution{
			Pattern:       domain.MovementPattern(p.Pattern),
			Contribution:  p.Contribution,
			RangeOfMotion: p.RangeOfMotion,
			Notes:         p.Notes,
		}
	}

	equipment := make([]domain.EquipmentType, len(req.Equipment))
	for i, eq := range req.Equipment {
		equipment[i] = domain.EquipmentType(eq)
	}

	exercise, err := h.exerciseService.CreateCustomExercise(userID, service.CreateExerciseRequest{
		Name:              req.Name,
		Description:       req.Description,
		PrimaryPatterns:   primaryPatterns,
		SecondaryPatterns: secondaryPatterns,
		Equipment:         equipment,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toExerciseResponse(exercise))
}

// POST /api/v1/exercises/from-template
func (h *ExerciseHandler) CreateFromTemplate(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	var req struct {
		TemplateID        string  `json:"template_id"`
		CustomName        *string `json:"custom_name"`
		CustomDescription *string `json:"custom_description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TemplateID == "" {
		http.Error(w, "Template ID is required", http.StatusBadRequest)
		return
	}

	customizations := service.ExerciseCustomizations{}
	if req.CustomName != nil {
		customizations.Name = *req.CustomName
	}
	if req.CustomDescription != nil {
		customizations.Description = *req.CustomDescription
	}

	exercise, err := h.exerciseService.CreateFromTemplate(
		userID,
		req.TemplateID,
		customizations,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toExerciseResponse(exercise))
}

// GET /api/v1/exercises/templates
func (h *ExerciseHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates := make([]TemplateResponse, 0, len(domain.ExerciseTemplates))

	for id, template := range domain.ExerciseTemplates {
		templates = append(templates, TemplateResponse{
			ID:          id,
			Name:        template.Name,
			Description: template.Description,
			Category:    template.Category,
			Patterns:    toPatternResponses(template.PrimaryPatterns),
			Equipment:   template.SuggestedEquipment,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"templates": templates,
		"count":     len(templates),
	})
}

// GET /api/v1/exercises
func (h *ExerciseHandler) ListExercises(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	// Parse query filters
	filters := service.ExerciseFilters{
		Search:      r.URL.Query().Get("search"),
		Pattern:     domain.MovementPattern(r.URL.Query().Get("pattern")),
		Equipment:   domain.EquipmentType(r.URL.Query().Get("equipment")),
		MuscleGroup: domain.MuscleGroup(r.URL.Query().Get("muscle")),
	}

	exercises, err := h.exerciseService.ListExercises(userID, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]ExerciseResponse, len(exercises))
	for i, ex := range exercises {
		responses[i] = toExerciseResponse(ex)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"exercises": responses,
		"count":     len(responses),
	})
}

// GET /api/v1/exercises/{id}
func (h *ExerciseHandler) GetExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID := domain.ExerciseID(r.PathValue("id"))
	userID := GetUserIDFromContext(r.Context())

	if exerciseID == "" {
		http.Error(w, "Exercise ID is required", http.StatusBadRequest)
		return
	}

	exercise, err := h.exerciseService.GetExercise(exerciseID, userID)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toExerciseResponse(exercise))
}

// PUT /api/v1/exercises/{id}
func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	exerciseID := domain.ExerciseID(r.PathValue("id"))

	if exerciseID == "" {
		http.Error(w, "Exercise ID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name              *string                      `json:"name"`
		Description       *string                      `json:"description"`
		PrimaryPatterns   []PatternContributionRequest `json:"primary_patterns,omitempty"`
		SecondaryPatterns []PatternContributionRequest `json:"secondary_patterns,omitempty"`
		Equipment         []string                     `json:"equipment,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateReq := service.UpdateExerciseRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	// Convert patterns if provided
	if req.PrimaryPatterns != nil {
		primaryPatterns := make([]domain.MovementPatternContribution, len(req.PrimaryPatterns))
		for i, p := range req.PrimaryPatterns {
			primaryPatterns[i] = domain.MovementPatternContribution{
				Pattern:       domain.MovementPattern(p.Pattern),
				Contribution:  p.Contribution,
				RangeOfMotion: p.RangeOfMotion,
				Notes:         p.Notes,
			}
		}
		updateReq.PrimaryPatterns = &primaryPatterns
	}

	if req.SecondaryPatterns != nil {
		secondaryPatterns := make([]domain.MovementPatternContribution, len(req.SecondaryPatterns))
		for i, p := range req.SecondaryPatterns {
			secondaryPatterns[i] = domain.MovementPatternContribution{
				Pattern:       domain.MovementPattern(p.Pattern),
				Contribution:  p.Contribution,
				RangeOfMotion: p.RangeOfMotion,
				Notes:         p.Notes,
			}
		}
		updateReq.SecondaryPatterns = &secondaryPatterns
	}

	if req.Equipment != nil {
		equipment := make([]domain.EquipmentType, len(req.Equipment))
		for i, eq := range req.Equipment {
			equipment[i] = domain.EquipmentType(eq)
		}
		updateReq.Equipment = &equipment
	}

	exercise, err := h.exerciseService.UpdateExercise(exerciseID, userID, updateReq)
	if err != nil {
		if err.Error() == "exercise not found" {
			http.Error(w, "Exercise not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: cannot modify system exercises" ||
			err.Error() == "unauthorized: exercise belongs to different user" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toExerciseResponse(exercise))
}

// DELETE /api/v1/exercises/{id}
func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	exerciseID := domain.ExerciseID(r.PathValue("id"))

	if exerciseID == "" {
		http.Error(w, "Exercise ID is required", http.StatusBadRequest)
		return
	}

	err := h.exerciseService.DeleteExercise(exerciseID, userID)
	if err != nil {
		if err.Error() == "exercise not found" {
			http.Error(w, "Exercise not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: cannot delete system exercises" ||
			err.Error() == "unauthorized: exercise belongs to different user" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/patterns
func (h *ExerciseHandler) ListPatterns(w http.ResponseWriter, r *http.Request) {
	patterns := make([]PatternInfoResponse, 0, len(domain.MovementPatternRegistry))

	for pattern, info := range domain.MovementPatternRegistry {
		patterns = append(patterns, PatternInfoResponse{
			Pattern:        string(pattern),
			Name:           info.Name,
			Plane:          string(info.Plane),
			PrimaryMuscles: info.PrimaryMuscles,
			JointActions:   info.JointActions,
			Description:    info.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"patterns": patterns,
		"count":    len(patterns),
	})
}

// Request/Response Types

type PatternContributionRequest struct {
	Pattern       string  `json:"pattern"`
	Contribution  float64 `json:"contribution"`
	RangeOfMotion string  `json:"range_of_motion"`
	Notes         string  `json:"notes"`
}

type ExerciseResponse struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	PrimaryPatterns   []PatternResponse      `json:"primary_patterns"`
	SecondaryPatterns []PatternResponse      `json:"secondary_patterns"`
	DerivedMuscles    DerivedMusclesResponse `json:"derived_muscles"`
	Equipment         []domain.EquipmentType `json:"equipment"`
	Source            string                 `json:"source"`
	IsCustom          bool                   `json:"is_custom"`
	CreatedAt         string                 `json:"created_at"`
	UpdatedAt         string                 `json:"updated_at"`
}

type PatternResponse struct {
	Pattern       string  `json:"pattern"`
	PatternName   string  `json:"pattern_name"`
	Contribution  float64 `json:"contribution"`
	RangeOfMotion string  `json:"range_of_motion"`
	Notes         string  `json:"notes"`
}

type DerivedMusclesResponse struct {
	Primary   []domain.MuscleGroup `json:"primary"`
	Secondary []domain.MuscleGroup `json:"secondary"`
}

type TemplateResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Patterns    []PatternResponse      `json:"patterns"`
	Equipment   []domain.EquipmentType `json:"equipment"`
}

type PatternInfoResponse struct {
	Pattern        string               `json:"pattern"`
	Name           string               `json:"name"`
	Plane          string               `json:"plane"`
	PrimaryMuscles []domain.MuscleGroup `json:"primary_muscles"`
	JointActions   []string             `json:"joint_actions"`
	Description    string               `json:"description"`
}

// Helper functions

func toExerciseResponse(ex domain.Exercise) ExerciseResponse {
	isCustom := ex.UserID != ""
	source := "system"
	if isCustom {
		source = "user"
	}
	return ExerciseResponse{
		ID:                string(ex.ID),
		Name:              ex.Name,
		Description:       ex.Description,
		PrimaryPatterns:   toPatternResponses(ex.PrimaryPatterns),
		SecondaryPatterns: toPatternResponses(ex.SecondaryPatterns),
		DerivedMuscles: DerivedMusclesResponse{
			Primary:   ex.GetPrimaryMuscles(),
			Secondary: ex.GetSecondaryMuscles(),
		},
		Equipment: ex.SuggestedEquipment,
		Source:    source,
		IsCustom:  isCustom,
		CreatedAt: ex.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: ex.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toPatternResponses(patterns []domain.MovementPatternContribution) []PatternResponse {
	responses := make([]PatternResponse, len(patterns))
	for i, p := range patterns {
		patternInfo := domain.MovementPatternRegistry[p.Pattern]
		responses[i] = PatternResponse{
			Pattern:       string(p.Pattern),
			PatternName:   patternInfo.Name,
			Contribution:  p.Contribution,
			RangeOfMotion: p.RangeOfMotion,
			Notes:         p.Notes,
		}
	}
	return responses
}
