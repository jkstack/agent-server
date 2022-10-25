package scriptengine

import (
	"encoding/base64"
	"errors"
	"fmt"
	"server/internal/api"
	"time"

	"github.com/jkstack/anet"
)

const (
	codeUnexpectedError = -65534
	codeTimeout         = -65535
)

func (e *Engine) run() bool {
	e.status = StatusRunning
	var cmd string
	var args []string
	switch e.args.Type {
	case TypeSh:
		cmd = "sh"
		args = append([]string{e.dir}, e.args.Args...)
	case TypeBash:
		cmd = "bash"
		args = append([]string{e.dir}, e.args.Args...)
	case TypePython:
		cmd = "python"
		args = append([]string{e.dir}, e.args.Args...)
	case TypePython3:
		cmd = "python3"
		args = append([]string{e.dir}, e.args.Args...)
	case TypeBat:
		cmd = e.dir
		args = e.args.Args
	case TypePowerShell:
		cmd = "powershell"
		args = append([]string{"-ExecutionPolicy", "remotesigned", e.dir}, e.args.Args...)
	case TypePhp:
		cmd = "php"
		args = append([]string{e.dir}, e.args.Args...)
	case TypeLua:
		cmd = "lua"
		args = append([]string{e.dir}, e.args.Args...)
	}
	taskID, err := e.cli.SendExecRun(cmd, args, e.args.Auth, e.args.User,
		e.args.WorkDir, e.args.Env, e.args.Timeout, e.dir)
	if err != nil {
		e.Err = err
		return false
	}

	var msg *anet.Msg
	select {
	case msg = <-e.cli.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		e.Err = ErrTimeout
		return false
	}

	switch {
	case msg.Type == anet.TypeError:
		e.Err = errors.New(msg.ErrorMsg)
		return false
	case msg.Type != anet.TypeExecd:
		e.Err = fmt.Errorf("invalid message type: %d", msg.Type)
		return false
	}

	if !msg.Execd.OK {
		e.Err = errors.New(msg.Execd.Msg)
		return false
	}

	e.Pid = msg.Execd.Pid

	ch := e.cli.ChanRead(taskID)
	done := time.After(time.Duration(e.args.Timeout) * time.Second)
	for {
		select {
		case msg := <-ch:
			switch msg.Type {
			case anet.TypeExecData:
				data, err := base64.StdEncoding.DecodeString(msg.ExecData.Data)
				if err != nil {
					e.Err = err
					return false
				}
				e.Err = e.OnData(data)
				if e.Err != nil {
					return false
				}
			case anet.TypeExecDone:
				e.Code = msg.ExecDone.Code
				return true
			default:
				e.Code = codeUnexpectedError
				e.Err = ErrUnexpectedError
				return false
			}
		case <-done:
			e.Code = codeTimeout
			e.Err = ErrTimeout
			return false
		}
	}
}
