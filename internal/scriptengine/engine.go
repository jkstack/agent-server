package scriptengine

import (
	"crypto/md5"
	"fmt"
	"server/internal/agent"
	"time"
)

type Type int
type Status int

const (
	TypeSh Type = iota
	TypeBash
	TypePython
	TypePython3
	TypeBat
	TypePowerShell
	TypePhp
	TypeLua
)

const (
	StatusStopWaiting Status = iota
	StatusUploading
	StatusRunning
	StatusDone
)

type Args struct {
	Data    string   // 脚本内容
	Type    Type     // 类型
	Args    []string // 参数
	Auth    string   // 提权方式
	User    string   // 运行身份
	WorkDir string   // 工作目录
	Env     []string // 环境变量
	Timeout int      // 超时时间
}

type Engine struct {
	cli      *agent.Agent
	args     Args
	fileName string
	status   Status
	Begin    time.Time
	End      time.Time
	Err      error
	dir      string
	Pid      int
	Code     int
	OnData   func([]byte) error
}

func New(cli *agent.Agent, args Args) *Engine {
	return &Engine{
		cli:      cli,
		args:     args,
		fileName: buildFileName(args),
		status:   StatusStopWaiting,
		OnData:   func([]byte) error { return nil },
	}
}

func (e *Engine) SetDataHandleFunc(fn func([]byte) error) {
	e.OnData = fn
}

func (e *Engine) Run() {
	e.Begin = time.Now()
	defer func() {
		if e.Err != nil {
			e.OnData([]byte(e.Err.Error()))
		}
		e.status = StatusDone
		e.End = time.Now()
	}()
	if !e.upload() {
		if e.Code == 0 {
			e.Code = codeUnexpectedError
		}
		return
	}
	if !e.run() {
		if e.Code == 0 {
			e.Code = codeUnexpectedError
		}
		return
	}
}

func buildFileName(args Args) string {
	enc := md5.Sum([]byte(args.Data))
	switch args.Type {
	case TypeSh, TypeBash:
		return fmt.Sprintf("%x.sh", enc)
	case TypePython, TypePython3:
		return fmt.Sprintf("%x.py", enc)
	case TypeBat:
		return fmt.Sprintf("%x.bat", enc)
	case TypePowerShell:
		return fmt.Sprintf("%x.ps1", enc)
	case TypePhp:
		return fmt.Sprintf("%x.php", enc)
	case TypeLua:
		return fmt.Sprintf("%x.lua", enc)
	default:
		return "unknown"
	}
}
