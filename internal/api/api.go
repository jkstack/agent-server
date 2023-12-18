package api

import (
	"time"
)

// RequestTimeout request timeout
const RequestTimeout = time.Minute

//	@title			agent-server
//	@version		TODO
//	@description	jkagent management server.

//	@contact.url	https://www.jkstack.com

//	@BasePath	/api

//	@tag.name			agents
//	@tag.description	agent相关操作接口
//	@tag.name			exec
//	@tag.description	exec-agent执行命令相关操作接口
//	@tag.name			file
//	@tag.description	exec-agent文件处理相关操作接口
//	@tag.name			foo
//	@tag.description	example-agent相关接口
//	@tag.name			info
//	@tag.description	服务器端信息相关接口
//	@tag.name			layout
//	@tag.description	编排运行相关接口
//	@tag.name			metrics
//	@tag.description	metrics-agent处理相关接口
//	@tag.name			script
//	@tag.description	脚本运行相关接口
//	@tag.name			ipmi
//	@tag.description	ipmi相关接口

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
