package domain

import (
	"context"
)

// Service ...
type Service interface {
	Put(ctx context.Context, accountID, userID, submissionID, runID string, extraPriorityValue int) error
	Delete(ctx context.Context, accountID, userID, submissionID, runID string) error
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
func (s *service) Put(ctx context.Context, accountID, userID, submissionID, runID string, extraPriorityValue int) error {
	id, err := NewExtraPriorityID(accountID, userID, submissionID, runID)
	if err != nil {
		return err
	}
	priority := &ExtraPriority{
		ID:                 id,
		AccountID:          accountID,
		UserID:             userID,
		SubmissionID:       submissionID,
		RunID:              runID,
		ExtraPriorityValue: extraPriorityValue,
	}
	return s.repo.Save(ctx, priority)
}

// Delete ...
func (s *service) Delete(ctx context.Context, accountID, userID, submissionID, runID string) error {
	id, err := NewExtraPriorityID(accountID, userID, submissionID, runID)
	if err != nil {
		return err
	}
	if _, err = s.repo.Get(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
