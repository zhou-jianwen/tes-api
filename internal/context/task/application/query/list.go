package query

import (
	"context"

	"github.com/GBA-BI/tes-api/pkg/consts"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/utils"
	"github.com/GBA-BI/tes-api/pkg/validator"
)

const defaultPageSize = 256

// ListQuery ...
type ListQuery struct {
	View      string `validate:"oneof=MINIMAL BASIC FULL"`
	PageSize  int    `validate:"gte=0,lte=2048"`
	PageToken *utils.PageToken
	Filter    *ListFilter
}

// ListFilter ...
type ListFilter struct {
	NamePrefix     string
	State          []string `validate:"dive,oneof=QUEUED INITIALIZING RUNNING COMPLETE SYSTEM_ERROR EXECUTOR_ERROR CANCELING CANCELED"`
	ClusterID      string
	WithoutCluster bool
}

func (q *ListQuery) setDefault() {
	if q.View == "" {
		q.View = consts.MinimalView
	}
	if q.PageSize == 0 {
		q.PageSize = defaultPageSize
	}
}

func (q *ListQuery) validate() error {
	if err := validator.Validate(q); err != nil {
		return err
	}
	if q.Filter == nil {
		return nil
	}
	if q.Filter.ClusterID != "" && q.Filter.WithoutCluster {
		return apperrors.NewInvalidError("cluster_id", "without_cluster")
	}
	return nil
}

// ListHandler ...
type ListHandler interface {
	Handle(ctx context.Context, query *ListQuery) ([]*Task, *utils.PageToken, error)
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
func (h *listHandler) Handle(ctx context.Context, query *ListQuery) ([]*Task, *utils.PageToken, error) {
	query.setDefault()
	if err := query.validate(); err != nil {
		return nil, nil, err
	}

	switch query.View {
	case consts.MinimalView:
		resMinimal, nextPageToken, err := h.readModel.ListMinimal(ctx, query.PageSize, query.PageToken, query.Filter)
		if err != nil {
			return nil, nil, err
		}
		res := make([]*Task, len(resMinimal))
		for index := range resMinimal {
			res[index] = &Task{TaskBasic: TaskBasic{TaskMinimal: *resMinimal[index]}}
		}
		return res, nextPageToken, nil
	case consts.BasicView:
		resBasic, nextPageToken, err := h.readModel.ListBasic(ctx, query.PageSize, query.PageToken, query.Filter)
		if err != nil {
			return nil, nil, err
		}
		res := make([]*Task, len(resBasic))
		for index := range resBasic {
			removeSystemLogs(resBasic[index])
			res[index] = &Task{TaskBasic: *resBasic[index]}
		}
		return res, nextPageToken, nil
	case consts.FullView:
		res, nextPageToken, err := h.readModel.ListFull(ctx, query.PageSize, query.PageToken, query.Filter)
		if err != nil {
			return nil, nil, err
		}
		return res, nextPageToken, nil
	default:
		return nil, nil, apperrors.NewInvalidError("view")
	}
}
