package layout

import (
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type nodeStatus struct {
	ID     string `json:"id" example:"exec-01" validate:"required"`                                  // 节点ID
	Begin  int64  `json:"begin" example:"1663816771" validate:"required"`                            // 开始时间
	End    int64  `json:"end" example:"1663816771"`                                                  // 仅当status=done时返回
	Status string `json:"status" example:"waiting" enums:"waiting,running,done" validate:"required"` // 当前状态
	Err    string `json:"err" example:"xxx not found"`                                               // 错误信息
}

type status struct {
	Begin int64        `json:"begin" example:"1663816771" validate:"required"` // 开始时间
	End   int64        `json:"end" example:"1663816771"`                       // 结束时间，仅当done=true时返回
	Done  bool         `json:"done" example:"false" validate:"required"`       // 该任务是否已执行完毕
	Index int          `json:"index" example:"0" validate:"required"`          // 当前正在运行的批次号
	Nodes []nodeStatus `json:"nodes"`                                          // 节点状态列表
}

// run 获取批量任务状态
// @ID /api/layout/status
// @Summary 获取批量任务状态
// @Tags layout
// @Accept  json
// @Produce json
// @Param   id   query string  true  "任务ID"
// @Success 200  {object}      api.Success{payload=status}
// @Router /layout/status/{id} [get]
func (h *Handler) status(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")

	h.RLock()
	t := h.tasks[id]
	h.RUnlock()

	if t == nil {
		g.Notfound("task")
	}

	var ret status
	ret.Begin = t.Begin.Unix()
	ret.End = t.End.Unix()
	ret.Done = t.Done
	ret.Index = t.Index
	for _, id := range t.IDS {
		t.RLock()
		err := t.NodeErrs[id]
		status := t.NodeStatus[id]
		begin := t.NodeBegin[id]
		end := t.NodeEnd[id]
		t.RUnlock()
		var errString string
		if err != nil {
			errString = err.Error()
		}
		ret.Nodes = append(ret.Nodes, nodeStatus{
			ID:     id,
			Begin:  begin.Unix(),
			End:    end.Unix(),
			Status: status.String(),
			Err:    errString,
		})
	}

	g.OK(ret)
}
