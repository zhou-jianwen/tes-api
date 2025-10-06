package sql

import (
	"context"
	"errors"

	applog "code.byted.org/epscp/go-common/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/domain"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
)

type repo struct {
	db *gorm.DB
}

// NewRepo ...
func NewRepo(ctx context.Context, db *gorm.DB) (domain.Repo, error) {
	if err := db.WithContext(ctx).AutoMigrate(&ExtraPriority{}); err != nil {
		return nil, err
	}
	return &repo{db: db}, nil
}

var _ domain.Repo = (*repo)(nil)

// Get ...
func (r *repo) Get(ctx context.Context, id string) (*domain.ExtraPriority, error) {
	var extraPriority ExtraPriority
	if err := r.db.WithContext(ctx).Model(&ExtraPriority{}).Where("`id` = ?", id).First(&extraPriority).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("extra priority", id)
		}
		applog.Errorw("failed to get extra priority", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return extraPriority.toDO(), nil
}

// Save ...
func (r *repo) Save(ctx context.Context, priority *domain.ExtraPriority) error {
	priorityPO := extraPriorityDOToPO(priority)
	if err := r.db.WithContext(ctx).Model(&ExtraPriority{}).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(priorityPO).Error; err != nil {
		applog.Errorw("failed to create extra priority", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}

// Delete ...
func (r *repo) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&ExtraPriority{}).Where("`id` = ?", id).Delete(&ExtraPriority{}).Error; err != nil {
		applog.Errorw("failed to delete extra priority", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}
