package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OK response api caller ok
func OK(g *gin.Context, payload any) {
	g.JSON(http.StatusOK, Success{
		Code:    http.StatusOK,
		Payload: payload,
	})
}

// ERR response api caller error
func ERR(g *gin.Context, code int, msg string) {
	g.JSON(http.StatusOK, Failure{
		Code: code,
		Msg:  msg,
	})
}

// PanicNotFound not found error
func PanicNotFound(what string) {
	panic(NotFound(what))
}

// PanicInvalidType invalid type error
func PanicInvalidType(want, got string) {
	panic(InvalidType{want: want, got: got})
}

// PanicTimeout timeout error
func PanicTimeout() {
	panic(Timeout{})
}

// HttpError response http error
func HttpError(g *gin.Context, code int, msg string) {
	g.String(code, msg)
}
