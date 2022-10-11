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

type temp struct {
	Name string  `json:"name" example:"coretemp_core_0" validate:"required"` // 名称
	Temp float64 `json:"temp" example:"38" validate:"required"`              // 温度
}

// static 获取节点的传感器温度数据
// @ID /api/metrics/dynamic/temps
// @Summary 获取节点的传感器温度数据
// @Tags metrics
// @Produce json
// @Param   id   path  string  true  "节点ID"
// @Success 200  {object}     api.Success{payload=[]temp}
// @Router /metrics/{id}/dynamic/temps [get]
func (h *Handler) dynamicTemps(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")

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
		anet.HMReqSensorsTemperatures,
	}, 0, nil)
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

	g.OK(transSensorsTemperatures(msg.HMDynamicRep.SensorsTemperatures))
}

func transSensorsTemperatures(temps []anet.HMSensorTemperature) []temp {
	ret := make([]temp, 0, len(temps))
	for _, t := range temps {
		ret = append(ret, temp{
			Name: t.Name,
			Temp: t.Temperature.Float(),
		})
	}
	return ret
}
