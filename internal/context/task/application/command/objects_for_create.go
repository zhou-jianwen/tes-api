package command

import (
	"time"

	"code.byted.org/epscp/vetes-api/internal/context/task/domain"
	"code.byted.org/epscp/vetes-api/pkg/consts"
	"code.byted.org/epscp/vetes-api/pkg/validator"
)

// CreateCommand ...
type CreateCommand struct {
	Name          string
	Description   string
	Inputs        []*Input  `validate:"dive"`
	Outputs       []*Output `validate:"dive"`
	Resources     *Resources
	Executors     []*Executor `validate:"gt=0,dive"`
	Volumes       []string
	Tags          map[string]string
	BioosInfo     *BioosInfo
	PriorityValue int
}

// Input ...
type Input struct {
	Name        string
	Description string
	Path        string `validate:"required,abspath"`
	Type        string `validate:"oneof=FILE DIRECTORY"`
	URL         string `validate:"required_without=Content,omitempty,uri"`
	Content     string `validate:"required_without=URL"`
}

// Output ...
type Output struct {
	Name        string
	Description string
	Path        string `validate:"required,abspath"`
	Type        string `validate:"oneof=FILE DIRECTORY"`
	URL         string `validate:"required,uri"`
}

// Resources ...
type Resources struct {
	CPUCores   int     `validate:"gte=0"`
	RamGB      float64 `validate:"gte=0"` // nolint
	DiskGB     float64 `validate:"gte=0"`
	BootDiskGB *int    `validate:"omitempty,gte=0"`
	GPU        *GPUResource
}

// GPUResource ...
type GPUResource struct {
	Count float64 `validate:"gt=0"`
	Type  string
}

// Executor ...
type Executor struct {
	Image   string   `validate:"required"`
	Command []string `validate:"gt=0,dive,required"`
	Workdir string   `validate:"omitempty,abspath"`
	Stdin   string   `validate:"omitempty,abspath"`
	Stdout  string   `validate:"omitempty,abspath"`
	Stderr  string   `validate:"omitempty,abspath"`
	Env     map[string]string
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

func (c *CreateCommand) setDefault() {}

func (c *CreateCommand) validate() error {
	return validator.Validate(c)
}

func (c *CreateCommand) toDO() *domain.Task {
	res := &domain.Task{
		TaskStatus: domain.TaskStatus{
			State:        consts.TaskQueued,
			CreationTime: time.Now().UTC().Truncate(time.Second),
		},
		Name:          c.Name,
		Description:   c.Description,
		Resources:     c.Resources.toDO(),
		Volumes:       c.Volumes,
		Tags:          c.Tags,
		BioosInfo:     c.BioosInfo.toDO(),
		PriorityValue: c.PriorityValue,
	}
	if len(c.Inputs) > 0 {
		res.Inputs = make([]*domain.Input, len(c.Inputs))
		for index, input := range c.Inputs {
			res.Inputs[index] = input.toDO()
		}
	}
	if len(c.Outputs) > 0 {
		res.Outputs = make([]*domain.Output, len(c.Outputs))
		for index, output := range c.Outputs {
			res.Outputs[index] = output.toDO()
		}
	}
	if len(c.Executors) > 0 {
		res.Executors = make([]*domain.Executor, len(c.Executors))
		for index, executor := range c.Executors {
			res.Executors[index] = executor.toDO()
		}
	}
	return res
}

func (i *Input) toDO() *domain.Input {
	if i == nil {
		return nil
	}
	return &domain.Input{
		Name:        i.Name,
		Description: i.Description,
		Path:        i.Path,
		Type:        i.Type,
		URL:         i.URL,
		Content:     i.Content,
	}
}

func (o *Output) toDO() *domain.Output {
	if o == nil {
		return nil
	}
	return &domain.Output{
		Name:        o.Name,
		Description: o.Description,
		Path:        o.Path,
		Type:        o.Type,
		URL:         o.URL,
	}
}

func (r *Resources) toDO() *domain.Resources {
	if r == nil {
		return nil
	}
	res := &domain.Resources{
		CPUCores:   r.CPUCores,
		RamGB:      r.RamGB,
		DiskGB:     r.DiskGB,
		BootDiskGB: r.BootDiskGB,
	}
	if r.GPU != nil {
		res.GPU = &domain.GPUResource{Count: r.GPU.Count, Type: r.GPU.Type}
	}
	return res
}

func (e *Executor) toDO() *domain.Executor {
	if e == nil {
		return nil
	}
	return &domain.Executor{
		Image:   e.Image,
		Command: e.Command,
		Workdir: e.Workdir,
		Stdin:   e.Stdin,
		Stdout:  e.Stdout,
		Stderr:  e.Stderr,
		Env:     e.Env,
	}
}

func (b *BioosInfo) toDO() *domain.BioosInfo {
	if b == nil {
		return nil
	}
	var meta *domain.BioosInfoMeta
	if b.Meta != nil {
		meta = &domain.BioosInfoMeta{
			AAIPassport:     b.Meta.AAIPassport,
			MountTOS:        b.Meta.MountTOS,
			BucketsAuthInfo: b.Meta.BucketsAuthInfo.toDO(),
		}
	}
	return &domain.BioosInfo{
		AccountID:    b.AccountID,
		UserID:       b.UserID,
		SubmissionID: b.SubmissionID,
		RunID:        b.RunID,
		Meta:         meta,
	}
}

func (b *BucketsAuthInfo) toDO() *domain.BucketsAuthInfo {
	if b == nil {
		return nil
	}
	res := &domain.BucketsAuthInfo{
		ReadOnly:  b.ReadOnly,
		ReadWrite: b.ReadWrite,
	}
	if len(b.External) > 0 {
		res.External = make([]*domain.ExternalBucketAuthInfo, len(b.External))
		for index, external := range b.External {
			res.External[index] = &domain.ExternalBucketAuthInfo{
				Bucket: external.Bucket,
				AK:     external.AK,
				SK:     external.SK,
			}
		}
	}
	return res
}
