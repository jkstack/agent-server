package agents

import (
	"server/internal/agent"
	"server/internal/api"
	"sort"

	"github.com/gin-gonic/gin"
)

type listArgs struct {
	Type string `form:"type"`
	OS   string `form:"os,default=all"`
	Page int    `form:"page,default=1" binding:"min=1"`
	Size int    `form:"size,default=20" binding:"min=10"`
}

// list 列出节点列表
// @ID /api/agents
// @Summary 获取节点列表
// @Tags agents
// @Accept  json
// @Produce json
// @Param   type  query string  false "节点类型,不指定则为所有" enums(example-agent,container-agent,metrics-agent,...)
// @Param   os    query string  false "节点操作系统,不指定则为所有" enums(linux,windows)
// @Param   page  query int     false "分页编号" default(1)  minimum(1)
// @Param   size  query int     false "每页数量" default(20) minimum(10)
// @Success 200   {object}      api.Success{payload=[]info}
// @Router /agents [get]
func (h *Handler) list(gin *gin.Context) {
	g := api.GetG(gin)

	var args listArgs
	if err := g.ShouldBindQuery(&args); err != nil {
		g.BadParam(err.Error())
		return
	}

	switch args.OS {
	case "linux", "windows":
	default:
		args.OS = ""
	}

	agents := g.GetAgents()

	ret := make([]info, 0, agents.Size())
	agents.Range(func(agent *agent.Agent) bool {
		want := true
		if len(args.Type) > 0 && agent.Type() != args.Type {
			want = false
		}
		if len(args.OS) > 0 && agent.Info().OS != args.OS {
			want = false
		}
		if want {
			a := agent.Info()
			ret = append(ret, info{
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
		return true
	})
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].ID < ret[j].ID
	})
	offset := (args.Page - 1) * args.Size
	if offset >= len(ret) {
		g.OK(nil)
		return
	}
	end := offset + args.Size
	if end > len(ret) {
		end = len(ret)
	}
	g.OK(ret[offset:end])
}
