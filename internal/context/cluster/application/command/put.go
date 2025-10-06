package command

import (
	"context"
	"time"

	"code.byted.org/epscp/vetes-api/internal/context/cluster/domain"
	"code.byted.org/epscp/vetes-api/pkg/validator"
)

// PutCommand ...
type PutCommand struct {
	ID       string `validate:"required"`
	Capacity *Capacity
	Limits   *Limits
}

// Capacity ...
type Capacity struct {
	Count       *int     `validate:"omitempty,gte=0"`
	CPUCores    *int     `validate:"omitempty,gte=0"`
	RamGB       *float64 `validate:"omitempty,gte=0"` // nolint
	DiskGB      *float64 `validate:"omitempty,gte=0"`
	GPUCapacity *GPUCapacity
}

// GPUCapacity ...
type GPUCapacity struct {
	GPU map[string]float64 `validate:"dive,gte=0"`
}

// Limits ...
type Limits struct {
	CPUCores *int     `validate:"omitempty,gte=0"`
	RamGB    *float64 `validate:"omitempty,gte=0"` // nolint
	GPULimit *GPULimit
}

// GPULimit ...
type GPULimit struct {
	GPU map[string]float64 `validate:"dive,gte=0"`
}

func (c *PutCommand) setDefault() {}

func (c *PutCommand) validate() error {
	return validator.Validate(c)
}

func (c *PutCommand) toDO() *domain.Cluster {
	res := &domain.Cluster{
		ID:                 c.ID,
		HeartbeatTimestamp: time.Now().UTC().Truncate(time.Second),
	}
	if c.Capacity != nil {
		res.Capacity = &domain.Capacity{
			Count:    c.Capacity.Count,
			CPUCores: c.Capacity.CPUCores,
			RamGB:    c.Capacity.RamGB,
			DiskGB:   c.Capacity.DiskGB,
		}
		if c.Capacity.GPUCapacity != nil {
			res.Capacity.GPUCapacity = &domain.GPUCapacity{
				GPU: c.Capacity.GPUCapacity.GPU,
			}
		}
	}
	if c.Limits != nil {
		res.Limits = &domain.Limits{
			CPUCores: c.Limits.CPUCores,
			RamGB:    c.Limits.RamGB,
		}
		if c.Limits.GPULimit != nil {
			res.Limits.GPULimit = &domain.GPULimit{
				GPU: c.Limits.GPULimit.GPU,
			}
		}
	}
	return res
}

// PutHandler ...
type PutHandler interface {
	Handle(ctx context.Context, cmd *PutCommand) error
}

type putHandler struct {
	svc domain.Service
}

var _ PutHandler = (*putHandler)(nil)

// NewPutHandler ...
func NewPutHandler(svc domain.Service) PutHandler {
	return &putHandler{svc: svc}
}

// Handle ...
func (h *putHandler) Handle(ctx context.Context, cmd *PutCommand) error {
	cmd.setDefault()
	if err := cmd.validate(); err != nil {
		return err
	}
	return h.svc.Put(ctx, cmd.toDO())
}
