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

type connection struct {
	Fd     uint32 `json:"fd" example:"8" validate:"required"`                                                                                                                                                                              // 句柄号
	Pid    int32  `json:"pid" example:"28093" validate:"required"`                                                                                                                                                                         // 所属进程ID
	Type   string `json:"type" example:"tcp4" enums:"tcp4,tcp6,udp4,udp6,unix,file" validate:"required"`                                                                                                                                   // 连接类型
	Local  string `json:"local,omitempty" example:"127.0.0.1:13081"`                                                                                                                                                                       // 本地地址
	Remote string `json:"remote,omitempty" example:"127.0.0.1:37470"`                                                                                                                                                                      // 远程地址
	Status string `json:"status" example:"ESTABLISHED" enums:",ESTABLISHED,SYN_SENT,SYN_RECEIVED,SYN_RECV,FIN_WAIT_1,FIN_WAIT_2,FIN_WAIT1,FIN_WAIT2,TIME_WAIT,CLOSE,CLOSED,CLOSE_WAIT,LAST_ACK,LISTEN,CLOSING,DELETE" validate:"required"` // 连接状态
}

// static 获取节点的连接列表数据
// @ID /api/metrics/dynamic/connections
// @description 1. 当指定kinds参数时获取的连接类型将会覆盖该agent节点配置文件中的类型设置
// @description 2. 当未指定kinds参数且该agent未配置task.conns.allow类型时默认返回所有类型的连接
// @Summary 获取节点的连接列表数据
// @Tags metrics
// @Accept  json
// @Produce json
// @Param   id    path string    true  "节点ID"
// @Param   kinds query []string false "获取连接类型" Enums(tcp,tcp4,tcp6,udp,udp4,udp6,unix)
// @Success 200   {object}       api.Success{payload=[]connection}
// @Router /metrics/{id}/dynamic/connections [get]
func (h *Handler) dynamicConnections(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	kinds := g.QueryArray("kinds")
	if len(kinds) == 1 && kinds[0] == "" {
		kinds = []string{}
	}

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.Notfound("agent")
		return
	}
	if cli.Type() != agent.TypeMetrics {
		g.InvalidType(agent.TypeMetrics, cli.Type())
	}

	taskID, err := cli.SendHMDynamicReq([]anet.HMDynamicReqType{
		anet.HMReqConnections,
	}, 0, kinds)
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

	g.OK(transDynamicConnections(msg.HMDynamicRep.Connections))
}

func transDynamicConnections(conns []anet.HMDynamicConnection) []connection {
	var ret []connection
	for _, conn := range conns {
		ret = append(ret, connection{
			Fd:     conn.Fd,
			Pid:    conn.Pid,
			Type:   conn.Type,
			Local:  conn.Local,
			Remote: conn.Remote,
			Status: conn.Status,
		})
	}
	return ret
}
