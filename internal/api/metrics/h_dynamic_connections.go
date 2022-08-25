package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
)

type connection struct {
}

// static 获取节点的连接列表数据
// @ID /api/metrics/dynamic/connections
// @Summary 获取节点的连接列表数据
// @Tags metrics
// @Produce json
// @Param   id   path string  true "节点ID"
// @Success 200  {object}     api.Success{payload=[]connection}
// @Router /metrics/{id}/dynamic/connections [get]
func (h *Handler) dynamicConnections(gin *gin.Context) {
}

func transDynamicConnections(conns []anet.HMDynamicConnection) []connection {
	return nil
}
