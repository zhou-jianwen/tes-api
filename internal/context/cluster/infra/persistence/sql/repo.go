package sql

import (
	"context"
	"errors"

	applog "code.byted.org/epscp/go-common/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"code.byted.org/epscp/vetes-api/internal/context/cluster/domain"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
)

type repo struct {
	db *gorm.DB
}

// NewRepo ...
func NewRepo(ctx context.Context, db *gorm.DB) (domain.Repo, error) {
	if err := db.WithContext(ctx).AutoMigrate(&Cluster{}); err != nil {
		return nil, err
	}
	return &repo{db: db}, nil
}

var _ domain.Repo = (*repo)(nil)

// Get ...
func (r *repo) Get(ctx context.Context, id string) (*domain.Cluster, error) {
	var cluster Cluster
	if err := r.db.WithContext(ctx).Model(&Cluster{}).Where("`id` = ?", id).First(&cluster).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("cluster", id)
		}
		applog.Errorw("failed to get cluster", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return cluster.toDO(), nil
}

// Save ...
func (r *repo) Save(ctx context.Context, cluster *domain.Cluster) error {
	clusterPO := clusterDOToPO(cluster)
	if err := r.db.WithContext(ctx).Model(&Cluster{}).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(clusterPO).Error; err != nil {
		applog.Errorw("failed to save cluster", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}

// Delete ...
func (r *repo) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&Cluster{}).Where("`id` = ?", id).Delete(&Cluster{}).Error; err != nil {
		applog.Errorw("failed to delete cluster", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}
