package agents

import (
	"time"

	"github.com/jkstack/anet"
	"github.com/prometheus/client_golang/prometheus"
)

func setValue(vec *prometheus.GaugeVec, id, t, tag string, n float64) {
	vec.With(prometheus.Labels{
		"id":         id,
		"agent_type": t,
		"tag":        tag,
	}).Set(n)
}

func (h *Handler) handleReport(id, t string, info *anet.AgentInfo) {
	h.mVersion.RLock()
	oldLabels := prometheus.Labels{
		"id":         id,
		"version":    h.oldVersion[id],
		"go_version": h.oldGoVersion[id],
	}
	h.mVersion.RUnlock()

	h.stAgentInfo.Delete(oldLabels)

	h.mVersion.Lock()
	h.oldVersion[id] = info.Version
	h.oldGoVersion[id] = info.GoVersion
	h.mVersion.Unlock()

	h.stAgentVersion.With(prometheus.Labels{
		"id":         id,
		"version":    info.Version,
		"go_version": info.GoVersion,
	}).Set(1)
	setValue(h.stAgentInfo, id, t, "cpu_usage", float64(info.CpuUsage))
	setValue(h.stAgentInfo, id, t, "memory_usage", float64(info.MemoryUsage))
	setValue(h.stAgentInfo, id, t, "threads", float64(info.Threads))
	setValue(h.stAgentInfo, id, t, "routines", float64(info.Routines))
	setValue(h.stAgentInfo, id, t, "startup", float64(info.Startup))
	setValue(h.stAgentInfo, id, t, "heap_in_use", float64(info.HeapInuse))
	setValue(h.stAgentInfo, id, t, "gc_0", info.GC["0"])
	setValue(h.stAgentInfo, id, t, "gc_0.25", info.GC["25"])
	setValue(h.stAgentInfo, id, t, "gc_0.5", info.GC["50"])
	setValue(h.stAgentInfo, id, t, "gc_0.75", info.GC["75"])
	setValue(h.stAgentInfo, id, t, "gc_1", info.GC["100"])
	setValue(h.stAgentInfo, id, t, "in_packets", float64(info.InPackets))
	setValue(h.stAgentInfo, id, t, "in_bytes", float64(info.InBytes))
	setValue(h.stAgentInfo, id, t, "out_packets", float64(info.OutPackets))
	setValue(h.stAgentInfo, id, t, "out_bytes", float64(info.OutBytes))
	setValue(h.stAgentInfo, id, t, "reconnect_count", float64(info.ReconnectCount))
	setValue(h.stAgentInfo, id, t, "read_chan_size", float64(info.ReadChanSize))
	setValue(h.stAgentInfo, id, t, "write_chan_size", float64(info.WriteChanSize))
	setValue(h.stAgentInfo, id, t, "report_time", float64(time.Now().Unix()))
}
