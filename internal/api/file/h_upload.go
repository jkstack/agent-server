package file

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"server/internal/agent"
	"server/internal/api"
	lutils "server/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/utils"
)

const uploadLimit = 1024 * 1024

// upload 上传文件
//	@ID			/api/file/upload
//	@Summary	上传文件
//	@Tags		file
//	@Accept		mpfd
//	@Produce	json
//	@Param		id			path		string	true	"节点ID"
//	@Param		dir			formData	string	true	"保存路径"
//	@Param		file		formData	file	true	"文件"
//	@Param		md5			formData	string	false	"md5校验码"
//	@Param		mod			formData	int		false	"文件权限（8进制）"	default(0644)
//	@Param		own_user	formData	string	false	"文件所属用户"
//	@Param		own_group	formData	string	false	"文件所属分组"
//	@Param		timeout		formData	int		false	"超时时间"	default(60)
//	@Success	200			{object}	api.Success{payload}
//	@Router		/file/{id}/upload [post]
func (h *Handler) upload(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	dir := g.PostForm("dir")
	if len(dir) == 0 {
		g.MissingParam("dir")
	}
	hdr, err := g.FormFile("file")
	if err != nil {
		g.BadParam("file:" + err.Error())
	}
	file, err := hdr.Open()
	if err != nil {
		g.BadParam("file:" + err.Error())
	}
	defer file.Close()
	md5sum := g.PostForm("md5")
	modStr := g.DefaultPostForm("mod", "0644")
	mod := 0644
	if len(modStr) > 0 {
		n, err := strconv.ParseInt(modStr, 8, 64)
		if err != nil {
			g.BadParam("mod:" + err.Error())
		}
		mod = int(n)
	}
	ownUser := g.PostForm("own_user")
	ownGroup := g.PostForm("own_group")
	timeoutStr := g.DefaultPostForm("timeout", "60")
	timeout := 60
	if len(timeoutStr) > 0 {
		n, err := strconv.ParseInt(timeoutStr, 10, 64)
		if err != nil {
			g.BadParam("timeout:" + err.Error())
		}
		timeout = int(n)
	}

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.HTTPError(http.StatusNotFound, "agent not found")
		return
	}
	if cli.Type() != agent.TypeExec {
		g.InvalidType(agent.TypeExec, cli.Type())
	}

	var srcMD5 [md5.Size]byte
	if len(md5sum) > 0 {
		encData, err := hex.DecodeString(md5sum)
		if err != nil {
			g.BadParam("md5:" + err.Error())
		}

		copy(srcMD5[:], encData)
	}

	var taskID string
	if hdr.Size <= uploadLimit {
		var data []byte
		data, err = io.ReadAll(file)
		utils.Assert(err)
		dstMD5 := md5.Sum(data)
		info := agent.UploadContext{
			Dir:      dir,
			Name:     hdr.Filename,
			Mod:      mod,
			OwnUser:  ownUser,
			OwnGroup: ownGroup,
			Size:     uint64(hdr.Size),
			Data:     data,
		}
		if len(md5sum) > 0 {
			if !bytes.Equal(srcMD5[:], dstMD5[:]) {
				g.BadParam("md5:invalid checksum")
			}
			info.Md5 = srcMD5
		} else {
			info.Md5 = md5.Sum(data)
		}
		taskID, err = cli.SendUpload(info, "")
	} else {
		var tmpDir string
		var dstMD5 [md5.Size]byte
		tmpDir, dstMD5, err = dumpFile(file, path.Join(h.cfg.CacheDir, "upload"))
		utils.Assert(err)
		defer os.Remove(tmpDir)
		if len(md5sum) > 0 && !bytes.Equal(srcMD5[:], dstMD5[:]) {
			g.BadParam("md5:invalid checksum")
			return
		}
		var token string
		token, err = utils.UUID(16, "0123456789abcdef")
		utils.Assert(err)
		taskID, err = lutils.TaskID()
		utils.Assert(err)
		uri := "/api/file/upload/" + taskID
		h.logUploadCache(taskID, tmpDir, token,
			time.Now().Add(time.Duration(timeout)*time.Second), true)
		defer h.removeUploadCache(taskID)
		taskID, err = cli.SendUpload(agent.UploadContext{
			Dir:      dir,
			Name:     hdr.Filename,
			Mod:      mod,
			OwnUser:  ownUser,
			OwnGroup: ownGroup,
			Size:     uint64(hdr.Size),
			Md5:      dstMD5,
			URI:      uri,
			Token:    token,
		}, taskID)
	}

	utils.Assert(err)
	defer cli.ChanClose(taskID)

	logging.Info("upload [%s] to %s on %s, task_id=%s",
		hdr.Filename, dir, id, taskID)

	var rep *anet.Msg
	select {
	case rep = <-cli.ChanRead(taskID):
	case <-time.After(time.Duration(timeout) * time.Second):
		g.Timeout()
	}

	switch {
	case rep.Type == anet.TypeError:
		g.ERR(http.StatusServiceUnavailable, rep.ErrorMsg)
		return
	case rep.Type != anet.TypeUploadRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", rep.Type))
		return
	}

	if !rep.UploadRep.OK {
		g.ERR(http.StatusServiceUnavailable, rep.UploadRep.ErrMsg)
		return
	}

	g.OK(nil)
}

func dumpFile(f multipart.File, dir string) (string, [md5.Size]byte, error) {
	var md [md5.Size]byte
	os.MkdirAll(dir, 0755)
	dst, err := os.CreateTemp(dir, "ul")
	if err != nil {
		return "", md, err
	}
	defer dst.Close()
	enc := md5.New()
	_, err = io.Copy(io.MultiWriter(dst, enc), f)
	if err != nil {
		return "", md, err
	}
	copy(md[:], enc.Sum(nil))
	return dst.Name(), md, nil
}
