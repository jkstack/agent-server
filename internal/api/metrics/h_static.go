package metrics

import (
	"fmt"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	runtime "github.com/jkstack/jkframe/utils"
)

type core struct {
	Processor int32   `json:"processor" example:"0"`                                     // 第几个核心
	Model     string  `json:"model" example:"Intel(R) Xeon(R) CPU E5-2620 v2 @ 2.10GHz"` // CPU型号
	Core      int32   `json:"core" example:"0"`                                          // 所在物理核上的编号
	Cores     int32   `json:"cores" example:"0"`                                         // 某块CPU上的编号
	Physical  int32   `json:"physical" example:"0"`                                      // 物理CPU编号
	Mhz       float64 `json:"mhz" example:"1600"`                                        // 核心频率
}

type disk struct {
	Model      string   `json:"model" example:"Samsung SSD 870 QVO 1TB"` // 品牌型号
	Total      uint64   `json:"total" example:"1000204886016"`           // 容量
	Type       string   `json:"type" example:"HDD" enums:"HDD,SSD"`      // 磁盘类型
	Partitions []string `json:"disks,omitempty" example:"/boot,/home"`   // 逻辑分区
}

type partition struct {
	Mount  string   `json:"mount" example:"/boot"`                    // linux为挂载路径如/run，windows为盘符如C:
	FSType string   `json:"fstype,omitempty" example:"NTFS"`          // 文件系统类型
	Opts   []string `json:"opts,omitempty" example:"rw,nosuid,nodev"` // 附加信息
	Total  uint64   `json:"total" example:"209666048"`                // 总容量
	INodes uint64   `json:"inodes" example:"4072701"`                 // inode数量
}

type intf struct {
	Index   int      `json:"index" example:"0"`                                        // 网卡下标
	Name    string   `json:"name" example:"lo"`                                        // 网卡名称
	Mtu     int      `json:"mtu" example:"1500"`                                       // 网卡mtu
	Flags   []string `json:"flags,omitempty" example:"UP,BROADCAST,RUNNING,MULTICAST"` // 网卡附加参数
	Mac     string   `json:"mac" example:"26:29:93:f2:84:22"`                          // 网卡mac地址
	Address []string `json:"addrs,omitempty" example:"192.168.1.100"`                  // 网卡上绑定的IP地址列表
}

type user struct {
	Name string `json:"name" example:"root"` // 用户名
	ID   string `json:"id" example:"0"`      // 用户ID
	GID  string `json:"gid" example:"0"`     // 用户组ID
}

type staticInfo struct {
	Time int64 `json:"time" example:"1661396019"` // 客户端时间
	Host struct {
		Name   string `json:"name" example:"DESKTOP-OQ8O3DR"` // 主机名
		UpTime int64  `json:"uptime" example:"1900"`          // 启动时长
	} `json:"host"`
	OS struct {
		Name            string `json:"name" example:"linux" enums:"linux,windows"`                             // 系统类型
		PlatformName    string `json:"platform_name" example:"centos" enums:"redhat,centos,debian,ubuntu,..."` // 系统名称
		PlatformVersion string `json:"platform_version" example:"7.7.1908"`                                    // 系统版本号
		Install         int64  `json:"install" example:"1661396019"`                                           // 系统安装时间
		Startup         int64  `json:"startup" example:"1661396019"`                                           // 系统启动时间
	} `json:"os"`
	Kernel struct {
		Version string `json:"version" example:"3.10.0-1062.el7.x86_64"`      // 内核版本
		Arch    string `json:"arch" example:"x86_64" enums:"x86_64,i386,..."` // 位数
	} `json:"kernel"`
	CPU struct {
		Physical int    `json:"physical" example:"2"` // 物理核心数
		Logical  int    `json:"logical" example:"4"`  // 逻辑核心数
		Cores    []core `json:"cores,omitempty"`      // 每个核心参数
	} `json:"cpu"`
	Memory struct {
		Physical uint64 `json:"physical" example:"33363566592"` // 物理内存大小
		Swap     uint64 `json:"swap" example:"8589930496"`      // swap内存大小
	} `json:"memory"`
	Disks      []disk      `json:"disks,omitempty"`               // 物理磁盘列表
	Partitions []partition `json:"partitions,omitempty"`          // 逻辑分区列表
	GateWay    string      `json:"gateway" example:"192.168.1.1"` // 网关地址
	Interface  []intf      `json:"interface,omitempty"`           // 网卡列表
	User       []user      `json:"user,omitempty"`                // 用户列表
}

