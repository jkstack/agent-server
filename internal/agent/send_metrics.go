package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

func (agent *Agent) SendHMStaticReq() (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeHMStaticReq
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}
