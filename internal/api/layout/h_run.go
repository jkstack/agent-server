package layout

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/scriptengine"
	iutils "server/internal/utils"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/utils"
)

const (
	onErrExit = iota
	onErrContinue
)

var errNotfound = errors.New("agent not found")
var errInvalidType = errors.New("invalid agent type")

type context struct {
	agents *agent.Agents
	taskID string
	ids    []string
	args   scriptengine.Args
	task   *task
}

// run 批量执行脚本
//	@ID				/api/layout/run
//	@description	分组示例:
//	@description	当ids参数为exec-01,exec-02,exec-03，group参数为1,1,2时，exec-01和exec-02节点为第一分组并发执行，当执行失败时若onerror参数为exit则不会执行后续分组中的exec-03节点任务
//	@Summary		批量执行脚本
//	@Tags			layout
//	@Accept			mpfd
//	@Produce		json
//	@Param			ids		formData	[]string					true	"节点ID列表"
//	@Param			group	formData	[]int						true	"分组列表"
//	@Param			file	formData	file						true	"文件"
//	@Param			type	formData	string						true	"脚本类型"	enums(sh,bash,python,python3,bat,powershell,php,lua)
//	@Param			args	formData	[]string					false	"参数"
//	@Param			md5		formData	string						false	"md5校验码"
//	@Param			auth	formData	string						false	"提权方式，仅linux有效"	enums(,sudo,su)
//	@Param			user	formData	string						false	"运行身份，仅linux有效"
//	@Param			workdir	formData	string						false	"工作目录"
//	@Param			env		formData	[]string					false	"环境变量"
//	@Param			timeout	formData	int							false	"单个节点的超时时间"		default(60)
//	@Param			onerror	formData	string						false	"执行失败时的后续操作"	enums(exit,continue)	default(exit)
//	@Success		200		{object}	api.Success{payload=string}	"payload为任务ID"
//	@Router			/layout/run [post]
func (h *Handler) run(gin *gin.Context) {
	g := api.GetG(gin)

	ids := g.PostFormArray("ids")
	if len(ids) == 0 {
		g.MissingParam("ids")
	}
	group := g.PostFormArray("group")
	if len(group) == 0 {
		g.MissingParam("group")
	}
	if len(ids) != len(group) {
		g.BadParam("invalid ids count and group count")
	}
	var index []int
	groups := make(map[int][]string) // group index => id
	for gi, gp := range group {
		n, err := strconv.ParseInt(gp, 10, 64)
		if err != nil {
			g.BadParam("group: " + err.Error())
		}
		if arr, ok := groups[int(n)]; ok {
			groups[int(n)] = append(arr, ids[gi])
		} else {
			groups[int(n)] = []string{ids[gi]}
			index = append(index, int(n))
		}
	}
	sort.Ints(index)
	hdr, err := g.FormFile("file")
	if err != nil {
		g.BadParam("file:" + err.Error())
	}
	file, err := hdr.Open()
	if err != nil {
		g.BadParam("file:" + err.Error())
	}
	defer file.Close()
	t := g.PostForm("type")
	args := g.PostFormArray("args")
	md5sum := g.PostForm("md5")
	auth := g.PostForm("auth")
	switch auth {
	case "", "sudo", "su":
	default:
		auth = ""
	}
	u := g.PostForm("user")
	workDir := g.PostForm("workdir")
	env := g.PostFormArray("env")
	timeout := 60
	str := g.PostForm("timeout")
	if len(str) > 0 {
		n, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			g.BadParam("timeout")
		}
		timeout = int(n)
	}
	str = g.PostForm("onerror")
	onErr := onErrExit
	if str == "continue" {
		onErr = onErrContinue
	}

	var srcMD5 [md5.Size]byte
	if len(md5sum) > 0 {
		encData, err := hex.DecodeString(md5sum)
		if err != nil {
			g.BadParam("md5:" + err.Error())
		}

		copy(srcMD5[:], encData)
	}

	data, err := io.ReadAll(file)
	utils.Assert(err)
	dstMD5 := md5.Sum(data)
	if len(md5sum) > 0 {
		if !bytes.Equal(srcMD5[:], dstMD5[:]) {
			g.BadParam("md5:invalid checksum")
		}
	}

	taskID, err := iutils.TaskID()
	utils.Assert(err)

	var sArgs scriptengine.Args
	sArgs.Data = string(data)
	switch t {
	case "sh":
		sArgs.Type = scriptengine.TypeSh
	case "bash":
		sArgs.Type = scriptengine.TypeBash
	case "python":
		sArgs.Type = scriptengine.TypePython
	case "python3":
		sArgs.Type = scriptengine.TypePython3
	case "bat":
		sArgs.Type = scriptengine.TypeBat
	case "powershell":
		sArgs.Type = scriptengine.TypePowerShell
	case "php":
		sArgs.Type = scriptengine.TypePhp
	case "lua":
		sArgs.Type = scriptengine.TypeLua
	default:
		g.BadParam("unsupported type")
	}
	sArgs.Args = args
	sArgs.Auth = auth
	sArgs.User = u
	sArgs.WorkDir = workDir
	sArgs.Env = env
	sArgs.Timeout = timeout

	agents := g.GetAgents()
	// pre check agent exists
	for _, i := range index {
		for _, id := range groups[i] {
			cli := agents.Get(id)
			if cli == nil {
				g.Notfound("agent: " + id)
				return
			}
			if cli.Type() != agent.TypeExec {
				g.InvalidType(agent.TypeExec, cli.Type())
			}
		}
	}

	tk := newTask(taskID, ids, group)
	h.Lock()
	h.tasks[taskID] = tk
	h.Unlock()

	go func(agents *agent.Agents, t *task) {
		ctx := context{
			agents: agents,
			taskID: taskID,
			args:   sArgs,
			task:   t,
		}
		for _, i := range index {
			t.Index = i
			ctx.ids = groups[i]
			ok := h.batch(ctx)
			if !ok && onErr == onErrExit {
				break
			}
		}
		t.Done = true
		t.End = time.Now()
	}(agents, tk)

	g.OK(taskID)
}

func (h *Handler) batch(ctx context) bool {
	ok := true
	run := func(wg *sync.WaitGroup, id string) bool {
		defer wg.Done()
		defer ctx.task.OnDone(id)
		ctx.task.OnRunning(id)
		cli := ctx.agents.Get(id)
		if cli == nil {
			ok = false
			ctx.task.OnErr(id, errNotfound)
			return false
		}
		if cli.Type() != agent.TypeExec {
			ok = false
			ctx.task.OnErr(id, errInvalidType)
			return false
		}
		e := scriptengine.New(cli, ctx.args)
		e.SetDataHandleFunc(func(b []byte) error {
			return nil
		})
		e.Run()
		if e.Err != nil {
			ok = false
			ctx.task.OnErr(id, e.Err)
			return false
		}
		if e.Code != 0 {
			ok = false
			ctx.task.OnErr(id, errors.New("")) // TODO: error message
		}
		return true
	}
	var wg sync.WaitGroup
	wg.Add(len(ctx.ids))
	for _, id := range ctx.ids {
		go run(&wg, id)
	}
	wg.Wait()
	return ok
}
