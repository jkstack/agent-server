package rpa

import (
	"server/internal/agent"
	"sync"

	"github.com/jkstack/anet"
)

// Server rpa server
type Server struct {
	sync.RWMutex
	UnimplementedRpaServer
	agents  *agent.Agents
	jobs    map[string]string                // agent id => task id
	ctrlRep map[string]chan *anet.RPACtrlRep // task id => response
}

// NewGRPC create rpa server
func NewGRPC(agents *agent.Agents) *Server {
	return &Server{
		agents:  agents,
		jobs:    make(map[string]string),
		ctrlRep: make(map[string]chan *anet.RPACtrlRep),
	}
}

// OnClose on agent close
func (svr *Server) OnClose(id string) {
	svr.Lock()
	defer svr.Unlock()
	if tid, ok := svr.jobs[id]; ok {
		delete(svr.jobs, id)
		delete(svr.ctrlRep, tid)
	}
}
