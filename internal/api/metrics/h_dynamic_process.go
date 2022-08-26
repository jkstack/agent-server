package metrics

import (
	"fmt"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	runtime "github.com/jkstack/jkframe/utils"
)

type process struct {
	ID            int32    `json:"id" example:"5167"`                       // 进程ID
	ParentID      int32    `json:"parent_id" example:"1"`                   // 父进程ID
	User          string   `json:"user" example:"root"`                     // 用户
	CpuUsage      float64  `json:"cpu_usage" example:"3.2"`                 // CPU使用率
	RssMemory     uint64   `json:"rss" example:"37654"`                     // 物理内存数
	VirtualMemory uint64   `json:"vms" example:"26754"`                     // 虚拟内存数
	SwapMemory    uint64   `json:"swap" example:"0"`                        // swap内存数
	MemoryUsage   float64  `json:"memory_usage" example:"6.4"`              // 内存使用率
	Cmd           []string `json:"cmd,omitempty" example:"/usr/bin/zsh,-i"` // 命令行
	Listen        []uint32 `json:"listen,omitempty" example:"8080,9090"`    // 监听端口
	Connections   int      `json:"conns" example:"16"`                      // 连接数
}

// static 获取节点的进程列表数据
// @ID /api/metrics/dynamic/process
// @Summary 获取节点的所有进程列表数据
// @Tags metrics
// @Produce json
// @Param   id   path  string  true  "节点ID"
// @Param   top  query integer false "数量限制"
// @Success 200  {object}     api.Success{payload=[]process}
// @Router /metrics/{id}/dynamic/process [get]
func (h *Handler) dynamicProcess(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	topStr := g.DefaultQuery("top", "0")
	top, _ := strconv.ParseInt(topStr, 10, 64)

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
		anet.HMReqProcess,
	}, int(top))
	runtime.Assert(err)
	defer cli.ChanClose(id)

	var msg *anet.Msg
	select {
	case msg = <-cli.ChanRead(taskID):
	case <-time.After(time.Minute):
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

	g.OK(transDynamicProcess(msg.HMDynamicRep.Process))
}

func transDynamicProcess(input []anet.HMDynamicProcess) []process {
	var ret []process
	for _, p := range input {
		ret = append(ret, process{
			ID:            p.ID,
			ParentID:      p.ParentID,
			User:          p.User,
			CpuUsage:      p.CpuUsage.Float(),
			RssMemory:     p.RssMemory,
			VirtualMemory: p.VirtualMemory,
			SwapMemory:    p.SwapMemory,
			MemoryUsage:   p.MemoryUsage.Float(),
			Cmd:           p.Cmd,
			Listen:        p.Listen,
			Connections:   p.Connections,
		})
	}
	return ret
}
