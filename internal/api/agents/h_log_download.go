package agents

import (
	"fmt"
	"net/http"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/utils"
)

type downloadArgs struct {
	Files   []string `form:"files"`
	Timeout int      `form:"timeout,default=600" binding:"min=1"`
}

// info 获取某个节点下的日志文件列表
//
//	@ID			/api/agents/log/download
//	@Summary	下载某个agent下的日志文件
//	@Tags		agents
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string		true	"节点ID"
//	@Param		files	query		[]string	true	"文件列表"
//	@Param		timeout	query		int			false	"超时时间"	default(600)	minimum(1)
//	@Success	200		{string}	string		"文件内容"
//	@Failure	500		{string}	string		"出错原因"
//	@Failure	503		{string}	string		"出错原因"
//	@Router		/agents/{id}/log/download [get]
func (h *Handler) logDownload(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	var args downloadArgs
	if err := g.ShouldBindQuery(&args); err != nil {
		g.BadParam(err.Error())
		return
	}
	if len(args.Files) == 0 {
		g.MissingParam("files")
		return
	}

	agents := g.GetAgents()

	agent := agents.Get(id)
	if agent == nil {
		g.Notfound("agent")
		return
	}

	taskID, err := agent.SendLogDownload(args.Files)
	utils.Assert(err)

	var rep *anet.Msg
	select {
	case rep = <-agent.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		g.Timeout()
	}

	switch {
	case rep.Type == anet.TypeError:
		g.ERR(http.StatusServiceUnavailable, rep.ErrorMsg)
		return
	case rep.Type != anet.TypeLogDownloadInfo:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", rep.Type))
		return
	}

	if !rep.LogDownloadInfo.OK {
		g.String(http.StatusServiceUnavailable, rep.LogDownloadInfo.ErrMsg)
		return
	}

	f, err := newCacheFile(h.cacheDir, rep.LogDownloadInfo.MD5)
	if err != nil {
		logging.Error("download failed, task_id=%s, err=%v", taskID, err)
		g.HTTPError(http.StatusInternalServerError, err.Error())
		return
	}
	defer f.Remove()

	left := rep.LogDownloadInfo.Size
	after := time.After(time.Duration(args.Timeout) * time.Second)
	for {
		var msg *anet.Msg
		select {
		case msg = <-agent.ChanRead(taskID):
		case <-after:
			g.Timeout()
			return
		}
		switch msg.Type {
		case anet.TypeLogDownloadData:
			n, err := f.Write(msg.LogDownloadData.Offset, msg.LogDownloadData.Data)
			if err != nil {
				logging.Error("download failed, task_id=%s, err=%v", taskID, err)
				g.HTTPError(http.StatusInternalServerError, err.Error())
				return
			}
			left -= uint64(n)
			if left == 0 {
				if !f.check() {
					g.HTTPError(http.StatusInternalServerError, "invalid checksum")
					return
				}
				g.FileAttachment(f.Name(), "log.zip")
				return
			}
		case anet.TypeDownloadError:
			g.HTTPError(http.StatusServiceUnavailable, msg.DownloadError.Msg)
			return
		}
	}
}
