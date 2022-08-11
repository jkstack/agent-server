package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/utils"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
)

const channelBuffer = 10000

type Agent struct {
	sync.RWMutex
	t      string
	parent *Agents
	info   anet.ComePayload
	remote *websocket.Conn
	// runtime
	chRead   chan *anet.Msg
	chWrite  chan *anet.Msg
	taskRead map[string]chan *anet.Msg
}

func (agent *Agent) Close() {
	logging.Info("agent %s connection closed", agent.info.ID)
	if agent.remote != nil {
		agent.remote.Close()
	}
	close(agent.chRead)
	close(agent.chWrite)
}

func (agent *Agent) remoteAddr() string {
	return agent.remote.RemoteAddr().String()
}

func (agent *Agent) read(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		utils.Recover(fmt.Sprintf("agent.read(%s)", agent.remoteAddr()))
		cancel()
	}()
	agent.remote.SetReadDeadline(time.Time{})
	send := func(taskID string, ch chan *anet.Msg, msg *anet.Msg) {
		defer func() {
			if err := recover(); err != nil {
				logging.Error("write to channel %s timeout", taskID)
			}
		}()
		select {
		case <-ctx.Done():
			return
		case ch <- msg:
		case <-time.After(10 * time.Second):
			return
		}
	}
	for {
		_, data, err := agent.remote.ReadMessage()
		if err != nil {
			logging.Error("agent.read(%s): %v", agent.remoteAddr(), err)
			return
		}

		agent.parent.stInPackets.Inc()
		agent.parent.stInBytes.Add(float64(len(data)))

		var msg anet.Msg
		err = json.Unmarshal(data, &msg)
		if err != nil {
			logging.Error("agent.read.unmarshal(%s): %v", agent.remoteAddr(), err)
			return
		}

		ch := agent.chRead
		if len(msg.TaskID) > 0 {
			agent.RLock()
			ch = agent.taskRead[msg.TaskID]
			agent.RUnlock()
			if ch == nil {
				// logging.Error("response channel %s not found", msg.TaskID)
				continue
			}
		}
		send(msg.TaskID, ch, &msg)
	}
}

func (agent *Agent) write(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		utils.Recover(fmt.Sprintf("agent.write(%s)", agent.remoteAddr()))
		cancel()
	}()
	send := func(msg *anet.Msg, i int) bool {
		data, err := json.Marshal(msg)
		if err != nil {
			logging.Error("agent.write.marshal(%s) %d times: %v", agent.remoteAddr(), i, err)
			return false
		}
		err = agent.remote.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			logging.Error("agent.write(%s) %d times: %v", agent.remoteAddr(), i, err)
			return false
		}
		agent.parent.stOutPackets.Inc()
		agent.parent.stOutBytes.Add(float64(len(data)))
		return true
	}
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-agent.chWrite:
			if msg == nil {
				return
			}
			if msg.Important {
				for i := 0; i < 10; i++ {
					if send(msg, i+1) {
						break
					}
				}
				continue
			}
			send(msg, 1)
		}
	}
}

func (agent *Agent) ID() string {
	return agent.info.ID
}

func (agent *Agent) Type() string {
	return agent.info.Name
}

func (agent *Agent) Info() anet.ComePayload {
	return agent.info
}

func (agent *Agent) ChanRead(id string) <-chan *anet.Msg {
	agent.RLock()
	defer agent.RUnlock()
	return agent.taskRead[id]
}

func (agent *Agent) ChanClose(id string) {
	agent.Lock()
	defer agent.Unlock()
	if ch := agent.taskRead[id]; ch != nil {
		close(ch)
		delete(agent.taskRead, id)
	}
}

func (agent *Agent) Unknown() <-chan *anet.Msg {
	return agent.chRead
}
