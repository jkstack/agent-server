package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

// SendIPMIDeviceInfo send ipmi device_info command
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

// SendIPMISensorList send ipmi sensor_list command
func (agent *Agent) SendIPMISensorList(req *anet.IPMICommonRequest) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeIPMISensorListReq
	msg.IPMICommonReq = req
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}
