package layout

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

type Handler struct {
	sync.RWMutex
	cfg   *conf.Configure
	tasks map[string]*task
}

func New() *Handler {
	return &Handler{
		tasks: make(map[string]*task),
	}
}

func (h *Handler) Module() string {
	return "layout"
}

func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.cfg = cfg
}

func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodPost, "/run"): h.run,
	}
}

func (h *Handler) OnConnect(*agent.Agent) {
}

func (h *Handler) OnClose(string) {
}

func (h *Handler) OnMessage(*agent.Agent, *anet.Msg) {
}
