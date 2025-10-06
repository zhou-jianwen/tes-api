package domain

import "context"

// Repo ...
type Repo interface {
	Get(ctx context.Context, id string) (*Cluster, error)
	Save(ctx context.Context, cluster *Cluster) error
	Delete(ctx context.Context, id string) error
}
