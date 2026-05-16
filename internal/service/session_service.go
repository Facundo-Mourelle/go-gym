package service

import (
	"fmt"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/domain/resistance"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type SessionService struct {
	sessionRepo     repository.SessionRepository
	exerciseRepo    repository.ExerciseRepository
	equipmentRepo   repository.EquipmentRepository
	workoutRepo     repository.WorkoutRepository
	profileRegistry *resistance.Registry
}

func NewSessionService(
	sessionRepo repository.SessionRepository,
	exerciseRepo repository.ExerciseRepository,
	equipmentRepo repository.EquipmentRepository,
	workoutRepo repository.WorkoutRepository,
	profileRegistry *resistance.Registry,
) *SessionService {
	return &SessionService{
		sessionRepo:     sessionRepo,
		exerciseRepo:    exerciseRepo,
		equipmentRepo:   equipmentRepo,
		workoutRepo:     workoutRepo,
		profileRegistry: profileRegistry,
	}
}

func (s *SessionService) StartSession(
	userID string,
	req StartSessionRequest,
) (StartSessionResponse, error) {

	if req.WorkoutPlanID != nil {
		workout, err := s.workoutRepo.FindByID(*req.WorkoutPlanID)
		if err != nil {
			return StartSessionResponse{}, fmt.Errorf("workout %v does not exist", req.WorkoutPlanID)
		}
		if workout.UserID != userID {
			return StartSessionResponse{}, fmt.Errorf("workout %v does not belong to user %v", req.WorkoutPlanID, userID)
		}
	}

	workoutPlanID := domain.WorkoutPlanID("")
	if req.WorkoutPlanID != nil {
		workoutPlanID = *req.WorkoutPlanID
	}

	session := domain.Session{
		ID:            domain.SessionID(generateID()),
		UserID:        userID,
		WorkoutPlanID: workoutPlanID,
		StartedAt:     time.Now(),
		CompletedAt:   time.Time{},
		PerformedSets: make([]domain.PerformedSet, 0),
		Notes:         req.Notes,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return StartSessionResponse{}, fmt.Errorf("failed to create session: %w", err)
	}

	return StartSessionResponse{
		SessionID:     session.ID,
		StartedAt:     session.StartedAt,
		WorkoutPlanID: &session.WorkoutPlanID,
	}, nil
}

func (s *SessionService) RecordSet(
	sessionID domain.SessionID,
	req RecordSetRequest,
) (RecordSetResponse, error) {

	_, err := s.exerciseRepo.FindByID(req.ExerciseID)
	if err != nil {
		return RecordSetResponse{}, fmt.Errorf("exercise not found: %w", err)
	}

	equipment, err := s.equipmentRepo.FindByID(req.EquipmentID)
	if err != nil {
		return RecordSetResponse{}, fmt.Errorf("equipment not found: %w", err)
	}

	profileID := equipment.ResistanceProfileID
	if profileID == "" {
		profileID = "free_weight"
	}
	profile, err := s.profileRegistry.Get(profileID)
	if err != nil {
		return RecordSetResponse{}, fmt.Errorf("resistance profile not found: %w", err)
	}

	effectiveLoad, err := profile.CalculateEffectiveLoad(req.RawLoad)
	if err != nil {
		return RecordSetResponse{}, fmt.Errorf("failed to calculate effective load: %w", err)
	}

	workoutExerciseID := domain.WorkoutExerciseID("")
	if req.WorkoutExerciseID != nil {
		workoutExerciseID = *req.WorkoutExerciseID
	}

	performedSet := domain.PerformedSet{
		ID:                domain.PerformedSetID(generateID()),
		WorkoutExerciseID: workoutExerciseID,
		ExerciseID:        req.ExerciseID,
		SetNumber:         req.SetNumber,
		Reps:              req.Reps,
		RepsInReserve:     req.RepsInReserve,
		RawLoad:           req.RawLoad,
		EquipmentID:       req.EquipmentID,
		EffectiveLoad:     effectiveLoad,
		PerformedAt:       time.Now(),
	}

	if err := s.sessionRepo.AddPerformedSet(sessionID, performedSet); err != nil {
		return RecordSetResponse{}, fmt.Errorf("failed to save set: %w", err)
	}

	volume := float64(req.Reps) * effectiveLoad

	return RecordSetResponse{
		PerformedSetID: performedSet.ID,
		EffectiveLoad:  effectiveLoad,
		Volume:         volume,
		PerformedAt:    performedSet.PerformedAt,
	}, nil
}

