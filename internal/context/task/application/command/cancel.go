package command

import (
	"context"

	"github.com/GBA-BI/tes-api/internal/context/task/domain"
	"github.com/GBA-BI/tes-api/pkg/validator"
)

// CancelCommand ...
type CancelCommand struct {
	ID string `validate:"required"`
}

func (c *CancelCommand) setDefault() {}

func (c *CancelCommand) validate() error {
	return validator.Validate(c)
}

// CancelHandler ...
type CancelHandler interface {
	Handle(ctx context.Context, cmd *CancelCommand) error
}

type cancelHandler struct {
	svc domain.Service
}

var _ CancelHandler = (*cancelHandler)(nil)

// NewCancelHandler ...
func NewCancelHandler(svc domain.Service) CancelHandler {
	return &cancelHandler{svc: svc}
}

// Handle ...
func (h *cancelHandler) Handle(ctx context.Context, cmd *CancelCommand) error {
	cmd.setDefault()
	if err := cmd.validate(); err != nil {
		return err
	}
	return h.svc.Cancel(ctx, cmd.ID)
}
