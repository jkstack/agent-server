package exec

import (
	"server/internal/agent"
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type info struct {
	TaskID string `json:"task_id" example:"20221008-00001-7a390a60f759aab5" validate:"required"` // 任务ID
	Pid    int    `json:"pid" example:"11186" validate:"required"`                               // 进程号
	Begin  int64  `json:"begin" example:"1665219359" validate:"required"`                        // 启动时间
}

// ps 列出正在运行中的任务
// @ID /api/exec/ps
// @Summary 列出正在运行中的任务
// @Tags exec
// @Accept  json
// @Produce json
// @Param   id   path string  true  "节点ID"
// @Success 200  {object}     api.Success{payload=[]info}
// @Router /exec/{id}/ps [get]
func (h *Handler) ps(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.Notfound("agent")
		return
	}
	if cli.Type() != agent.TypeExec {
		g.InvalidType(agent.TypeExec, cli.Type())
		return
	}

	ts := h.getTasks(id)
	if ts == nil {
		g.OK([]info{})
		return
	}

	var ret []info
	ts.list(func(t *task) {
		if t.isDone() {
			return
		}
		ret = append(ret, info{
			TaskID: t.id,
			Pid:    t.pid,
			Begin:  t.begin.Unix(),
		})
	})

	g.OK(ret)
}