func (s *SessionService) UpdateSet(
	sessionID domain.SessionID,
	setID domain.PerformedSetID,
	req UpdateSetRequest,
) (RecordSetResponse, error) {

	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return RecordSetResponse{}, fmt.Errorf("session not found: %w", err)
	}

	if !session.CompletedAt.IsZero() {
		return RecordSetResponse{}, fmt.Errorf("cannot update sets in completed session")
	}

	setIndex := -1
	for i, set := range session.PerformedSets {
		if set.ID == setID {
			setIndex = i
			break
		}
	}

	if setIndex == -1 {
		return RecordSetResponse{}, fmt.Errorf("set not found in session")
	}

	set := session.PerformedSets[setIndex]

	if req.RawLoad != nil || req.EquipmentID != nil {
		equipmentID := set.EquipmentID
		if req.EquipmentID != nil {
			equipmentID = *req.EquipmentID
		}

		equipment, err := s.equipmentRepo.FindByID(equipmentID)
		if err != nil {
			return RecordSetResponse{}, fmt.Errorf("equipment not found: %w", err)
		}

		profileID := equipment.ResistanceProfileID
		if profileID == "" {
			profileID = "free_weight"
		}
		profile, err := s.profileRegistry.Get(profileID)
		if err != nil {
			return RecordSetResponse{}, fmt.Errorf("resistance profile not found: %w", err)
		}

		rawLoad := set.RawLoad
		if req.RawLoad != nil {
			rawLoad = *req.RawLoad
		}

		effectiveLoad, err := profile.CalculateEffectiveLoad(rawLoad)
		if err != nil {
			return RecordSetResponse{}, fmt.Errorf("failed to calculate effective load: %w", err)
		}

		set.RawLoad = rawLoad
		set.EquipmentID = equipmentID
		set.EffectiveLoad = effectiveLoad
	}

	if req.Reps != nil {
		set.Reps = *req.Reps
	}

	if req.RepsInReserve != nil {
		set.RepsInReserve = *req.RepsInReserve
	}

	// Update the set in the session
	session.PerformedSets[setIndex] = set

	if err := s.sessionRepo.Update(session); err != nil {
		return RecordSetResponse{}, fmt.Errorf("failed to update session: %w", err)
	}

	volume := float64(set.Reps) * set.EffectiveLoad

	return RecordSetResponse{
		PerformedSetID: set.ID,
		EffectiveLoad:  set.EffectiveLoad,
		Volume:         volume,
		PerformedAt:    set.PerformedAt,
	}, nil
}

func (s *SessionService) DeleteSet(
	sessionID domain.SessionID,
	setID domain.PerformedSetID,
) error {

	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	if !session.CompletedAt.IsZero() {
		return fmt.Errorf("cannot delete sets from completed session")
	}

	newSets := make([]domain.PerformedSet, 0, len(session.PerformedSets)-1)
	found := false

	for _, set := range session.PerformedSets {
		if set.ID == setID {
			found = true
			continue
		}
		newSets = append(newSets, set)
	}

	if !found {
		return fmt.Errorf("set not found in session")
	}

	session.PerformedSets = newSets

	return s.sessionRepo.Update(session)
}

func (s *SessionService) CompleteSession(
	sessionID domain.SessionID,
	notes string,
) (CompleteSessionResponse, error) {

	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return CompleteSessionResponse{}, fmt.Errorf("session not found: %w", err)
	}

	if !session.CompletedAt.IsZero() {
		return CompleteSessionResponse{}, fmt.Errorf("session already completed")
	}

	now := time.Now()

	if err := s.sessionRepo.Complete(sessionID, now); err != nil {
		return CompleteSessionResponse{}, fmt.Errorf("failed to complete session: %w", err)
	}

	if notes != "" {
		session.Notes = notes
		if err := s.sessionRepo.Update(session); err != nil {
			return CompleteSessionResponse{}, fmt.Errorf("failed to update notes: %w", err)
		}
	}

	totalSets := session.SessionVolume()

	duration := now.Sub(session.StartedAt)

	return CompleteSessionResponse{
		SessionID:   sessionID,
		CompletedAt: now,
		Duration:    duration,
		TotalSets:   totalSets,
	}, nil
}

