package handlers

import (
	"time"

	applog "github.com/GBA-BI/tes-api/pkg/log"

	"github.com/GBA-BI/tes-api/internal/context/task/application/command"
	"github.com/GBA-BI/tes-api/internal/context/task/application/query"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

func (r *CreateTaskRequest) toDTO() *command.CreateCommand {
	if r == nil {
		return nil
	}
	res := &command.CreateCommand{
		Name:          r.Name,
		Description:   r.Description,
		Resources:     r.Resources.toDTO(),
		Volumes:       r.Volumes,
		Tags:          r.Tags,
		BioosInfo:     r.BioosInfo.toDTO(),
		PriorityValue: r.PriorityValue,
	}
	if len(r.Inputs) > 0 {
		res.Inputs = make([]*command.Input, len(r.Inputs))
		for index, input := range r.Inputs {
			res.Inputs[index] = input.toDTO()
		}
	}
	if len(r.Outputs) > 0 {
		res.Outputs = make([]*command.Output, len(r.Outputs))
		for index, output := range r.Outputs {
			res.Outputs[index] = output.toDTO()
		}
	}
	if len(r.Executors) > 0 {
		res.Executors = make([]*command.Executor, len(r.Executors))
		for index, executor := range r.Executors {
			res.Executors[index] = executor.toDTO()
		}
	}
	return res
}

func (i *Input) toDTO() *command.Input {
	if i == nil {
		return nil
	}
	return &command.Input{
		Name:        i.Name,
		Description: i.Description,
		Path:        i.Path,
		Type:        i.Type,
		URL:         i.URL,
		Content:     i.Content,
	}
}

func (o *Output) toDTO() *command.Output {
	if o == nil {
		return nil
	}
	return &command.Output{
		Name:        o.Name,
		Description: o.Description,
		Path:        o.Path,
		Type:        o.Type,
		URL:         o.URL,
	}
}

func (r *Resources) toDTO() *command.Resources {
	if r == nil {
		return nil
	}
	res := &command.Resources{
		CPUCores:   r.CPUCores,
		RamGB:      r.RamGB,
		DiskGB:     r.DiskGB,
		BootDiskGB: r.BootDiskGB,
	}
	if r.GPU != nil {
		res.GPU = &command.GPUResource{
			Count: r.GPU.Count,
			Type:  r.GPU.Type,
		}
	}
	return res
}

func (e *Executor) toDTO() *command.Executor {
	if e == nil {
		return nil
	}
	return &command.Executor{
		Image:   e.Image,
		Workdir: e.Workdir,
		Command: e.Command,
		Stdin:   e.Stdin,
		Stdout:  e.Stdout,
		Stderr:  e.Stderr,
		Env:     e.Env,
	}
}

func (b *BioosInfo) toDTO() *command.BioosInfo {
	if b == nil {
		return nil
	}
	var meta *command.BioosInfoMeta
	if b.Meta != nil {
		meta = &command.BioosInfoMeta{
			AAIPassport:     b.Meta.AAIPassport,
			MountTOS:        b.Meta.MountTOS,
			BucketsAuthInfo: b.Meta.BucketsAuthInfo.toDTO(),
		}
	}
	return &command.BioosInfo{
		AccountID:    b.AccountID,
		UserID:       b.UserID,
		SubmissionID: b.SubmissionID,
		RunID:        b.RunID,
		Meta:         meta,
	}
}

func (b *BucketsAuthInfo) toDTO() *command.BucketsAuthInfo {
	if b == nil {
		return nil
	}
	res := &command.BucketsAuthInfo{
		ReadOnly:  b.ReadOnly,
		ReadWrite: b.ReadWrite,
	}
	if len(b.External) > 0 {
		res.External = make([]*command.ExternalBucketAuthInfo, len(b.External))
		for index, external := range b.External {
			res.External[index] = &command.ExternalBucketAuthInfo{
				Bucket: external.Bucket,
				AK:     external.AK,
				SK:     external.SK,
			}
		}
	}
	return res
}

func (r *ListTasksRequest) toDTO() (*query.ListQuery, error) {
	if r == nil {
		return nil, nil
	}
	pageToken, err := utils.ParsePageToken(r.PageToken)
	if err != nil {
		return nil, err
	}
	return &query.ListQuery{
		View:      r.View,
		PageSize:  r.PageSize,
		PageToken: pageToken,
		Filter: &query.ListFilter{
			NamePrefix:     r.NamePrefix,
			State:          r.State,
			ClusterID:      r.ClusterID,
			WithoutCluster: r.WithoutCluster,
		},
	}, nil
}

func (r *GetTaskRequest) toDTO() *query.GetQuery {
	if r == nil {
		return nil
	}
	return &query.GetQuery{ID: r.ID, View: r.View}
}

func (r *CancelTaskRequest) toDTO() *command.CancelCommand {
	if r == nil {
		return nil
	}
	return &command.CancelCommand{ID: r.ID}
}

func (r *UpdateTaskRequest) toDTO() (*command.UpdateCommand, error) {
	if r == nil {
		return nil, nil
	}
	res := &command.UpdateCommand{
		ID:        r.ID,
		ClusterID: r.ClusterID,
		State:     r.State,
	}
	var err error
	if len(r.Logs) > 0 {
		res.Logs = make([]*command.TaskLog, len(r.Logs))
		for index, taskLog := range r.Logs {
			if res.Logs[index], err = taskLog.toDTO(); err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

func (t *TaskLog) toDTO() (*command.TaskLog, error) {
	if t == nil {
		return nil, nil
	}
	res := &command.TaskLog{
		ClusterID:  t.ClusterID,
		SystemLogs: t.SystemLogs,
	}
	if t.StartTime != nil && *t.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, *t.StartTime)
		if err != nil {
			applog.Errorw("parse startTime of taskLog", "err", err)
			return nil, apperrors.NewInvalidError("start_time")
		}
		res.StartTime = &startTime
	}
	if t.EndTime != nil && *t.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, *t.EndTime)
		if err != nil {
			applog.Errorw("parse endTime of taskLog", "err", err)
			return nil, apperrors.NewInvalidError("end_time")
		}
		res.EndTime = &endTime
	}
	var err error
	if len(t.Logs) > 0 {
		res.Logs = make([][]*command.ExecutorLog, len(t.Logs))
		for i, executorLogs := range t.Logs {
			if len(executorLogs) > 0 {
				res.Logs[i] = make([]*command.ExecutorLog, len(executorLogs))
				for j, executorLog := range executorLogs {
					if res.Logs[i][j], err = executorLog.toDTO(); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return res, nil
}

func (e *ExecutorLog) toDTO() (*command.ExecutorLog, error) {
	if e == nil {
		return nil, nil
	}
	res := &command.ExecutorLog{
		ExecutorID: e.ExecutorID,
	}
	if e.StartTime != nil && *e.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, *e.StartTime)
		if err != nil {
			applog.Errorw("parse startTime of executorLog", "err", err)
			return nil, apperrors.NewInvalidError("start_time")
		}
		res.StartTime = &startTime
	}
	if e.EndTime != nil && *e.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, *e.EndTime)
		if err != nil {
			applog.Errorw("parse endTime of executorLog", "err", err)
			return nil, apperrors.NewInvalidError("end_time")
		}
		res.EndTime = &endTime
	}
	return res, nil
}

