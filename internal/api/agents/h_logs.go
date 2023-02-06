package agents

import "github.com/gin-gonic/gin"

type fileInfo struct {
	Name    string `json:"name" example:"metrics-agent.log" validate:"required"` // 文件名
	Size    uint64 `json:"size" example:"155" validate:"required"`               // 文件大小
	ModTime int64  `json:"mod_time" example:"1663816771" validate:"required"`    // 修改时间
}

// info 获取某个节点下的日志文件列表
//
//	@ID			/api/agents/logs
//	@Summary	获取某个节点下的日志文件列表
//	@Tags		agents
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"节点ID"
//	@Success	200	{object}	api.Success{payload=fileInfo}
//	@Router		/agents/{id}/logs [get]
func (h *Handler) logs(gin *gin.Context) {
}
