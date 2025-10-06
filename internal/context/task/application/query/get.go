package query

import (
	"context"

	"github.com/GBA-BI/tes-api/pkg/consts"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/validator"
)

// GetQuery ...
type GetQuery struct {
	ID   string `validate:"required"`
	View string `validate:"oneof=MINIMAL BASIC FULL"`
}

func (q *GetQuery) setDefault() {
	if q.View == "" {
		q.View = consts.MinimalView
	}
}

func (q *GetQuery) validate() error {
	return validator.Validate(q)
}

// GetHandler ...
type GetHandler interface {
	Handle(ctx context.Context, query *GetQuery) (*Task, error)
}

type getHandler struct {
	readModel ReadModel
}

var _ GetHandler = (*getHandler)(nil)

// NewGetHandler ...
func NewGetHandler(readModel ReadModel) GetHandler {
	return &getHandler{readModel: readModel}
}

// Handle ...
func (h *getHandler) Handle(ctx context.Context, query *GetQuery) (*Task, error) {
	query.setDefault()
	if err := query.validate(); err != nil {
		return nil, err
	}

	switch query.View {
	case consts.MinimalView:
		resMinimal, err := h.readModel.GetMinimal(ctx, query.ID)
		if err != nil {
			return nil, err
		}
		return &Task{TaskBasic: TaskBasic{TaskMinimal: *resMinimal}}, nil
	case consts.BasicView:
		resBasic, err := h.readModel.GetBasic(ctx, query.ID)
		if err != nil {
			return nil, err
		}
		removeSystemLogs(resBasic)
		return &Task{TaskBasic: *resBasic}, nil
	case consts.FullView:
		res, err := h.readModel.GetFull(ctx, query.ID)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, apperrors.NewInvalidError("view")
	}
}

func removeSystemLogs(task *TaskBasic) {
	if task == nil {
		return
	}
	for index := range task.Logs {
		if task.Logs[index] != nil {
			task.Logs[index].SystemLogs = nil
		}
	}
}
