package rpa

import "github.com/gin-gonic/gin"

type inSelectorArgs struct {
	CallBack string `json:"callback" example:"https://www.baidu.com"` // 回调地址
}

// inSelector 进入元素选择器状态
//
//	@ID						/api/rpa/in_selector
//	@Description.markdown	in_selector.md
//	@Summary				进入元素选择器状态
//	@Tags					rpa
//	@Accept					json
//	@Produce				json
//	@Param					id		path		string			true	"节点ID"
//	@Param					args	body		inSelectorArgs	true	"需启动的任务列表"
//	@Success				200		{object}	api.Success
//	@Router					/rpa/{id}/in_selector [post]
func (h *Handler) inSelector(gin *gin.Context) {
}
