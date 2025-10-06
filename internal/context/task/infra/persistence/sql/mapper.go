package sql

import (
	"github.com/GBA-BI/tes-api/internal/context/task/application/query"
	"github.com/GBA-BI/tes-api/internal/context/task/domain"
)

func (t *TaskState) toDTO() *query.TaskMinimal {
	if t == nil {
		return nil
	}
	return &query.TaskMinimal{ID: t.ID, State: t.State}
}

func (t *TaskBasic) toDTO() *query.TaskBasic {
	if t == nil {
		return nil
	}
	res := &query.TaskBasic{
		TaskMinimal:   *t.TaskState.toDTO(),
		Name:          t.Name,
		Description:   t.Description,
		Resources:     t.Resources.toDTO(),
		Volumes:       t.Volumes,
		Tags:          t.Tags,
		CreationTime:  t.CreationTime,
		BioosInfo:     t.BioosInfo.toDTO(),
		PriorityValue: t.PriorityValue,
	}
	if len(t.Executors) > 0 {
		res.Executors = make([]*query.Executor, len(t.Executors))
		for index, executor := range t.Executors {
			res.Executors[index] = executor.toDTO()
		}
	}
	if len(t.Logs) > 0 {
		res.Logs = make([]*query.TaskLog, len(t.Logs))
		for index, log := range t.Logs {
			res.Logs[index] = log.toDTO()
		}
	}
	if t.ClusterID == nil {
		res.ClusterID = ""
	} else {
		res.ClusterID = *t.ClusterID
	}
	return res
}

func (t *Task) toDTO() *query.Task {
	if t == nil {
		return nil
	}
	res := &query.Task{
		TaskBasic: *t.TaskBasic.toDTO(),
	}
	if len(t.Inputs) > 0 {
		res.Inputs = make([]*query.Input, len(t.Inputs))
		for index, input := range t.Inputs {
			res.Inputs[index] = input.toDTO()
		}
	}
	if len(t.Outputs) > 0 {
		res.Outputs = make([]*query.Output, len(t.Outputs))
		for index, output := range t.Outputs {
			res.Outputs[index] = output.toDTO()
		}
	}
	return res
}

func (i *Input) toDTO() *query.Input {
	if i == nil {
		return nil
	}
	return &query.Input{
		Name:        i.Name,
		Description: i.Description,
		Path:        i.Path,
		Type:        i.Type,
		URL:         i.URL,
		Content:     i.Content,
	}
}

func (o *Output) toDTO() *query.Output {
	if o == nil {
		return nil
	}
	return &query.Output{
		Name:        o.Name,
		Description: o.Description,
		Path:        o.Path,
		Type:        o.Type,
		URL:         o.URL,
	}
}

func (r *Resources) toDTO() *query.Resources {
	if r == nil {
		return nil
	}
	res := &query.Resources{
		CPUCores:   r.CPUCores,
		RamGB:      r.RamGB,
		DiskGB:     r.DiskGB,
		BootDiskGB: r.BootDiskGB,
	}
	if r.GPUType != nil && r.GPUCount != nil {
		res.GPU = &query.GPUResource{Count: *r.GPUCount, Type: *r.GPUType}
	}
	return res
}

func (e *Executor) toDTO() *query.Executor {
	if e == nil {
		return nil
	}
	return &query.Executor{
		Image:   e.Image,
		Command: e.Command,
		Workdir: e.Workdir,
		Stdin:   e.Stdin,
		Stdout:  e.Stdout,
		Stderr:  e.Stderr,
		Env:     e.Env,
	}
}

func (t *TaskLog) toDTO() *query.TaskLog {
	if t == nil {
		return nil
	}
	res := &query.TaskLog{
		ClusterID:  t.ClusterID,
		StartTime:  t.StartTime,
		EndTime:    t.EndTime,
		SystemLogs: t.SystemLogs,
	}
	if len(t.Logs) > 0 {
		res.Logs = make([][]*query.ExecutorLog, len(t.Logs))
		for i, executorLogs := range t.Logs {
			if len(executorLogs) > 0 {
				res.Logs[i] = make([]*query.ExecutorLog, len(executorLogs))
				for j, executorLog := range executorLogs {
					res.Logs[i][j] = executorLog.toDTO()
				}
			}
		}
	}
	return res
}

func (e *ExecutorLog) toDTO() *query.ExecutorLog {
	if e == nil {
		return nil
	}
	return &query.ExecutorLog{
		ExecutorID: e.ExecutorID,
		StartTime:  e.StartTime,
		EndTime:    e.EndTime,
	}
}

func (b *BioosInfo) toDTO() *query.BioosInfo {
	if b == nil {
		return nil
	}
	var meta *query.BioosInfoMeta
	if b.Meta != nil {
		meta = &query.BioosInfoMeta{
			AAIPassport:     b.Meta.AAIPassport,
			MountTOS:        b.Meta.MountTOS,
			BucketsAuthInfo: b.Meta.BucketsAuthInfo.toDTO(),
		}
	}
	return &query.BioosInfo{
		AccountID:    b.AccountID,
		UserID:       b.UserID,
		SubmissionID: b.SubmissionID,
		RunID:        b.RunID,
		Meta:         meta,
	}
}

func (b *BucketsAuthInfo) toDTO() *query.BucketsAuthInfo {
	if b == nil {
		return nil
	}
	res := &query.BucketsAuthInfo{
		ReadOnly:  b.ReadOnly,
		ReadWrite: b.ReadWrite,
	}
	if len(b.External) > 0 {
		res.External = make([]*query.ExternalBucketAuthInfo, len(b.External))
		for index, external := range b.External {
			res.External[index] = &query.ExternalBucketAuthInfo{
				Bucket: external.Bucket,
				AK:     external.AK,
				SK:     external.SK,
			}
		}
	}
	return res
}

