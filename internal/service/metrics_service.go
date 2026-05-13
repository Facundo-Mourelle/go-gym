package service

import (
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/domain/calculator"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type MetricsService struct {
	sessionRepo    repository.SessionRepository
	exerciseRepo   repository.ExerciseRepository
	equipmentRepo  repository.EquipmentRepository
	volumeCalc     *calculator.VolumeCalculator
	progressCalc   *calculator.ProgressCalculator
}

type VolumeMetricsResponse struct {
	TotalSets       int
	MuscleBreakdown []MuscleVolumeBreakdown
}

type MuscleVolumeBreakdown struct {
	MuscleGroup domain.MuscleGroup
	Volume      float64
}

func NewMetricsService(
	sessionRepo repository.SessionRepository,
	exerciseRepo repository.ExerciseRepository,
	equipmentRepo repository.EquipmentRepository,
	volumeCalc *calculator.VolumeCalculator,
	progressCalc *calculator.ProgressCalculator,
) *MetricsService {
	return &MetricsService{
		sessionRepo:   sessionRepo,
		exerciseRepo:  exerciseRepo,
		equipmentRepo: equipmentRepo,
		volumeCalc:    volumeCalc,
		progressCalc:  progressCalc,
	}
}

func (m *MetricsService) GetVolumeMetrics(
	userID string,
	startDate, endDate time.Time,
) (VolumeMetricsResponse, error) {

	// Fetch all sessions with performed sets in date range
	sessions, err := m.sessionRepo.FindByUserAndDateRange(userID, startDate, endDate)
	if err != nil {
		return VolumeMetricsResponse{}, err
	}

	exercises, err := m.exerciseRepo.FindAll()
	if err != nil {
		return VolumeMetricsResponse{}, err
	}

	exerciseMap := make(map[domain.ExerciseID]domain.Exercise)
	for _, ex := range exercises {
		exerciseMap[ex.ID] = ex
	}

	totalSets := 0
	muscleVolumes := make(map[domain.MuscleGroup]float64)

	for _, session := range sessions {
		if session.CompletedAt.IsZero() {
			continue // Skip incomplete sessions
		}

		totalSets += session.SessionVolume()

		mgVolumes := m.volumeCalc.CalculateMuscleGroupVolume(session, exerciseMap)
		for muscle, volume := range mgVolumes {
			muscleVolumes[muscle] += volume
		}
	}

	// Convert to response format
	muscleBreakdown := make([]MuscleVolumeBreakdown, 0, len(muscleVolumes))
	for muscle, volume := range muscleVolumes {
		muscleBreakdown = append(muscleBreakdown, MuscleVolumeBreakdown{
			MuscleGroup: muscle,
			Volume:      volume,
		})
	}

	return VolumeMetricsResponse{
		TotalSets:       totalSets,
		MuscleBreakdown: muscleBreakdown,
	}, nil
}

func (m *MetricsService) GetExerciseProgress(
	userID string,
	exerciseID domain.ExerciseID,
	startDate, endDate *time.Time,
) (ExerciseProgressResponse, error) {

	sessions, err := m.sessionRepo.FindByUser(userID)
	if err != nil {
		return ExerciseProgressResponse{}, err
	}

	exercises, err := m.exerciseRepo.FindAll()
	if err != nil {
		return ExerciseProgressResponse{}, err
	}

	equipmentList, err := m.equipmentRepo.FindAll()
	if err != nil {
		return ExerciseProgressResponse{}, err
	}

	equipmentMap := make(map[domain.EquipmentID]domain.EquipmentType, len(equipmentList))
	for _, eq := range equipmentList {
		equipmentMap[eq.ID] = eq.Type
	}

	exerciseMap := make(map[domain.ExerciseID]domain.Exercise)
	for _, ex := range exercises {
		exerciseMap[ex.ID] = ex
	}

	formula := calculator.Epley1RM{}
	progressCalc := calculator.NewProgressCalculator(formula)
	dataPoints := progressCalc.GenerateSetProgressData(
		exerciseID,
		sessions,
		exerciseMap,
		startDate,
		endDate,
	)

	summary := progressCalc.CalculateProgressSummary(dataPoints)

	// Group data by set number for frontend
	dataBySetNumber := make(map[int][]SetDataPointResponse)
	availableSetNumbers := progressCalc.GetAvailableSetNumbers(dataPoints)

	for _, setNumber := range availableSetNumbers {
		setData := progressCalc.GetDataPointsBySetNumber(dataPoints, setNumber)
		dataBySetNumber[setNumber] = toSetDataPointResponses(setData, equipmentMap)
	}
	return ExerciseProgressResponse{
		ExerciseID:      exerciseID,
		ExerciseName:    summary.ExerciseName,
		AllDataPoints:   toSetDataPointResponses(dataPoints, equipmentMap),
		DataBySetNumber: dataBySetNumber,
		Summary:         toProgressSummaryResponse(summary),
	}, nil
}

