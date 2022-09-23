package exec

import (
	"io"
	"net/http"
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

// pty 获取返回内容
// @ID /api/exec/pty
// @Summary 获取返回内容
// @Tags exec
// @Produce json
// @Param   id  path string true   "节点ID"
// @Param   pid path string true   "进程号"
// @Success 200 {string}    string "输出内容"
// @Failure 404 {string}    string "<what> not found"
// @Failure 500 {string}    string "错误内容"
// @Router /exec/{id}/pty/{pid} [get]
func (h *Handler) pty(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	pid := g.GetInt("pid")

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.HttpError(http.StatusNotFound, "agent not found")
		return
	}

	h.RLock()
	task := h.tasks[pid]
	h.RUnlock()

	if task == nil {
		g.HttpError(http.StatusNotFound, "task not found")
		return
	}

	data, err := io.ReadAll(task.cache)
	if err != nil {
		g.HttpError(http.StatusInternalServerError, err.Error())
		return
	}

	g.HttpData(data)
}
