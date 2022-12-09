package script

import (
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/conf"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/stat"
)

// Handler api handler
type Handler struct {
	sync.RWMutex
	cfg *conf.Configure
}

// New create api handler
func New() *Handler {
	return &Handler{}
}

// Module get module name
func (h *Handler) Module() string {
	return "script"
}

// Init initialize module
func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.cfg = cfg
}

// HandleFuncs get funcs
func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodPost, "/:id/run"): h.run,
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
