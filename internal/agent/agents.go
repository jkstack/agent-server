package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/stat"
)

type Agents struct {
	sync.RWMutex
	data         map[string]*Agent
	stInPackets  *stat.Counter
	stOutPackets *stat.Counter
	stInBytes    *stat.Counter
	stOutBytes   *stat.Counter
}

func NewAgents(stats *stat.Mgr) *Agents {
	agents := &Agents{
		data:         make(map[string]*Agent),
		stInPackets:  stats.NewCounter("in_packets"),
		stOutPackets: stats.NewCounter("out_packets"),
		stInBytes:    stats.NewCounter("in_bytes"),
		stOutBytes:   stats.NewCounter("out_bytes"),
	}
	go agents.print()
	return agents
}

// New new agent
func (agents *Agents) New(conn *websocket.Conn, come *anet.ComePayload) (*Agent, <-chan struct{}) {
	cli := &Agent{
		t:        come.Name,
		parent:   agents,
		info:     *come,
		remote:   conn,
		chRead:   make(chan *anet.Msg, channelBuffer),
		chWrite:  make(chan *anet.Msg, channelBuffer),
		taskRead: make(map[string]chan *anet.Msg, channelBuffer/10),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go cli.read(ctx, cancel, conn.RemoteAddr().String())
	go cli.write(ctx, cancel, conn.RemoteAddr().String())
	agents.Add(cli)
	return cli, ctx.Done()
}

// Add add agent
func (agents *Agents) Add(agent *Agent) {
	agents.Lock()
	defer agents.Unlock()
	if old := agents.data[agent.info.ID]; old != nil {
		old.Close()
	}
	agents.data[agent.info.ID] = agent
}

// Remove remote agent
func (agents *Agents) Remove(id string) {
	agents.Lock()
	defer agents.Unlock()
	if agent := agents.data[id]; agent != nil {
		agent.Close()
	}
	delete(agents.data, id)
}

// Get get agent by id
func (agents *Agents) Get(id string) *Agent {
	agents.RLock()
	defer agents.RUnlock()
	return agents.data[id]
}

// Range list agents
func (agents *Agents) Range(cb func(*Agent) bool) {
	agents.RLock()
	defer agents.RUnlock()
	for _, c := range agents.data {
		next := cb(c)
		if !next {
			return
		}
	}
}

// Size get agents count
func (agents *Agents) Size() int {
	return len(agents.data)
}

func (agents *Agents) print() {
	var logs []string
	agents.RLock()
	for _, cli := range agents.data {
		if len(cli.chWrite) > 0 || len(cli.chRead) > 0 {
			logs = append(logs, fmt.Sprintf("agent %s: write chan=%d, read chan=%d",
				cli.ID(), len(cli.chWrite), len(cli.chRead)))
		}
	}
	agents.RUnlock()
	if len(logs) > 0 {
		logging.Info(strings.Join(logs, "\n"))
	}
}

func (agents *Agents) Prefix(str string) []*Agent {
	var ret []*Agent
	agents.RLock()
	for id, cli := range agents.data {
		if strings.HasPrefix(id, str) {
			ret = append(ret, cli)
		}
	}
	agents.RUnlock()
	return ret
}

func (agents *Agents) Contains(str string) []*Agent {
	var ret []*Agent
	agents.RLock()
	for id, cli := range agents.data {
		if strings.Contains(id, str) {
			ret = append(ret, cli)
		}
	}
	agents.RUnlock()
	return ret
}

func (agents *Agents) All() []*Agent {
	var ret []*Agent
	agents.RLock()
	for _, cli := range agents.data {
		ret = append(ret, cli)
	}
	agents.RUnlock()
	return ret
}
