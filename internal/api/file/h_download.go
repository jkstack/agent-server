package file

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"server/internal/agent"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/utils"
)

type downloadArgs struct {
	Dir     string `form:"dir"`
	Timeout int    `form:"timeout,default=600" binding:"min=1"`
}

// download 下载文件
// @ID /api/file/download
// @Summary 下载文件
// @Tags file
// @Accept  json
// @Produce plain
// @Param   id      path  string true  "节点ID"
// @Param   dir     query string true  "文件路径"
// @Param   timeout query int    false "超时时间" default(600)  minimum(1)
// @Success 200 {string}  string "输出内容"
// @Failure 404 {string}  string "file not found"
// @Failure 500 {string}  string "出错原因"
// @Failure 503 {string}  string "出错原因"
// @Router /file/{id}/download [get]
func (h *Handler) download(gin *gin.Context) {
	g := api.GetG(gin)

	var args downloadArgs
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

	taskID, err := cli.SendFileDownload(args.Dir)
	utils.Assert(err)
	defer cli.ChanClose(taskID)

	var rep *anet.Msg
	select {
	case rep = <-cli.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		g.Timeout()
	}

	switch {
	case rep.Type == anet.TypeError:
		g.ERR(http.StatusServiceUnavailable, rep.ErrorMsg)
		return
	case rep.Type != anet.TypeDownloadRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", rep.Type))
		return
	}

	if !rep.DownloadRep.OK {
		logging.Error("download [%s] on %s failed, task_id=%s, msg=%s", args.Dir, id, taskID, rep.DownloadRep.ErrMsg)
		g.HTTPError(http.StatusServiceUnavailable, rep.DownloadRep.ErrMsg)
		return
	}
	logging.Info("download [%s] on %s success, task_id=%s, size=%d, md5=%x", args.Dir, id, taskID,
		rep.DownloadRep.Size, rep.DownloadRep.MD5)

	f, err := tmpFile(h.cfg.CacheDir, rep.DownloadRep.Size)
	if err != nil {
		logging.Error("download [%s] on %s failed, task_id=%s, err=%v", args.Dir, id, taskID, err)
		g.HTTPError(http.StatusInternalServerError, err.Error())
		return
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	left := rep.DownloadRep.Size
	after := time.After(time.Duration(args.Timeout) * time.Second)
	for {
		var msg *anet.Msg
		select {
		case msg = <-cli.ChanRead(taskID):
		case <-after:
			g.Timeout()
			return
		}
		switch msg.Type {
		case anet.TypeDownloadData:
			n, err := writeFile(f, msg.DownloadData)
			if err != nil {
				logging.Error("download [%s] on %s, task_id=%s, err=%v", args.Dir, id, taskID, err)
				g.HTTPError(http.StatusInternalServerError, err.Error())
				return
			}
			left -= uint64(n)
			if left == 0 {
				serveFile(g, f, id, args.Dir, taskID, rep.DownloadRep.MD5)
				return
			}
		case anet.TypeDownloadError:
			g.HTTPError(http.StatusServiceUnavailable, msg.DownloadError.Msg)
			return
		}
	}
}

func tmpFile(cacheDir string, size uint64) (*os.File, error) {
	tmp := path.Join(cacheDir, "download")
	os.MkdirAll(tmp, 0755)
	f, err := os.CreateTemp(tmp, "dl")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = fillFile(f, size)
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, err
	}
	f.Close()
	f, err = os.OpenFile(f.Name(), os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, err
	}
	return f, err
}

func serveFile(g *api.GContext, f *os.File, id, dir, taskID string, src [md5.Size]byte) {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		logging.Error("download [%s] on %s, task_id=%s, err=%v", dir, id, taskID, err)
		g.HTTPError(http.StatusInternalServerError, err.Error())
		return
	}
	dst, err := md5From(f)
	if err != nil {
		logging.Error("download [%s] on %s, task_id=%s, err=%v", dir, id, taskID, err)
		g.HTTPError(http.StatusInternalServerError, err.Error())
		return
	}
	if !bytes.Equal(dst[:], src[:]) {
		logging.Error("download [%s] on %s, task_id=%s, invalid md5checksum, src=%x, dst=%x",
			dir, id, taskID, src, dst)
		g.HTTPError(http.StatusInternalServerError, "invalid checksum")
		return
	}
	f.Close()
	g.File(f.Name())
}
