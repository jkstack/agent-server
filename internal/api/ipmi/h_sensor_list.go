package ipmi

import (
	"fmt"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/utils"
)

type sensorUnitValue struct {
	Percent bool   `json:"percent" validate:"required" default:"false"`  // 是否是百分比
	Base    string `json:"base" validate:"required" default:"degrees C"` // 基本单位
	Op      string `json:"op,omitempty" default:"/"`                     // 乘或除
	Mod     string `json:"mod,omitempty" default:"minute"`               // 第二项单位
}

type sensorCritical struct {
	NonCritical    *float64 `json:"non_critical,omitempty" default:"60"`     // 恢复数值
	Critical       *float64 `json:"critical,omitempty" default:"70"`         // 告警数值
	NonRecoverable *float64 `json:"non_recoverable,omitempty" default:"100"` // 严重告警数值
}

type sensorValue struct {
	Unit    sensorUnitValue `json:"unit" validate:"required"`                 // 当前传感器的单位信息
	Current float64         `json:"current" validate:"required" default:"23"` // 当前数值
	Lower   *sensorCritical `json:"lower,omitempty"`                          // 最低告警数值
	Upper   *sensorCritical `json:"upper,omitempty"`                          // 最高告警数值
}

type sensorInfo struct {
	ID       uint16       `json:"id" validate:"required" default:"0"`              // 传感器序号
	SensorID uint8        `json:"sensor_id" validate:"required" default:"44"`      // 传感器ID
	EntityID string       `json:"entity_id" validate:"required" default:"12.1"`    // 实体ID
	Name     string       `json:"name" validate:"required" default:"Ambient Temp"` // 传感器名称
	Type     string       `json:"type" validate:"required" default:"Temperature"`  // 传感器类型
	Discrete bool         `json:"discrete" validate:"required" default:"false"`    // 是否是离散传感器
	Values   *sensorValue `json:"values,omitempty"`                                // 传感器数值
}

// sensorList 获取服务器的传感器列表
//
//	@ID				/api/ipmi/sensors
//	@Description	常量定义：https://docs.jkservice.org/dp/others/ipmi/defs/
//	@Summary		获取服务器的传感器列表
//	@Tags			ipmi
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string								true	"节点ID"
//	@Param			addr	query		string								true	"IPMI地址"
//	@Param			user	query		string								true	"IPMI用户名"
//	@Param			pass	query		string								true	"IPMI密码"
//	@Param			mode	query		string								false	"IPMI连接模式"	enums(lan,lanplus,auto)	default(auto)
//	@Success		200		{object}	api.Success{payload=[]sensorInfo}	"传感器列表"
//	@Router			/ipmi/{id}/sensors [get]
func (h *Handler) sensorList(gin *gin.Context) {
	g := api.GetG(gin)

	var args commonArgs
	if err := g.ShouldBindQuery(&args); err != nil {
		g.BadParam(err.Error())
		return
	}

	switch args.Mode {
	case "lan", "lanplus":
	default:
		args.Mode = "auto"
	}

	id := g.Param("id")

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.HTTPError(http.StatusNotFound, "agent not found")
		return
	}
	if cli.Type() != agent.TypeIPMI {
		g.InvalidType(agent.TypeIPMI, cli.Type())
		return
	}

	taskID, err := cli.SendIPMISensorList(args.toRequest())
	utils.Assert(err)
	defer cli.ChanClose(taskID)

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
	case msg.Type != anet.TypeIPMISensorListRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	if !msg.IPMISensorList.OK {
		g.ERR(http.StatusServiceUnavailable, msg.IPMISensorList.Msg)
		return
	}

	g.OK(h.buildSensorList(msg.IPMISensorList))
}

func (h *Handler) buildSensorList(list *anet.IPMISensorList) []sensorInfo {
	sensors := make([]sensorInfo, 0, len(list.List))
	build := func(n *utils.Float64P2) *float64 {
		if n != nil {
			dn := float64(*n)
			return &dn
		}
		return nil
	}
	for _, s := range list.List {
		var values *sensorValue
		if s.Values != nil {
			values = new(sensorValue)
			values.Unit.Percent = s.Values.Unit.Percent
			values.Unit.Base = s.Values.Unit.Base
			values.Unit.Op = s.Values.Unit.Op
			values.Unit.Mod = s.Values.Unit.Mod
			values.Current = float64(s.Values.Current)
			if s.Values.Lower != nil {
				values.Lower = new(sensorCritical)
				values.Lower.NonCritical = build(s.Values.Lower.NonCritical)
				values.Lower.Critical = build(s.Values.Lower.Critical)
				values.Lower.NonRecoverable = build(s.Values.Lower.NonRecoverable)
			}
			if s.Values.Upper != nil {
				values.Upper = new(sensorCritical)
				values.Upper.NonCritical = build(s.Values.Upper.NonCritical)
				values.Upper.Critical = build(s.Values.Upper.Critical)
				values.Upper.NonRecoverable = build(s.Values.Upper.NonRecoverable)
			}
		}
		sensors = append(sensors, sensorInfo{
			ID:       s.ID,
			SensorID: s.SensorID,
			EntityID: s.EntityID,
			Name:     s.Name,
			Type:     s.Type,
			Discrete: s.Discrete,
			Values:   values,
		})
	}
	return sensors
}
