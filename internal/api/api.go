package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate swag init -g api.go

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
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg,omitempty" example:"错误内容"`
}

func OK(g *gin.Context, payload any) {
	g.JSON(http.StatusOK, Success{
		Code:    http.StatusOK,
		Payload: payload,
	})
}

func ERR(g *gin.Context, code int, msg string) {
	g.JSON(http.StatusOK, Failure{
		Code: code,
		Msg:  msg,
	})
}

func HttpError(g *gin.Context, code int, msg string) {
	g.String(code, msg)
}
