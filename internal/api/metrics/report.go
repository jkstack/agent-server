package metrics

import (
	"github.com/Shopify/sarama"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/prometheus/client_golang/prometheus"
)

func (h *Handler) handleReport(id string, data *anet.HMAgentStatus) {
	running := make(map[string]int)
	for _, job := range data.Jobs {
		running[job] = 1
	}
	for _, job := range allJobs {
		labels := prometheus.Labels{
			"id":   id,
			"name": job,
		}
		h.stJobs.Delete(labels)
		h.stJobs.With(labels).Set(float64(running[job]))
		h.stBytesSent.Delete(labels)
		h.stBytesSent.With(labels).Set(float64(data.ReportBytes[job]))
		h.stCounts.Delete(labels)
		h.stCounts.With(labels).Set(float64(data.ReportCount[job]))
	}
	h.stWarning.With(prometheus.Labels{"id": id}).Set(float64(data.Warnings))
}

func HandleReportError(cli sarama.AsyncProducer) {
	for err := range cli.Errors() {
		key, err := err.Msg.Key.Encode()
		if err != nil {
			logging.Error("encode message key: %v", err)
			continue
		}
		logging.Error("metrics send for [%s]: %v", string(key), err)
	}
}