func (t *TaskStatus) toDO() *domain.TaskStatus {
	if t == nil {
		return nil
	}
	res := &domain.TaskStatus{
		ID:           t.ID,
		State:        t.State,
		CreationTime: t.CreationTime,
	}
	if len(t.Logs) > 0 {
		res.Logs = make([]*domain.TaskLog, len(t.Logs))
		for index, log := range t.Logs {
			res.Logs[index] = log.toDO()
		}
	}
	if t.ClusterID == nil {
		res.ClusterID = ""
	} else {
		res.ClusterID = *t.ClusterID
	}
	res.StatusResourceVersion = t.StatusResourceVersion
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
	return &domain.ExecutorLog{
		ExecutorID: e.ExecutorID,
		StartTime:  e.StartTime,
		EndTime:    e.EndTime,
	}
}

func taskStatusDOToPO(taskStatus *domain.TaskStatus) *TaskStatus {
	if taskStatus == nil {
		return nil
	}
	res := &TaskStatus{
		TaskState: TaskState{
			ID:    taskStatus.ID,
			State: taskStatus.State,
		},
		CreationTime: taskStatus.CreationTime,
		ClusterID:    &taskStatus.ClusterID,
	}
	if len(taskStatus.Logs) > 0 {
		res.Logs = make([]*TaskLog, len(taskStatus.Logs))
		for index, log := range taskStatus.Logs {
			res.Logs[index] = taskLogDOToPO(log)
		}
	}
	res.StatusResourceVersion = taskStatus.StatusResourceVersion
	return res
}

func taskDOToPO(task *domain.Task) *Task {
	if task == nil {
		return nil
	}
	res := &Task{
		TaskBasic: TaskBasic{
			TaskStatus:    *taskStatusDOToPO(&task.TaskStatus),
			Name:          task.Name,
			Description:   task.Description,
			Resources:     resourcesDOToPO(task.Resources),
			Volumes:       task.Volumes,
			Tags:          task.Tags,
			BioosInfo:     bioosInfoDOToPO(task.BioosInfo),
			PriorityValue: task.PriorityValue,
		},
	}

	if len(task.Inputs) > 0 {
		res.Inputs = make([]*Input, len(task.Inputs))
		for index, input := range task.Inputs {
			res.Inputs[index] = inputDOToPO(input)
		}
	}
	if len(task.Outputs) > 0 {
		res.Outputs = make([]*Output, len(task.Outputs))
		for index, output := range task.Outputs {
			res.Outputs[index] = outputDOToPO(output)
		}
	}
	if len(task.Executors) > 0 {
		res.Executors = make([]*Executor, len(task.Executors))
		for index, executor := range task.Executors {
			res.Executors[index] = executorDOToPO(executor)
		}
	}
	return res
}

func taskLogDOToPO(log *domain.TaskLog) *TaskLog {
	if log == nil {
		return nil
	}
	res := &TaskLog{
		ClusterID:  log.ClusterID,
		StartTime:  log.StartTime,
		EndTime:    log.EndTime,
		SystemLogs: log.SystemLogs,
	}
	if len(log.Logs) > 0 {
		res.Logs = make([][]*ExecutorLog, len(log.Logs))
		for i, executorLogs := range log.Logs {
			if len(executorLogs) > 0 {
				res.Logs[i] = make([]*ExecutorLog, len(executorLogs))
				for j, executorLog := range executorLogs {
					res.Logs[i][j] = executorLogDOToPO(executorLog)
				}
			}
		}
	}
	return res
}

func executorLogDOToPO(log *domain.ExecutorLog) *ExecutorLog {
	if log == nil {
		return nil
	}
	return &ExecutorLog{
		ExecutorID: log.ExecutorID,
		StartTime:  log.StartTime,
		EndTime:    log.EndTime,
	}
}

func inputDOToPO(input *domain.Input) *Input {
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

func outputDOToPO(output *domain.Output) *Output {
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

func resourcesDOToPO(resources *domain.Resources) *Resources {
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
		res.GPUCount = &resources.GPU.Count
		res.GPUType = &resources.GPU.Type
	}
	return res
}

func executorDOToPO(executor *domain.Executor) *Executor {
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

func bioosInfoDOToPO(info *domain.BioosInfo) *BioosInfo {
	if info == nil {
		return nil
	}
	var meta *BioosInfoMeta
	if info.Meta != nil {
		meta = &BioosInfoMeta{
			AAIPassport:     info.Meta.AAIPassport,
			MountTOS:        info.Meta.MountTOS,
			BucketsAuthInfo: bucketsAuthInfoDOToPO(info.Meta.BucketsAuthInfo),
		}
	}
	return &BioosInfo{
		AccountID:    info.AccountID,
		UserID:       info.UserID,
		SubmissionID: info.SubmissionID,
		RunID:        info.RunID,
		Meta:         meta,
	}
}

func bucketsAuthInfoDOToPO(info *domain.BucketsAuthInfo) *BucketsAuthInfo {
	if info == nil {
		return nil
	}
	res := &BucketsAuthInfo{
		ReadOnly:  info.ReadOnly,
		ReadWrite: info.ReadWrite,
	}
	if len(info.External) > 0 {
		res.External = make([]*ExternalBucketAuthInfo, len(info.External))
		for index, external := range info.External {
			res.External[index] = &ExternalBucketAuthInfo{
				Bucket: external.Bucket,
				AK:     external.AK,
				SK:     external.SK,
			}
		}
	}
	return res
}
