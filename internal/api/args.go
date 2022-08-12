package api

import (
	"bytes"
	"io"
)

// ShouldBindQuery ShouldBindQuery of gin
func (ctx *GContext) ShouldBindQuery(obj any) error {
	return ctx.g.ShouldBindQuery(obj)
}

// Param Param of gin
func (ctx *GContext) Param(key string) string {
	return ctx.g.Param(key)
}

// ShouldBindJSON ShouldBindJSON of gin
func (ctx *GContext) ShouldBindJSON(obj any) error {
	return ctx.g.ShouldBindJSON(obj)
}

// RequestBody get request body
func (ctx *GContext) RequestBody() string {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(ctx.g.Request.Body)
	if err != nil {
		panic(err)
	}
	ctx.g.Request.Body = io.NopCloser(&buf)
	return string(buf.String())
}
