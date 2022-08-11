package app

import (
	"fmt"
	"io"
	"net/http"
	"server/internal/agent"
	"server/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/logging"
)

func bindRecovery(g *gin.RouterGroup) {
	g.Use(gin.CustomRecoveryWithWriter(io.Discard, func(g *gin.Context, err any) {
		switch err := err.(type) {
		case api.MissingParam:
			api.ERR(g, http.StatusBadRequest, err.Error())
		case api.BadParam:
			api.ERR(g, http.StatusBadRequest, err.Error())
		case api.NotFound:
			api.ERR(g, http.StatusNotFound, err.Error())
		case api.Timeout:
			api.ERR(g, http.StatusBadGateway, err.Error())
		case api.Notfound:
			api.ERR(g, http.StatusNotFound, err.Error())
		default:
			api.ERR(g, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			logging.Error("err: %v", err)
		}
	}))
}

func (app *App) reg(g *gin.RouterGroup, module string, route api.Route, cb func(*agent.Agents, *gin.Context)) {
	g.Handle(route.Method, route.Uri, func(g *gin.Context) {
		if app.blocked {
			api.ERR(g, http.StatusServiceUnavailable, "rate limit")
			return
		}
		counter := app.stats.NewCounter("api_counter_" + module + "_" + route.MetricName)
		counter.Inc()
		tick := app.stats.NewTick("api_pref_" + module + "_" + route.MetricName)
		defer tick.Close()
		cb(app.agents, g)
	})
}
