package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/google/uuid"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
	"github.com/hertz-contrib/requestid"
	"github.com/hertz-contrib/swagger"
	swaggerfiles "github.com/swaggo/files"

	"github.com/GBA-BI/tes-api/internal/apiserver/middlewares/hertz"
	"github.com/GBA-BI/tes-api/pkg/consts"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	appserver "github.com/GBA-BI/tes-api/pkg/server"
	"github.com/GBA-BI/tes-api/pkg/serviceinfo"
	"github.com/GBA-BI/tes-api/pkg/utils"
	"github.com/GBA-BI/tes-api/pkg/version"
)

func setupHTTPServer(opts *appserver.HTTPOptions, registers ...appserver.RouteRegister) *server.Hertz {
	httpOptions := []config.Option{
		server.WithHostPorts(fmt.Sprintf(":%d", opts.Port)),
		server.WithMaxRequestBodySize(opts.MaxRequestBodySize),
		server.WithTracer(prometheus.NewServerTracer(fmt.Sprintf(":%d", opts.MetricsPort), "/metrics", prometheus.WithEnableGoCollector(true))),
	}

	httpServer := server.Default(httpOptions...)
	setupMiddlewares(httpServer)
	setupRouter(httpServer)
	for _, r := range registers {
		r.AddRoute(httpServer)
	}
	return httpServer
}

func setupMiddlewares(h *server.Hertz) {
	h.Use(
		requestid.New(
			requestid.WithGenerator(func(_ context.Context, _ *app.RequestContext) string {
				return uuid.New().String()
			}),
			// set custom header for request id
			requestid.WithCustomHeaderStrKey(consts.XRequestIDKey),
		),
		hertz.Logger(),
	)
}

func setupRouter(h *server.Hertz) {
	h.GET("/ping", PingHandler)
	h.GET("/version", VersionHandler)
	h.GET(consts.Ga4ghAPIPrefix+"/service-info", ServiceInfoHandler)
	url := swagger.URL("/swagger/doc.json") // The url pointing to API definition
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerfiles.Handler, url))

	h.NoRoute(func(_ context.Context, ctx *app.RequestContext) {
		utils.WriteHertzErrorResponse(ctx, apperrors.NewHertzRouteNotFoundError(ctx))
	})
}

// PingHandler ping handler
//
//	@Summary		ping
//	@Description	ping
//	@Produce		application/json
//	@Router			/ping [get]
//	@Success		200
func PingHandler(_ context.Context, ctx *app.RequestContext) {
	ctx.JSON(http.StatusOK, map[string]string{
		"ping": "pong",
	})
}

// VersionHandler version handler
//
//	@Summary		version
//	@Description	version
//	@Produce		application/json
//	@Router			/version [get]
//	@Success		200	{object}	version.Info
func VersionHandler(_ context.Context, ctx *app.RequestContext) {
	ctx.JSON(http.StatusOK, version.Get())
}

// ServiceInfoHandler ga4gh serviceInfo handler
//
//	@Summary		service-info
//	@Description	ga4gh service-info
//	@Produce		application/json
//	@Router			/api/ga4gh/tes/v1/service-info [get]
//	@Success		200	{object}	serviceinfo.ServiceInfo
func ServiceInfoHandler(_ context.Context, ctx *app.RequestContext) {
	ctx.JSON(http.StatusOK, serviceinfo.DefaultServiceInfo)
}
