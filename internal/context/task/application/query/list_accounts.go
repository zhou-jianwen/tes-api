package query

import "context"

// ListAccountsHandler ...
type ListAccountsHandler interface {
	Handle(ctx context.Context) ([]*AccountInfo, error)
}

type listAccountsHandler struct {
	readModel ReadModel
}

var _ ListAccountsHandler = (*listAccountsHandler)(nil)

// NewListAccountsHandler ...
func NewListAccountsHandler(readModel ReadModel) ListAccountsHandler {
	return &listAccountsHandler{readModel: readModel}
}

func (h *listAccountsHandler) Handle(ctx context.Context) ([]*AccountInfo, error) {
	return h.readModel.ListAccounts(ctx)
}
