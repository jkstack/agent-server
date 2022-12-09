package api

import (
	"time"
)

// RequestTimeout request timeout
const RequestTimeout = 10 * time.Second

// @title       agent-server
// @version     TODO
// @description jkagent management server.

// @contact.url  https://www.jkstack.com

// @BasePath /api

// Success response success
type Success struct {
	Code          int   `json:"code" example:"200" validate:"required"`  // 状态码
	Payload       any   `json:"payload,omitempty"`                       // 内容
	ExecutionTime int64 `json:"extime" example:"70" validate:"required"` // 耗时(毫秒)
}

// Failure response failure
type Failure struct {
	Code          int    `json:"code" example:"1" validate:"required"`                         // 状态码
	Msg           string `json:"msg,omitempty" example:"错误内容" validate:"required"`             // 错误内容
	RequestID     string `json:"reqid" example:"20220812-00000001-2bf6c4" validate:"required"` // 请求ID
	ExecutionTime int64  `json:"extime" example:"70" validate:"required"`                      // 耗时(毫秒)
}

// Route route object
type Route struct {
	Method string
	URI    string
}

// MakeRoute create route object
func MakeRoute(method, uri string) Route {
	return Route{
		Method: method,
		URI:    uri,
	}
}
