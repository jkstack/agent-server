package script

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"path/filepath"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/scriptengine"
	iutils "server/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/cache/l2cache"
	"github.com/jkstack/jkframe/utils"
)

type runPayload struct {
	TaskID string `json:"id" example:"20220922-00001-4bc99720760771f6" validate:"required"` // 任务ID
	Begin  int64  `json:"begin" example:"1663816771" validate:"required"`                   // 开始时间戳
	End    int64  `json:"end,omitempty" example:"1663816771"`                               // 任务结束时间
	Code   int    `json:"code" example:"0" validate:"required"`                             // 任务的返回状态
	Data   string `json:"data,omitempty"`                                                   // 返回内容（base64编码）
}

// run 执行脚本
//	@ID			/api/script/run
//	@Summary	执行脚本
//	@Tags		script
//	@Accept		mpfd
//	@Produce	json
//	@Param		id		path		string		true	"节点ID"
//	@Param		file	formData	file		true	"文件"
//	@Param		type	formData	string		true	"脚本类型"	enums(sh,bash,python,python3,bat,powershell,php,lua)
//	@Param		args	formData	[]string	false	"参数"
//	@Param		md5		formData	string		false	"md5校验码"
//	@Param		auth	formData	string		false	"提权方式，仅linux有效"	enums(,sudo,su)
//	@Param		user	formData	string		false	"运行身份，仅linux有效"
//	@Param		workdir	formData	string		false	"工作目录"
//	@Param		env		formData	[]string	false	"环境变量"
//	@Param		timeout	formData	int			false	"超时时间"	default(60)
//	@Success	200		{object}	api.Success{payload=runPayload}
//	@Router		/script/{id}/run [post]
func (h *Handler) run(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")

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

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.Notfound("agent")
		return
	}
	if cli.Type() != agent.TypeExec {
		g.InvalidType(agent.TypeExec, cli.Type())
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

	cache, err := l2cache.New(102400, filepath.Join(h.cfg.CacheDir, "script"))
	utils.Assert(err)
	defer cache.Close()

	e := scriptengine.New(cli, sArgs)
	e.SetDataHandleFunc(func(data []byte) error {
		_, err = cache.Write(data)
		return err
	})
	e.Run()

	var ret runPayload
	ret.TaskID = taskID
	ret.Begin = e.Begin.Unix()
	ret.End = e.End.Unix()
	ret.Code = e.Code
	data, err = io.ReadAll(cache)
	utils.Assert(err)
	ret.Data = base64.StdEncoding.EncodeToString(data)

	g.OK(ret)
}
