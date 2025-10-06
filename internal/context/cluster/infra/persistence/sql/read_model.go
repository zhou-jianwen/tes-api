package sql

import (
	"context"

	applog "code.byted.org/epscp/go-common/log"
	"gorm.io/gorm"

	"code.byted.org/epscp/vetes-api/internal/context/cluster/application/query"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
)

type readModel struct {
	db *gorm.DB
}

// NewReadModel ...
func NewReadModel(ctx context.Context, db *gorm.DB) (query.ReadModel, error) {
	if err := db.WithContext(ctx).AutoMigrate(&Cluster{}); err != nil {
		return nil, err
	}
	return &readModel{db: db}, nil
}

var _ query.ReadModel = (*readModel)(nil)

// List ...
func (r *readModel) List(ctx context.Context, filter *query.ListFilter) ([]*query.Cluster, error) {
	db := r.db.WithContext(ctx).Model(&Cluster{})
	db = listFilter(db, filter)

	clusters := make([]*Cluster, 0)
	if err := db.Find(&clusters).Error; err != nil {
		applog.Errorw("failed to list clusters", "err", err)
		return nil, apperrors.NewInternalError(err)
	}

	res := make([]*query.Cluster, 0, len(clusters))
	for _, cluster := range clusters {
		res = append(res, cluster.toDTO())
	}
	return res, nil
}

func listFilter(db *gorm.DB, filter *query.ListFilter) *gorm.DB {
	return db
}
