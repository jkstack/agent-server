package scriptengine

import (
	"crypto/md5"
	"errors"
	"fmt"
	"server/internal/agent"
	"time"

	"github.com/jkstack/anet"
)

func (e *Engine) upload() bool {
	e.status = StatusUploading
	info := agent.UploadContext{
		Dir:  "$$TMP$$",
		Name: e.fileName,
		Mod:  0644,
		Size: uint64(len(e.args.Data)),
		Md5:  md5.Sum([]byte(e.args.Data)),
		Data: []byte(e.args.Data),
	}
	taskID, err := e.cli.SendUpload(info, "")
	if err != nil {
		e.Err = err
		return false
	}
	defer e.cli.ChanClose(taskID)

	var rep *anet.Msg
	select {
	case rep = <-e.cli.ChanRead(taskID):
	case <-time.After(time.Duration(e.args.Timeout) * time.Second):
		e.Err = ErrTimeout
		return false
	}

	switch {
	case rep.Type == anet.TypeError:
		e.Err = errors.New(rep.ErrorMsg)
		return false
	case rep.Type != anet.TypeUploadRep:
		e.Err = fmt.Errorf("invalid message type: %d", rep.Type)
		return false
	}

	if !rep.UploadRep.OK {
		e.Err = errors.New(rep.UploadRep.ErrMsg)
		return false
	}
	e.dir = rep.UploadRep.Dir
	return true
}
