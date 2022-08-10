package agent

import (
	"net/http"
	"server/internal/agent"

	"github.com/gin-gonic/gin"
)

// list 列出主机列表
// @ID /api/agent/list
// @Produce json
// @Router /agent/list [get]
func (h *Handler) list(agents *agent.Agents, g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{
		"a": "b",
	})
}
