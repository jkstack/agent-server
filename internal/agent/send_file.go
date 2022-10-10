package agent

import (
	"server/internal/utils"

	"github.com/jkstack/anet"
)

func (agent *Agent) SendFileLs(dir string) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeLsReq
	msg.LSReq = &anet.LsReq{Dir: dir}
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}

func (agent *Agent) SendFileDownload(dir string) (string, error) {
	id, err := utils.TaskID()
	if err != nil {
		return "", err
	}
	var msg anet.Msg
	msg.Type = anet.TypeDownloadReq
	msg.DownloadReq = &anet.DownloadReq{Dir: dir}
	msg.TaskID = id
	agent.Lock()
	agent.taskRead[id] = make(chan *anet.Msg)
	agent.Unlock()
	agent.chWrite <- &msg
	return id, nil
}
