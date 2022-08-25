package metrics

import (
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/conf"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/stat"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Module() string {
	return "metrics"
}

func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
}

func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodGet, "/:id/static"):              h.static,
		api.MakeRoute(http.MethodGet, "/:id/dynamic"):             h.dynamic,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/usage"):       h.dynamicUsage,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/process"):     h.dynamicProcess,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/connections"): h.dynamicConnections,
	}
}

func (h *Handler) OnConnect(*agent.Agent) {
}

func (h *Handler) OnClose(string) {
}

func (h *Handler) OnMessage(*agent.Agent, *anet.Msg) {
}
