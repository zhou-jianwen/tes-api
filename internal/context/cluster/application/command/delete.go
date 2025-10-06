package command

import (
	"context"

	"code.byted.org/epscp/vetes-api/internal/context/cluster/domain"
	"code.byted.org/epscp/vetes-api/pkg/validator"
)

// DeleteCommand ...
type DeleteCommand struct {
	ID string `validate:"required"`
}

func (c *DeleteCommand) setDefault() {}

func (c *DeleteCommand) validate() error {
	return validator.Validate(c)
}

// DeleteHandler ...
type DeleteHandler interface {
	Handle(ctx context.Context, cmd *DeleteCommand) error
}

type deleteHandler struct {
	svc domain.Service
}

var _ DeleteHandler = (*deleteHandler)(nil)

// NewDeleteHandler ...
func NewDeleteHandler(svc domain.Service) DeleteHandler {
	return &deleteHandler{svc: svc}
}

// Handle ...
func (h *deleteHandler) Handle(ctx context.Context, cmd *DeleteCommand) error {
	cmd.setDefault()
	if err := cmd.validate(); err != nil {
		return err
	}
	return h.svc.Delete(ctx, cmd.ID)
}