func (r *GatherTasksResourcesRequest) toDTO() *query.GatherQuery {
	return &query.GatherQuery{Filter: &query.GatherFilter{
		State:       r.State,
		ClusterID:   r.ClusterID,
		WithCluster: r.WithCluster,
		AccountID:   r.AccountID,
		UserID:      r.UserID,
	}}
}

func taskDTOToVO(task *query.Task) *Task {
	if task == nil {
		return nil
	}
	res := &Task{
		ID:            task.ID,
		State:         task.State,
		Name:          task.Name,
		Description:   task.Description,
		Resources:     resourcesDTOToVO(task.Resources),
		Volumes:       task.Volumes,
		Tags:          task.Tags,
		BioosInfo:     bioosInfoDTOToVO(task.BioosInfo),
		PriorityValue: task.PriorityValue,
		ClusterID:     task.ClusterID,
	}
	if !task.CreationTime.IsZero() {
		res.CreationTime = task.CreationTime.Format(time.RFC3339)
	}
	if len(task.Inputs) > 0 {
		res.Inputs = make([]*Input, len(task.Inputs))
		for index, input := range task.Inputs {
			res.Inputs[index] = inputDTOToVO(input)
		}
	}
	if len(task.Outputs) > 0 {
		res.Outputs = make([]*Output, len(task.Outputs))
		for index, output := range task.Outputs {
			res.Outputs[index] = outputDTOToVO(output)
		}
	}
	if len(task.Executors) > 0 {
		res.Executors = make([]*Executor, len(task.Executors))
		for index, executor := range task.Executors {
			res.Executors[index] = executorDTOToVO(executor)
		}
	}
	if len(task.Logs) > 0 {
		res.Logs = make([]*TaskLog, len(task.Logs))
		for index, taskLog := range task.Logs {
			res.Logs[index] = taskLogDTOToVO(taskLog)
		}
	}
	return res
}

func inputDTOToVO(input *query.Input) *Input {
	if input == nil {
		return nil
	}
	return &Input{
		Name:        input.Name,
		Description: input.Description,
		Path:        input.Path,
		Type:        input.Type,
		URL:         input.URL,
		Content:     input.Content,
	}
}

func outputDTOToVO(output *query.Output) *Output {
	if output == nil {
		return nil
	}
	return &Output{
		Name:        output.Name,
		Description: output.Description,
		Path:        output.Path,
		Type:        output.Type,
		URL:         output.URL,
	}
}