func toSetDataPointResponses(dataPoints []calculator.SetDataPoint, equipmentMap map[domain.EquipmentID]domain.EquipmentType) []SetDataPointResponse {
	responses := make([]SetDataPointResponse, len(dataPoints))
	for i, dp := range dataPoints {
		eqType := equipmentMap[dp.EquipmentID]
		responses[i] = SetDataPointResponse{
			Date:          dp.Date,
			SessionID:     string(dp.SessionID),
			SetID:         string(dp.SetID),
			SetNumber:     dp.SetNumber,
			Score:         dp.Score,
			Reps:          dp.Reps,
			EffectiveLoad: dp.EffectiveLoad,
			RawLoad:       dp.RawLoad,
			RepsInReserve: dp.RepsInReserve,
			EquipmentType: string(eqType),
		}
	}
	return responses
}

func toProgressSummaryResponse(summary calculator.ProgressSummary) ProgressSummaryResponse {
	setSummaries := make(map[int]SetNumberSummaryResponse)

	for setNum, setSummary := range summary.SetSummaries {
		setSummaries[setNum] = SetNumberSummaryResponse{
			SetNumber:          setSummary.SetNumber,
			DataPointCount:     setSummary.DataPointCount,
			BestScore:          setSummary.BestScore,
			WorstScore:         setSummary.WorstScore,
			AverageScore:       setSummary.AverageScore,
			LatestScore:        setSummary.LatestScore,
			FirstScore:         setSummary.FirstScore,
			TrendDirection:     string(setSummary.TrendDirection),
			TrendStrength:      setSummary.TrendStrength,
			Slope:              setSummary.Slope,
			TotalImprovement:   setSummary.TotalImprovement,
			PercentImprovement: setSummary.PercentImprovement,
		}
	}
	return ProgressSummaryResponse{
		ExerciseID:   string(summary.ExerciseID),
		ExerciseName: summary.ExerciseName,
		DateRange: DateRangeResponse{
			StartDate: summary.DateRange.StartDate,
			EndDate:   summary.DateRange.EndDate,
		},
		SetSummaries:        setSummaries,
		OverallBestScore:    summary.OverallBestScore,
		OverallAverageScore: summary.OverallAverageScore,
		OverallLatestScore:  summary.OverallLatestScore,
		TotalDataPoints:     summary.TotalDataPoints,
	}
}

type ExerciseProgressResponse struct {
	ExerciseID   domain.ExerciseID `json:"exercise_id"`
	ExerciseName string            `json:"exercise_name"`

	AllDataPoints []SetDataPointResponse `json:"all_data_points"`

	DataBySetNumber map[int][]SetDataPointResponse `json:"data_by_set_number"`

	Summary ProgressSummaryResponse `json:"summary"`
}

type SetDataPointResponse struct {
	Date          time.Time `json:"date"`
	SessionID     string    `json:"session_id"`
	SetID         string    `json:"set_id"`
	SetNumber     int       `json:"set_number"`
	Reps          int       `json:"reps"`
	Score         float64   `json:"score"`
	EffectiveLoad float64   `json:"effective_load"`
	RawLoad       float64   `json:"weight"`
	RepsInReserve int       `json:"rir"`
	EquipmentType string    `json:"equipment_type"`
}

type ProgressSummaryResponse struct {
	ExerciseID          string                           `json:"exercise_id"`
	ExerciseName        string                           `json:"exercise_name"`
	DateRange           DateRangeResponse                `json:"date_range"`
	SetSummaries        map[int]SetNumberSummaryResponse `json:"set_summaries"`
	OverallBestScore    float64                          `json:"overall_best_score"`
	OverallAverageScore float64                          `json:"overall_average_score"`
	OverallLatestScore  float64                          `json:"overall_latest_score"`
	TotalDataPoints     int                              `json:"total_data_points"`
}

type SetNumberSummaryResponse struct {
	SetNumber          int     `json:"set_number"`
	DataPointCount     int     `json:"data_point_count"`
	BestScore          float64 `json:"best_score"`
	WorstScore         float64 `json:"worst_score"`
	AverageScore       float64 `json:"average_score"`
	LatestScore        float64 `json:"latest_score"`
	FirstScore         float64 `json:"first_score"`
	TrendDirection     string  `json:"trend_direction"`
	TrendStrength      float64 `json:"trend_strength"`
	Slope              float64 `json:"slope"`
	TotalImprovement   float64 `json:"total_improvement"`
	PercentImprovement float64 `json:"percent_improvement"`
}

type DateRangeResponse struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type SessionProgressResponse struct {
	SessionID    string    `json:"session_id"`
	Date         time.Time `json:"date"`
	SetCount     int       `json:"set_count"`
	AverageScore float64   `json:"average_score"`
	BestSetScore float64   `json:"best_set_score"`
}
