package exec

import (
	"io"
	"net/http"
	"server/internal/api"
	"strconv"

	"github.com/gin-gonic/gin"
)

// pty 获取返回内容
//	@ID			/api/exec/pty
//	@Summary	获取返回内容
//	@Tags		exec
//	@Accept		json
//	@Produce	plain
//	@Param		id	path		string	true	"节点ID"
//	@Param		pid	path		int		true	"进程号"
//	@Success	200	{string}	string	"输出内容"
//	@Failure	404	{string}	string	"\<what\> not found"
//	@Failure	500	{string}	string	"出错原因"
//	@Router		/exec/{id}/pty/{pid} [get]
func (h *Handler) pty(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	pid, err := strconv.ParseInt(g.Param("pid"), 10, 64)
	if err != nil {
		g.BadParam("pid")
		return
	}

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.HTTPError(http.StatusNotFound, "agent not found")
		return
	}

	task := h.getTasksOrCreate(id).get(int(pid))
	if task == nil {
		g.HTTPError(http.StatusNotFound, "task not found")
		return
	}

	data, err := io.ReadAll(task.cache)
	if err != nil {
		g.HTTPError(http.StatusInternalServerError, err.Error())
		return
	}

	g.HTTPData(data)
}
