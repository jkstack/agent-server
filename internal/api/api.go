package api

import (
	"server/internal/agent"
	"time"

	"github.com/gin-gonic/gin"
)

const RequestTimeout = 10 * time.Second

//go:generate swag init -g api.go -o ../../docs

// @title       agent-server
// @description jkagent management server.

// @contact.url  https://www.jkstack.com

// @BasePath /api

// Success response success
type Success struct {
	Code    int `json:"code" example:"0"`
	Payload any `json:"payload,omitempty"`
}

// Failure response failure
type Failure struct {
	Code      int    `json:"code" example:"1"`
	Msg       string `json:"msg,omitempty" example:"错误内容"`
	RequestID string `json:"reqid,omitempty" example:"20220812-00000001-2bf6c4"`
}

// GetRequestID get current request id
func GetRequestID(g *gin.Context) string {
	if id, ok := g.Get(KeyRequestID); ok {
		return id.(string)
	}
	return "00000000"
}

type Route struct {
	Method string
	Uri    string
}

func MakeRoute(method, uri string) Route {
	return Route{
		Method: method,
		Uri:    uri,
	}
}

func GetAgents(g *gin.Context) *agent.Agents {
	agents, _ := g.Get(KeyAgents)
	return agents.(*agent.Agents)
}
