package metrics

import (
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/conf"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/stat"
	"github.com/prometheus/client_golang/prometheus"
)

var allJobs = []string{"static", "usage", "process", "conns", "temps"}

type jobStatus struct {
	running   bool
	interval  uint64
	bytesSent uint64
	count     uint64
}

type jobs [5]jobStatus

type Handler struct {
	sync.RWMutex
	stJobs    *prometheus.GaugeVec
	stWarning *prometheus.GaugeVec
	cli       sarama.AsyncProducer
	topic     string
	jobs      map[string]jobs
}

func New() *Handler {
	return &Handler{
		jobs: make(map[string]jobs),
	}
}

func (h *Handler) Module() string {
	return "metrics"
}

func (h *Handler) Init(cfg *conf.Configure, mgr *stat.Mgr) {
	h.stJobs = mgr.RawVec("metrics_jobs", []string{"id", "name",
		"interval", "bytes_sent", "report_count"})
	h.stWarning = mgr.RawVec("metrics_warning", []string{"id"})
	go h.updateJobs()
	h.cli = cfg.MetricsCli
	h.topic = cfg.Metrics.Kafka.Topic
	if h.cli != nil {
		go HandleReportError(h.cli)
	}
}

func (h *Handler) HandleFuncs() map[api.Route]func(*gin.Context) {
	return map[api.Route]func(*gin.Context){
		api.MakeRoute(http.MethodGet, "/:id/static"):              h.static,
		api.MakeRoute(http.MethodGet, "/:id/dynamic"):             h.dynamic,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/usage"):       h.dynamicUsage,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/process"):     h.dynamicProcess,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/connections"): h.dynamicConnections,
		api.MakeRoute(http.MethodGet, "/:id/dynamic/temps"):       h.dynamicTemps,
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
