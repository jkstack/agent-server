package api

import "github.com/jkstack/jkframe/logging"

// LogInfo write info log
func (ctx *GContext) LogInfo(fmt string, a ...interface{}) {
	logging.Info("iii [%s] "+fmt, a...)
}

// LogWarning write warning log
func (ctx *GContext) LogWarning(fmt string, a ...interface{}) {
	logging.Warning("www [%s] "+fmt, a...)
}

// LogDebug write debug log
func (ctx *GContext) LogDebug(fmt string, a ...interface{}) {
	logging.Debug("ddd [%s] "+fmt, a...)
}

// LogError write error log
func (ctx *GContext) LogError(fmt string, a ...interface{}) {
	logging.Error("eee [%s] "+fmt, a...)
}
