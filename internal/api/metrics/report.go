package metrics

import (
	"github.com/jkstack/anet"
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
		h.stJobs.With(labels).Set(float64(running[job]))
		h.stBytesSent.With(labels).Set(float64(data.ReportBytes[job]))
		h.stCounts.With(labels).Set(float64(data.ReportCount[job]))
	}
	h.stWarning.With(prometheus.Labels{"id": id}).Set(float64(data.Warnings))
}
