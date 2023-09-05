package rpa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	uts "server/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/utils"
)

var callbackCli = &http.Client{
	Timeout: 10 * time.Second,
}

type inSelectorArgs struct {
	CallBack  string `json:"callback" example:"https://www.baidu.com"` // 回调地址
	RequestID string `json:"requestId" example:"20230728..."`          // 请求ID
}

// inSelector 进入元素选择器状态
//
//	@ID						/api/rpa/in_selector
//	@Description.markdown	in_selector.md
//	@Summary				进入元素选择器状态
//	@Tags					rpa
//	@Accept					json
//	@Produce				json
//	@Param					id		path		string			true	"节点ID"
//	@Param					args	body		inSelectorArgs	true	"需启动的任务列表"
//	@Success				200		{object}	api.Success
//	@Router					/rpa/{id}/in_selector [post]
func (h *Handler) inSelector(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	var args inSelectorArgs
	if err := g.ShouldBindJSON(&args); err != nil {
		g.BadParam(err.Error())
		return
	}

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.Notfound("agent")
		return
	}
	if cli.Type() != agent.TypeRPA {
		g.InvalidType(agent.TypeRPA, cli.Type())
	}

	taskID, err := cli.SendRpaInSelector()
	utils.Assert(err)
	defer cli.ChanClose(taskID)

	var msg *anet.Msg
	select {
	case msg = <-cli.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		g.Timeout()
	}

	switch {
	case msg.Type == anet.TypeError:
		g.ERR(http.StatusServiceUnavailable, msg.ErrorMsg)
		return
	case msg.Type != anet.TypeRPASelectorRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	if msg.RPASelectorRep.Code != 1 {
		g.ERR(http.StatusServiceUnavailable, msg.RPASelectorRep.Msg)
		return
	}

	go h.waitSelectorResult(cli, taskID, args.CallBack, args.RequestID)

	g.OK(nil)
}

func (h *Handler) waitSelectorResult(cli *agent.Agent, taskID, callback, requestID string) {
	defer uts.Recover("wait selector result")
	defer cli.ChanClose(taskID)
	begin := time.Now()
	timeout := time.After(time.Hour)
	select {
	case <-timeout:
		h.callback(cli, callback, begin, &anet.Msg{
			RPASelectorResult: &anet.RPASelectorResult{
				Code: 65535,
				Msg:  "execution timeout",
			},
		}, requestID)
	case msg := <-cli.ChanRead(taskID):
		logging.Info("rpa selector result: %d %s", msg.RPASelectorResult.Code, msg.RPASelectorResult.Msg)
		if msg.RPASelectorResult.Code == 1 {
			h.logImage(taskID, msg.RPASelectorResult.Image)
		}
		h.callback(cli, callback, begin, msg, requestID)
	}
}

func (h *Handler) callback(cli *agent.Agent, url string, begin time.Time, msg *anet.Msg, requestID string) {
	var buf bytes.Buffer
	var body struct {
		ID        string `json:"id"`
		Timestamp int64  `json:"ts"`
		Cost      int    `json:"cost"`
		Code      int    `json:"code"`
		Msg       string `json:"msg"`
		Content   string `json:"content,omitempty"`
		ImageURI  string `json:"imageUri,omitempty"`
		RequestID string `json:"requestId,omitempty"`
	}
	body.ID = cli.ID()
	body.Timestamp = time.Now().Unix()
	body.Cost = int(time.Since(begin).Seconds())
	body.Code = msg.RPASelectorResult.Code
	body.Msg = msg.RPASelectorResult.Msg
	if msg.RPASelectorResult.Code == 1 {
		body.Content = msg.RPASelectorResult.Content
		body.ImageURI = fmt.Sprintf("/api/rpa/files/%s", msg.TaskID)
	}
	body.RequestID = requestID
	utils.Assert(json.NewEncoder(&buf).Encode(body))
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	utils.Assert(err)
	req.Header.Set("Content-Type", "application/json")
	rep, err := callbackCli.Do(req)
	utils.Assert(err)
	defer rep.Body.Close()
	if rep.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(rep.Body)
		logging.Error("callback failed: %s", string(data))
	}
}

func (h *Handler) logImage(taskID, image string) {
	data, err := base64.StdEncoding.DecodeString(image)
	utils.Assert(err)
	h.newCache(taskID, data)
}
