package rpa

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

type selectorValidateArgs struct {
	Content string `json:"content"` // 验证内容
}

// selectorValidate 元素选择器结果验证
//
//	@ID			/api/rpa/selector_validate
//	@Summary	元素选择器结果验证
//	@Tags		rpa
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"节点ID"
//	@Param		args	body		selectorValidateArgs	true	"需启动的任务列表"
//	@Success	200		{object}	api.Success
//	@Router		/rpa/{id}/selector_validate [post]
func (h *Handler) selectorValidate(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	var args selectorValidateArgs
	if err := g.ShouldBindJSON(&args); err != nil {
		g.BadParam(err.Error())
		return
	}

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.Notfound("agent")
		return
	}
	if cli.Type() != agent.TypeRPA {
		g.InvalidType(agent.TypeRPA, cli.Type())
	}

	taskID, err := cli.SendRpaSelectorValidate(args.Content)
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
	case msg.Type != anet.TypeRPASelectorValidateRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	if msg.RPASelectorValidateRep.OK {
		g.OK(nil)
		return
	}

	g.ERR(http.StatusServiceUnavailable, msg.RPASelectorValidateRep.Msg)
}
