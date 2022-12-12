package info

import (
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type serverInfo struct {
	ID         string `json:"id" example:"cluster-01" validate:"required"`     // 集群ID
	Version    string `json:"version" example:"1.0.0" validate:"required"`     // 服务端版本号
	Agents     int    `json:"agents" example:"10" validate:"required"`         // 当前连接的agent数量
	IsBlocking bool   `json:"is_blocking" example:"false" validate:"required"` // 是否处于限流状态
}

// server 获取当前服务器状态
// @ID /api/info/server
// @Summary 获取当前服务器状态
// @Tags info
// @Accept  json
// @Produce json
// @Success 200  {object}  api.Success{payload=serverInfo}
// @Router /info/server [get]
func (h *Handler) server(gin *gin.Context) {
	g := api.GetG(gin)

	agents := g.GetAgents()

	var info serverInfo
	info.ID = h.id
	info.Version = h.version
	info.Agents = agents.Size()
	info.IsBlocking = *h.isBlocking

	g.OK(info)
}
