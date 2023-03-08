package rpa

import (
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (svr *Server) Run(args *RunArgs, stream Rpa_RunServer) error {
	agentID := args.GetId()
	agent := svr.agents.Get(agentID)
	if agent == nil {
		return grpc.Errorf(codes.NotFound, "agent not found")
	}
	taskID, err := agent.SendRpaRun(args.GetUrl(), args.GetIsDebug())
	if err != nil {
		return err
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
			return grpc.ErrClientConnClosing
		}
		msg := <-ch
		switch msg.Type {
		case anet.TypeRPAControlRep:
			chRep <- msg.RPACtrlRep
		case anet.TypeRPALog:
			err := stream.Send(&Log{Data: string(*msg.RPALog)})
			if err != nil {
				logging.Error("send log error: %v", err)
				return nil
			}
		}
	}
}
