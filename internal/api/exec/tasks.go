package exec

import (
	"sync"
	"time"
)

type tasks struct {
	sync.RWMutex
	id   string        // agent id
	data map[int]*task // process id => task
}

func newTasks(id string) *tasks {
	ts := &tasks{
		id:   id,
		data: make(map[int]*task),
	}
	go ts.clean()
	return ts
}

func (ts *tasks) add(t *task) {
	ts.Lock()
	ts.data[t.pid] = t
	ts.Unlock()
}

func (ts *tasks) get(pid int) *task {
	ts.RLock()
	defer ts.RUnlock()
	return ts.data[pid]
}

func (ts *tasks) list(fn func(*task)) {
	tasks := make([]*task, 0, len(ts.data))
	ts.RLock()
	for _, t := range ts.data {
		tasks = append(tasks, t)
	}
	ts.RUnlock()
	for _, t := range tasks {
		fn(t)
	}
}

func (ts *tasks) clean() {
	fetch := func() []*task {
		ret := make([]*task, 0, len(ts.data))
		now := time.Now()
		ts.RLock()
		defer ts.RUnlock()
		for _, t := range ts.data {
			if now.After(t.clean) {
				ret = append(ret, t)
			}
		}
		return ret
	}
	for {
		for _, task := range fetch() {
			task.close()
			ts.Lock()
			delete(ts.data, task.pid)
			ts.Unlock()
		}
		time.Sleep(time.Minute)
	}
}
