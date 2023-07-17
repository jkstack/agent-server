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

// New create rpa server
func NewGRPC(agents *agent.Agents) *Server {
	return &Server{
		agents:  agents,
		jobs:    make(map[string]string),
		ctrlRep: make(map[string]chan *anet.RPACtrlRep),
	}
}