func resourcesDTOToVO(resources *query.Resources) *Resources {
	if resources == nil {
		return nil
	}
	res := &Resources{
		CPUCores:   resources.CPUCores,
		RamGB:      resources.RamGB,
		DiskGB:     resources.DiskGB,
		BootDiskGB: resources.BootDiskGB,
	}
	if resources.GPU != nil {
		res.GPU = &GPUResource{Count: resources.GPU.Count, Type: resources.GPU.Type}
	}
	return res
}

func executorDTOToVO(executor *query.Executor) *Executor {
	if executor == nil {
		return nil
	}
	return &Executor{
		Image:   executor.Image,
		Command: executor.Command,
		Workdir: executor.Workdir,
		Stdin:   executor.Stdin,
		Stdout:  executor.Stdout,
		Stderr:  executor.Stderr,
		Env:     executor.Env,
	}
}

func taskLogDTOToVO(taskLog *query.TaskLog) *TaskLog {
	if taskLog == nil {
		return nil
	}
	res := &TaskLog{
		ClusterID:  taskLog.ClusterID,
		SystemLogs: taskLog.SystemLogs,
	}
	if taskLog.StartTime != nil && !taskLog.StartTime.IsZero() {
		res.StartTime = utils.Point(taskLog.StartTime.Format(time.RFC3339))
	}
	if taskLog.EndTime != nil && !taskLog.EndTime.IsZero() {
		res.EndTime = utils.Point(taskLog.EndTime.Format(time.RFC3339))
	}
	if len(taskLog.Logs) > 0 {
		res.Logs = make([][]*ExecutorLog, len(taskLog.Logs))
		for i, executorLogs := range taskLog.Logs {
			if len(executorLogs) > 0 {
				res.Logs[i] = make([]*ExecutorLog, len(executorLogs))
				for j, executorLog := range executorLogs {
					res.Logs[i][j] = executorLogDTOToVO(executorLog)
				}
			}
		}
	}
	return res
}

func executorLogDTOToVO(executorLog *query.ExecutorLog) *ExecutorLog {
	if executorLog == nil {
		return nil
	}
	res := &ExecutorLog{
		ExecutorID: executorLog.ExecutorID,
	}
	if executorLog.StartTime != nil && !executorLog.StartTime.IsZero() {
		res.StartTime = utils.Point(executorLog.StartTime.Format(time.RFC3339))
	}
	if executorLog.EndTime != nil && !executorLog.EndTime.IsZero() {
		res.EndTime = utils.Point(executorLog.EndTime.Format(time.RFC3339))
	}
	return res
}

func bioosInfoDTOToVO(bioosInfo *query.BioosInfo) *BioosInfo {
	if bioosInfo == nil {
		return nil
	}
	var meta *BioosInfoMeta
	if bioosInfo.Meta != nil {
		meta = &BioosInfoMeta{
			AAIPassport:     bioosInfo.Meta.AAIPassport,
			MountTOS:        bioosInfo.Meta.MountTOS,
			BucketsAuthInfo: bucketsAuthInfoDTOToVO(bioosInfo.Meta.BucketsAuthInfo),
		}
	}
	return &BioosInfo{
		AccountID:    bioosInfo.AccountID,
		UserID:       bioosInfo.UserID,
		SubmissionID: bioosInfo.SubmissionID,
		RunID:        bioosInfo.RunID,
		Meta:         meta,
	}
}

func bucketsAuthInfoDTOToVO(bucketsAuthInfo *query.BucketsAuthInfo) *BucketsAuthInfo {
	if bucketsAuthInfo == nil {
		return nil
	}
	res := &BucketsAuthInfo{
		ReadOnly:  bucketsAuthInfo.ReadOnly,
		ReadWrite: bucketsAuthInfo.ReadWrite,
	}
	if len(bucketsAuthInfo.External) > 0 {
		res.External = make([]*ExternalBucketAuthInfo, len(bucketsAuthInfo.External))
		for index, externalBucketAuthInfo := range bucketsAuthInfo.External {
			res.External[index] = &ExternalBucketAuthInfo{
				Bucket: externalBucketAuthInfo.Bucket,
				AK:     externalBucketAuthInfo.AK,
				SK:     externalBucketAuthInfo.SK,
			}
		}
	}
	return res
}

func tasksResourcesDTOToVO(resources *query.TasksResources) *GatherTasksResourcesResponse {
	return &GatherTasksResourcesResponse{
		Count:    resources.Count,
		CPUCores: resources.CPUCores,
		RamGB:    resources.RamGB,
		DiskGB:   resources.DiskGB,
		GPU:      resources.GPU,
	}
}

func taskAccountInfoDTOToVO(accountInfo *query.AccountInfo) *AccountInfo {
	return &AccountInfo{
		AccountID: accountInfo.AccountID,
		UserIDs:   accountInfo.UserIDs,
	}
}
