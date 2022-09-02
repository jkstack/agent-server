package agents

import (
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/conf"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/stat"
	"github.com/prometheus/client_golang/prometheus"
)

type Handler struct {
	stAgentVersion *prometheus.GaugeVec
	stAgentInfo    *prometheus.GaugeVec
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Module() string {
	return "agents"
}

func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.stAgentVersion = mgr.RawVec("agent_version", []string{"id", "version", "go_version"})
	h.stAgentInfo = mgr.RawVec("agent_info", []string{"id", "agent_type", "tag"})
}

func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodGet, ""):     h.list,
		api.MakeRoute(http.MethodGet, "/:id"): h.info,
	}
}

func (h *Handler) OnConnect(*agent.Agent) {
}

func (h *Handler) OnClose(string) {
}

func (h *Handler) OnMessage(cli *agent.Agent, msg *anet.Msg) {
	if msg.Type != anet.TypeAgentInfo {
		return
	}
	h.handleReport(cli.ID(), cli.Type(), msg.AgentInfo)
}
