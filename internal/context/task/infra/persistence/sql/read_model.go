package sql

import (
	"context"
	"errors"
	"fmt"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"gorm.io/gorm"

	"github.com/GBA-BI/tes-api/internal/context/task/application/query"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

type readModel struct {
	db *gorm.DB
}

// NewReadModel ...
func NewReadModel(ctx context.Context, db *gorm.DB) (query.ReadModel, error) {
	if err := db.WithContext(ctx).AutoMigrate(&Task{}); err != nil {
		return nil, err
	}
	return &readModel{db: db}, nil
}

var _ query.ReadModel = (*readModel)(nil)

// ListMinimal ...
func (r *readModel) ListMinimal(ctx context.Context, pageSize int, pageToken *utils.PageToken, filter *query.ListFilter) ([]*query.TaskMinimal, *utils.PageToken, error) {
	db := r.db.WithContext(ctx).Model(&Task{}).Order("`id`")
	db = listFilter(db, filter)

	db = db.Limit(pageSize)
	if pageToken != nil {
		db = db.Where("`id` > ?", pageToken.LastID)
	}

	taskStates := make([]*TaskState, 0)
	if err := db.Find(&taskStates).Error; err != nil {
		applog.Errorw("failed to list taskStates", "err", err)
		return nil, nil, apperrors.NewInternalError(err)
	}
	// maybe remains more tasks
	var nextPageToken *utils.PageToken
	if len(taskStates) == pageSize {
		nextPageToken = &utils.PageToken{LastID: taskStates[len(taskStates)-1].ID}
	}

	res := make([]*query.TaskMinimal, 0, len(taskStates))
	for _, taskState := range taskStates {
		res = append(res, taskState.toDTO())
	}
	return res, nextPageToken, nil
}

// ListBasic ...
func (r *readModel) ListBasic(ctx context.Context, pageSize int, pageToken *utils.PageToken, filter *query.ListFilter) ([]*query.TaskBasic, *utils.PageToken, error) {
	db := r.db.WithContext(ctx).Model(&Task{}).Order("`id`")
	db = listFilter(db, filter)

	db = db.Limit(pageSize)
	if pageToken != nil {
		db = db.Where("`id` > ?", pageToken.LastID)
	}

	taskBasics := make([]*TaskBasic, 0)
	if err := db.Find(&taskBasics).Error; err != nil {
		applog.Errorw("failed to list taskBasics", "err", err)
		return nil, nil, apperrors.NewInternalError(err)
	}
	// maybe remains more tasks
	var nextPageToken *utils.PageToken
	if len(taskBasics) == pageSize {
		nextPageToken = &utils.PageToken{LastID: taskBasics[len(taskBasics)-1].ID}
	}

	res := make([]*query.TaskBasic, 0, len(taskBasics))
	for _, taskBasic := range taskBasics {
		res = append(res, taskBasic.toDTO())
	}
	return res, nextPageToken, nil
}

// ListFull ...
func (r *readModel) ListFull(ctx context.Context, pageSize int, pageToken *utils.PageToken, filter *query.ListFilter) ([]*query.Task, *utils.PageToken, error) {
	db := r.db.WithContext(ctx).Model(&Task{}).Order("`id`")
	db = listFilter(db, filter)

	db = db.Limit(pageSize)
	if pageToken != nil {
		db = db.Where("`id` > ?", pageToken.LastID)
	}

	tasks := make([]*Task, 0)
	if err := db.Find(&tasks).Error; err != nil {
		applog.Errorw("failed to list tasks", "err", err)
		return nil, nil, apperrors.NewInternalError(err)
	}
	// maybe remains more tasks
	var nextPageToken *utils.PageToken
	if len(tasks) == pageSize {
		nextPageToken = &utils.PageToken{LastID: tasks[len(tasks)-1].ID}
	}

	res := make([]*query.Task, 0, len(tasks))
	for _, task := range tasks {
		res = append(res, task.toDTO())
	}
	return res, nextPageToken, nil
}

// GetMinimal ...
func (r *readModel) GetMinimal(ctx context.Context, id string) (*query.TaskMinimal, error) {
	var taskState TaskState
	if err := r.db.WithContext(ctx).Model(&Task{}).Where("`id` = ?", id).First(&taskState).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("task", id)
		}
		applog.Errorw("failed to get taskState", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return taskState.toDTO(), nil
}

// GetBasic ...
func (r *readModel) GetBasic(ctx context.Context, id string) (*query.TaskBasic, error) {
	var taskBasic TaskBasic
	if err := r.db.WithContext(ctx).Model(&Task{}).Where("`id` = ?", id).First(&taskBasic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("task", id)
		}
		applog.Errorw("failed to get taskBasic", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return taskBasic.toDTO(), nil
}

// GetFull ...
func (r *readModel) GetFull(ctx context.Context, id string) (*query.Task, error) {
	var task Task
	if err := r.db.WithContext(ctx).Model(&Task{}).Where("`id` = ?", id).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("task", id)
		}
		applog.Errorw("failed to get task", "err", err)
		return nil, apperrors.NewInternalError(err)
	}

	return task.toDTO(), nil
}

