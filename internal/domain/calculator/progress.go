package calculator

import (
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"math"
	"sort"
	"time"
)

type SetDataPoint struct {
	Date         time.Time
	SessionID    domain.SessionID
	SetID        domain.PerformedSetID
	SetNumber    int // accounts for volume intra-session
	ExerciseID   domain.ExerciseID
	ExerciseName string
	// Raw performance data
	Reps          int
	EffectiveLoad float64
	RawLoad       float64
	EquipmentID   domain.EquipmentID

	Score float64 // Computed by scoring formula
}

type ScoringFormula interface {
	Calculate(reps int, effectiveLoad float64) float64
	Name() string
	Description() string
}

type ProgressCalculator struct {
	scoringFormula ScoringFormula
}

func NewProgressCalculator(formula ScoringFormula) *ProgressCalculator {
	return &ProgressCalculator{
		scoringFormula: formula,
	}
}

func (p *ProgressCalculator) GenerateSetProgressData(
	exerciseID domain.ExerciseID,
	sessions []domain.Session,
	exercises map[domain.ExerciseID]domain.Exercise,
	startDate, endDate *time.Time,
) []SetDataPoint {

	dataPoints := make([]SetDataPoint, 0)
	exercise, exerciseExists := exercises[exerciseID]

	for _, session := range sessions {
		if session.CompletedAt.IsZero() {
			continue
		}

		if startDate != nil && session.CompletedAt.Before(*startDate) {
			continue
		}

		if endDate != nil && session.CompletedAt.After(*endDate) {
			continue
		}

		for _, set := range session.PerformedSets {
			if set.ExerciseID != exerciseID {
				continue
			}

			score := p.scoringFormula.Calculate(set.Reps, set.EffectiveLoad)

			exerciseName := string(exerciseID)
			if exerciseExists {
				exerciseName = exercise.Name
			}

			dataPoint := SetDataPoint{
				Date:          set.PerformedAt,
				SessionID:     session.ID,
				SetID:         set.ID,
				SetNumber:     set.SetNumber,
				ExerciseID:    set.ExerciseID,
				ExerciseName:  exerciseName,
				Reps:          set.Reps,
				EffectiveLoad: set.EffectiveLoad,
				RawLoad:       set.RawLoad,
				EquipmentID:   set.EquipmentID,
				Score:         score,
			}

			dataPoints = append(dataPoints, dataPoint)
		}
	}

	sort.Slice(dataPoints, func(i, j int) bool {
		return dataPoints[i].Date.Before(dataPoints[j].Date)
	})

	return dataPoints
}

// Progress is individual to set number in between sessions
type ProgressSummary struct {
	ExerciseID   domain.ExerciseID
	ExerciseName string
	DateRange    DateRange

	// Per-set-number summaries
	SetSummaries map[int]SetNumberSummary

	// Overall statistics (across all sets)
	OverallBestScore    float64
	OverallAverageScore float64
	OverallLatestScore  float64
	TotalDataPoints     int
}

type SetNumberSummary struct {
	SetNumber      int
	DataPointCount int

	// Score statistics
	BestScore    float64
	WorstScore   float64
	AverageScore float64
	LatestScore  float64
	FirstScore   float64

	// Trend analysis
	TrendDirection TrendDirection
	TrendStrength  float64 // R^2 value
	Slope          float64 // Score improvement per session

	// Improvement metrics
	TotalImprovement   float64 // Latest - First
	PercentImprovement float64 // (Latest - First) / First * 100
}

type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}

type TrendDirection string

const (
	TrendImproving TrendDirection = "improving"
	TrendDeclining TrendDirection = "declining"
	TrendStable    TrendDirection = "stable"
)

func (p *ProgressCalculator) CalculateProgressSummary(
	dataPoints []SetDataPoint,
) ProgressSummary {

	if len(dataPoints) == 0 {
		return ProgressSummary{
			SetSummaries: make(map[int]SetNumberSummary),
		}
	}

	summary := ProgressSummary{
		ExerciseID:   dataPoints[0].ExerciseID,
		ExerciseName: dataPoints[0].ExerciseName,
		DateRange: DateRange{
			StartDate: dataPoints[0].Date,
			EndDate:   dataPoints[len(dataPoints)-1].Date,
		},
		SetSummaries:    make(map[int]SetNumberSummary),
		TotalDataPoints: len(dataPoints),
	}

	setGroups := p.groupBySetNumber(dataPoints)
	// Calculate summary for each set number
	overallScoreSum := 0.0
	overallBestScore := 0.0
	overallLatestScore := 0.0

	for setNumber, setDataPoints := range setGroups {
		setSummary := p.calculateSetNumberSummary(setNumber, setDataPoints)
		summary.SetSummaries[setNumber] = setSummary

		overallScoreSum += setSummary.AverageScore * float64(setSummary.DataPointCount)

		if setSummary.BestScore > overallBestScore {
			overallBestScore = setSummary.BestScore
		}

		// Latest score from the highest set number
		if setNumber >= len(setGroups) {
			overallLatestScore = setSummary.LatestScore
		}
	}

	summary.OverallBestScore = overallBestScore
	summary.OverallAverageScore = overallScoreSum / float64(len(dataPoints))
	summary.OverallLatestScore = overallLatestScore

	return summary
}

