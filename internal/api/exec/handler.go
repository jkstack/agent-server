package exec

import (
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/conf"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/stat"
)

type Handler struct {
	sync.RWMutex
	cfg   *conf.Configure
	tasks map[int]*task
}

func New() *Handler {
	h := &Handler{
		tasks: make(map[int]*task),
	}
	go h.clean()
	return h
}

func (h *Handler) Module() string {
	return "exec"
}

func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.cfg = cfg
}

func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodPost, "/:id/run"):        h.run,
		api.MakeRoute(http.MethodGet, "/:id/status/:pid"): h.status,
		api.MakeRoute(http.MethodGet, "/:id/pty/:pid"):    h.pty,
	}
}

func (h *Handler) OnConnect(*agent.Agent) {
}

func (h *Handler) OnClose(string) {
}

func (h *Handler) OnMessage(*agent.Agent, *anet.Msg) {
}

func (h *Handler) clean() {
	fetch := func() []*task {
		ret := make([]*task, 0, len(h.tasks))
		now := time.Now()
		h.RLock()
		defer h.RUnlock()
		for _, t := range h.tasks {
			if now.After(t.clean) {
				ret = append(ret, t)
			}
		}
		return ret
	}
	for {
		for _, task := range fetch() {
			task.close()
			h.Lock()
			delete(h.tasks, task.pid)
			h.Unlock()
		}
		time.Sleep(time.Minute)
	}
}
