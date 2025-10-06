package domain

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/GBA-BI/tes-api/pkg/consts"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

// GenTaskID returns task-[0-9a-f](8), which keep the same as TESK
func GenTaskID() string {
	randBytes := make([]byte, 4)
	rand.Read(randBytes)
	return fmt.Sprintf("task-%x", randBytes)
}

// Task ...
type Task struct {
	TaskStatus
	Name          string
	Description   string
	Inputs        []*Input
	Outputs       []*Output
	Resources     *Resources
	Executors     []*Executor
	Volumes       []string
	Tags          map[string]string
	BioosInfo     *BioosInfo
	PriorityValue int
}

// TaskStatus contains fields not specified by CreateTask
type TaskStatus struct {
	ID           string
	State        string
	Logs         []*TaskLog
	CreationTime time.Time
	ClusterID    string

	StatusResourceVersion int
}

// Input ...
type Input struct {
	Name        string
	Description string
	Path        string
	Type        string
	URL         string
	Content     string
}

// Output ...
type Output struct {
	Name        string
	Description string
	Path        string
	Type        string
	URL         string
}

// Resources ...
type Resources struct {
	CPUCores   int
	RamGB      float64 // nolint
	DiskGB     float64
	BootDiskGB *int
	GPU        *GPUResource
}

// GPUResource ...
type GPUResource struct {
	Count float64
	Type  string
}

// Executor ...
type Executor struct {
	Image   string
	Command []string
	Workdir string
	Stdin   string
	Stdout  string
	Stderr  string
	Env     map[string]string
}

// TaskLog ...
type TaskLog struct {
	ClusterID  string
	Logs       [][]*ExecutorLog
	StartTime  *time.Time
	EndTime    *time.Time
	SystemLogs []string
}

// ExecutorLog ...
type ExecutorLog struct {
	ExecutorID string
	StartTime  *time.Time
	EndTime    *time.Time
}

// BioosInfo ...
type BioosInfo struct {
	AccountID    string
	UserID       string
	SubmissionID string
	RunID        string
	Meta         *BioosInfoMeta
}

// BioosInfoMeta ...
type BioosInfoMeta struct {
	AAIPassport     *string
	MountTOS        *bool
	BucketsAuthInfo *BucketsAuthInfo
}

// BucketsAuthInfo ...
type BucketsAuthInfo struct {
	ReadOnly  []string
	ReadWrite []string
	External  []*ExternalBucketAuthInfo
}

// ExternalBucketAuthInfo ...
type ExternalBucketAuthInfo struct {
	Bucket string
	AK     string
	SK     string
}

var executingStates = map[string]struct{}{
	consts.TaskQueued:       {},
	consts.TaskInitializing: {},
	consts.TaskRunning:      {},
}

var finishedStates = map[string]struct{}{
	consts.TaskCanceled:      {},
	consts.TaskComplete:      {},
	consts.TaskExecutorError: {},
	consts.TaskSystemError:   {},
}

// UpdateState ...
func (t *TaskStatus) UpdateState(newState string) error {
	if t.State == newState {
		return nil
	}

	if _, ok := finishedStates[t.State]; ok {
		return apperrors.NewCannotExecError("finished job state cannot be updated")
	}

	if t.State == consts.TaskCanceling {
		if _, ok := executingStates[newState]; ok {
			return apperrors.NewCannotExecError("CANCELING job state cannot be updated back to executing")
		}
	}

	if t.State != consts.TaskCanceling && newState == consts.TaskCanceled {
		return apperrors.NewCannotExecError("only job state is CANCELING, it can be changed to CANCELED")
	}

	t.State = newState
	return nil
}

// Cancel ...
func (t *TaskStatus) Cancel() error {
	return t.UpdateState(consts.TaskCanceling)
}

// UpdateClusterID ...
func (t *TaskStatus) UpdateClusterID(clusterID string) error {
	if t.ClusterID == clusterID {
		return nil
	}
	if t.State != consts.TaskQueued {
		return apperrors.NewCannotExecError("only QUEUED job cluster_id may be changed")
	}
	t.ClusterID = clusterID
	return nil
}

