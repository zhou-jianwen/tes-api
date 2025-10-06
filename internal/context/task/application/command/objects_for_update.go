package command

import (
	"time"

	"code.byted.org/epscp/vetes-api/internal/context/task/domain"
	"code.byted.org/epscp/vetes-api/pkg/validator"
)

// UpdateCommand ...
type UpdateCommand struct {
	ID        string `validate:"required"`
	ClusterID *string
	State     *string    `validate:"omitempty,oneof=QUEUED INITIALIZING RUNNING COMPLETE SYSTEM_ERROR EXECUTOR_ERROR CANCELING CANCELED"`
	Logs      []*TaskLog `validate:"unique=ClusterID,dive"`
}

// TaskLog ...
type TaskLog struct {
	ClusterID  string           `validate:"required"`
	Logs       [][]*ExecutorLog `validate:"dive,unique=ExecutorID,dive"`
	StartTime  *time.Time
	EndTime    *time.Time
	SystemLogs []string
}

// ExecutorLog ...
type ExecutorLog struct {
	ExecutorID string `validate:"required"`
	StartTime  *time.Time
	EndTime    *time.Time
}

func (c *UpdateCommand) setDefault() {}

func (c *UpdateCommand) validate() error {
	return validator.Validate(c)
}

type taskLogs []*TaskLog

func (l taskLogs) toDO() []*domain.TaskLog {
	if len(l) == 0 {
		return nil
	}
	res := make([]*domain.TaskLog, len(l))
	for index, log := range l {
		res[index] = log.toDO()
	}
	return res
}

func (t *TaskLog) toDO() *domain.TaskLog {
	if t == nil {
		return nil
	}
	res := &domain.TaskLog{
		ClusterID:  t.ClusterID,
		StartTime:  t.StartTime,
		EndTime:    t.EndTime,
		SystemLogs: t.SystemLogs,
	}
	if len(t.Logs) > 0 {
		res.Logs = make([][]*domain.ExecutorLog, len(t.Logs))
		for i, executorLogs := range t.Logs {
			if len(executorLogs) > 0 {
				res.Logs[i] = make([]*domain.ExecutorLog, len(executorLogs))
				for j, executorLog := range executorLogs {
					res.Logs[i][j] = executorLog.toDO()
				}
			}
		}
	}
	return res
}

func (e *ExecutorLog) toDO() *domain.ExecutorLog {
	if e == nil {
		return nil
	}
	return &domain.ExecutorLog{ExecutorID: e.ExecutorID, StartTime: e.StartTime, EndTime: e.EndTime}
}
