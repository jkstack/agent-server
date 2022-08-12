package foo

import (
	"fmt"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"time"

	"github.com/jkstack/anet"
	runtime "github.com/jkstack/jkframe/utils"
)

// info 调用example类型的agent
// @ID /api/foo
// @Description 调用example类型的agent
// @Tags foo
// @Produce json
// @Param   id    path string  true "节点ID"
// @Success 200   {object}     api.Success
// @Router /foo/{id} [get]
func (h *Handler) foo(ctx *api.GContext, agents *agent.Agents) {
	id := ctx.Param("id")

	cli := agents.Get(id)
	if cli == nil {
		ctx.NotFound("agent")
		return
	}
	if cli.Type() != agent.TypeExample {
		ctx.InvalidType(agent.TypeExample, cli.Type())
		return
	}

	taskID, err := cli.SendFoo()
	runtime.Assert(err)
	defer cli.ChanClose(id)

	var msg *anet.Msg
	select {
	case msg = <-cli.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		ctx.Timeout()
	}

	switch {
	case msg.Type == anet.TypeError:
		ctx.ERR(http.StatusServiceUnavailable, msg.ErrorMsg)
		return
	case msg.Type != anet.TypeBar:
		ctx.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	ctx.OK(nil)
}
