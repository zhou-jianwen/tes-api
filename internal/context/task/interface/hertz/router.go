package hertz

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"

	"code.byted.org/epscp/vetes-api/internal/context/task/application"
	"code.byted.org/epscp/vetes-api/internal/context/task/interface/hertz/handlers"
	"code.byted.org/epscp/vetes-api/pkg/consts"
	appserver "code.byted.org/epscp/vetes-api/pkg/server"
)

type register struct {
	svc *application.TaskService
}

// NewRouterRegister ...
func NewRouterRegister(taskService *application.TaskService) appserver.RouteRegister {
	return &register{
		svc: taskService,
	}
}

func (r *register) AddRoute(h route.IRouter) {
	taskGA4GH := h.Group(consts.Ga4ghAPIPrefix + "/tasks")

	taskGA4GH.POST("", func(c context.Context, ctx *app.RequestContext) {
		handlers.CreateTask(c, ctx, r.svc.TaskCommands.Create)
	})

	taskGA4GH.GET("", func(c context.Context, ctx *app.RequestContext) {
		handlers.ListTasks(c, ctx, r.svc.TaskQueries.List)
	})

	taskGA4GH.GET("/:id", func(c context.Context, ctx *app.RequestContext) {
		handlers.GetTask(c, ctx, r.svc.TaskQueries.Get)
	})

	// colon in "/{id}:cancel" conflict with hertz
	taskGA4GH.POST("/:idWithCancel", func(c context.Context, ctx *app.RequestContext) {
		handlers.CancelTask(c, ctx, r.svc.TaskCommands.Cancel)
	})

	taskOther := h.Group(consts.OtherAPIPrefix + "/tasks")

	taskOther.PATCH("/:id", func(c context.Context, ctx *app.RequestContext) {
		handlers.UpdateTask(c, ctx, r.svc.TaskCommands.Update)
	})

	taskOther.GET("/resources", func(c context.Context, ctx *app.RequestContext) {
		handlers.GatherTasksResources(c, ctx, r.svc.TaskQueries.Gather)
	})

	taskOther.GET("/accounts", func(c context.Context, ctx *app.RequestContext) {
		handlers.ListTasksAccounts(c, ctx, r.svc.TaskQueries.ListAccounts)
	})
}
