package handlers

import (
	"context"
	"strings"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/GBA-BI/tes-api/internal/context/task/application/command"
	"github.com/GBA-BI/tes-api/internal/context/task/application/query"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

// CreateTask create task
//
//	@Summary		create task
//	@Description	create task
//	@Tags			task
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/api/ga4gh/tes/v1/tasks [post]
//	@Param			request	body		CreateTaskRequest	true	"create task request"
//	@Success		200		{object}	CreateTaskResponse
//	@Failure		400		{object}	apperrors.AppError	"invalid param"
//	@Failure		500		{object}	apperrors.AppError	"internal system error"
func CreateTask(c context.Context, ctx *app.RequestContext, handler command.CreateHandler) {
	var req CreateTaskRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	id, err := handler.Handle(c, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}

	resp := &CreateTaskResponse{ID: id}
	utils.WriteHertzOKResponse(ctx, resp)
}

// ListTasks list tasks
//
//	@Summary		list tasks
//	@Description	list tasks
//	@Tags			task
//	@Produce		application/json
//	@Router			/api/ga4gh/tes/v1/tasks [get]
//	@Param			name_prefix		query		string		false	"query name prefix"
//	@Param			page_size		query		int			false	"query page size"	maximum(2048)	default(256)
//	@Param			page_token		query		string		false	"query page token"
//	@Param			view			query		string		false	"query view"	Enums(MINIMAL,BASIC,FULL)	default(MINIMAL)
//	@Param			state			query		[]string	false	"query state array"
//	@Param			cluster_id		query		string		false	"query cluster id"
//	@Param			without_cluster	query		bool		false	"query without cluster"
//	@Success		200				{object}	ListTasksResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListTasks(c context.Context, ctx *app.RequestContext, handler query.ListHandler) {
	var req ListTasksRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	qry, err := req.toDTO()
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	tasks, nextPageToken, err := handler.Handle(c, qry)
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}

	items := make([]*Task, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, taskDTOToVO(task))
	}
	resp := &ListTasksResponse{
		Tasks:         items,
		NextPageToken: utils.GenPageToken(nextPageToken),
	}
	utils.WriteHertzOKResponse(ctx, resp)
}

// GetTask get task
//
//	@Summary		get task
//	@Description	get task by id
//	@Tags			task
//	@Produce		application/json
//	@Router			/api/ga4gh/tes/v1/tasks/{id} [get]
//	@Param			id		path		string	true	"get task id"
//	@Param			view	query		string	false	"query view"	Enums(MINIMAL,BASIC,FULL)	default(MINIMAL)
//	@Success		200		{object}	GetTaskResponse
//	@Failure		400		{object}	apperrors.AppError	"invalid param"
//	@Failure		404		{object}	apperrors.AppError	"not found"
//	@Failure		500		{object}	apperrors.AppError	"internal system error"
func GetTask(c context.Context, ctx *app.RequestContext, handler query.GetHandler) {
	var req GetTaskRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	task, err := handler.Handle(c, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	utils.WriteHertzOKResponse(ctx, &GetTaskResponse{Task: taskDTOToVO(task)})
}

// CancelTask cancel task
//
//	@Summary		cancel task
//	@Description	cancel task by id
//	@Tags			task
//	@Produce		application/json
//	@Router			/api/ga4gh/tes/v1/tasks/{id}:cancel [post]
//	@Param			id	path		string	true	"cancel task id"
//	@Success		200	{object}	CancelTaskResponse
//	@Failure		400	{object}	apperrors.AppError	"invalid param or cannot execute"
//	@Failure		404	{object}	apperrors.AppError	"not found"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func CancelTask(c context.Context, ctx *app.RequestContext, handler command.CancelHandler) {
	// colon in "/{id}:cancel" conflict with hertz
	idWithCancel := ctx.Param("idWithCancel")
	if !strings.HasSuffix(idWithCancel, ":cancel") {
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzRouteNotFoundError(ctx))
		return
	}
	id := idWithCancel[:len(idWithCancel)-7] // remove ":cancel"

	var req CancelTaskRequest
	req.ID = id
	// no bind

	if err := handler.Handle(c, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	utils.WriteHertzOKResponse(ctx, &CancelTaskResponse{})
}

// UpdateTask update task
//
//	@Summary		update task
//	@Description	update task by id
//	@Tags			task
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/api/v1/tasks/{id} [patch]
//	@Param			id		path		string				true	"update task id"
//	@Param			request	body		UpdateTaskRequest	true	"update task request"
//	@Success		200		{object}	UpdateTaskResponse
//	@Failure		400		{object}	apperrors.AppError	"invalid param or cannot execute"
//	@Failure		404		{object}	apperrors.AppError	"not found"
//	@Failure		500		{object}	apperrors.AppError	"internal system error"
func UpdateTask(c context.Context, ctx *app.RequestContext, handler command.UpdateHandler) {
	var req UpdateTaskRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	cmd, err := req.toDTO()
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}

	if err = handler.Handle(c, cmd); err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	utils.WriteHertzOKResponse(ctx, &UpdateTaskResponse{})
}

// GatherTasksResources gather tasks resources
//
//	@Summary		gather tasks resources
//	@Description	gather tasks resources
//	@Tags			task
//	@Produce		application/json
//	@Router			/api/v1/tasks/resources [get]
//	@Param			state			query		[]string	false	"query state array"
//	@Param			cluster_id		query		string		false	"query cluster id"
//	@Param			with_cluster	query		bool		false	"query with cluster"
//	@Param			account_id		query		string		false	"query account id"
//	@Param			user_id			query		string		false	"query user id"
//	@Success		200				{object}	GatherTasksResourcesResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func GatherTasksResources(c context.Context, ctx *app.RequestContext, handler query.GatherHandler) {
	var req GatherTasksResourcesRequest
	if err := ctx.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzBindError(err))
		return
	}

	res, err := handler.Handle(c, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	utils.WriteHertzOKResponse(ctx, tasksResourcesDTOToVO(res))
}

// ListTasksAccounts list tasks accounts
//
//	@Summary		list tasks accounts
//	@Description	list tasks accounts
//	@Tags			task
//	@Produce		application/json
//	@Router			/api/v1/tasks/accounts [get]
//	@Success		200	{object}	ListTasksAccountsResponse
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func ListTasksAccounts(c context.Context, ctx *app.RequestContext, handler query.ListAccountsHandler) {
	accountInfos, err := handler.Handle(c)
	if err != nil {
		utils.WriteHertzErrorResponse(ctx, err)
		return
	}
	items := make([]*AccountInfo, 0, len(accountInfos))
	for _, accountInfo := range accountInfos {
		items = append(items, taskAccountInfoDTOToVO(accountInfo))
	}
	utils.WriteHertzOKResponse(ctx, &ListTasksAccountsResponse{Accounts: items})
}