// GatherResources ...
func (r *readModel) GatherResources(ctx context.Context, filter *query.GatherFilter) (*query.TasksResources, error) {
	db := r.db.WithContext(ctx).Model(&Task{})
	db = gatherFilter(db, filter)

	res := &query.TasksResources{}

	var normalResources struct {
		Count    int     `gorm:"column:count"`
		CPUCores int     `gorm:"column:cpu_cores"`
		RamGB    float64 `gorm:"column:ram_gb"` // nolint
		DiskGB   float64 `gorm:"column:disk_gb"`
	}
	if err := db.Select("COUNT(*) AS `count`, SUM(`cpu_cores`) AS `cpu_cores`, SUM(`ram_gb`) AS `ram_gb`, SUM(`disk_gb`) AS `disk_gb`").Find(&normalResources).Error; err != nil {
		applog.Errorw("failed to gather tasks normal resources", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	res.Count = normalResources.Count
	res.CPUCores = normalResources.CPUCores
	res.RamGB = normalResources.RamGB
	res.DiskGB = normalResources.DiskGB

	var gpuResources []*struct {
		Type  string  `gorm:"column:gpu_type"`
		Count float64 `gorm:"column:gpu_count"`
	}
	if err := db.Select("`gpu_type`, SUM(`gpu_count`) AS `gpu_count`").Group("gpu_type").
		Where("`gpu_type` IS NOT NULL AND `gpu_count` IS NOT NULL").Find(&gpuResources).Error; err != nil {
		applog.Errorw("failed to gather tasks gpu resources", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	if len(gpuResources) > 0 {
		res.GPU = make(map[string]float64, len(gpuResources))
		for _, gpuResource := range gpuResources {
			res.GPU[gpuResource.Type] = gpuResource.Count
		}
	}

	return res, nil
}

// ListAccounts ...
func (r *readModel) ListAccounts(ctx context.Context) ([]*query.AccountInfo, error) {
	db := r.db.WithContext(ctx).Model(&Task{})

	var accountWithUsers []*struct {
		AccountID string `gorm:"column:account_id"`
		UserID    string `gorm:"column:user_id"`
	}
	if err := db.Distinct("`account_id`", "`user_id`").Find(&accountWithUsers).Error; err != nil {
		applog.Errorw("failed to list task accounts", "err", err)
		return nil, apperrors.NewInternalError(err)
	}

	accountsMap := make(map[string][]string)
	for _, accountWithUser := range accountWithUsers {
		if accountWithUser.AccountID == "" {
			continue
		}
		if _, ok := accountsMap[accountWithUser.AccountID]; !ok {
			accountsMap[accountWithUser.AccountID] = make([]string, 0)
		}
		accountsMap[accountWithUser.AccountID] = append(accountsMap[accountWithUser.AccountID], accountWithUser.UserID)
	}

	res := make([]*query.AccountInfo, 0, len(accountsMap))
	for accountID, userIDs := range accountsMap {
		res = append(res, &query.AccountInfo{
			AccountID: accountID,
			UserIDs:   userIDs,
		})
	}
	return res, nil
}

func listFilter(db *gorm.DB, filter *query.ListFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if filter.NamePrefix != "" {
		db = db.Where("`name` LIKE ?", fmt.Sprintf("%s%%", utils.EscapeLikeSpecialChars(filter.NamePrefix)))
	}
	if len(filter.State) > 0 {
		db = db.Where("`state` IN ?", filter.State)
	}
	if filter.ClusterID != "" {
		db = db.Where("`cluster_id` = ?", filter.ClusterID)
	}
	if filter.WithoutCluster {
		db = db.Where("`cluster_id` = ''")
	}
	return db
}

func gatherFilter(db *gorm.DB, filter *query.GatherFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if len(filter.State) > 0 {
		db = db.Where("`state` IN ?", filter.State)
	}
	if filter.ClusterID != "" {
		db = db.Where("`cluster_id` = ?", filter.ClusterID)
	} else if filter.WithCluster {
		db = db.Where("`cluster_id` <> ''")
	}
	if filter.AccountID != "" {
		db = db.Where("`account_id` = ?", filter.AccountID)
	}
	if filter.UserID != "" {
		db = db.Where("`user_id` = ?", filter.UserID)
	}
	return db
}
