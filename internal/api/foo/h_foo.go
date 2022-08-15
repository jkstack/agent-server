package foo

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

// info 调用example类型的agent
// @ID /api/foo
// @Description 调用example类型的agent
// @Tags foo
// @Produce json
// @Param   id    path string  true "节点ID"
// @Success 200   {object}     api.Success
// @Router /foo/{id} [get]
func (h *Handler) foo(g *gin.Context) {
	id := g.Param("id")

	agents := api.GetAgents(g)

	cli := agents.Get(id)
	if cli == nil {
		api.PanicNotFound("agent")
	}
	if cli.Type() != agent.TypeExample {
		api.PanicInvalidType(agent.TypeExample, cli.Type())
	}

	taskID, err := cli.SendFoo()
	runtime.Assert(err)
	defer cli.ChanClose(id)

	var msg *anet.Msg
	select {
	case msg = <-cli.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		api.PanicTimeout()
	}

	switch {
	case msg.Type == anet.TypeError:
		api.ERR(g, http.StatusServiceUnavailable, msg.ErrorMsg)
		return
	case msg.Type != anet.TypeBar:
		api.ERR(g, http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	api.OK(g, nil)
}
