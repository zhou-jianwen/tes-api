package hertz

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"

	"code.byted.org/epscp/vetes-api/internal/context/cluster/application"
	"code.byted.org/epscp/vetes-api/internal/context/cluster/interface/hertz/handlers"
	"code.byted.org/epscp/vetes-api/pkg/consts"
	appserver "code.byted.org/epscp/vetes-api/pkg/server"
)

type register struct {
	svc *application.ClusterService
}

// NewRouterRegister ...
func NewRouterRegister(clusterService *application.ClusterService) appserver.RouteRegister {
	return &register{
		svc: clusterService,
	}
}

// AddRoute ...
func (r *register) AddRoute(h route.IRouter) {
	cluster := h.Group(consts.OtherAPIPrefix + "/clusters")

	cluster.PUT("/:id", func(c context.Context, ctx *app.RequestContext) {
		handlers.PutCluster(c, ctx, r.svc.ClusterCommands.Put)
	})

	cluster.GET("", func(c context.Context, ctx *app.RequestContext) {
		handlers.ListClusters(c, ctx, r.svc.ClusterQueries.List)
	})

	cluster.DELETE("/:id", func(c context.Context, ctx *app.RequestContext) {
		handlers.DeleteCluster(c, ctx, r.svc.ClusterCommands.Delete)
	})
}
