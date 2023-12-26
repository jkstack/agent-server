package rpa

import (
	"encoding/json"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Run run rpa workflow
func (svr *Server) Run(args *RunArgs, stream Rpa_RunServer) error {
	logging.Info("rpa run: %s", args.String())
	agentID := args.GetId()
	agent := svr.agents.Get(agentID)
	if agent == nil {
		return svr.runResponse(codes.NotFound, "agent not found", "")
	}
	svr.RLock()
	_, ok := svr.jobs[agentID]
	svr.RUnlock()
	if ok {
		return svr.runResponse(codes.Unavailable, "agent is busy", "")
	}
	taskID, err := agent.SendRpaRun(args.GetUrl(), args.GetIsDebug())
	if err != nil {
		return svr.runResponse(codes.Unavailable, "send message: "+err.Error(), "")
	}
	defer agent.ChanClose(taskID)
	chRep := make(chan *anet.RPACtrlRep, 1)
	svr.Lock()
	svr.jobs[agentID] = taskID
	svr.ctrlRep[taskID] = chRep
	svr.Unlock()
	defer func() {
		close(chRep)
		svr.Lock()
		delete(svr.jobs, agentID)
		delete(svr.ctrlRep, taskID)
		svr.Unlock()
	}()
	for {
		ch := agent.ChanRead(taskID)
		if ch == nil {
			return svr.runResponse(codes.Unavailable, "agent is closed", "")
		}
		msg := <-ch
		if msg == nil {
			return svr.runResponse(codes.Unavailable, "agent is closed", "")
		}
		switch msg.Type {
		case anet.TypeRPAControlRep:
			chRep <- msg.RPACtrlRep
		case anet.TypeRPALog:
			err := stream.Send(&Log{Data: string(*msg.RPALog)})
			if err != nil {
				logging.Error("send log error: %v", err)
				return svr.runResponse(codes.Internal, "send log error: "+err.Error(), "")
			}
		case anet.TypeRPAFinish:
			payload := msg.RPAFinish
			var code int
			if payload.Code == 1 {
				code = 0
			} else {
				code = 1
			}
			logging.Info("rpa finish(%d): %s", payload.Code, payload.Msg)
			return svr.runResponse(codes.Code(code), payload.Msg, payload.Data)
		}
	}
}

func (svr *Server) runResponse(code codes.Code, msg, data string) error {
	var rep struct {
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}
	rep.Msg = msg
	rep.Data = data
	payload, err := json.Marshal(rep)
	if err != nil {
		return status.Errorf(codes.Internal, `{"msg":"`+err.Error()+`","data":""}`)
	}
	return status.Errorf(code, "%s", string(payload))
}
