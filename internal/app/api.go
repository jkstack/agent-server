package app

import (
	"fmt"
	"net/http"
	"server/internal/agent"
	lapi "server/internal/api"
	"strings"

	"github.com/jkstack/jkframe/api"
	"github.com/jkstack/jkframe/logging"
)

func (app *App) reg(uri string, cb func(*agent.Agents, *api.Context)) {
	http.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		if app.blocked {
			http.Error(w, "rate limit", http.StatusServiceUnavailable)
			return
		}
		statName := strings.TrimSuffix(uri, "/")
		statName = strings.ReplaceAll(statName, "/", ":")
		counter := app.stats.NewCounter("api_counter" + statName)
		counter.Inc()
		tick := app.stats.NewTick("api_pref" + statName)
		defer tick.Close()
		ctx := api.NewContext(w, r)
		defer func() {
			if err := recover(); err != nil {
				switch err := err.(type) {
				case api.MissingParam:
					ctx.ERR(http.StatusBadRequest, err.Error())
				case api.BadParam:
					ctx.ERR(http.StatusBadRequest, err.Error())
				case api.NotFound:
					ctx.ERR(http.StatusNotFound, err.Error())
				case api.Timeout:
					ctx.ERR(http.StatusBadGateway, err.Error())
				case lapi.Notfound:
					ctx.ERR(http.StatusNotFound, err.Error())
				default:
					ctx.ERR(http.StatusInternalServerError, fmt.Sprintf("%v", err))
					logging.Error("err: %v", err)
				}
			}
		}()
		cb(app.agents, ctx)
	})
}
