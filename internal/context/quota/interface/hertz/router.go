package hertz

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"

	"code.byted.org/epscp/vetes-api/internal/context/quota/application"
	"code.byted.org/epscp/vetes-api/internal/context/quota/interface/hertz/handlers"
	"code.byted.org/epscp/vetes-api/pkg/consts"
	appserver "code.byted.org/epscp/vetes-api/pkg/server"
)

type register struct {
	svc *application.QuotaService
}

// NewRouterRegister ...
func NewRouterRegister(quotaService *application.QuotaService) appserver.RouteRegister {
	return &register{
		svc: quotaService,
	}
}

// AddRoute ...
func (r *register) AddRoute(h route.IRouter) {
	quota := h.Group(consts.OtherAPIPrefix + "/quota")

	quota.GET("", func(c context.Context, ctx *app.RequestContext) {
		handlers.GetQuota(c, ctx, r.svc.QuotaQueries.Get)
	})

	quota.PUT("", func(c context.Context, ctx *app.RequestContext) {
		handlers.PutQuota(c, ctx, r.svc.QuotaCommands.Put)
	})

	quota.DELETE("", func(c context.Context, ctx *app.RequestContext) {
		handlers.DeleteQuota(c, ctx, r.svc.QuotaCommands.Delete)
	})
}
