package app

import (
	"fmt"
	"net/http"
	"server/internal/api"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/stat"
)

func handleRecovery(g *gin.Context, err any) {
	pref, _ := g.Get(api.KeyPerformance)
	pref.(*stat.Tick).Close()
	if err == nil {
		return
	}
	ctx := api.GetG(g)
	switch err := err.(type) {
	case api.MissingParam:
		ctx.ERR(http.StatusBadRequest, err.Error())
	case api.BadParam:
		ctx.ERR(http.StatusBadRequest, err.Error())
	case api.Timeout:
		ctx.ERR(http.StatusGatewayTimeout, err.Error())
	case api.Notfound:
		ctx.ERR(http.StatusNotFound, err.Error())
	case api.InvalidType:
		ctx.ERR(http.StatusFailedDependency, err.Error())
	default:
		// TODO: trace log
		ctx.ERR(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		logging.Error("err: %v", err)
	}
}

func (app *App) ratelimit(g *gin.Context) {
	if app.blocked {
		g.Abort()
		api.GetG(g).ERR(http.StatusServiceUnavailable, "rate limit")
	}
}

func (app *App) point(g *gin.Context) {
	uri := strings.ReplaceAll(g.FullPath(), "/", "_")
	counter := app.stats.NewCounter("api_counter_" + uri + ":" + g.Request.Method)
	counter.Inc()
	tick := app.stats.NewTick("api_pref_" + uri + ":" + g.Request.Method)
	g.Set(api.KeyPerformance, tick)
}

func (app *App) context(g *gin.Context) {
	ctx := api.New(g, app.agents)
	g.Set(api.KeyGContext, ctx)
}