// static 获取节点的静态数据
// @ID /api/metrics/static
// @Summary 获取节点的静态数据
// @Tags metrics
// @Produce json
// @Param   id   path string  true "节点ID"
// @Success 200  {object}     api.Success{payload=staticInfo}
// @Router /metrics/{id}/static [get]
func (h *Handler) static(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")

	agents := g.GetAgents()

	cli := agents.Get(id)
	if cli == nil {
		g.NotFound("agent")
		return
	}
	if cli.Type() != agent.TypeMetrics {
		g.InvalidType(agent.TypeMetrics, cli.Type())
	}

	taskID, err := cli.SendHMStaticReq()
	runtime.Assert(err)
	defer cli.ChanClose(id)

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
	case msg.Type != anet.TypeHMStaticRep:
		g.ERR(http.StatusInternalServerError, fmt.Sprintf("invalid message type: %d", msg.Type))
		return
	}

	g.OK(transStatic(msg.HMStatic))
}

func transStatic(payload *anet.HMStaticPayload) staticInfo {
	var ret staticInfo
	ret.Time = payload.Time.Unix()
	fillStaticHost(&ret, payload)
	fillStaticCpu(&ret, payload)
	fillStaticMemory(&ret, payload)
	fillStaticDisk(&ret, payload)
	fillStaticNetwork(&ret, payload)
	fillStaticUser(&ret, payload)
	return ret
}

func fillStaticHost(ret *staticInfo, payload *anet.HMStaticPayload) {
	ret.Host.Name = payload.Host.Name
	ret.Host.UpTime = int64(payload.Host.UpTime.Seconds())
	ret.OS.Name = payload.OS.Name
	ret.OS.PlatformName = payload.OS.PlatformName
	ret.OS.PlatformVersion = payload.OS.PlatformVersion
	ret.OS.Install = payload.OS.Install.Unix()
	ret.OS.Startup = payload.OS.Startup.Unix()
	ret.Kernel.Version = payload.Kernel.Version
	ret.Kernel.Arch = payload.Kernel.Arch
}

func fillStaticCpu(ret *staticInfo, payload *anet.HMStaticPayload) {
	ret.CPU.Physical = payload.CPU.Physical
	ret.CPU.Logical = payload.CPU.Logical
	for _, c := range payload.CPU.Cores {
		ret.CPU.Cores = append(ret.CPU.Cores, core{
			Processor: c.Processor,
			Model:     c.Model,
			Core:      c.Core,
			Cores:     c.Cores,
			Physical:  c.Physical,
			Mhz:       c.Mhz.Float(),
		})
	}
}

func fillStaticMemory(ret *staticInfo, payload *anet.HMStaticPayload) {
	ret.Memory.Physical = payload.Memory.Physical
	ret.Memory.Swap = payload.Memory.Swap
}

func fillStaticDisk(ret *staticInfo, payload *anet.HMStaticPayload) {
	for _, d := range payload.Disks {
		ret.Disks = append(ret.Disks, disk{
			Model:      d.Model,
			Total:      d.Total,
			Type:       d.Type,
			Partitions: d.Partitions,
		})
	}
	for _, p := range payload.Partitions {
		ret.Partitions = append(ret.Partitions, partition{
			Mount:  p.Name,
			FSType: p.FSType,
			Opts:   p.Opts,
			Total:  p.Total,
			INodes: p.INodes,
		})
	}
}

func fillStaticNetwork(ret *staticInfo, payload *anet.HMStaticPayload) {
	ret.GateWay = payload.GateWay
	for _, i := range payload.Interface {
		ret.Interface = append(ret.Interface, intf{
			Index:   i.Index,
			Name:    i.Name,
			Mtu:     i.Mtu,
			Flags:   i.Flags,
			Mac:     i.Mac,
			Address: i.Address,
		})
	}
}

func fillStaticUser(ret *staticInfo, payload *anet.HMStaticPayload) {
	for _, u := range payload.User {
		ret.User = append(ret.User, user{
			Name: u.Name,
			ID:   u.ID,
			GID:  u.GID,
		})
	}
}
