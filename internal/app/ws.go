package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	runtime "github.com/jkstack/jkframe/utils"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

// handleWS agent连接处理接口
func (app *App) handleWS(g *gin.Context, mods []handler) {
	if !app.connectLimit.Allow() {
		api.HttpError(g, http.StatusServiceUnavailable, "rate limit")
		return
	}
	onConnect := make(chan *agent.Agent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case cli := <-onConnect:
			for _, mod := range mods {
				mod.OnConnect(cli)
			}
		case <-ctx.Done():
			return
		}
	}()
	cli := app.agent(g.Writer, g.Request, onConnect, cancel)
	go func() {
		for {
			select {
			case msg := <-cli.Unknown():
				if msg == nil {
					return
				}
				for _, mod := range mods {
					mod.OnMessage(cli, msg)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	if cli != nil {
		<-ctx.Done()
		app.stAgentCount.Dec()
		logging.Info("agent %s connection closed", cli.ID())
		for _, mod := range mods {
			mod.OnClose(cli.ID())
		}
	}
}

func (app *App) agent(w http.ResponseWriter, r *http.Request,
	onConnect chan *agent.Agent, cancel context.CancelFunc) *agent.Agent {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("upgrade websocket: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	come, err := app.waitCome(conn)
	if err != nil {
		logging.Error("wait come message(%s): %v", conn.RemoteAddr().String(), err)
		return nil
	}
	if app.handshake(conn, come) {
		app.stAgentCount.Inc()
		logging.Info("agent %s connection on, type=%s, os=%s, arch=%s, mac=%s",
			come.ID, come.Name, come.OS, come.Arch, come.MAC)
		cli := app.agents.New(conn, come, cancel)
		app.agents.Add(cli)
		onConnect <- cli
		return cli
	}
	return nil
}

func (app *App) waitCome(conn *websocket.Conn) (*anet.ComePayload, error) {
	conn.SetReadDeadline(time.Now().Add(time.Minute))
	var msg anet.Msg
	err := conn.ReadJSON(&msg)
	if err != nil {
		return nil, err
	}
	if msg.Type != anet.TypeCome {
		return nil, errors.New("invalid come message")
	}
	return msg.Come, nil
}

func (app *App) handshake(conn *websocket.Conn, come *anet.ComePayload) (ok bool) {
	var errMsg string
	defer func() {
		var rep anet.Msg
		rep.Type = anet.TypeHandshake
		rep.Important = true
		if ok {
			rep.Handshake = &anet.HandshakePayload{
				OK: true,
				ID: come.ID,
				// TODO: redirect
			}
		} else {
			rep.Handshake = &anet.HandshakePayload{
				OK:  false,
				Msg: errMsg,
			}
		}
		data, err := json.Marshal(rep)
		if err != nil {
			logging.Error("build handshake message: %v", err)
			return
		}
		conn.WriteMessage(websocket.TextMessage, data)
	}()
	app.connectLock.Lock()
	defer app.connectLock.Unlock()
	if len(come.ID) == 0 {
		id, err := runtime.UUID(16, "0123456789abcdef")
		if err != nil {
			errMsg = fmt.Sprintf("generate agent id: %v", err)
			logging.Error(errMsg)
			return false
		}
		come.ID = fmt.Sprintf("agent-%s-%s", time.Now().Format("20060102"), id)
	} else if app.agents.Get(come.ID) != nil {
		errMsg = "agent id conflict"
		return false
	}
	return true
}
