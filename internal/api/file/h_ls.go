package file

import (
	"fmt"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/utils"
)

type lsArgs struct {
	Dir string `form:"dir"`
}

type info struct {
	Name    string `json:"name" example:"bin" validate:"required"`            // 文件名
	Auth    uint64 `json:"auth" example:"0777" validate:"required"`           // 文件权限（10进制）
	User    string `json:"user,omitempty" example:"root"`                     // 所属用户
	Group   string `json:"group,omitempty" example:"root"`                    // 所属组
	Size    uint64 `json:"size" example:"155" validate:"required"`            // 文件大小
	ModTime int64  `json:"mod_time" example:"1663816771" validate:"required"` // 更新时间
	IsDir   bool   `json:"is_dir" example:"true" validate:"required"`         // 是否是目录
	IsLink  bool   `json:"is_link" example:"false" validate:"required"`       // 是否是软链
	LinkDir string `json:"link_dir,omitempty" example:"/usr/bin"`             // 连接路径
}

// ls 查询文件列表
// @ID /api/file/ls
// @Summary 查询文件列表
// @Tags file
// @Accept  json
// @Produce json
// @Param   id   path  string true "节点ID"
// @Param   dir  query string true "查询路径"
// @Success 200  {object}     api.Success{payload=[]info}
// @Router /file/{id}/ls [get]
func (h *Handler) ls(gin *gin.Context) {
	g := api.GetG(gin)

	var args lsArgs
	if err := g.ShouldBindQuery(&args); err != nil {
		g.BadParam(err.Error())
		return
	}

	id := g.Param("id")

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.HTTPError(http.StatusNotFound, "agent not found")
		return
	}
	if cli.Type() != agent.TypeExec {
		g.InvalidType(agent.TypeExec, cli.Type())
	}

	taskID, err := cli.SendFileLs(args.Dir)
	utils.Assert(err)
	defer cli.ChanClose(taskID)

	var msg *anet.Msg
	select {
	case msg = <-cli.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		g.Timeout()
	}

	switch {
	case msg.Type == anet.TypeError:
		g.ERR(http.StatusServiceUnavailable, msg.ErrorMsg)
		return
	case msg.Type != anet.TypeLsRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	if !msg.LSRep.OK {
		g.ERR(http.StatusServiceUnavailable, msg.LSRep.ErrMsg)
		return
	}

	ret := make([]info, len(msg.LSRep.Files))
	for i, file := range msg.LSRep.Files {
		var ix info
		ix.Name = file.Name
		ix.Auth = uint64(file.Mod)
		ix.User = file.User
		ix.Group = file.Group
		ix.Size = file.Size
		ix.ModTime = file.ModTime.Unix()
		ix.IsDir = file.Mod.IsDir()
		ix.IsLink = file.IsLink
		ix.LinkDir = file.LinkDir
		ret[i] = ix
	}

	g.OK(ret)
}
