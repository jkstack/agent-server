package metrics

import (
	"fmt"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	runtime "github.com/jkstack/jkframe/utils"
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
	Mount      string  `json:"name" example:"/"`              // linux为挂载路径如/run，windows为盘符如C:
	Used       uint64  `json:"used" example:"16920992"`       // 已使用字节数
	Free       uint64  `json:"free" example:"232815064"`      // 剩余字节数
	Usage      float64 `json:"usage" example:"6.27"`          // 磁盘使用率
	InodeUsed  uint64  `json:"inode_used" example:"778282"`   // inode已使用数量
	InodeFree  uint64  `json:"inode_free" example:"15998934"` // inode剩余数量
	InodeUsage float64 `json:"inode_usage" example:"5.64"`    // inode使用率
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
	g := api.GetG(gin)

	id := g.Param("id")

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.NotFound("agent")
		return
	}
	if cli.Type() != agent.TypeMetrics {
		g.InvalidType(agent.TypeMetrics, cli.Type())
	}

	taskID, err := cli.SendHMDynamicReq([]anet.HMDynamicReqType{
		anet.HMReqUsage, anet.HMReqProcess, anet.HMReqConnections,
	})
	runtime.Assert(err)
	defer cli.ChanClose(id)

	var msg *anet.Msg
	select {
	case msg = <-cli.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		g.Timeout()
	}

	switch {
	case msg.Type == anet.TypeError:
		g.ERR(http.StatusServiceUnavailable, msg.ErrorMsg)
		return
	case msg.Type != anet.TypeHMDynamicRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	g.OK(transDynamicUsage(msg.HMDynamicRep.Usage))
}

func transDynamicUsage(input *anet.HMDynamicUsage) *usage {
	var ret usage
	ret.Cpu.Usage = input.Cpu.Usage
	ret.Memory.Used = input.Memory.Used
	ret.Memory.Free = input.Memory.Free
	ret.Memory.Available = input.Memory.Available
	ret.Memory.Total = input.Memory.Total
	ret.Memory.Usage = input.Memory.Usage
	ret.Memory.SwapUsed = input.Memory.SwapUsed
	ret.Memory.SwapFree = input.Memory.SwapFree
	ret.Memory.SwapTotal = input.Memory.SwapTotal
	for _, part := range input.Partitions {
		ret.Partitions = append(ret.Partitions, partitionUsage{
			Mount:      part.Name,
			Used:       part.Used,
			Free:       part.Free,
			Usage:      part.Usage,
			InodeUsed:  part.InodeUsed,
			InodeFree:  part.InodeFree,
			InodeUsage: part.InodeUsage,
		})
	}
	for _, intf := range input.Interface {
		ret.Interface = append(ret.Interface, interfaceUsage{
			Name:        intf.Name,
			BytesSent:   intf.BytesSent,
			BytesRecv:   intf.BytesRecv,
			PacketsSent: intf.PacketsSent,
			PacketsRecv: intf.PacketsRecv,
		})
	}
	return &ret
}