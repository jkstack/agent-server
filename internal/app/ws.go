package app

import (
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

var errIDConflict = errors.New("agent id conflict")

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

func dispatchMessage(done <-chan struct{}, cli *agent.Agent, mods []handler) {
	for {
		select {
		case msg := <-cli.Unknown():
			if msg == nil {
				return
			}
			for _, mod := range mods {
				mod.OnMessage(cli, msg)
			}
		case <-done:
			return
		}
	}
}

func (app *App) onConnect(id string, mods []handler) {
	for _, mod := range mods {
		mod.OnClose(id)
	}
}

// handleWS agent连接处理接口
func (app *App) handleWS(g *gin.Context, mods []handler) {
	if !app.connectLimit.Allow() {
		api.HttpError(g, http.StatusServiceUnavailable, "rate limit")
		return
	}
	conn, err := upgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		logging.Error("upgrade websocket: %v", err)
		api.HttpError(g, http.StatusInternalServerError, err.Error())
		return
	}
	defer conn.Close()
	come, err := app.waitCome(conn)
	if err != nil {
		logging.Error("wait come message(%s): %v", conn.RemoteAddr().String(), err)
		return
	}
	if ok, err := app.responseHandshake(conn, come); !ok {
		logging.Error("response handshake failed, agent_id=%s, src_ip=%s: %v",
			come.ID, conn.RemoteAddr().String(), err)
		return
	}

	app.stAgentCount.Inc()
	logging.Info("agent %s connection on, type=%s, os=%s, arch=%s, mac=%s",
		come.ID, come.Name, come.OS, come.Arch, come.MAC)

	cli, done := app.agents.New(conn, come)
	defer cli.Close()
	app.onConnect(cli.ID(), mods)
	go dispatchMessage(done, cli, mods)

	<-done

	app.stAgentCount.Dec()
	logging.Info("agent %s connection closed", cli.ID())
	for _, mod := range mods {
		mod.OnClose(cli.ID())
	}
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

func (app *App) responseHandshake(conn *websocket.Conn, come *anet.ComePayload) (ok bool, err error) {
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
			return false, err
		}
		come.ID = fmt.Sprintf("agent-%s-%s", time.Now().Format("20060102"), id)
	} else if app.agents.Get(come.ID) != nil {
		errMsg = "agent id conflict"
		return false, errIDConflict
	}
	return true, nil
}
