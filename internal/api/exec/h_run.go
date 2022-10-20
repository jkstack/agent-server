package exec

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"
	"server/internal/agent"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/utils"
)

type runArgs struct {
	Result  bool     `json:"result" example:"true" default:"true"`                      // 是否直接返回数据
	Cmd     string   `json:"cmd" example:"echo" validate:"required" binding:"required"` // 命令
	Args    []string `json:"args" example:"hello world"`                                // 参数
	Auth    string   `json:"auth" example:"sudo" enums:",sudo,su" `                     // 提权方式，仅linux有效
	User    string   `json:"user" example:"nobody"`                                     // 运行身份，仅linux有效
	WorkDir string   `json:"workdir" example:"/var/log"`                                // 工作目录
	Env     []string `json:"env" example:"DB_HOST=127.0.0.1,DB_NAME=jkstack"`           // 环境变量
	Timeout int      `json:"timeout" example:"3600" default:"60"`                       // 超时时间
}

type result struct {
	Running bool   `json:"running" example:"true" validate:"required"` // 任务是否还在执行中
	Code    int    `json:"code" example:"0" validate:"required"`       // 任务的返回状态
	End     int64  `json:"end,omitempty" example:"1663816771"`         // 任务结束时间
	Data    string `json:"data,omitempty"`                             // 返回内容（base64编码）
}

type runPayload struct {
	TaskID string  `json:"id" example:"20220922-00001-4bc99720760771f6" validate:"required"` // 任务ID
	Begin  int64   `json:"begin" example:"1663816771" validate:"required"`                   // 开始时间戳
	Pid    int     `json:"pid" example:"9655" validate:"required"`                           // 进程ID
	Result *result `json:"stauts,omitempty"`                                                 // 任务状态，仅当result参数为true时返回
}

// run 执行命令或脚本
// @ID /api/exec/run
// @Summary 执行命令或脚本
// @Tags exec
// @Accept  json
// @Produce json
// @Param   id   path string  true  "节点ID"
// @Param   args body runArgs true "需启动的任务列表"
// @Success 200  {object}     api.Success{payload=runPayload}
// @Router /exec/{id}/run [post]
func (h *Handler) run(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")
	args := runArgs{
		Result:  true,
		Timeout: 60,
	}
	if err := g.ShouldBindJson(&args); err != nil {
		g.BadParam(err.Error())
		return
	}

	switch args.Auth {
	case "", "sudo", "su":
	default:
		args.Auth = ""
	}

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.Notfound("agent")
		return
	}
	if cli.Type() != agent.TypeExec {
		g.InvalidType(agent.TypeExec, cli.Type())
	}

	taskID, err := cli.SendExecRun(
		args.Cmd, args.Args,
		args.Auth, args.User,
		args.WorkDir, args.Env,
		args.Timeout,
	)
	utils.Assert(err)

	var msg *anet.Msg
	select {
	case msg = <-cli.ChanRead(taskID):
	case <-time.After(api.RequestTimeout):
		g.Timeout()
	}

	switch {
	case msg.Type == anet.TypeError:
		g.ERR(http.StatusServiceUnavailable, msg.ErrorMsg)
		return
	case msg.Type != anet.TypeExecd:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	task, err := newTask(cli, taskID, msg.Execd.Pid, filepath.Join(h.cfg.CacheDir, "exec"),
		time.Duration(args.Timeout)*time.Second, msg.Execd.Time)
	if err != nil {
		logging.Error("can not create task %s: %v", taskID, err)
		g.ERR(http.StatusInternalServerError, "can not create task")
		return
	}

	if args.Result {
		defer task.close()
		task.wait()
		data, err := task.data()
		if err != nil {
			logging.Error("can not load data for task %s: %v", taskID, err)
			g.ERR(http.StatusInternalServerError, "can not load data")
			return
		}
		g.OK(runPayload{
			TaskID: taskID,
			Begin:  msg.Execd.Time.Unix(),
			Pid:    msg.Execd.Pid,
			Result: &result{
				Running: false,
				Code:    task.code,
				End:     task.end.Unix(),
				Data:    base64.StdEncoding.EncodeToString(data),
			},
		})
		return
	}

	h.getTasksOrCreate(id).add(task)

	g.OK(runPayload{
		TaskID: taskID,
		Begin:  msg.Execd.Time.Unix(),
		Pid:    msg.Execd.Pid,
	})
}
