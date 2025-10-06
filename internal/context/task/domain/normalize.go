package domain

// Normalizer ...
type Normalizer interface {
	Normalize(task *Task) error
}
