package exec

import (
	"encoding/base64"
	"io"
	"server/internal/agent"
	"time"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/cache/l2cache"
	"github.com/jkstack/jkframe/logging"
)

const (
	codeUnexpectedError = -65534
	codeTimeout         = -65535
)

type task struct {
	remote *agent.Agent
	id     string
	pid    int
	clean  time.Time
	// runtime
	cache    *l2cache.Cache
	done     chan struct{}
	doneFlag bool
	code     int
	begin    time.Time
	end      time.Time
}

func newTask(remote *agent.Agent, id string, pid int, cacheDir string,
	timeout time.Duration, begin time.Time) (*task, error) {
	cache, err := l2cache.New(102400, cacheDir)
	if err != nil {
		return nil, err
	}
	t := &task{
		remote: remote,
		id:     id,
		clean:  time.Now().Add(timeout + time.Hour),
		done:   make(chan struct{}),
		cache:  cache,
		begin:  begin,
	}
	go t.recv(remote.ChanRead(id), time.After(timeout))
	return t, nil
}

func (t *task) close() {
	t.cache.Close()
	close(t.done)
	t.remote.ChanClose(t.id)
}

func (t *task) recv(ch <-chan *anet.Msg, done <-chan time.Time) {
	defer func() {
		if t.end.IsZero() {
			t.end = time.Now()
		}
		select {
		case t.done <- struct{}{}:
		default:
		}
		t.doneFlag = true
	}()
	for {
		select {
		case msg := <-ch:
			switch msg.Type {
			case anet.TypeExecData:
				data, err := base64.StdEncoding.DecodeString(msg.ExecData.Data)
				if err != nil {
					logging.Error("can not decode exec data for task %s: %v", t.id, err)
					return
				}
				_, err = t.cache.Write(data)
				if err != nil {
					logging.Error("can not write data for task %s: %v", t.id, err)
					return
				}
			case anet.TypeExecDone:
				t.code = msg.ExecDone.Code
				t.end = msg.ExecDone.Time
				return
			default:
				t.code = codeUnexpectedError
				logging.Error("invalid message type want %d or %d got %d",
					anet.TypeExecData, anet.TypeExecDone, msg.Type)
				return
			}
		case <-done:
			t.code = codeTimeout
			return
		}
	}
}

func (t *task) wait() {
	<-t.done
}

func (t *task) data() ([]byte, error) {
	return io.ReadAll(t.cache)
}
