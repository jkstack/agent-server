package exec

import (
	"server/internal/api"
	"strconv"

	"github.com/gin-gonic/gin"
)

type status struct {
	Begin int64 `json:"begin" example:"1663816771" validate:"required"` // 任务开始时间
	End   int64 `json:"end,omitempty" example:"1663816771"`             // 结束时间，仅当running=false时返回
	Done  bool  `json:"done" example:"false" validate:"required"`       // 该任务是否已执行完毕
	Code  int   `json:"code" example:"0"`                               // 运行结果状态码
}

// status 获取运行状态
// @ID /api/exec/status
// @Summary 获取运行状态
// @Tags exec
// @Produce json
// @Param   id  path string true "节点ID"
// @Param   pid path int    true "进程号"
// @Success 200 {object}    api.Success{payload=status}
// @Router /exec/{id}/status/{pid} [get]
func (h *Handler) status(gin *gin.Context) {
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

	h.RLock()
	task := h.tasks[int(pid)]
	h.RUnlock()

	if task == nil {
		g.NotFound("pid")
		return
	}

	ret := status{
		Begin: task.begin.Unix(),
		Done:  task.doneFlag,
	}

	if ret.Done {
		ret.End = task.end.Unix()
		ret.Code = task.code
	}

	g.OK(ret)
}
