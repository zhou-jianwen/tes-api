package command

import "github.com/GBA-BI/tes-api/internal/context/task/domain"

// Commands ...
type Commands struct {
	Create CreateHandler
	Cancel CancelHandler
	Update UpdateHandler
}

// NewCommands ...
func NewCommands(svc domain.Service) *Commands {
	return &Commands{
		Create: NewCreateHandler(svc),
		Cancel: NewCancelHandler(svc),
		Update: NewUpdateHandler(svc),
	}
}
