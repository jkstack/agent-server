package api

import "net/http"

// OK response api caller ok
func (ctx *GContext) OK(payload any) {
	ctx.g.JSON(http.StatusOK, Success{
		Code:      http.StatusOK,
		Payload:   payload,
		RequestID: ctx.reqID,
	})
}

// ERR response api caller error
func (ctx *GContext) ERR(code int, msg string) {
	ctx.g.JSON(http.StatusOK, Failure{
		Code:      code,
		Msg:       msg,
		RequestID: ctx.reqID,
	})
}

// NotFound not found error
func (ctx *GContext) NotFound(what string) {
	panic(NotFound(what))
}

// HttpError response http error
func (ctx *GContext) HttpError(code int, msg string) {
	ctx.g.String(code, msg)
}
