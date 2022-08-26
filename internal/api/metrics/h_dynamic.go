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

type dynamicInfo struct {
	Usage       *usage       `json:"usage,omitempty"`       // usage数据
	Process     []process    `json:"process,omitempty"`     // 进程列表
	Connections []connection `json:"connections,omitempty"` // 连接列表
}

// static 获取节点的所有动态数据
// @ID /api/metrics/dynamic
// @Summary 获取节点的所有动态数据
// @Tags metrics
// @Produce json
// @Param   id   path string  true "节点ID"
// @Success 200  {object}     api.Success{payload=dynamicInfo}
// @Router /metrics/{id}/dynamic [get]
func (h *Handler) dynamic(gin *gin.Context) {
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

	var ret dynamicInfo
	ret.Usage = transDynamicUsage(msg.HMDynamicRep.Usage)
	ret.Process = transDynamicProcess(msg.HMDynamicRep.Process)
	ret.Connections = transDynamicConnections(msg.HMDynamicRep.Connections)
	g.OK(ret)
}