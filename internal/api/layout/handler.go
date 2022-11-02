package layout

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
	tasks map[string]*task
}

func New() *Handler {
	h := &Handler{
		tasks: make(map[string]*task),
	}
	go h.clear()
	return h
}

func (h *Handler) Module() string {
	return "layout"
}

func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.cfg = cfg
}

func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodPost, "/run"):       h.run,
		api.MakeRoute(http.MethodGet, "/status/:id"): h.status,
	}
}

func (h *Handler) OnConnect(*agent.Agent) {
}

func (h *Handler) OnClose(string) {
}

func (h *Handler) OnMessage(*agent.Agent, *anet.Msg) {
}

func (h *Handler) clear() {
	tick := func() {
		list := make([]*task, 0, len(h.tasks))
		h.RLock()
		for _, t := range h.tasks {
			if t.Done && time.Since(t.End).Hours() >= 1 {
				list = append(list, t)
			}
		}
		h.RUnlock()

		for _, t := range list {
			h.Lock()
			delete(h.tasks, t.ID)
			h.Unlock()
		}
	}
	for {
		time.Sleep(time.Minute)
		tick()
	}
}
