package command

import (
	"context"

	"code.byted.org/epscp/vetes-api/internal/context/task/domain"
)

// UpdateHandler ...
type UpdateHandler interface {
	Handle(ctx context.Context, cmd *UpdateCommand) error
}

type updateHandler struct {
	svc domain.Service
}

var _ UpdateHandler = (*updateHandler)(nil)

// NewUpdateHandler ...
func NewUpdateHandler(svc domain.Service) UpdateHandler {
	return &updateHandler{svc: svc}
}

// Handle ...
func (h *updateHandler) Handle(ctx context.Context, cmd *UpdateCommand) error {
	cmd.setDefault()
	if err := cmd.validate(); err != nil {
		return err
	}
	return h.svc.Update(ctx, cmd.ID, cmd.State, cmd.ClusterID, taskLogs(cmd.Logs).toDO())
}
