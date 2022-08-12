package app

import (
	"fmt"
	"io"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/logging"
)

func bindRecovery(g *gin.RouterGroup) {
	g.Use(gin.CustomRecoveryWithWriter(io.Discard, func(g *gin.Context, err any) {
		if err == nil {
			return
		}
		ac, _ := g.Get("X-GContext")
		ctx := ac.(*api.GContext)
		switch err := err.(type) {
		case api.MissingParam:
			ctx.ERR(http.StatusBadRequest, err.Error())
		case api.BadParam:
			ctx.ERR(http.StatusBadRequest, err.Error())
		case api.NotFound:
			ctx.ERR(http.StatusNotFound, err.Error())
		case api.Timeout:
			ctx.ERR(http.StatusBadGateway, err.Error())
		case api.Notfound:
			ctx.ERR(http.StatusNotFound, err.Error())
		default:
			ctx.ERR(http.StatusInternalServerError, fmt.Sprintf("%v", err))
			logging.Error("err: %v", err)
		}
	}))
}

func (app *App) regAPI(g *gin.RouterGroup, module string, route api.Route, cb func(*api.GContext, *agent.Agents)) {
	logging.Info("route => %s /api/%s%s", route.Method, module, route.Uri)
	g.Handle(route.Method, route.Uri, func(g *gin.Context) {
		ctx := api.New(g)
		if app.blocked {
			ctx.ERR(http.StatusServiceUnavailable, "rate limit")
			return
		}
		counter := app.stats.NewCounter("api_counter_" + module + "_" + route.MetricName + ":" + route.Method)
		counter.Inc()
		tick := app.stats.NewTick("api_pref_" + module + "_" + route.MetricName + ":" + route.Method)
		defer tick.Close()
		body := ctx.RequestBody()
		if len(body) > 0 {
			logging.Info("REQ [%s] <= %s \"%s\"\n%s",
				ctx.ReqID(), ctx.RemoteIP(), ctx.Request(), body)
		} else {
			logging.Info("REQ [%s] <= %s \"%s\"",
				ctx.ReqID(), ctx.RemoteIP(), ctx.Request())
		}
		begin := time.Now()
		cb(ctx, app.agents)
		logging.Info("REP [%s] => %s %d %f",
			ctx.ReqID(), ctx.RemoteIP(), ctx.ContentLength(), time.Since(begin).Seconds())
	})
}
