package domain

import (
	"context"
)

// Service ...
type Service interface {
	Put(ctx context.Context, cluster *Cluster) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repo
}

var _ Service = (*service)(nil)

// NewService ...
func NewService(repo Repo) Service {
	return &service{repo: repo}
}

// Put ...
func (s *service) Put(ctx context.Context, cluster *Cluster) error {
	return s.repo.Save(ctx, cluster)
}

// Delete ...
func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.Get(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
