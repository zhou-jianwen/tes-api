package domain

import "context"

// Repo ...
type Repo interface {
	Create(ctx context.Context, task *Task) error
	GetStatus(ctx context.Context, id string) (*TaskStatus, error)
	UpdateStatus(ctx context.Context, taskStatus *TaskStatus) (bool, error)
	CheckIDExist(ctx context.Context, id string) (bool, error)
}
