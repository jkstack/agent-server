package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jkstack/jkframe/logging"
)

// OK response api caller ok
func (ctx *GContext) OK(payload any) {
	ctx.JSON(http.StatusOK, Success{
		Code:          http.StatusOK,
		Payload:       payload,
		ExecutionTime: time.Since(ctx.begin).Milliseconds(),
	})
}

// ERR response api caller error
func (ctx *GContext) ERR(code int, msg string) {
	ctx.JSON(http.StatusOK, Failure{
		Code:          code,
		Msg:           msg,
		RequestID:     ctx.reqID,
		ExecutionTime: time.Since(ctx.begin).Milliseconds(),
	})
}

// ErrAndLog response api caller error and log arguments
func (ctx *GContext) ErrAndLog(code int, msg string) {
	ctx.JSON(http.StatusOK, Failure{
		Code:      code,
		Msg:       msg,
		RequestID: ctx.reqID,
	})
	format := "REQUEST ERROR --- [%s] ---> code: %d, msg: %s"
	args := []interface{}{ctx.reqID, code, msg}
	format += "\n=> uri: %s"
	args = append(args, ctx.Request.RequestURI)
	if ctx.qryArgs != nil {
		format += "\n=> query: %s"
		args = append(args, fmt.Sprintf("%#v", ctx.qryArgs))
	}
	if ctx.reqBody != nil {
		format += "\n=> body: %s"
		args = append(args, fmt.Sprintf("%#v", ctx.reqBody))
	}
	format += "\n"
	logging.Error(format, args...)
}

// InvalidType invalid type error
func (ctx *GContext) InvalidType(want, got string) {
	panic(InvalidType{want: want, got: got})
}

// Timeout timeout error
func (ctx *GContext) Timeout() {
	panic(Timeout{})
}

// Notfound not found error
func (ctx *GContext) Notfound(what string) {
	panic(Notfound(what))
}

// BadParam bad param error
func (ctx *GContext) BadParam(param string) {
	panic(BadParam(param))
}

// MissingParam missing param error
func (ctx *GContext) MissingParam(what string) {
	panic(MissingParam(what))
}

// HTTPError response http error
func (ctx *GContext) HTTPError(code int, msg string) {
	ctx.String(code, msg)
}

// HTTPData response http data
func (ctx *GContext) HTTPData(data []byte) {
	ct := http.DetectContentType(data)
	ctx.Data(http.StatusOK, ct, data)
}
