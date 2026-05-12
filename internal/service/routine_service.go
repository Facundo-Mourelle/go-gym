package service

import (
	"fmt"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type CreateRoutineRequest struct {
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	MovementPatterns []domain.MovementPattern `json:"movement_patterns"`
}

type RoutineResponse struct {
	ID               string              `json:"id"`
	UserID           string              `json:"user_id"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	MovementPatterns []domain.MovementPattern `json:"movement_patterns"`
	IsPreset         bool                `json:"is_preset"`
	CreatedAt        string              `json:"created_at"`
	UpdatedAt        string              `json:"updated_at"`
}

type RoutineSummaryResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	PatternCount int    `json:"pattern_count"`
	IsPreset     bool   `json:"is_preset"`
}

type RoutineService struct {
	routineRepo repository.RoutineRepository
}

func NewRoutineService(routineRepo repository.RoutineRepository) *RoutineService {
	return &RoutineService{
		routineRepo: routineRepo,
	}
}

func (s *RoutineService) CreateRoutine(userID string, req CreateRoutineRequest) (RoutineResponse, error) {
	now := time.Now()
	routine := domain.Routine{
		ID:               domain.RoutineID(generateID()),
		UserID:           userID,
		Name:             req.Name,
		Description:      req.Description,
		MovementPatterns: req.MovementPatterns,
		IsPreset:         false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.routineRepo.Create(routine); err != nil {
		return RoutineResponse{}, fmt.Errorf("failed to create routine: %w", err)
	}

	return s.toRoutineResponse(routine), nil
}

func (s *RoutineService) GetRoutine(id domain.RoutineID, userID string) (RoutineResponse, error) {
		routine, err := s.routineRepo.FindByID(id)
	if err != nil {
		return RoutineResponse{}, fmt.Errorf("routine not found: %w", err)
	}

	if routine.UserID != userID {
		return RoutineResponse{}, fmt.Errorf("unauthorized: routine belongs to different user")
	}

	return s.toRoutineResponse(routine), nil
}

func (s *RoutineService) ListRoutines(userID string) ([]RoutineSummaryResponse, error) {
	routines, err := s.routineRepo.FindByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routines: %w", err)
	}

	responses := make([]RoutineSummaryResponse, len(routines))
	for i, routine := range routines {
		responses[i] = RoutineSummaryResponse{
			ID:           string(routine.ID),
			Name:         routine.Name,
			Description:  routine.Description,
			PatternCount: len(routine.MovementPatterns),
			IsPreset:     routine.IsPreset,
		}
	}

	return responses, nil
}

func (s *RoutineService) DeleteRoutine(id domain.RoutineID, userID string) error {
	routine, err := s.routineRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("routine not found: %w", err)
	}

	if routine.UserID != userID {
		return fmt.Errorf("unauthorized: routine belongs to different user")
	}

	return s.routineRepo.Delete(id)
}

func (s *RoutineService) SeedStarterRoutines(userID string) error {
	now := time.Now()

	routines := []domain.Routine{
		{
			ID:          domain.RoutineID(generateID()),
			UserID:      userID,
			Name:        "Push Day",
			Description: "Upper body pushing movements - chest, shoulders, triceps",
			MovementPatterns: []domain.MovementPattern{
				domain.HorizontalPush,
				domain.VerticalPush,
				domain.ShoulderAbduction,
			},
			IsPreset:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          domain.RoutineID(generateID()),
			UserID:      userID,
			Name:        "Pull Day",
			Description: "Upper body pulling movements - back, biceps",
			MovementPatterns: []domain.MovementPattern{
				domain.HorizontalPull,
				domain.VerticalPull,
				domain.ShoulderFlexion,
				domain.ShoulderExtension,
			},
			IsPreset:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          domain.RoutineID(generateID()),
			UserID:      userID,
			Name:        "Legs Day",
			Description: "Lower body movements - quads, glutes, hamstrings",
			MovementPatterns: []domain.MovementPattern{
				domain.SquatPattern,
				domain.HipHinge,
				domain.KneeFixedExtension,
				domain.HipAdduction,
			},
			IsPreset:  true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, routine := range routines {
		if err := s.routineRepo.Create(routine); err != nil {
			return err
		}
	}

	return nil
}

func (s *RoutineService) toRoutineResponse(routine domain.Routine) RoutineResponse {
	return RoutineResponse{
		ID:               string(routine.ID),
		UserID:           routine.UserID,
		Name:             routine.Name,
		Description:      routine.Description,
		MovementPatterns: routine.MovementPatterns,
		IsPreset:         routine.IsPreset,
		CreatedAt:        routine.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        routine.UpdatedAt.Format(time.RFC3339),
	}
}