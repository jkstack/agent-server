package agents

import (
	"fmt"
	"net/http"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/utils"
)

type fileInfo struct {
	Name    string `json:"name" example:"metrics-agent.log" validate:"required"` // 文件名
	Size    uint64 `json:"size" example:"155" validate:"required"`               // 文件大小
	ModTime int64  `json:"mod_time" example:"1663816771" validate:"required"`    // 修改时间
}

// info 获取某个节点下的日志文件列表
//
//	@ID			/api/agents/logs
//	@Summary	获取某个节点下的日志文件列表
//	@Tags		agents
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"节点ID"
//	@Success	200	{object}	api.Success{payload=[]fileInfo}
//	@Router		/agents/{id}/logs [get]
func (h *Handler) logs(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")

	agents := g.GetAgents()

	agent := agents.Get(id)
	if agent == nil {
		g.Notfound("agent")
		return
	}

	taskID, err := agent.SendLogLs()
	utils.Assert(err)

	var msg *anet.Msg
	select {
	case msg = <-agent.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		g.Timeout()
	}

	switch {
	case msg.Type == anet.TypeError:
		g.ERR(http.StatusServiceUnavailable, msg.ErrorMsg)
		return
	case msg.Type != anet.TypeLogLsRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	var files []fileInfo
	for _, file := range msg.LsLog.Files {
		files = append(files, fileInfo{
			Name:    file.Name,
			Size:    file.Size,
			ModTime: file.ModTime.Unix(),
		})
	}

	g.OK(files)
}
