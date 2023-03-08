package rpa

import (
	"server/internal/agent"
	"sync"

	"github.com/jkstack/anet"
)

type Server struct {
	sync.RWMutex
	UnimplementedRpaServer
	agents  *agent.Agents
	jobs    map[string]string                // agent id => task id
	ctrlRep map[string]chan *anet.RPACtrlRep // task id => response
}

func New(agents *agent.Agents) *Server {
	return &Server{
		agents:  agents,
		jobs:    make(map[string]string),
		ctrlRep: make(map[string]chan *anet.RPACtrlRep),
	}
}
