package ipmi

import (
	"fmt"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/utils"
)

type commonArgs struct {
	Addr string `form:"addr"`
	User string `form:"user"`
	Pass string `form:"pass"`
	Mode string `form:"mode"`
}

func (args commonArgs) toRequest() *anet.IPMICommonRequest {
	return &anet.IPMICommonRequest{
		Addr:      args.Addr,
		Username:  args.User,
		Password:  args.Pass,
		Interface: args.Mode,
	}
}

type deviceInfo struct {
	OEM             string `json:"oem" validate:"required" default:"IBM"`              // 生产厂商
	FirmwareVersion string `json:"firmware_version" validate:"required" default:"0.0"` // 固件版本
	IPMIVersion     string `json:"ipmi_version" validate:"required" default:"2.0"`     // IPMI版本
}

// deviceInfo 获取服务器的设备信息
//
//	@ID			/api/ipmi/device
//	@Summary	获取服务器的设备信息
//	@Tags		ipmi
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string							true	"节点ID"
//	@Param		addr	query		string							true	"IPMI地址"
//	@Param		user	query		string							true	"IPMI用户名"
//	@Param		pass	query		string							true	"IPMI密码"
//	@Param		mode	query		string							false	"IPMI连接模式"	enums(lan,lanplus,auto) default(auto)
//	@Success	200		{object}	api.Success{payload=deviceInfo}	"服务器IPMI信息"
//	@Router		/ipmi/{id}/device [get]
func (h *Handler) deviceInfo(gin *gin.Context) {
	g := api.GetG(gin)

	var args commonArgs
	if err := g.ShouldBindQuery(&args); err != nil {
		g.BadParam(err.Error())
		return
	}

	tmp := strings.SplitN(args.Addr, ":", 2)
	if len(tmp) == 1 {
		args.Addr = fmt.Sprintf("%s:623", tmp[0])
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

	taskID, err := cli.SendIPMIDeviceInfo(args.toRequest())
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
	case msg.Type != anet.TypeIPMIDeviceInfoRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	if !msg.IPMIDeviceInfo.OK {
		g.ERR(http.StatusServiceUnavailable, msg.IPMIDeviceInfo.Msg)
		return
	}

	info := msg.IPMIDeviceInfo
	g.OK(deviceInfo{
		OEM:             info.OEM,
		FirmwareVersion: info.FirmwareVersion,
		IPMIVersion:     info.IPMIVersion,
	})
}
