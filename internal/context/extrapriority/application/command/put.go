package command

import (
	"context"

	"github.com/GBA-BI/tes-api/internal/context/extrapriority/domain"
	"github.com/GBA-BI/tes-api/pkg/validator"
)

// PutCommand ...
type PutCommand struct {
	AccountID          string
	UserID             string
	SubmissionID       string
	RunID              string
	ExtraPriorityValue int `validate:"required"`
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
	return h.svc.Put(ctx, cmd.AccountID, cmd.UserID, cmd.SubmissionID, cmd.RunID, cmd.ExtraPriorityValue)
}
