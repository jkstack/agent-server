package api

import (
	"strings"

	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/utils"
)

type logLevel int

const (
	logDebug logLevel = iota
	logInfo
	logWarning
	logError
)

type logItem struct {
	level  logLevel
	format string
	data   []interface{}
	trace  []string
}

func (ctx *GContext) log(level logLevel, format string, args ...interface{}) {
	item := logItem{
		level:  level,
		format: format,
		data:   args,
	}
	if level == logError {
		item.trace = utils.Trace("  ")
	}
	ctx.muLogs.Lock()
	ctx.logs = append(ctx.logs, item)
	ctx.muLogs.Unlock()
}

// LogDebug write debug log
func (ctx *GContext) LogDebug(format string, args ...interface{}) {
	ctx.log(logDebug, format, args...)
}

// LogInfo write info log
func (ctx *GContext) LogInfo(format string, args ...interface{}) {
	ctx.log(logInfo, format, args...)
}

// LogWarning write warning log
func (ctx *GContext) LogWarning(format string, args ...interface{}) {
	ctx.log(logWarning, format, args...)
}

// LogError write error log
func (ctx *GContext) LogError(format string, args ...interface{}) {
	ctx.log(logError, format, args...)
}

// ShowLog print trace log
func (ctx *GContext) ShowLog() {
	for _, item := range ctx.logs {
		format := "REQUEST --- [%s] ---> " + item.format
		args := append([]interface{}{ctx.reqID}, item.data...)
		switch item.level {
		case logDebug:
			logging.Debug(format, args...)
		case logInfo:
			logging.Info(format, args...)
		case logWarning:
			logging.Warning(format, args...)
		case logError:
			logging.Printf("[ERROR]"+format+"\n%s",
				append(args, strings.Join(item.trace, "\n"))...)
		}
	}
}
