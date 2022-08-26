package info

import (
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type serverInfo struct {
	Version    string `json:"version" example:"1.0.0"`     // 服务端版本号
	Agents     int    `json:"agents" example:"10"`         // 当前连接的agent数量
	IsBlocking bool   `json:"is_blocking" example:"false"` // 是否处于限流状态
}

// server 获取当前服务器状态
// @ID /api/info/server
// @Summary 获取当前服务器状态
// @Tags info
// @Produce json
// @Success 200  {object}  api.Success{payload=serverInfo}
// @Router /info/server [get]
func (h *Handler) server(gin *gin.Context) {
	g := api.GetG(gin)

	agents := g.GetAgents()

	var info serverInfo
	info.Version = h.version
	info.Agents = agents.Size()
	info.IsBlocking = *h.isBlocking

	g.OK(info)
}