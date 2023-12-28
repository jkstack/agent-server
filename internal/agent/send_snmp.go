package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

// SendSNMPList send snmp list command
func (agent *Agent) SendSNMPList(req *anet.SNMPReq) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeSNMPListReq
	msg.SNMPReq = req
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}
