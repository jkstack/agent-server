package api

import (
	"fmt"
	"server/internal/agent"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/utils"
)

type GContext struct {
	*gin.Context
	agents *agent.Agents

	begin   time.Time
	reqID   string
	qryArgs any
	reqBody any

	muLogs sync.Mutex
	logs   []logItem
}

func GetG(g *gin.Context) *GContext {
	return g.MustGet(KeyGContext).(*GContext)
}

var number uint64

const defaultUID = "ffffffff"

func New(g *gin.Context, agents *agent.Agents) *GContext {
	next := atomic.AddUint64(&number, 1)
	uid, err := utils.UUID(8, "0123456789abcdef")
	if err != nil {
		logging.Error("generate uid for request %d failed, reset to default", next)
		uid = defaultUID
	}
	return &GContext{
		Context: g,
		agents:  agents,
		begin:   time.Now(),
		reqID: fmt.Sprintf("%s-%08d-%s",
			time.Now().Format("20060102"), next%99999999, uid),
	}
}

func (ctx *GContext) GetAgents() *agent.Agents {
	return ctx.agents
}

func (ctx *GContext) ShouldBindQuery(obj any) error {
	err := ctx.Context.ShouldBindQuery(obj)
	ctx.qryArgs = obj
	return err
}

func (ctx *GContext) ShouldBindJson(obj any) error {
	err := ctx.Context.ShouldBindJSON(obj)
	ctx.reqBody = obj
	return err
}

func (ctx *GContext) PostFormArray(key string) []string {
	value := ctx.Context.PostFormArray(key)
	if len(value) == 1 && len(value[0]) == 0 {
		value = nil
	}
	return value
}
