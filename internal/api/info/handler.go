package info

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
	version    string
	isBlocking *bool
}

func New(version string, isBlocking *bool) *Handler {
	return &Handler{
		version:    version,
		isBlocking: isBlocking,
	}
}

func (h *Handler) Module() string {
	return "info"
}

func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
}

func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodGet, "/server"): h.server,
	}
}

func (h *Handler) OnConnect(*agent.Agent) {
}

func (h *Handler) OnClose(string) {
}

func (h *Handler) OnMessage(*agent.Agent, *anet.Msg) {
}
