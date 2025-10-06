package handlers

import (
	"context"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/GBA-BI/tes-api/internal/context/quota/application/command"
	"github.com/GBA-BI/tes-api/internal/context/quota/application/query"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

// GetQuota get quota
//
//	@Summary		get quota
//	@Description	get quota
//	@Tags			quota
//	@Produce		application/json
//	@Router			/api/v1/quota [get]
//	@Param			global		query		bool	false	"query global quota"
//	@Param			account_id	query		string	false	"query account quota"
//	@Param			user_id		query		string	false	"query user quota"
//	@Success		200			{object}	GetQuotaResponse
//	@Failure		400			{object}	apperrors.AppError	"invalid param"
//	@Failure		404			{object}	apperrors.AppError	"not found"
//	@Failure		500			{object}	apperrors.AppError	"internal system error"
func GetQuota(c context.Context, ctx *app.RequestContext, handler query.GetHandler) {
	var req GetQuotaRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	quota, err := handler.Handle(c, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	utils.WriteHertzOKResponse(ctx, quotaDTOToVO(quota))
}

// PutQuota create or update quota
//
//	@Summary		put quota
//	@Description	put quota
//	@Tags			quota
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/api/v1/quota [put]
//	@Param			request	body		PutQuotaRequest	true	"put quota request"
//	@Success		200		{object}	PutQuotaResponse
//	@Failure		400		{object}	apperrors.AppError	"invalid param"
//	@Failure		500		{object}	apperrors.AppError	"internal system error"
func PutQuota(c context.Context, ctx *app.RequestContext, handler command.PutHandler) {
	var req PutQuotaRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(c, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}

	resp := &PutQuotaResponse{}
	utils.WriteHertzOKResponse(ctx, resp)
}

// DeleteQuota delete quota
//
//	@Summary		delete quota
//	@Description	delete quota
//	@Tags			quota
//	@Produce		application/json
//	@Router			/api/v1/quota [delete]
//	@Param			global		query		bool	false	"query global quota"
//	@Param			account_id	query		string	false	"query account quota"
//	@Param			user_id		query		string	false	"query user quota"
//	@Success		200			{object}	DeleteQuotaResponse
//	@Failure		400			{object}	apperrors.AppError	"invalid param"
//	@Failure		404			{object}	apperrors.AppError	"not found"
//	@Failure		500			{object}	apperrors.AppError	"internal system error"
func DeleteQuota(c context.Context, ctx *app.RequestContext, handler command.DeleteHandler) {
	var req DeleteQuotaRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(c, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}

	resp := &DeleteQuotaResponse{}
	utils.WriteHertzOKResponse(ctx, resp)
}
