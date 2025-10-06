package sql

import (
	"context"
	"errors"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"gorm.io/gorm"

	"github.com/GBA-BI/tes-api/internal/context/task/domain"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

type repo struct {
	db *gorm.DB
}

// NewRepo ...
func NewRepo(ctx context.Context, db *gorm.DB) (domain.Repo, error) {
	if err := db.WithContext(ctx).AutoMigrate(&Task{}); err != nil {
		return nil, err
	}
	return &repo{db: db}, nil
}

var _ domain.Repo = (*repo)(nil)

// Create ...
func (r *repo) Create(ctx context.Context, task *domain.Task) error {
	taskPO := taskDOToPO(task)
	if err := r.db.WithContext(ctx).Model(&Task{}).Create(taskPO).Error; err != nil {
		applog.Errorw("failed to create task", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}

// GetStatus ...
func (r *repo) GetStatus(ctx context.Context, id string) (*domain.TaskStatus, error) {
	var taskStatus TaskStatus
	if err := r.db.WithContext(ctx).Model(&Task{}).Where("`id` = ?", id).First(&taskStatus).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("task", id)
		}
		applog.Errorw("failed to get taskStatus", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return taskStatus.toDO(), nil
}

// UpdateStatus ...
func (r *repo) UpdateStatus(ctx context.Context, taskStatus *domain.TaskStatus) (bool, error) {
	taskStatusPO := taskStatusDOToPO(taskStatus)
	oldStatusResourceVersion := taskStatusPO.StatusResourceVersion
	taskStatusPO.StatusResourceVersion = oldStatusResourceVersion + 1
	res := r.db.WithContext(ctx).Model(&Task{}).Where("`id` = ?", taskStatusPO.ID).
		Where("`status_resource_version` = ?", oldStatusResourceVersion).
		Updates(taskStatusPO)
	if err := res.Error; err != nil {
		applog.Errorw("failed to update taskStatus", "err", err)
		return false, apperrors.NewInternalError(err)
	}
	return res.RowsAffected > 0, nil
}

// CheckIDExist ...
func (r *repo) CheckIDExist(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Task{}).Where("`id` = ?", id).Count(&count).Error; err != nil {
		applog.Errorw("failed to count tasks", "err", err)
		return false, apperrors.NewInternalError(err)
	}
	return count > 0, nil
}
