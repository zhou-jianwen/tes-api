package query

import (
	"context"

	"code.byted.org/epscp/vetes-api/pkg/validator"
)

// ListQuery ...
type ListQuery struct {
	Filter *ListFilter
}

// ListFilter ...
type ListFilter struct {
}

func (q *ListQuery) setDefault() {}

func (q *ListQuery) validate() error {
	return validator.Validate(q)
}

// ListHandler ...
type ListHandler interface {
	Handle(ctx context.Context, query *ListQuery) ([]*Cluster, error)
}

type listHandler struct {
	readModel ReadModel
}

var _ ListHandler = (*listHandler)(nil)

// NewListHandler ...
func NewListHandler(readModel ReadModel) ListHandler {
	return &listHandler{readModel: readModel}
}

// Handle ...
func (h *listHandler) Handle(ctx context.Context, query *ListQuery) ([]*Cluster, error) {
	query.setDefault()
	if err := query.validate(); err != nil {
		return nil, err
	}
	return h.readModel.List(ctx, query.Filter)
}
