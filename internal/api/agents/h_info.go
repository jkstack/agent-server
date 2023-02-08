package agents

import (
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type info struct {
	ID       string `json:"id" example:"ct-01" validate:"required"`                                  // agent-id
	Type     string `json:"type" example:"container-agent" validate:"required"`                      // agent类型
	Version  string `json:"version" example:"1.0.0" validate:"required"`                             // agent版本号
	IP       string `json:"ip" example:"127.0.0.1" validate:"required"`                              // ip地址
	MAC      string `json:"mac" example:"00:15:5d:c9:e0:17" validate:"required"`                     // mac地址
	OS       string `json:"os" example:"linux" enums:"windows,linux" validate:"required"`            // 操作系统类型
	Platform string `json:"platform" example:"debian" enums:"debian,centos,..." validate:"required"` // 操作系统名称
	Arch     string `json:"arch" example:"x86_64" enums:"i386,x86_64,..." validate:"required"`       // 操作系统位数
	IsBusy   bool   `json:"is_busy" example:"false" validate:"required"`                             // 是否处于忙碌状态
}

// info 获取某个节点信息
//	@ID			/api/agents/info
//	@Summary	获取某个节点信息
//	@Tags		agents
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"节点ID"
//	@Success	200	{object}	api.Success{payload=info}
//	@Router		/agents/{id} [get]
func (h *Handler) info(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")

	agents := g.GetAgents()

	agent := agents.Get(id)
	if agent == nil {
		g.Notfound("agent")
		return
	}
	a := agent.Info()
	g.OK(info{
		ID:       agent.ID(),
		Type:     agent.Type(),
		Version:  a.Version,
		IP:       a.IP.String(),
		MAC:      a.MAC,
		OS:       a.OS,
		Platform: a.Platform,
		Arch:     a.Arch,
		IsBusy:   false, // TODO
	})
}
