package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

// SendFileLs send ls command
func (agent *Agent) SendIPMIDeviceInfo(req *anet.IPMICommonRequest) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeIPMIDeviceInfoReq
	msg.IPMICommonReq = req
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}
