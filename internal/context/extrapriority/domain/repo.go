package domain

import "context"

// Repo ...
type Repo interface {
	Get(ctx context.Context, id string) (*ExtraPriority, error)
	Save(ctx context.Context, priority *ExtraPriority) error
	Delete(ctx context.Context, id string) error
}
