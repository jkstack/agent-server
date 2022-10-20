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
// @Summary 调用example类型的agent
// @Tags foo
// @Accept  json
// @Produce json
// @Param   id    path string  true "节点ID"
// @Success 200   {object}     api.Success
// @Router /foo/{id} [get]
func (h *Handler) foo(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.Notfound("agent")
	}
	if cli.Type() != agent.TypeExample {
		g.InvalidType(agent.TypeExample, cli.Type())
	}

	taskID, err := cli.SendFoo()
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
	case msg.Type != anet.TypeBar:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	g.OK(nil)
}
