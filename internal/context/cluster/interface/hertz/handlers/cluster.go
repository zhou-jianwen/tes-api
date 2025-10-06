package handlers

import (
	"context"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/GBA-BI/tes-api/internal/context/cluster/application/command"
	"github.com/GBA-BI/tes-api/internal/context/cluster/application/query"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

// PutCluster reports cluster info
//
//	@Summary		put cluster
//	@Description	put cluster
//	@Tags			cluster
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/api/v1/clusters/{id} [put]
//	@Param			id		path		string				true	"put cluster id"
//	@Param			request	body		PutClusterRequest	true	"put cluster request"
//	@Success		200		{object}	PutClusterResponse
//	@Failure		400		{object}	apperrors.AppError	"invalid param"
//	@Failure		500		{object}	apperrors.AppError	"internal system error"
func PutCluster(c context.Context, ctx *app.RequestContext, handler command.PutHandler) {
	var req PutClusterRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(c, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}

	resp := &PutClusterResponse{}
	utils.WriteHertzOKResponse(ctx, resp)
}

// ListClusters lists cluster
//
//	@Summary		list clusters
//	@Description	list clusters
//	@Tags			cluster
//	@Produce		application/json
//	@Router			/api/v1/clusters [get]
//	@Success		200	{object}	ListClustersResponse
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func ListClusters(c context.Context, ctx *app.RequestContext, handler query.ListHandler) {
	var req ListClustersRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	clusters, err := handler.Handle(c, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	resp := make(ListClustersResponse, 0, len(clusters))
	for _, cluster := range clusters {
		resp = append(resp, clusterDTOToVO(cluster))
	}
	utils.WriteHertzOKResponse(ctx, resp)
}

// DeleteCluster delete cluster
//
//	@Summary		delete cluster
//	@Description	delete cluster
//	@Tags			cluster
//	@Produce		application/json
//	@Router			/api/v1/clusters/{id} [delete]
//	@Param			id	path		string	true	"delete cluster id"
//	@Success		200	{object}	DeleteClusterResponse
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		404	{object}	apperrors.AppError	"not found"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func DeleteCluster(c context.Context, ctx *app.RequestContext, handler command.DeleteHandler) {
	var req DeleteClusterRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(c, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}

	resp := &DeleteClusterResponse{}
	utils.WriteHertzOKResponse(ctx, resp)
}
