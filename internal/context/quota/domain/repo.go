package domain

import "context"

// Repo ...
type Repo interface {
	Get(ctx context.Context, id string) (*Quota, error)
	Save(ctx context.Context, quota *Quota) error
	Delete(ctx context.Context, id string) error
}
