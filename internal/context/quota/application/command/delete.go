package command

import (
	"context"

	"code.byted.org/epscp/vetes-api/internal/context/quota/domain"
	"code.byted.org/epscp/vetes-api/pkg/validator"
)

// DeleteCommand ...
type DeleteCommand struct {
	Global    bool
	AccountID string
	UserID    string
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
	return &deleteHandler{
		svc: svc,
	}
}

// Handle ...
func (h *deleteHandler) Handle(ctx context.Context, cmd *DeleteCommand) error {
	cmd.setDefault()
	if err := cmd.validate(); err != nil {
		return err
	}
	return h.svc.Delete(ctx, cmd.Global, cmd.AccountID, cmd.UserID)
}
