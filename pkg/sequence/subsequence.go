package sequence

// Subsequence ...
type Subsequence struct {
	sequence Sequence
	score    float64
}

// Score ...
func (s Subsequence) Score() float64 {
	return s.score
}

func (s Subsequence) String() string {
	return s.sequence.String()
}
