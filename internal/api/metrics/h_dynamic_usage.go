package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
)

type usage struct {
	Cpu struct {
		Usage float64 `json:"usage" example:"2.3"` // CPU使用率(百分比)
	} `json:"cpu"`
	Memory struct {
		Used      uint64  `json:"used" example:"1595712"`       // 已使用字节数
		Free      uint64  `json:"free" example:"94896"`         // 剩余字节数
		Available uint64  `json:"available" example:"201332"`   // 可用字节数
		Total     uint64  `json:"total" example:"1956784"`      // 总字节数
		Usage     float64 `json:"usage" example:"1.2"`          // 内存使用率(百分比)
		SwapUsed  uint64  `json:"swap_used" example:"1146264"`  // swap已使用字节数
		SwapFree  uint64  `json:"swap_free" example:"7242344"`  // swap剩余字节数
		SwapTotal uint64  `json:"swap_total" example:"8388608"` // swap总字节数
	} `json:"memory"`
	Partitions []partitionUsage `json:"partitions,omitempty"` // 分区
	Interface  []interfaceUsage `json:"interface,omitempty"`  // 网卡
}

type partitionUsage struct {
	Name  string  `json:"name,omitempty"`  // linux为挂载路径如/run，windows为盘符如C:
	Used  uint64  `json:"used,omitempty"`  // 已使用字节数
	Free  uint64  `json:"free,omitempty"`  // 剩余字节数
	Usage float64 `json:"usage,omitempty"` // 磁盘使用率
}

type interfaceUsage struct {
	Name        string `json:"name" example:"eth0"`          // 网卡名称
	BytesSent   uint64 `json:"bytes_sent" example:"6162729"` // 发送字节数
	BytesRecv   uint64 `json:"bytes_recv" example:"24422"`   // 接收字节数
	PacketsSent uint64 `json:"packets_sent" example:"5699"`  // 发送数据包数量
	PacketsRecv uint64 `json:"packets_recv" example:"4399"`  // 接收数据包数量
}

// static 获取节点的usage动态数据
// @ID /api/metrics/dynamic/usage
// @Summary 获取节点的usage动态数据
// @Tags metrics
// @Produce json
// @Param   id   path string  true "节点ID"
// @Success 200  {object}     api.Success{payload=usage}
// @Router /metrics/{id}/dynamic/usage [get]
func (h *Handler) dynamicUsage(gin *gin.Context) {
}

func transDynamicUsage(usage *anet.HMDynamicUsage) *usage {
	return nil
}
