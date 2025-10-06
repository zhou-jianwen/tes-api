package command

import "github.com/GBA-BI/tes-api/internal/context/quota/domain"

// Commands ...
type Commands struct {
	Put    PutHandler
	Delete DeleteHandler
}

// NewCommands ...
func NewCommands(svc domain.Service) *Commands {
	return &Commands{
		Put:    NewPutHandler(svc),
		Delete: NewDeleteHandler(svc),
	}
}
