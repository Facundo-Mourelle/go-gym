package service

import (
	"fmt"
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
	"slices"
	"strings"
	"time"
)

type ExerciseService struct {
	exerciseRepo repository.ExerciseRepository
}

func NewExerciseService(repo repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{
		exerciseRepo: repo,
	}
}

func (s *ExerciseService) CreateCustomExercise(
	userID string,
	req CreateExerciseRequest,
) (domain.Exercise, error) {

	builder := domain.NewExerciseBuilder().
		WithName(req.Name).
		WithDescription(req.Description).
		WithUserID(userID)

	for _, pattern := range req.PrimaryPatterns {
		builder.WithPrimaryPattern(
			pattern.Pattern,
			pattern.Contribution,
			pattern.RangeOfMotion,
			pattern.Notes,
		)
	}

	for _, pattern := range req.SecondaryPatterns {
		builder.WithSecondaryPattern(
			pattern.Pattern,
			pattern.Contribution,
			pattern.RangeOfMotion,
			pattern.Notes,
		)
	}

	if len(req.Equipment) > 0 {
		builder.WithEquipment(req.Equipment...)
	}

	// Build and validate
	exercise, err := builder.Build()
	if err != nil {
		return domain.Exercise{}, err
	}

	// Save to database
	if err := s.exerciseRepo.Create(exercise); err != nil {
		return domain.Exercise{}, err
	}

	return exercise, nil
}

func (s *ExerciseService) CreateFromTemplate(
	userID string,
	templateID string,
	customizations ExerciseCustomizations,
) (domain.Exercise, error) {

	builder, err := domain.CreateFromTemplate(templateID, userID)
	if err != nil {
		return domain.Exercise{}, err
	}

	if customizations.Name != "" {
		builder.WithName(customizations.Name)
	}

	if customizations.Description != "" {
		builder.WithDescription(customizations.Description)
	}

	exercise, err := builder.Build()
	if err != nil {
		return domain.Exercise{}, err
	}

	if err := s.exerciseRepo.Create(exercise); err != nil {
		return domain.Exercise{}, err
	}

	return exercise, nil
}

func (s *ExerciseService) GetExercise(exerciseID domain.ExerciseID) (domain.Exercise, error) {
	exercise, err := s.exerciseRepo.FindByID(exerciseID)
	if err != nil {
		return domain.Exercise{}, fmt.Errorf("exercise not found")
	}
	return exercise, nil
}

func (s *ExerciseService) ListExercises(
	userID string,
	filters ExerciseFilters,
) ([]domain.Exercise, error) {

	exercises, err := s.exerciseRepo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	// Apply filters
	if filters.Search != "" {
		exercises = filterByName(exercises, filters.Search)
	}

	if filters.Pattern != "" {
		exercises = filterByPattern(exercises, filters.Pattern)
	}

	if filters.Equipment != "" {
		exercises = filterByEquipment(exercises, filters.Equipment)
	}

	if filters.MuscleGroup != "" {
		exercises = filterByMuscle(exercises, filters.MuscleGroup)
	}

	return exercises, nil
}

func (s *ExerciseService) UpdateExercise(
	exerciseID domain.ExerciseID,
	userID string,
	req UpdateExerciseRequest,
) (domain.Exercise, error) {

	exercise, err := s.exerciseRepo.FindByID(exerciseID)
	if err != nil {
		return domain.Exercise{}, fmt.Errorf("exercise not found")
	}

	if exercise.UserID != userID {
		return domain.Exercise{}, fmt.Errorf("unauthorized: exercise belongs to different user")
	}

	if req.Name != nil {
		exercise.Name = *req.Name
	}

	if req.Description != nil {
		exercise.Description = *req.Description
	}

	if req.PrimaryPatterns != nil {
		exercise.PrimaryPatterns = *req.PrimaryPatterns
	}

	if req.SecondaryPatterns != nil {
		exercise.SecondaryPatterns = *req.SecondaryPatterns
	}

	if req.Equipment != nil {
		exercise.SuggestedEquipment = *req.Equipment
	}

	exercise.UpdatedAt = time.Now()

	if err := exercise.Validate(); err != nil {
		return domain.Exercise{}, err
	}

	if err := s.exerciseRepo.Update(exercise); err != nil {
		return domain.Exercise{}, err
	}

	return exercise, nil
}

func (s *ExerciseService) DeleteExercise(
	exerciseID domain.ExerciseID,
	userID string,
) error {

	exercise, err := s.exerciseRepo.FindByID(exerciseID)
	if err != nil {
		return fmt.Errorf("exercise not found")
	}

	if exercise.UserID != userID {
		return fmt.Errorf("unauthorized: exercise belongs to different user")
	}

	return s.exerciseRepo.Delete(exerciseID)
}

type CreateExerciseRequest struct {
	Name              string
	Description       string
	PrimaryPatterns   []domain.MovementPatternContribution
	SecondaryPatterns []domain.MovementPatternContribution
	Equipment         []domain.EquipmentType
}

type UpdateExerciseRequest struct {
	Name              *string
	Description       *string
	PrimaryPatterns   *[]domain.MovementPatternContribution
	SecondaryPatterns *[]domain.MovementPatternContribution
	Equipment         *[]domain.EquipmentType
}

type ExerciseCustomizations struct {
	Name        string
	Description string
}

type ExerciseFilters struct {
	Search      string
	Pattern     domain.MovementPattern
	Equipment   domain.EquipmentType
	MuscleGroup domain.MuscleGroup
}

func filterByPattern(exercises []domain.Exercise, pattern domain.MovementPattern) []domain.Exercise {
	result := make([]domain.Exercise, 0)
	for _, ex := range exercises {
		for _, p := range ex.PrimaryPatterns {
			if p.Pattern == pattern {
				result = append(result, ex)
				break
			}
		}
	}
	return result
}

func filterByEquipment(exercises []domain.Exercise, equipment domain.EquipmentType) []domain.Exercise {
	result := make([]domain.Exercise, 0)
	for _, ex := range exercises {
		if slices.Contains(ex.SuggestedEquipment, equipment) {
			result = append(result, ex)
		}
	}
	return result
}

func filterByMuscle(exercises []domain.Exercise, muscle domain.MuscleGroup) []domain.Exercise {
	result := make([]domain.Exercise, 0)
	for _, ex := range exercises {
		primaryMuscles := ex.GetPrimaryMuscles()
		if slices.Contains(primaryMuscles, muscle) {
			result = append(result, ex)
		}
	}
	return result
}

func filterByName(exercises []domain.Exercise, search string) []domain.Exercise {
	result := make([]domain.Exercise, 0)
	searchLower := strings.ToLower(search)
	for _, ex := range exercises {
		if strings.Contains(strings.ToLower(ex.Name), searchLower) {
			result = append(result, ex)
		}
	}
	return result
}
