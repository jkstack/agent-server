package rpa

import (
	"context"

	"github.com/jkstack/anet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (svr *Server) Control(ctx context.Context, args *ControlArgs) (*ControlResponse, error) {
	agentID := args.GetId()
	agent := svr.agents.Get(agentID)
	if agent == nil {
		return nil, grpc.Errorf(codes.NotFound, "agent not found")
	}
	svr.RLock()
	taskID := svr.jobs[agentID]
	svr.RUnlock()
	if len(taskID) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "task not found")
	}
	var status int
	switch args.GetSt() {
	case ControlArgs_pause:
		status = anet.RPAPause
	case ControlArgs_stop:
		status = anet.RPAStop
	case ControlArgs_continue:
		status = anet.RPAContinue
	}
	err := agent.SendRpaCtrl(taskID, status)
	if err != nil {
		return nil, err
	}
	svr.RLock()
	ch := svr.ctrlRep[taskID]
	svr.RUnlock()
	if ch == nil {
		return nil, grpc.Errorf(codes.Unavailable, "channel not found")
	}
	rep := <-ch
	if rep == nil {
		return nil, grpc.Errorf(codes.Unavailable, "channel closed")
	}
	return &ControlResponse{
		Ok:  rep.OK,
		Msg: rep.Msg,
	}, nil
}
