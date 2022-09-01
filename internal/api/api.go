package api

import (
	"time"
)

const RequestTimeout = 10 * time.Second

//go:generate swag init -g api.go -o ../../docs

// @title       agent-server
// @version     TODO
// @description jkagent management server.

// @contact.url  https://www.jkstack.com

// @BasePath /api

// Success response success
type Success struct {
	Code    int `json:"code" example:"200" validate:"required"`
	Payload any `json:"payload,omitempty"`
}

// Failure response failure
type Failure struct {
	Code      int    `json:"code" example:"1" validate:"required"`
	Msg       string `json:"msg,omitempty" example:"错误内容" validate:"required"`
	RequestID string `json:"reqid" example:"20220812-00000001-2bf6c4" validate:"required"`
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
