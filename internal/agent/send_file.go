package agent

import (
	"crypto/md5"
	"os"
	"server/internal/utils"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/compress"
)

func (agent *Agent) SendFileLs(dir string) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeLsReq
	msg.LSReq = &anet.LsReq{Dir: dir}
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}

func (agent *Agent) SendFileDownload(dir string) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeDownloadReq
	msg.DownloadReq = &anet.DownloadReq{Dir: dir}
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}

type UploadContext struct {
	Dir      string
	Name     string
	Mod      int
	OwnUser  string
	OwnGroup string
	Size     uint64
	Md5      [md5.Size]byte
	Md5Check bool
	// send data
	Data []byte
	// send file
	Uri   string
	Token string
}

func (agent *Agent) SendUpload(ctx UploadContext, id string) (string, error) {
	if len(id) == 0 {
		var err error
		id, err = utils.TaskID()
		if err != nil {
			return "", err
		}
	}
	var msg anet.Msg
	msg.Type = anet.TypeUpload
	msg.TaskID = id
	msg.Upload = &anet.Upload{
		Dir:      ctx.Dir,
		Name:     ctx.Name,
		Mod:      os.FileMode(ctx.Mod),
		OwnUser:  ctx.OwnUser,
		OwnGroup: ctx.OwnGroup,
		Size:     ctx.Size,
		MD5:      ctx.Md5,
	}
	if len(ctx.Data) > 0 {
		msg.Upload.Data = compress.Compress(ctx.Data)
	} else if len(ctx.Uri) > 0 {
		msg.Upload.URI = ctx.Uri
		msg.Upload.Token = ctx.Token
	}
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}
