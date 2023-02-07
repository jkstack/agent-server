package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

// SendLogLs send ls log command
func (agent *Agent) SendLogLs() (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeLogLsReq
	msg.TaskID = id
	agent.chWrite <- &msg
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	return id, nil
}
