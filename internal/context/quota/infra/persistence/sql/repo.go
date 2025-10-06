package sql

import (
	"context"
	"errors"

	applog "code.byted.org/epscp/go-common/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"code.byted.org/epscp/vetes-api/internal/context/quota/domain"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
)

type repo struct {
	db *gorm.DB
}

// NewRepo ...
func NewRepo(ctx context.Context, db *gorm.DB) (domain.Repo, error) {
	if err := db.WithContext(ctx).AutoMigrate(&Quota{}); err != nil {
		return nil, err
	}
	return &repo{db: db}, nil
}

var _ domain.Repo = (*repo)(nil)

// Get ...
func (r *repo) Get(ctx context.Context, id string) (*domain.Quota, error) {
	var quota Quota
	if err := r.db.WithContext(ctx).Model(&Quota{}).Where("`id` = ?", id).First(&quota).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("quota", id)
		}
		applog.Errorw("failed to get quota", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return quota.toDO(), nil
}

// Save ...
func (r *repo) Save(ctx context.Context, quota *domain.Quota) error {
	quotaPO := quotaDOToPO(quota)
	if err := r.db.WithContext(ctx).Model(&Quota{}).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(quotaPO).Error; err != nil {
		applog.Errorw("failed to save quota", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}

// Delete ...
func (r *repo) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&Quota{}).Where("`id` = ?", id).Delete(&Quota{}).Error; err != nil {
		applog.Errorw("failed to delete quota", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}
