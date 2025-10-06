package command

import "code.byted.org/epscp/vetes-api/internal/context/cluster/domain"

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
