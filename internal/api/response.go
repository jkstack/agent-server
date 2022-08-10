package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OK(g *gin.Context, payload any) {
	g.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"payload": payload,
	})
}

func ERR(g *gin.Context, code int, msg string) {
	g.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func HttpError(g *gin.Context, code int, msg string) {
	g.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
	})
}