func (s *SessionService) GetSession(
	sessionID domain.SessionID,
	userID string,
) (SessionDetailResponse, error) {

	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return SessionDetailResponse{}, fmt.Errorf("session not found: %w", err)
	}

	if session.UserID != userID {
		return SessionDetailResponse{}, fmt.Errorf("unauthorized")
	}

	exerciseMap := make(map[domain.ExerciseID]domain.Exercise)
	for _, set := range session.PerformedSets {
		if _, exists := exerciseMap[set.ExerciseID]; !exists {
			exercise, err := s.exerciseRepo.FindByID(set.ExerciseID)
			if err != nil {
				// Exercise might have been deleted, skip
				continue
			}
			exerciseMap[set.ExerciseID] = exercise
		}
	}

	exerciseGroups := s.groupSetsByExercise(session.PerformedSets, exerciseMap)

	totalSets := session.SessionVolume()

	var duration *time.Duration
	var completedAt *time.Time
	if !session.CompletedAt.IsZero() {
		d := session.CompletedAt.Sub(session.StartedAt)
		duration = &d
		t := session.CompletedAt
		completedAt = &t
	}

	var wpID *domain.WorkoutPlanID
	if session.WorkoutPlanID != "" {
		id := session.WorkoutPlanID
		wpID = &id
	}

	var totalVolume float64
	for _, g := range exerciseGroups {
		totalVolume += g.TotalVolume
	}

	return SessionDetailResponse{
		SessionID:      session.ID,
		WorkoutPlanID:  wpID,
		StartedAt:      session.StartedAt,
		CompletedAt:    completedAt,
		Duration:       duration,
		Notes:          session.Notes,
		ExerciseGroups: exerciseGroups,
		TotalSets:      totalSets,
		TotalVolume:    totalVolume,
	}, nil
}

func (s *SessionService) groupSetsByExercise(
	sets []domain.PerformedSet,
	exercises map[domain.ExerciseID]domain.Exercise,
) []ExerciseGroupResponse {

	groups := make(map[domain.ExerciseID][]domain.PerformedSet)
	exerciseOrder := make([]domain.ExerciseID, 0)

	for _, set := range sets {
		if _, exists := groups[set.ExerciseID]; !exists {
			exerciseOrder = append(exerciseOrder, set.ExerciseID)
			groups[set.ExerciseID] = make([]domain.PerformedSet, 0)
		}
		groups[set.ExerciseID] = append(groups[set.ExerciseID], set)
	}

	// Convert to response format
	result := make([]ExerciseGroupResponse, 0, len(groups))

	for _, exerciseID := range exerciseOrder {
		sets := groups[exerciseID]
		exercise, exists := exercises[exerciseID]

		exerciseName := string(exerciseID)
		if exists {
			exerciseName = exercise.Name
		}

		var totalVolume float64
		for _, s := range sets {
			totalVolume += float64(s.Reps) * s.EffectiveLoad
		}

		setResponses := make([]PerformedSetResponse, len(sets))
		for i, set := range sets {
			setVolume := float64(set.Reps) * set.EffectiveLoad
		setResponses[i] = PerformedSetResponse{
				SetID:         set.ID,
				SetNumber:     set.SetNumber,
				Reps:          set.Reps,
				RepsInReserve: set.RepsInReserve,
				RawLoad:       set.RawLoad,
				EffectiveLoad: set.EffectiveLoad,
				Volume:        setVolume,
				EquipmentID:   set.EquipmentID,
				PerformedAt:   set.PerformedAt,
			}
		}

		result = append(result, ExerciseGroupResponse{
			ExerciseID:   exerciseID,
			ExerciseName: exerciseName,
			Sets:         setResponses,
			TotalVolume:  totalVolume,
			SetCount:     len(sets),
		})
	}

	return result
}

