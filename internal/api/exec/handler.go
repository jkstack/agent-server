package exec

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
	cfg   *conf.Configure
	tasks map[string]*tasks // agent id => tasks
}

// New create api handler
func New() *Handler {
	return &Handler{
		tasks: make(map[string]*tasks),
	}
}

// Module get module name
func (h *Handler) Module() string {
	return "exec"
}

// Init initialize module
func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.cfg = cfg
}

// HandleFuncs get funcs
func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodPost, "/:id/run"):         h.run,
		api.MakeRoute(http.MethodGet, "/:id/status/:pid"):  h.status,
		api.MakeRoute(http.MethodGet, "/:id/pty/:pid"):     h.pty,
		api.MakeRoute(http.MethodDelete, "/:id/kill/:pid"): h.kill,
		api.MakeRoute(http.MethodGet, "/:id/ps"):           h.ps,
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

func (h *Handler) getTasksOrCreate(id string) *tasks {
	h.Lock()
	defer h.Unlock()
	if ts, ok := h.tasks[id]; ok {
		return ts
	}
	ts := newTasks(id)
	h.tasks[id] = ts
	return ts
}

func (h *Handler) getTasks(id string) *tasks {
	h.RLock()
	defer h.RUnlock()
	return h.tasks[id]
}
