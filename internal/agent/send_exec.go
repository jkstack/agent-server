package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

func (agent *Agent) SendExecRun(
	cmd string, args []string,
	auth, user string,
	workDir string, env []string,
	timeout int) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeExec
	msg.Exec = &anet.ExecPayload{
		Cmd:     cmd,
		Args:    args,
		Auth:    auth,
		User:    user,
		WorkDir: workDir,
		Env:     env,
		Timeout: timeout,
	}
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}
