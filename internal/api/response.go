package api

import (
	"net/http"
)

// OK response api caller ok
func (ctx *GContext) OK(payload any) {
	ctx.JSON(http.StatusOK, Success{
		Code:    http.StatusOK,
		Payload: payload,
	})
}

// ERR response api caller error
func (ctx *GContext) ERR(code int, msg string) {
	ctx.JSON(http.StatusOK, Failure{
		Code: code,
		Msg:  msg,
	})
}

// NotFound not found error
func (ctx *GContext) NotFound(what string) {
	panic(Notfound(what))
}

// InvalidType invalid type error
func (ctx *GContext) InvalidType(want, got string) {
	panic(InvalidType{want: want, got: got})
}

// Timeout timeout error
func (ctx *GContext) Timeout() {
	panic(Timeout{})
}

// HttpError response http error
func (ctx *GContext) HttpError(code int, msg string) {
	ctx.String(code, msg)
}
