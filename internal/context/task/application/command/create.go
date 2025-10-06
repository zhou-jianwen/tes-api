package command

import (
	"context"

	"code.byted.org/epscp/vetes-api/internal/context/task/domain"
)

// CreateHandler ...
type CreateHandler interface {
	Handle(ctx context.Context, cmd *CreateCommand) (string, error)
}

type createHandler struct {
	svc domain.Service
}

var _ CreateHandler = (*createHandler)(nil)

// NewCreateHandler ...
func NewCreateHandler(svc domain.Service) CreateHandler {
	return &createHandler{svc: svc}
}

// Handle ...
func (h *createHandler) Handle(ctx context.Context, cmd *CreateCommand) (string, error) {
	cmd.setDefault()
	if err := cmd.validate(); err != nil {
		return "", err
	}
	return h.svc.Create(ctx, cmd.toDO())
}
