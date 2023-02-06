package agents

import "github.com/gin-gonic/gin"

// info 获取某个节点下的日志文件列表
//
//	@ID			/api/agents/log
//	@Summary	下载某个agent下的日志文件
//	@Tags		agents
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string		true	"节点ID"
//	@Param		files	query		[]string	true	"文件列表"
//	@Success	200		{string}	string		"文件内容"
//	@Failure	500		{string}	string		"出错原因"
//	@Failure	503		{string}	string		"出错原因"
//	@Router		/agents/{id}/log/download [get]
func (h *Handler) logDownload(gin *gin.Context) {
}
