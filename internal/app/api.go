package app

import (
	"fmt"
	"net/http"
	"server/internal/api"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/stat"
	"github.com/jkstack/jkframe/utils"
)

func handleRecovery(g *gin.Context, err any) {
	pref, _ := g.Get(api.KeyPerformance)
	pref.(*stat.Tick).Close()
	if err == nil {
		return
	}
	switch err := err.(type) {
	case api.MissingParam:
		api.ERR(g, http.StatusBadRequest, err.Error())
	case api.BadParam:
		api.ERR(g, http.StatusBadRequest, err.Error())
	case api.NotFound:
		api.ERR(g, http.StatusNotFound, err.Error())
	case api.Timeout:
		api.ERR(g, http.StatusGatewayTimeout, err.Error())
	case api.Notfound:
		api.ERR(g, http.StatusNotFound, err.Error())
	case api.InvalidType:
		api.ERR(g, http.StatusFailedDependency, err.Error())
	default:
		// TODO: trace log
		api.ERR(g, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		logging.Error("err: %v", err)
	}
}

func (app *App) ratelimit(g *gin.Context) {
	if app.blocked {
		g.Abort()
		api.ERR(g, http.StatusServiceUnavailable, "rate limit")
	}
}

func (app *App) point(g *gin.Context) {
	uri := strings.ReplaceAll(g.FullPath(), "/", "_")
	counter := app.stats.NewCounter("api_counter_" + uri + ":" + g.Request.Method)
	counter.Inc()
	tick := app.stats.NewTick("api_pref_" + uri + ":" + g.Request.Method)
	g.Set("X-PERFORMANCE", tick)
}

var number uint64

const defaultUID = "ffffffff"

func (app *App) bind(g *gin.Context) {
	// agents
	g.Set(api.KeyAgents, app.agents)

	// request-id
	next := atomic.AddUint64(&number, 1)
	uid, err := utils.UUID(8, "0123456789abcdef")
	if err != nil {
		logging.Error("generate uid for request %d failed, reset to default", next)
		uid = defaultUID
	}
	g.Set(api.KeyRequestID, fmt.Sprintf("%s-%08d-%s",
		time.Now().Format("20060102"), next%99999999, uid))
}
