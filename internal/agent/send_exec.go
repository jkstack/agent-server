package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

// SendExecRun send exec task
func (agent *Agent) SendExecRun(
	cmd string, args []string,
	auth, user string,
	workDir string, env []string,
	timeout int, deferRm ...string) (string, error) {
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
	if len(deferRm) > 0 {
		msg.Exec.DeferRM = deferRm[0]
	}
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}

// SendExecKill send kill command
func (agent *Agent) SendExecKill(pid int) error {
	id, err := utils.TaskID()
	if err != nil {
		return err
	}
	var msg anet.Msg
	msg.Type = anet.TypeExecKill
	msg.ExecKill = &anet.ExecKill{
		Pid: pid,
	}
	msg.TaskID = id
	agent.chWrite <- &msg
	return nil
}
