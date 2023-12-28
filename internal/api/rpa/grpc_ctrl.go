package rpa

import (
	"context"

	"github.com/jkstack/anet"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Control send control packet to rpa agent
func (svr *Server) Control(ctx context.Context, args *ControlArgs) (*ControlResponse, error) {
	agentID := args.GetId()
	agent := svr.agents.Get(agentID)
	if agent == nil {
		return nil, status.Errorf(codes.NotFound, "agent not found")
	}
	svr.RLock()
	taskID := svr.jobs[agentID]
	svr.RUnlock()
	if len(taskID) == 0 {
		return nil, status.Errorf(codes.NotFound, "task not found")
	}
	var st int
	switch args.GetSt() {
	case ControlArgs_Pause:
		st = anet.RPAPause
	case ControlArgs_Stop:
		st = anet.RPAStop
	case ControlArgs_Resume:
		st = anet.RPAContinue
	}
	err := agent.SendRpaCtrl(taskID, st)
	if err != nil {
		return nil, err
	}
	svr.RLock()
	ch := svr.ctrlRep[taskID]
	svr.RUnlock()
	if ch == nil {
		return nil, status.Errorf(codes.Unavailable, "channel not found")
	}
	rep := <-ch
	if rep == nil {
		return nil, status.Errorf(codes.Unavailable, "channel closed")
	}
	return &ControlResponse{
		Ok:  rep.OK,
		Msg: rep.Msg,
	}, nil
}
