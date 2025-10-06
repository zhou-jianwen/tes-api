package command

import (
	"context"

	"github.com/GBA-BI/tes-api/internal/context/quota/domain"
	"github.com/GBA-BI/tes-api/pkg/validator"
)

// PutCommand ...
type PutCommand struct {
	Global        bool
	AccountID     string
	UserID        string
	ResourceQuota *ResourceQuota
}

// ResourceQuota ...
type ResourceQuota struct {
	Count    *int     `validate:"omitempty,gte=0"`
	CPUCores *int     `validate:"omitempty,gte=0"`
	RamGB    *float64 `validate:"omitempty,gte=0"` // nolint
	DiskGB   *float64 `validate:"omitempty,gte=0"`
	GPUQuota *GPUQuota
}

// GPUQuota ...
type GPUQuota struct {
	GPU map[string]float64 `validate:"dive,gte=0"`
}

func (r *ResourceQuota) toDO() *domain.ResourceQuota {
	if r == nil {
		return nil
	}
	return &domain.ResourceQuota{
		Count:    r.Count,
		CPUCores: r.CPUCores,
		RamGB:    r.RamGB,
		DiskGB:   r.DiskGB,
		GPUQuota: r.GPUQuota.toDO(),
	}
}

func (g *GPUQuota) toDO() *domain.GPUQuota {
	if g == nil {
		return nil
	}
	return &domain.GPUQuota{GPU: g.GPU}
}

func (c *PutCommand) setDefault() {}

func (c *PutCommand) validate() error {
	return validator.Validate(c)
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
	return h.svc.Put(ctx, cmd.Global, cmd.AccountID, cmd.UserID, cmd.ResourceQuota.toDO())
}
