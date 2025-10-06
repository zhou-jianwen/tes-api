package hertz

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/requestid"

	"github.com/GBA-BI/tes-api/pkg/consts"
)

// Logger log middleware for hertz.
func Logger() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		if !bytes.HasPrefix(ctx.Request.Path(), []byte(consts.Ga4ghAPIPrefix)) &&
			!bytes.HasPrefix(ctx.Request.Path(), []byte(consts.OtherAPIPrefix)) {
			return
		}

		startTime := time.Now()
		ctx.Next(c)
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		requestID := requestid.Get(ctx)
		reqMethod := string(ctx.Request.Method())
		reqURI := string(ctx.Request.RequestURI())
		statusCode := ctx.Response.StatusCode()
		clientIP := ctx.ClientIP()

		applog.CtxInfow(c, "process request",
			"status_code", statusCode,
			"req_action", extractAction(reqMethod, reqURI),
			"latency_time", latencyTime,
			"client_ip", clientIP,
			"req_method", reqMethod,
			"req_uri", reqURI,
			"request_id", requestID)
	}
}

func extractAction(reqMethod, reqURI string) string {
	uri, err := url.ParseRequestURI(reqURI)
	if err != nil {
		applog.Errorw("invalid requestURI", "err", err)
		return "<invalid action>"
	}
	path := uri.Path
	switch {
	case serviceInfoRegexp.MatchString(path) && reqMethod == http.MethodGet:
		return "GetServiceInfo"
	case listTasksCreateTaskRegexp.MatchString(path):
		switch reqMethod {
		case http.MethodGet:
			return "ListTasks"
		case http.MethodPost:
			return "CreateTask"
		}
	case getTaskRegexp.MatchString(path) && reqMethod == http.MethodGet:
		return "GetTask"
	case cancelTaskRegexp.MatchString(path) && reqMethod == http.MethodPost:
		return "CancelTask"
	case updateTaskRegexp.MatchString(path) && reqMethod == http.MethodPatch:
		return "UpdateTask"
	case gatherTasksResourcesRegexp.MatchString(path) && reqMethod == http.MethodGet:
		return "GatherTasksResources"
	case listTasksAccountsRegexp.MatchString(path) && reqMethod == http.MethodGet:
		return "ListTasksAccounts"
	case putDeleteClusterRegexp.MatchString(path):
		switch reqMethod {
		case http.MethodPut:
			return "PutCluster"
		case http.MethodDelete:
			return "DeleteCluster"
		}
	case listClustersRegexp.MatchString(path) && reqMethod == http.MethodGet:
		return "ListClusters"
	case quotaRegexp.MatchString(path):
		switch reqMethod {
		case http.MethodGet:
			return "GetQuota"
		case http.MethodPut:
			return "PutQuota"
		case http.MethodDelete:
			return "DeleteQuota"
		}
	case extraPriorityRegexp.MatchString(path):
		switch reqMethod {
		case http.MethodGet:
			return "GetExtraPriority"
		case http.MethodPut:
			return "PutExtraPriority"
		case http.MethodDelete:
			return "DeleteExtraPriority"
		}
	}
	return "<invalid action>"
}

var (
	serviceInfoRegexp          = regexp.MustCompile(fmt.Sprintf("^%s/service-info$", consts.Ga4ghAPIPrefix))
	listTasksCreateTaskRegexp  = regexp.MustCompile(fmt.Sprintf("^%s/tasks$", consts.Ga4ghAPIPrefix))
	getTaskRegexp              = regexp.MustCompile(fmt.Sprintf("^%s/tasks/task-[a-z0-9]+$", consts.Ga4ghAPIPrefix))
	cancelTaskRegexp           = regexp.MustCompile(fmt.Sprintf("^%s/tasks/task-[a-z0-9]+:cancel$", consts.Ga4ghAPIPrefix))
	updateTaskRegexp           = regexp.MustCompile(fmt.Sprintf("^%s/tasks/task-[a-z0-9]+$", consts.OtherAPIPrefix))
	gatherTasksResourcesRegexp = regexp.MustCompile(fmt.Sprintf("^%s/tasks/resources$", consts.OtherAPIPrefix))
	listTasksAccountsRegexp    = regexp.MustCompile(fmt.Sprintf("^%s/tasks/accounts$", consts.OtherAPIPrefix))
	putDeleteClusterRegexp     = regexp.MustCompile(fmt.Sprintf("^%s/clusters/.+", consts.OtherAPIPrefix))
	listClustersRegexp         = regexp.MustCompile(fmt.Sprintf("^%s/clusters$", consts.OtherAPIPrefix))
	quotaRegexp                = regexp.MustCompile(fmt.Sprintf("^%s/quota$", consts.OtherAPIPrefix))
	extraPriorityRegexp        = regexp.MustCompile(fmt.Sprintf("^%s/extra_priority$", consts.OtherAPIPrefix))
)
