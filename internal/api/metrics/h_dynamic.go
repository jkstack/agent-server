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

type dynamicInfo struct {
	Usage       *usage       `json:"usage,omitempty"`       // usage数据
	Process     []process    `json:"process,omitempty"`     // 进程列表
	Connections []connection `json:"connections,omitempty"` // 连接列表
}

// static 获取节点的所有动态数据
// @ID /api/metrics/dynamic
// @description 1. 当指定top参数时将会获取CPU占用率最高的n个进程数据
// @description 2. 当指定kinds参数时获取的连接类型将会覆盖该agent节点配置文件中的类型设置
// @description 3. 当未指定kinds参数且该agent未配置task.conns.allow类型时默认返回所有类型的连接
// @Summary 获取节点的所有动态数据
// @Tags metrics
// @Produce json
// @Param   id    path  string   true  "节点ID"
// @Param   top   query integer  false "获取进程列表时的数量限制"
// @Param   kinds query []string false "获取连接类型" Enums(tcp,tcp4,tcp6,udp,udp4,udp6,unix)
// @Success 200   {object}       api.Success{payload=dynamicInfo}
// @Router /metrics/{id}/dynamic [get]
func (h *Handler) dynamic(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	topStr := g.DefaultQuery("top", "0")
	top, _ := strconv.ParseInt(topStr, 10, 64)
	kinds := g.QueryArray("kinds")
	if len(kinds) == 1 && kinds[0] == "" {
		kinds = []string{}
	}

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
		anet.HMReqUsage, anet.HMReqProcess,
		anet.HMReqConnections, anet.HMReqSensorsTemperatures,
	}, int(top), kinds)
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

	var ret dynamicInfo
	ret.Usage = transDynamicUsage(msg.HMDynamicRep.Usage)
	ret.Process = transDynamicProcess(msg.HMDynamicRep.Process)
	ret.Connections = transDynamicConnections(msg.HMDynamicRep.Connections)
	g.OK(ret)
}
