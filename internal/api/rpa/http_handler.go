package rpa

import (
	"net/http"
	"path/filepath"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/conf"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/cache/l2cache"
	"github.com/jkstack/jkframe/stat"
	"github.com/jkstack/jkframe/utils"
)

type cache struct {
	t     time.Time
	cache *l2cache.Cache
}

// Handler api handler
type Handler struct {
	sync.RWMutex
	cacheDir string
	files    map[string]*cache
}

// New create api handler
func New() *Handler {
	h := &Handler{
		files: make(map[string]*cache),
	}
	go h.clear()
	return h
}

// Module get module name
func (h *Handler) Module() string {
	return "rpa"
}

// Init initialize module
func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.cacheDir = cfg.CacheDir
}

// HandleFuncs get funcs
func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodPost, "/:id/in_selector"): h.inSelector,
		api.MakeRoute(http.MethodGet, "/files/:id"):        h.file,
	}
}

// OnConnect agent connect callback
func (h *Handler) OnConnect(*agent.Agent) {
}

// OnClose agent connection closed callback
func (h *Handler) OnClose(string) {
}

// OnMessage received agent message callback
func (h *Handler) OnMessage(agent *agent.Agent, msg *anet.Msg) {
}

func (h *Handler) newCache(taskID string, data []byte) {
	c, err := l2cache.New(102400, filepath.Join(h.cacheDir, "rpa"))
	utils.Assert(err)
	_, err = c.Write(data)
	utils.Assert(err)
	ch := &cache{
		t:     time.Now(),
		cache: c,
	}
	h.Lock()
	h.files[taskID] = ch
	h.Unlock()
}

func (h *Handler) clear() {
	fetch := func() []string {
		ret := make([]string, 0, len(h.files))
		h.RLock()
		defer h.RUnlock()
		for id, file := range h.files {
			if time.Since(file.t) > time.Hour {
				ret = append(ret, id)
			}
		}
		return ret
	}
	drop := func(id string) {
		h.Lock()
		defer h.Unlock()
		c, ok := h.files[id]
		if !ok {
			return
		}
		c.cache.Close()
		delete(h.files, id)
	}
	for {
		for _, id := range fetch() {
			drop(id)
		}
		time.Sleep(time.Minute)
	}
}
