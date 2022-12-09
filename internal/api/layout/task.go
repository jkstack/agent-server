package layout

import (
	"sync"
	"time"
)

type statusCode int

const (
	// StatusWaiting waiting status
	StatusWaiting statusCode = iota
	// StatusRunning running status
	StatusRunning
	// StatusDone done status
	StatusDone
)

func (code statusCode) String() string {
	switch code {
	case StatusWaiting:
		return "waiting"
	case StatusRunning:
		return "running"
	case StatusDone:
		return "done"
	default:
		return "unknown"
	}
}

type task struct {
	sync.RWMutex
	ID     string
	Begin  time.Time
	End    time.Time
	IDS    []string // agent id
	Groups []string // groups
	// runtime
	Index      int
	Done       bool
	NodeErrs   map[string]error
	NodeStatus map[string]statusCode
	NodeBegin  map[string]time.Time
	NodeEnd    map[string]time.Time
}

func newTask(taskID string, ids, groups []string) *task {
	status := make(map[string]statusCode)
	for _, id := range ids {
		status[id] = StatusWaiting
	}
	return &task{
		ID:         taskID,
		Begin:      time.Now(),
		IDS:        ids,
		Groups:     groups,
		NodeStatus: status,
		NodeBegin:  make(map[string]time.Time),
		NodeEnd:    make(map[string]time.Time),
	}
}

func (t *task) OnRunning(node string) {
	t.Lock()
	t.NodeStatus[node] = StatusRunning
	t.NodeBegin[node] = time.Now()
	t.Unlock()
}

func (t *task) OnErr(node string, err error) {
	t.Lock()
	defer t.Unlock()
	t.NodeErrs[node] = err
	t.NodeStatus[node] = StatusDone
	t.NodeEnd[node] = time.Now()
}

func (t *task) OnDone(node string) {
	t.Lock()
	t.NodeStatus[node] = StatusDone
	t.NodeEnd[node] = time.Now()
	t.Unlock()
}
