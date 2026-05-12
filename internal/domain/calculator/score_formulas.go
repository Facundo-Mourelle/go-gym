package calculator


type Epley1RM struct{}

func (e Epley1RM) Calculate(reps int, weight float64) float64 {
	if reps == 1 {
		return weight
	}
	return weight * (1.0 + float64(reps)/30.0)
}

func (e Epley1RM) Name() string {
	return "Epley 1RM"
}

func (e Epley1RM) Description() string {
	return "Standard Epley formula"
}
