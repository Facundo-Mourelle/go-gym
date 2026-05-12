package calculator_test

import (
	"testing"
	"github.com/Facundo-Mourelle/go-gym/internal/domain/calculator"
)

func TestEpley1RM_Calculate(t *testing.T) {
	formula := calculator.Epley1RM{}

	tests := []struct {
		name     string
		weight   float64
		reps     int
		expected float64
	}{
		{"1 rep max", 100, 1, 100},
		{"5 reps", 100, 5, 116.66666666666667},
		{"10 reps", 100, 10, 133.33333333333331},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formula.Calculate(tt.reps, tt.weight)
			// check with a small epsilon for float comparison
			if got < tt.expected-0.001 || got > tt.expected+0.001 {
				t.Errorf("Epley1RM.Calculate() = %v, want %v", got, tt.expected)
			}
		})
	}
}
