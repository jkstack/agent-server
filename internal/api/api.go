package api

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/utils"
)

const RequestTimeout = 10 * time.Second

//go:generate swag init -g api.go -o ../../docs

// @title       agent-server
// @description jkagent management server.

// @contact.url  https://www.jkstack.com

// @BasePath /api

// Success response success
type Success struct {
	Code      int    `json:"code" example:"0"`
	Payload   any    `json:"payload,omitempty"`
	RequestID string `json:"reqid" example:"20220812-00000001-2bf6c4"`
}

// Failure response failure
type Failure struct {
	Code      int    `json:"code" example:"1"`
	Msg       string `json:"msg,omitempty" example:"错误内容"`
	RequestID string `json:"reqid" example:"20220812-00000001-2bf6c4"`
}

// GContext context for gin
type GContext struct {
	g     *gin.Context
	reqID string
}

var number uint64

const defaultUID = "ffffffff"

// New create gin context
func New(g *gin.Context) *GContext {
	next := atomic.AddUint64(&number, 1)
	uid, err := utils.UUID(8, "0123456789abcdef")
	if err != nil {
		logging.Error("generate uid for request %d failed, reset to default", next)
		uid = defaultUID
	}
	ctx := &GContext{
		g: g,
		reqID: fmt.Sprintf("%s-%08d-%s",
			time.Now().Format("20060102"), next, uid),
	}
	g.Set("X-GContext", ctx)
	return ctx
}

// ReqID get request id
func (ctx *GContext) ReqID() string {
	return ctx.reqID
}

// RemoteIP get remote ip
func (ctx *GContext) RemoteIP() string {
	return ctx.g.RemoteIP()
}

// Request get request info
func (ctx *GContext) Request() string {
	return ctx.g.Request.Method + " " +
		ctx.g.Request.RequestURI + " " +
		ctx.g.Request.Proto
}

// ContentLength get response length
func (ctx *GContext) ContentLength() int {
	n := ctx.g.Writer.Size()
	if n < 0 {
		n = 0
	}
	return n
}
