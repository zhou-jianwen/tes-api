package domain

import (
	"context"
)

// Service ...
type Service interface {
	Create(ctx context.Context, task *Task) (string, error)
	Cancel(ctx context.Context, id string) error
	Update(ctx context.Context, id string, state, clusterID *string, logs []*TaskLog) error
}

type service struct {
	repo       Repo
	normalizer Normalizer
}

var _ Service = (*service)(nil)

// NewService ...
func NewService(repo Repo, normalizer Normalizer) Service {
	return &service{
		repo:       repo,
		normalizer: normalizer,
	}
}

// Create ...
func (s *service) Create(ctx context.Context, task *Task) (string, error) {
	var id string
	for {
		id = GenTaskID()
		exist, err := s.repo.CheckIDExist(ctx, id)
		if err != nil {
			return "", err
		}
		if !exist {
			break
		}
	}
	task.ID = id

	if err := s.normalizer.Normalize(task); err != nil {
		return "", err
	}

	return id, s.repo.Create(ctx, task)
}

// Cancel ...
func (s *service) Cancel(ctx context.Context, id string) error {
	for {
		updated, err := s.cancel(ctx, id)
		if err != nil {
			return err
		}
		if updated {
			return nil
		}
	}
}

func (s *service) cancel(ctx context.Context, id string) (bool, error) {
	taskStatus, err := s.repo.GetStatus(ctx, id)
	if err != nil {
		return false, err
	}
	if err = taskStatus.Cancel(); err != nil {
		return false, err
	}
	return s.repo.UpdateStatus(ctx, taskStatus)
}

// Update ...
func (s *service) Update(ctx context.Context, id string, state, clusterID *string, logs []*TaskLog) error {
	for {
		updated, err := s.update(ctx, id, state, clusterID, logs)
		if err != nil {
			return err
		}
		if updated {
			return nil
		}
	}
}

func (s *service) update(ctx context.Context, id string, state, clusterID *string, logs []*TaskLog) (bool, error) {
	taskStatus, err := s.repo.GetStatus(ctx, id)
	if err != nil {
		return false, err
	}

	if state != nil {
		if err = taskStatus.UpdateState(*state); err != nil {
			return false, err
		}
	}

	if clusterID != nil {
		if err = taskStatus.UpdateClusterID(*clusterID); err != nil {
			return false, err
		}
	}

	if len(logs) > 0 {
		if err = taskStatus.UpdateLogs(logs); err != nil {
			return false, err
		}
	}

	return s.repo.UpdateStatus(ctx, taskStatus)
}
