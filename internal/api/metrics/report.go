package metrics

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/prometheus/client_golang/prometheus"
)

func (h *Handler) handleReport(id string, data *anet.HMAgentStatus) {
	running := make(map[string]bool)
	for _, job := range data.Jobs {
		running[job] = true
	}
	var jobs jobs
	for i, job := range allJobs {
		var status jobStatus
		status.running = running[job]
		status.interval = 0 // TODO
		status.bytesSent = data.ReportBytes[job]
		status.count = data.ReportCount[job]
		jobs[i] = status
	}
	h.Lock()
	h.jobs[id] = jobs
	h.Unlock()
	h.stWarning.With(prometheus.Labels{"id": id}).Set(float64(data.Warnings))
}

func (h *Handler) updateJobs() {
	tick := func() {
		h.stJobs.Reset()
		h.RLock()
		defer h.RUnlock()
		for agentID, jobs := range h.jobs {
			for i, status := range jobs {
				running := 1.
				if !status.running {
					running = 0
				}
				h.stJobs.With(prometheus.Labels{
					"id":           agentID,
					"name":         allJobs[i],
					"interval":     fmt.Sprintf("%d", status.interval),
					"bytes_sent":   fmt.Sprintf("%d", status.bytesSent),
					"report_count": fmt.Sprintf("%d", status.count),
				}).Set(running)
			}
		}
	}
	for {
		time.Sleep(10 * time.Second)
		tick()
	}
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
