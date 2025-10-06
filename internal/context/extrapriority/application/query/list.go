package query

import (
	"context"

	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

// ListQuery ...
type ListQuery struct {
	Filter *ListFilter
}

// ListFilter ...
type ListFilter struct {
	AccountID    string
	SubmissionID string
	RunID        string
}

func (q *ListQuery) setDefault() {}

func (q *ListQuery) validate() error {
	if q.Filter == nil {
		return nil
	}
	nonEmptyCnt := 0
	if q.Filter.AccountID != "" {
		nonEmptyCnt++
	}
	if q.Filter.SubmissionID != "" {
		nonEmptyCnt++
	}
	if q.Filter.RunID != "" {
		nonEmptyCnt++
	}
	if nonEmptyCnt > 1 {
		return apperrors.NewInvalidError("only one query shall be non-empty")
	}
	return nil
}

// ListHandler ...
type ListHandler interface {
	Handle(ctx context.Context, query *ListQuery) ([]*ExtraPriority, error)
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
func (h *listHandler) Handle(ctx context.Context, query *ListQuery) ([]*ExtraPriority, error) {
	query.setDefault()
	if err := query.validate(); err != nil {
		return nil, err
	}
	res, err := h.readModel.List(ctx, query.Filter)
	if err != nil {
		return nil, err
	}
	return res, nil
}