func (p *ProgressCalculator) groupBySetNumber(dataPoints []SetDataPoint) map[int][]SetDataPoint {
	groups := make(map[int][]SetDataPoint)

	for _, dp := range dataPoints {
		if _, exists := groups[dp.SetNumber]; !exists {
			groups[dp.SetNumber] = make([]SetDataPoint, 0)
		}
		groups[dp.SetNumber] = append(groups[dp.SetNumber], dp)
	}

	return groups
}

func (p *ProgressCalculator) calculateSetNumberSummary(
	setNumber int,
	dataPoints []SetDataPoint,
) SetNumberSummary {

	if len(dataPoints) == 0 {
		return SetNumberSummary{SetNumber: setNumber}
	}

	summary := SetNumberSummary{
		SetNumber:      setNumber,
		DataPointCount: len(dataPoints),
		BestScore:      dataPoints[0].Score,
		WorstScore:     dataPoints[0].Score,
		FirstScore:     dataPoints[0].Score,
		LatestScore:    dataPoints[len(dataPoints)-1].Score,
	}

	totalScore := 0.0

	for _, dp := range dataPoints {
		totalScore += dp.Score

		if dp.Score > summary.BestScore {
			summary.BestScore = dp.Score
		}
		if dp.Score < summary.WorstScore {
			summary.WorstScore = dp.Score
		}
	}

	summary.AverageScore = totalScore / float64(len(dataPoints))

	summary.TotalImprovement = summary.LatestScore - summary.FirstScore
	if summary.FirstScore != 0 {
		summary.PercentImprovement = (summary.TotalImprovement / summary.FirstScore) * 100
	}

	// Calculate trend
	trend, strength, slope := p.calculateTrendForSetNumber(dataPoints)
	summary.TrendDirection = trend
	summary.TrendStrength = strength
	summary.Slope = slope

	return summary
}

func (p *ProgressCalculator) calculateTrendForSetNumber(
	dataPoints []SetDataPoint,
) (TrendDirection, float64, float64) {

	if len(dataPoints) < 3 {
		return TrendStable, 0.0, 0.0
	}

	// Linear regression: y = mx + b where y = score, x = index (time proxy)
	n := float64(len(dataPoints))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, dp := range dataPoints {
		x := float64(i)
		y := dp.Score

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n
	// Calculate R² to determine trend strength
	meanY := sumY / n
	ssTotal := 0.0
	ssResidual := 0.0

	for i, dp := range dataPoints {
		x := float64(i)
		predicted := slope*x + intercept
		ssTotal += math.Pow(dp.Score-meanY, 2)
		ssResidual += math.Pow(dp.Score-predicted, 2)
	}

	rSquared := 1 - (ssResidual / ssTotal)
	if math.IsNaN(rSquared) || ssTotal == 0 {
		rSquared = 0
	}

	// Determine direction based on slope
	threshold := 0.01
	var direction TrendDirection

	if slope > threshold {
		direction = TrendImproving
	} else if slope < -threshold {
		direction = TrendDeclining
	} else {
		direction = TrendStable
	}

	return direction, math.Abs(rSquared), slope
}

func (p *ProgressCalculator) GetDataPointsBySetNumber(
	dataPoints []SetDataPoint,
	setNumber int,
) []SetDataPoint {
	filtered := make([]SetDataPoint, 0)

	for _, dp := range dataPoints {
		if dp.SetNumber == setNumber {
			filtered = append(filtered, dp)
		}
	}

	return filtered
}

func (p *ProgressCalculator) GetAvailableSetNumbers(
	dataPoints []SetDataPoint,
) []int {
	setNumberMap := make(map[int]bool)

	for _, dp := range dataPoints {
		setNumberMap[dp.SetNumber] = true
	}

	setNumbers := make([]int, 0, len(setNumberMap))
	for setNumber := range setNumberMap {
		setNumbers = append(setNumbers, setNumber)
	}

	sort.Ints(setNumbers)
	return setNumbers
}
