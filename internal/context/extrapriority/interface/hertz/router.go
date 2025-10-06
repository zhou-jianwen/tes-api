package hertz

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"

	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/application"
	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/interface/hertz/handlers"
	"code.byted.org/epscp/vetes-api/pkg/consts"
	appserver "code.byted.org/epscp/vetes-api/pkg/server"
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
