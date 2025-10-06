package domain

import (
	"context"
	"fmt"

	"github.com/GBA-BI/tes-api/pkg/consts"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

// Service ...
type Service interface {
	GetOrDefault(ctx context.Context, global bool, accountID, userID string) (*Quota, error)
	Put(ctx context.Context, global bool, accountID, userID string, resourceQuota *ResourceQuota) error
	Delete(ctx context.Context, global bool, accountID, userID string) error
}

type service struct {
	repo Repo
}

var _ Service = (*service)(nil)

// NewService ...
func NewService(repo Repo) Service {
	return &service{repo: repo}
}

// GetOrDefault ...
func (s *service) GetOrDefault(ctx context.Context, global bool, accountID, userID string) (*Quota, error) {
	id, err := NewQuotaID(global, accountID, userID)
	if err != nil {
		return nil, err
	}

	quota, err := s.repo.Get(ctx, id)
	if err == nil {
		return quota, nil
	}
	if !apperrors.IsCode(err, apperrors.NotFoundCode) {
		return nil, err
	}

	mayGetDefault := !global && accountID != "" && accountID != consts.DefaultQuotaAccountID && userID == ""
	if !mayGetDefault {
		return nil, err
	}

	defaultID, _ := NewQuotaID(false, consts.DefaultQuotaAccountID, "")
	quota, err = s.repo.Get(ctx, defaultID)
	if err != nil {
		if apperrors.IsCode(err, apperrors.NotFoundCode) {
			return nil, apperrors.NewNotFoundError("quota", fmt.Sprintf("%s and default", id))
		}
		return nil, err
	}
	return quota, nil
}

// Put ...
func (s *service) Put(ctx context.Context, global bool, accountID, userID string, resourceQuota *ResourceQuota) error {
	id, err := NewQuotaID(global, accountID, userID)
	if err != nil {
		return err
	}
	return s.repo.Save(ctx, &Quota{
		ID:            id,
		AccountID:     accountID,
		UserID:        userID,
		ResourceQuota: resourceQuota,
	})
}

// Delete ...
func (s *service) Delete(ctx context.Context, global bool, accountID, userID string) error {
	id, err := NewQuotaID(global, accountID, userID)
	if err != nil {
		return err
	}
	if _, err = s.repo.Get(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
