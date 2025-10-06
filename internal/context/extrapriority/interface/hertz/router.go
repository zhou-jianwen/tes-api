package hertz

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"

	"github.com/GBA-BI/tes-api/internal/context/extrapriority/application"
	"github.com/GBA-BI/tes-api/internal/context/extrapriority/interface/hertz/handlers"
	"github.com/GBA-BI/tes-api/pkg/consts"
	appserver "github.com/GBA-BI/tes-api/pkg/server"
)

type register struct {
	svc *application.ExtraPriorityService
}

// NewRouterRegister ...
func NewRouterRegister(extraPriorityService *application.ExtraPriorityService) appserver.RouteRegister {
	return &register{
		svc: extraPriorityService,
	}
}

// AddRoute ...
func (r *register) AddRoute(h route.IRouter) {
	extraPriority := h.Group(consts.OtherAPIPrefix + "/extra_priority")

	extraPriority.PUT("", func(c context.Context, ctx *app.RequestContext) {
		handlers.PutExtraPriority(c, ctx, r.svc.ExtraPriorityCommands.Put)
	})

	extraPriority.GET("", func(c context.Context, ctx *app.RequestContext) {
		handlers.ListExtraPriority(c, ctx, r.svc.ExtraPriorityQueries.List)
	})

	extraPriority.DELETE("", func(c context.Context, ctx *app.RequestContext) {
		handlers.DeleteExtraPriority(c, ctx, r.svc.ExtraPriorityCommands.Delete)
	})
}
