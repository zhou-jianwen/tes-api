package handlers

import (
	"context"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/GBA-BI/tes-api/internal/context/extrapriority/application/command"
	"github.com/GBA-BI/tes-api/internal/context/extrapriority/application/query"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

// PutExtraPriority create or update extra priority on tasks
//
//	@Summary		create or update tasks extra priority
//	@Description	create or update extra priority on tasks
//	@Tags			priority
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/api/v1/extra_priority [put]
//	@Param			account_id		query		string					false	"query account id"
//	@Param			user_id			query		string					false	"query user id"
//	@Param			submission_id	query		string					false	"query submission id"
//	@Param			run_id			query		string					false	"query run id"
//	@Param			request			body		PutExtraPriorityRequest	true	"put tasks extra priority request"
//	@Success		200				{object}	PutExtraPriorityResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func PutExtraPriority(c context.Context, ctx *app.RequestContext, handler command.PutHandler) {
	var req PutExtraPriorityRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(c, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	utils.WriteHertzOKResponse(ctx, &PutExtraPriorityResponse{})
}

// ListExtraPriority list extra priority on tasks
//
//	@Summary		list tasks extra priority
//	@Description	list extra priority on tasks
//	@Tags			priority
//	@Produce		application/json
//	@Router			/api/v1/extra_priority [get]
//	@Param			account_id		query		string	false	"query account id"
//	@Param			submission_id	query		string	false	"query submission id"
//	@Param			run_id			query		string	false	"query run id"
//	@Success		200				{object}	ListExtraPriorityResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListExtraPriority(c context.Context, ctx *app.RequestContext, handler query.ListHandler) {
	var req ListExtraPriorityRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	extraPriorities, err := handler.Handle(c, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	resp := make(ListExtraPriorityResponse, 0, len(extraPriorities))
	for _, extraPriority := range extraPriorities {
		resp = append(resp, extraPriorityDTOToVO(extraPriority))
	}
	utils.WriteHertzOKResponse(ctx, resp)
}

// DeleteExtraPriority delete extra priority on tasks
//
//	@Summary		delete tasks extra priority
//	@Description	delete extra priority on tasks
//	@Tags			priority
//	@Produce		application/json
//	@Router			/api/v1/extra_priority [delete]
//	@Param			account_id		query		string	false	"query account id"
//	@Param			user_id			query		string	false	"query user id"
//	@Param			submission_id	query		string	false	"query submission id"
//	@Param			run_id			query		string	false	"query run id"
//	@Success		200				{object}	DeleteExtraPriorityResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		404				{object}	apperrors.AppError	"not found"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func DeleteExtraPriority(c context.Context, ctx *app.RequestContext, handler command.DeleteHandler) {
	var req DeleteExtraPriorityRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(c, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	utils.WriteHertzOKResponse(ctx, &DeleteExtraPriorityResponse{})
}
