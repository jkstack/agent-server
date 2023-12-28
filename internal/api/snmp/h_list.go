package snmp

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

// list 获取snmp清单
//
//	@ID			/api/snmp/list
//	@Summary	获取snmp清单
//	@Tags		snmp
//	@Accept		json
//	@Produce	json
//	@Param		id			path		string									true	"节点ID"
//	@Param		host		query		string									true	"SNMP服务地址"
//	@Param		community	query		string									true	"SNMP community"
//	@Param		oid			query		string									false	"SNMP oid"
//	@Success	200			{object}	api.Success{payload=[]anet.SNMPItem}	"snmp列表"
//	@Router		/snmp/{id}/list [get]
func (h *Handler) list(gin *gin.Context) {
	g := api.GetG(gin)

	var args struct {
		Host      string `form:"host"`
		Community string `form:"community"`
		OID       string `form:"oid"`
	}
	if err := g.ShouldBindQuery(&args); err != nil {
		g.BadParam(err.Error())
		return
	}

	id := g.Param("id")

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.HTTPError(http.StatusNotFound, "agent not found")
		return
	}
	if cli.Type() != agent.TypeSNMP {
		g.InvalidType(agent.TypeSNMP, cli.Type())
		return
	}

	taskID, err := cli.SendSNMPList(&anet.SNMPReq{
		Host:      args.Host,
		Community: args.Community,
		OID:       args.OID,
	})
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
	case msg.Type != anet.TypeSNMPListRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	if !msg.SNMPRep.OK {
		g.ERR(http.StatusServiceUnavailable, msg.SNMPRep.Msg)
		return
	}

	g.OK(msg.SNMPRep.Items)
}
