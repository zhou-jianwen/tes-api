package query

import (
	"context"

	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
	"code.byted.org/epscp/vetes-api/pkg/validator"
)

// GatherQuery ...
type GatherQuery struct {
	Filter *GatherFilter
}

// GatherFilter ...
type GatherFilter struct {
	State       []string `validate:"dive,oneof=QUEUED INITIALIZING RUNNING COMPLETE SYSTEM_ERROR EXECUTOR_ERROR CANCELING CANCELED"`
	ClusterID   string
	WithCluster bool
	AccountID   string
	UserID      string
}

func (q *GatherQuery) setDefault() {}

func (q *GatherQuery) validate() error {
	if err := validator.Validate(q); err != nil {
		return err
	}
	if q.Filter == nil {
		return nil
	}
	if q.Filter.AccountID == "" && q.Filter.UserID != "" {
		return apperrors.NewInvalidError("empty account_id with non-empty user_id")
	}
	return nil
}

// GatherHandler ...
type GatherHandler interface {
	Handle(ctx context.Context, query *GatherQuery) (*TasksResources, error)
}

type gatherHandler struct {
	readModel ReadModel
}

var _ GatherHandler = (*gatherHandler)(nil)

// NewGatherHandler ...
func NewGatherHandler(readModel ReadModel) GatherHandler {
	return &gatherHandler{readModel: readModel}
}

// Handle ...
func (h *gatherHandler) Handle(ctx context.Context, query *GatherQuery) (*TasksResources, error) {
	query.setDefault()
	if err := query.validate(); err != nil {
		return nil, err
	}
	return h.readModel.GatherResources(ctx, query.Filter)
}
