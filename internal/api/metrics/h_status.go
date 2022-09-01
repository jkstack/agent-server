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

type job struct {
	Name     string `json:"name" example:"static" enums:"static,usage,process,conns"` // 任务名称
	Interval int    `json:"interval" example:"5"`                                     // 间隔时间
}

type status struct {
	Jobs       []job    `json:"jobs"`                                                                   // 正在运行的任务列表
	AllowConns []string `json:"allow_conns" example:"tcp,udp" enums:"tcp,tcp4,tcp6,udp,udp4,udp6,unix"` // 采集的连接类型
}

// status 获取节点自动采集状态
// @ID /api/metrics/status
// @Summary 获取节点自动采集状态
// @Tags metrics
// @Produce json
// @Param   id   path string  true "节点ID"
// @Success 200  {object}     api.Success{payload=status}
// @Router /metrics/{id}/status [get]
func (h *Handler) status(gin *gin.Context) {
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

	taskID, err := cli.SendHMQueryStatus()
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
	case msg.Type != anet.TypeHMCollectStatus:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	var st status
	st.AllowConns = msg.HMCollectStatus.ConnsAllow
	for _, j := range msg.HMCollectStatus.Jobs {
		st.Jobs = append(st.Jobs, job{
			Name:     j.Name,
			Interval: int(j.Interval.Seconds()),
		})
	}
	g.OK(st)
}
