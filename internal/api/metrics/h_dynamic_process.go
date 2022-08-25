package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
)

type process struct {
}

// static 获取节点的进程列表数据
// @ID /api/metrics/dynamic/process
// @Summary 获取节点的所有进程列表数据
// @Tags metrics
// @Produce json
// @Param   id   path string  true "节点ID"
// @Success 200  {object}     api.Success{payload=[]process}
// @Router /metrics/{id}/dynamic/process [get]
func (h *Handler) dynamicProcess(gin *gin.Context) {
}

func transDynamicProcess(process []anet.HMDynamicProcess) []process {
	return nil
}
