package exec

import (
	"server/internal/agent"
	"server/internal/api"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/utils"
)

// kill 结束进程
// @ID /api/exec/kill
// @Summary 结束进程
// @Tags exec
// @Produce json
// @Param   id   path string  true  "节点ID"
// @Param   pid  path int     true  "进程号"
// @Success 200  {object}     api.Success
// @Router /exec/{id}/kill/{pid} [delete]
func (h *Handler) kill(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	pid, err := strconv.ParseInt(g.Param("pid"), 10, 64)
	if err != nil {
		api.BadParamErr("pid")
		return
	}

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.NotFound("agent")
		return
	}
	if cli.Type() != agent.TypeExec {
		g.InvalidType(agent.TypeExec, cli.Type())
	}

	err = cli.SendExecKill(int(pid))
	utils.Assert(err)

	g.OK(nil)
}
