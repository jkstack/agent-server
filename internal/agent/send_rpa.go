package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

// SendRpaRun send rpa run command
func (agent *Agent) SendRpaRun(url string, isDebug bool) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeRPARun
	msg.RPARun = &anet.RPARunArgs{
		URL:     url,
		IsDebug: isDebug,
	}
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}

// SendRpaCtrl send rpa control command
func (agent *Agent) SendRpaCtrl(taskID string, status int) error {
	var msg anet.Msg
	msg.Type = anet.TypeRPAControlReq
	msg.RPACtrlReq = &anet.RPACtrlReq{
		Status: status,
	}
	msg.TaskID = taskID
	agent.chWrite <- &msg
	return nil
}

// SendRpaInSelector send rpa in selector command
func (agent *Agent) SendRpaInSelector() (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeRPASelectorReq
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}

// SendRpaSelectorValidate send rpa in selector validate request
func (agent *Agent) SendRpaSelectorValidate(content string) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeRPASelectorValidateReq
	msg.TaskID = id
	msg.RPASelectorValidateReq = &anet.RPASelectorValidateReq{
		Content: content,
	}
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}