// UpdateLogs ...
func (t *TaskStatus) UpdateLogs(newLogs []*TaskLog) error {
	t.Logs = mergeTaskLogs(t.Logs, newLogs)

	for _, taskLog := range t.Logs {
		if err := validateTaskLog(taskLog, t.CreationTime); err != nil {
			return err
		}
	}
	return nil
}

func mergeTaskLogs(old, new []*TaskLog) []*TaskLog {
	for _, newLog := range new {
		if newLog == nil {
			continue
		}
		existMatch := false
		for index, oldLog := range old {
			if newLog.ClusterID == oldLog.ClusterID {
				old[index] = mergeTaskLog(old[index], newLog)
				existMatch = true
				break
			}
		}
		if !existMatch {
			old = append(old, newLog)
		}
	}
	return old
}

func mergeTaskLog(old, new *TaskLog) *TaskLog {
	if new == nil {
		return old
	}
	old.Logs = mergeExecutorLogsSlice(old.Logs, new.Logs)
	if new.StartTime != nil {
		old.StartTime = new.StartTime
	}
	if new.EndTime != nil {
		old.EndTime = new.EndTime
	}
	old.SystemLogs = mergeSystemLogs(old.SystemLogs, new.SystemLogs)
	return old
}

func mergeExecutorLogsSlice(old, new [][]*ExecutorLog) [][]*ExecutorLog {
	less := min(len(old), len(new))
	for index := 0; index < less; index++ {
		old[index] = mergeExecutorLogs(old[index], new[index])
	}
	if len(old) < len(new) {
		old = append(old, new[len(old):]...)
	}
	return old
}

func mergeExecutorLogs(old, new []*ExecutorLog) []*ExecutorLog {
	for _, newLog := range new {
		if newLog == nil {
			continue
		}
		existMatch := false
		for index, oldLog := range old {
			if newLog.ExecutorID == oldLog.ExecutorID {
				old[index] = mergeExecutorLog(old[index], newLog)
				existMatch = true
				break
			}
		}
		if !existMatch {
			old = append(old, newLog)
		}
	}
	return old
}

func mergeExecutorLog(old, new *ExecutorLog) *ExecutorLog {
	if new == nil {
		return old
	}
	if new.StartTime != nil {
		old.StartTime = new.StartTime
	}
	if new.EndTime != nil {
		old.EndTime = new.EndTime
	}
	return old
}

func mergeSystemLogs(old, new []string) []string {
	less := min(len(old), len(new))
	for index := 0; index < less; index++ {
		if new[index] != "" {
			old[index] = new[index]
		}
	}
	if len(old) < len(new) {
		old = append(old, new[len(old):]...)
	}
	return old
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func validateTaskLog(taskLog *TaskLog, creationTime time.Time) error {
	if err := validateTime(taskLog.StartTime, taskLog.EndTime); err != nil {
		return err
	}
	if taskLog.StartTime != nil && taskLog.StartTime.Before(creationTime) {
		return apperrors.NewInvalidError("task start_time before task creation_time")
	}
	for _, executorLogs := range taskLog.Logs {
		for _, executorLog := range executorLogs {
			if err := validateExecutorLog(executorLog, taskLog.StartTime, taskLog.EndTime); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateExecutorLog(executorLog *ExecutorLog, taskStartTime, taskEndTime *time.Time) error {
	if err := validateTime(executorLog.StartTime, executorLog.EndTime); err != nil {
		return err
	}
	if executorLog.StartTime != nil {
		if taskStartTime == nil {
			return apperrors.NewInvalidError("empty task start_time with non-empty executor start_time")
		} else if executorLog.StartTime.Before(*taskStartTime) {
			return apperrors.NewInvalidError("executor start_time before task start_time")
		}
	}
	if taskEndTime != nil {
		if executorLog.EndTime == nil {
			return apperrors.NewInvalidError("empty executor end_time with non-empty task end_time")
		} else if executorLog.EndTime.After(*taskEndTime) {
			return apperrors.NewInvalidError("executor end_time after task end_time")
		}
	}
	return nil
}

func validateTime(startTime, endTime *time.Time) error {
	if endTime == nil {
		return nil
	}
	if startTime == nil {
		return apperrors.NewInvalidError("empty start_time with non-empty end_time")
	} else if startTime.After(*endTime) {
		return apperrors.NewInvalidError("start_time after end_time")
	}
	return nil
}
