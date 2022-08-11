package agent

import (
	"server/internal/agent"
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

type info struct {
	ID       string `json:"id" example:"agent_id"`
	Type     string `json:"type" example:"agent类型"`
	Version  string `json:"version" example:"agent版本号"`
	IP       string `json:"ip" example:"ip地址"`
	MAC      string `json:"mac" example:"mac地址"`
	OS       string `json:"os" example:"操作系统类型" enums:"windows,linux"`
	Platform string `json:"platform" example:"操作系统名称" enums:"debian,centos,..."`
	Arch     string `json:"arch" example:"操作系统位数" enums:"i386,x86_64,..."`
}

// info 获取某个节点信息
// @ID /api/agent/info
// @Description 获取某个节点信息
// @Produce json
// @Param   id    path string  true "节点ID,不指定则列出所有节点"
// @Success 200   {object}     api.Success{payload=info}
// @Router /agent/info/{id} [get]
func (h *Handler) info(agents *agent.Agents, g *gin.Context) {
	id := g.Param("id")
	agent := agents.Get(id)
	a := agent.Info()
	api.OK(g, info{
		ID:       agent.ID(),
		Type:     agent.Type(),
		Version:  a.Version,
		IP:       a.IP.String(),
		MAC:      a.MAC,
		OS:       a.OS,
		Platform: a.Platform,
		Arch:     a.Arch,
	})
}
