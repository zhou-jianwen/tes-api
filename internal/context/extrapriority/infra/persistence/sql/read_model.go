package sql

import (
	"context"

	applog "code.byted.org/epscp/go-common/log"
	"gorm.io/gorm"

	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/application/query"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
)

type readModel struct {
	db *gorm.DB
}

// NewReadModel ...
func NewReadModel(ctx context.Context, db *gorm.DB) (query.ReadModel, error) {
	if err := db.WithContext(ctx).AutoMigrate(&ExtraPriority{}); err != nil {
		return nil, err
	}
	return &readModel{db: db}, nil
}

var _ query.ReadModel = (*readModel)(nil)

// List ...
func (r *readModel) List(ctx context.Context, filter *query.ListFilter) ([]*query.ExtraPriority, error) {
	db := r.db.WithContext(ctx).Model(&ExtraPriority{})
	db = listFilter(db, filter)

	extraPriorities := make([]*ExtraPriority, 0)
	if err := db.Find(&extraPriorities).Error; err != nil {
		applog.Errorw("failed to list extraPriorities", "err", err)
		return nil, apperrors.NewInternalError(err)
	}

	res := make([]*query.ExtraPriority, 0, len(extraPriorities))
	for _, extraPriority := range extraPriorities {
		res = append(res, extraPriority.toDTO())
	}
	return res, nil
}

func listFilter(db *gorm.DB, filter *query.ListFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if filter.AccountID != "" {
		db = db.Where("`account_id` = ?", filter.AccountID)
	}
	if filter.SubmissionID != "" {
		db = db.Where("`submission_id` = ?", filter.SubmissionID)
	}
	if filter.RunID != "" {
		db = db.Where("`run_id` = ?", filter.RunID)
	}
	return db
}
