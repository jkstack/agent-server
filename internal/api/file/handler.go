package file

import (
	"net/http"
	"os"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/conf"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/stat"
)

type uploadInfo struct {
	token   string
	dir     string
	rm      bool
	timeout time.Time
}

// Handler api handler
type Handler struct {
	sync.RWMutex
	cfg         *conf.Configure
	uploadCache map[string]*uploadInfo
}

// New create api handler
func New() *Handler {
	h := &Handler{
		uploadCache: make(map[string]*uploadInfo),
	}
	go h.clean()
	return h
}

// Module get module name
func (h *Handler) Module() string {
	return "file"
}

// Init initialize module
func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.cfg = cfg
}

// HandleFuncs get funcs
func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodGet, "/:id/ls"):       h.ls,
		api.MakeRoute(http.MethodGet, "/:id/download"): h.download,
		api.MakeRoute(http.MethodPost, "/:id/upload"):  h.upload,
		api.MakeRoute(http.MethodGet, "/upload/:id"):   h.uploadHandle,
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

func (h *Handler) clean() {
}

func (h *Handler) logUploadCache(taskID, dir, token string,
	deadline time.Time, rm bool) {
	h.Lock()
	h.uploadCache[taskID] = &uploadInfo{
		token:   token,
		dir:     dir,
		rm:      rm,
		timeout: deadline,
	}
	h.Unlock()
	logging.Info("log cache: %s", taskID)
}

func (h *Handler) removeUploadCache(id string) bool {
	h.Lock()
	cache := h.uploadCache[id]
	delete(h.uploadCache, id)
	h.Unlock()
	if cache == nil {
		return false
	}
	logging.Info("removed cache: %s", id)
	if cache.rm {
		os.Remove(cache.dir)
	}
	return true
}
