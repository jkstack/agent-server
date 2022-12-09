package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

// SendHMStaticReq send static data get request
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

// SendHMDynamicReq send dynamic data get request
func (agent *Agent) SendHMDynamicReq(req []anet.HMDynamicReqType,
	top int, kind []string) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeHMDynamicReq
	msg.TaskID = id
	msg.HMDynamicReq = &anet.HMDynamicReq{
		Req:        req,
		Top:        top,
		AllowConns: kind,
	}
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}

// SendHMQueryStatus send query status request
func (agent *Agent) SendHMQueryStatus() (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeHMQueryCollect
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}

// SendHMChangeStatus send change status task
func (agent *Agent) SendHMChangeStatus(jobs []string) error {
	id, err := utils.TaskID()
	if err != nil {
		return err
	}
	var msg anet.Msg
	msg.Type = anet.TypeHMChangeCollectStatus
	msg.TaskID = id
	msg.HMChangeStatus = &anet.HMChangeReportStatus{
		Jobs: jobs,
	}
	agent.chWrite <- &msg
	return nil
}
