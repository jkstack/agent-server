package metrics

import (
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/conf"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/stat"
	"github.com/prometheus/client_golang/prometheus"
)

var allJobs = []string{"static", "usage", "process", "conns"}

type Handler struct {
	stJobs      *prometheus.GaugeVec
	stWarning   *prometheus.GaugeVec
	stBytesSent *prometheus.GaugeVec
	stCounts    *prometheus.GaugeVec
	cli         sarama.AsyncProducer
	topic       string
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Module() string {
	return "metrics"
}

func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.stJobs = mgr.RawVec("metrics_jobs", []string{"id", "name"})
	h.stWarning = mgr.RawVec("metrics_warning", []string{"id"})
	h.stBytesSent = mgr.RawVec("metrics_bytes_sent", []string{"id", "name"})
	h.stCounts = mgr.RawVec("metrics_report_count", []string{"id", "name"})
	h.cli = cfg.MetricsCli
	h.topic = cfg.Metrics.Topic
}

func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodGet, "/:id/static"):              h.static,
		api.MakeRoute(http.MethodGet, "/:id/dynamic"):             h.dynamic,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/usage"):       h.dynamicUsage,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/process"):     h.dynamicProcess,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/connections"): h.dynamicConnections,
		api.MakeRoute(http.MethodGet, "/:id/status"):              h.getStatus,
		api.MakeRoute(http.MethodPut, "/:id/status"):              h.setStatus,
		api.MakeRoute(http.MethodPut, "/status"):                  h.batchSetStatus,
	}
}

func (h *Handler) OnConnect(*agent.Agent) {
}

func (h *Handler) OnClose(string) {
}

func (h *Handler) OnMessage(agent *agent.Agent, msg *anet.Msg) {
	switch msg.Type {
	case anet.TypeHMStaticRep:
		logging.Debug("agent [%s] report static info", agent.ID())
		h.saveStaticData(agent.ID(), msg.HMStatic)
	case anet.TypeHMDynamicRep:
		logging.Debug("agent [%s] report dynamic info", agent.ID())
		h.saveDynamicData(agent.ID(), msg.HMDynamicRep)
	case anet.TypeHMReportAgentStatus:
		h.handleReport(agent.ID(), msg.HMAgentStatus)
	}
}
