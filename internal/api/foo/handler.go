package foo

import (
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/conf"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/stat"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Module() string {
	return "foo"
}

func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
}

func (h *Handler) HandleFuncs() map[api.Route]func(*api.GContext, *agent.Agents) {
	return map[api.Route]func(*api.GContext, *agent.Agents){
		api.MakeRoute(http.MethodGet, "/:id", "foo"): h.foo,
	}
}

func (h *Handler) OnConnect(*agent.Agent) {
}

func (h *Handler) OnClose(string) {
}

func (h *Handler) OnMessage(*agent.Agent, *anet.Msg) {
}