func (s *SessionService) ListUserSessions(
	userID string,
	limit int,
) ([]SessionSummaryResponse, error) {

	sessions, err := s.sessionRepo.FindByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sessions: %w", err)
	}

	// Sort by most recent first (already sorted by repo, but double check)
	if limit > 0 && len(sessions) > limit {
		sessions = sessions[:limit]
	}

	responses := make([]SessionSummaryResponse, len(sessions))

	for i, session := range sessions {
		var duration *time.Duration
		var completedAt *time.Time
		if !session.CompletedAt.IsZero() {
			d := session.CompletedAt.Sub(session.StartedAt)
			duration = &d
			t := session.CompletedAt
			completedAt = &t
		}

		var wpID *domain.WorkoutPlanID
		if session.WorkoutPlanID != "" {
			id := session.WorkoutPlanID
			wpID = &id
		}

		responses[i] = SessionSummaryResponse{
			SessionID:     session.ID,
			WorkoutPlanID: wpID,
			StartedAt:     session.StartedAt,
			CompletedAt:   completedAt,
			Duration:      duration,
			TotalSets:     len(session.PerformedSets),
		}
	}

	return responses, nil
}

type StartSessionRequest struct {
	WorkoutPlanID *domain.WorkoutPlanID `json:"workout_plan_id,omitempty"`
	Notes         string                `json:"notes"`
}

type StartSessionResponse struct {
	SessionID     domain.SessionID      `json:"session_id"`
	StartedAt     time.Time             `json:"started_at"`
	WorkoutPlanID *domain.WorkoutPlanID `json:"workout_plan_id,omitempty"`
}

type RecordSetRequest struct {
	WorkoutExerciseID *domain.WorkoutExerciseID
	ExerciseID        domain.ExerciseID
	SetNumber         int
	Reps              int
	RepsInReserve     int
	RawLoad           float64
	EquipmentID       domain.EquipmentID
	RestSeconds       int
	Notes             string
}

type RecordSetResponse struct {
	PerformedSetID domain.PerformedSetID `json:"performed_set_id"`
	EffectiveLoad  float64              `json:"effective_load"`
	Volume         float64              `json:"volume"`
	PerformedAt    time.Time            `json:"performed_at"`
}

type UpdateSetRequest struct {
	Reps          *int
	RawLoad       *float64
	EquipmentID   *domain.EquipmentID
	RestSeconds   *int
	RepsInReserve *int
	Notes         *string
}

type CompleteSessionResponse struct {
	SessionID   domain.SessionID `json:"session_id"`
	CompletedAt time.Time        `json:"completed_at"`
	Duration    time.Duration    `json:"duration"`
	TotalSets   int              `json:"total_sets"`
}

type SessionDetailResponse struct {
	SessionID      domain.SessionID        `json:"session_id"`
	WorkoutPlanID  *domain.WorkoutPlanID   `json:"workout_plan_id,omitempty"`
	StartedAt      time.Time               `json:"started_at"`
	CompletedAt    *time.Time              `json:"completed_at"`
	Duration       *time.Duration          `json:"duration"`
	Notes          string                  `json:"notes"`
	ExerciseGroups []ExerciseGroupResponse `json:"exercise_groups"`
	TotalSets      int                     `json:"total_sets"`
	TotalVolume    float64                 `json:"total_volume"`
}

type ExerciseGroupResponse struct {
	ExerciseID   domain.ExerciseID      `json:"exercise_id"`
	ExerciseName string                 `json:"exercise_name"`
	Sets         []PerformedSetResponse `json:"sets"`
	TotalVolume  float64                `json:"total_volume"`
	SetCount     int                    `json:"set_count"`
}

type PerformedSetResponse struct {
	SetID             domain.PerformedSetID `json:"set_id"`
	WorkoutExerciseID domain.WorkoutExerciseID `json:"workout_exercise_id,omitempty"`
	ExerciseID        domain.ExerciseID        `json:"exercise_id"`
	SetNumber         int                      `json:"set_number"`
	Reps              int                      `json:"reps"`
	RepsInReserve     int                      `json:"reps_in_reserve"`
	RawLoad           float64                  `json:"raw_load"`
	EquipmentID       domain.EquipmentID       `json:"equipment_id"`
	EffectiveLoad     float64                  `json:"effective_load"`
	Volume            float64                  `json:"volume"`
	PerformedAt       time.Time                `json:"performed_at"`
}

type SessionSummaryResponse struct {
	SessionID     domain.SessionID     `json:"session_id"`
	WorkoutPlanID *domain.WorkoutPlanID `json:"workout_plan_id,omitempty"`
	StartedAt     time.Time            `json:"started_at"`
	CompletedAt   *time.Time           `json:"completed_at"`
	Duration      *time.Duration       `json:"duration"`
	TotalSets     int                  `json:"total_sets"`
}

// Helper function to generate IDs
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
