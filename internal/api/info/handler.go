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

// Handler api handler
type Handler struct {
	version    string
	isBlocking *bool
}

// New create api handler
func New(version string, isBlocking *bool) *Handler {
	return &Handler{
		version:    version,
		isBlocking: isBlocking,
	}
}

// Module get module name
func (h *Handler) Module() string {
	return "info"
}

// Init initialize module
func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
}

// HandleFuncs get funcs
func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodGet, "/server"): h.server,
	}
}

// OnConnect agent connect callback
func (h *Handler) OnConnect(*agent.Agent) {
}

// OnClose agent connection closed callback
func (h *Handler) OnClose(string) {
}

// OnMessage received agent message callback
func (h *Handler) OnMessage(*agent.Agent, *anet.Msg) {
}
