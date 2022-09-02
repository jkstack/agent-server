package metrics

import (
	"encoding/json"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func percent(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) * 100. / float64(b)
}

func (h *Handler) saveStaticData(agentID string, data *anet.HMStaticPayload) {
	var static StaticData
	static.AgentId = agentID
	static.Time = timestamppb.New(data.Time)
	static.HostName = data.Host.Name
	static.Uptime = uint64(data.Host.UpTime.Seconds())
	static.OsName = data.OS.Name
	static.PlatformName = data.OS.PlatformName
	static.PlatformVersion = data.OS.PlatformVersion
	static.Install = timestamppb.New(data.OS.Install)
	static.Startup = timestamppb.New(data.OS.Startup)
	static.KernelVersion = data.Kernel.Version
	static.Arch = data.Kernel.Arch
	static.PhysicalCpu = uint64(data.CPU.Physical)
	static.LogicalCpu = uint64(data.CPU.Logical)
	for _, core := range data.CPU.Cores {
		static.Cores = append(static.Cores, &StaticCore{
			Processor: uint32(core.Processor),
			Model:     core.Model,
			Core:      uint32(core.Core),
			Cores:     uint32(core.Cores),
			Physical:  uint32(core.Physical),
			Mhz:       float32(core.Mhz),
		})
	}
	for _, disk := range data.Disks {
		var t StaticDiskDiskType
		switch strings.ToLower(disk.Type) {
		case "hdd":
			t = StaticDisk_hdd
		case "fdd":
			t = StaticDisk_fdd
		case "odd":
			t = StaticDisk_odd
		default:
			t = StaticDisk_unknown
		}
		static.Disks = append(static.Disks, &StaticDisk{
			Model:      disk.Model,
			Total:      disk.Total,
			Type:       t,
			Partitions: disk.Partitions,
		})
	}
	for _, part := range data.Partitions {
		static.Partitions = append(static.Partitions, &StaticPartition{
			Mount:   part.Name,
			Type:    part.FSType,
			Options: part.Opts,
			Total:   part.Total,
			Inodes:  part.INodes,
		})
	}
	static.Gateway = data.GateWay
	for _, intf := range data.Interface {
		static.Interfaces = append(static.Interfaces, &StaticInterface{
			Index: uint64(intf.Index),
			Name:  intf.Name,
			Mtu:   uint32(intf.Mtu),
			Flags: intf.Flags,
			Addrs: intf.Address,
			Mac:   intf.Mac,
		})
	}
	for _, user := range data.User {
		static.Users = append(static.Users, &StaticUser{
			Name: user.Name,
			Id:   user.ID,
			Gid:  user.GID,
		})
	}
	jBuf, _ := json.Marshal(data)
	pBuf, _ := proto.Marshal(&static)
	logging.Info("static data from [%s], json=%s, proto=%s, saved=%.02f%%",
		agentID,
		humanize.IBytes(uint64(len(jBuf))),
		humanize.IBytes(uint64(len(pBuf))),
		percent(len(jBuf)-len(pBuf), len(jBuf)))
}

func (h *Handler) saveDynamicData(agentID string, data *anet.HMDynamicRep) {
	var dynamic DynamicData
	dynamic.AgentId = agentID
	dynamic.Begin = timestamppb.New(data.Begin)
	dynamic.End = timestamppb.New(data.End)
	dynamic.HasUsage = data.Usage != nil
	if dynamic.HasUsage {
		dynamic.Usage = new(DynamicUsage)
		dynamic.Usage.CpuUsage = float32(data.Usage.Cpu.Usage)
		dynamic.Usage.MemoryUsed = data.Usage.Memory.Used
		dynamic.Usage.MemoryFree = data.Usage.Memory.Free
		dynamic.Usage.MemoryAvailable = data.Usage.Memory.Available
		dynamic.Usage.MemoryUsage = float32(data.Usage.Memory.Usage)
		dynamic.Usage.SwapUsed = data.Usage.Memory.SwapUsed
		dynamic.Usage.SwapFree = data.Usage.Memory.SwapFree
		for _, parts := range data.Usage.Partitions {
			dynamic.Usage.Partitions = append(dynamic.Usage.Partitions, &DynamicPartition{
				Mount:      parts.Name,
				Used:       parts.Used,
				Free:       parts.Free,
				Usage:      float32(parts.Usage),
				InodeUsed:  parts.InodeUsed,
				InodeFree:  parts.InodeFree,
				InodeUsage: float32(parts.InodeUsage),
			})
		}
		for _, intf := range data.Usage.Interface {
			dynamic.Usage.Interfaces = append(dynamic.Usage.Interfaces, &DynamicInterface{
				Name:        intf.Name,
				BytesSent:   intf.BytesSent,
				BytesRecv:   intf.BytesRecv,
				PacketsSent: intf.PacketsSent,
				PacketsRecv: intf.PacketsRecv,
			})
		}
	}
	for _, process := range data.Process {
		dynamic.Processes = append(dynamic.Processes, &DynamicProcess{
			Id:          uint32(process.ID),
			ParentId:    uint32(process.ParentID),
			User:        process.User,
			CpuUsage:    float32(process.CpuUsage),
			Rss:         process.RssMemory,
			Vms:         process.VirtualMemory,
			Swap:        process.SwapMemory,
			MemoryUsage: float32(process.MemoryUsage),
			Cmd:         process.Cmd,
			Listen:      process.Listen,
			Connections: uint64(process.Connections),
		})
	}
	for _, conn := range data.Connections {
		var t DynamicConnectionConnectionType
		switch strings.ToLower(conn.Type) {
		case "tcp4":
			t = DynamicConnection_tcp4
		case "tcp6":
			t = DynamicConnection_tcp6
		case "udp4":
			t = DynamicConnection_udp4
		case "udp6":
			t = DynamicConnection_udp6
		case "unix":
			t = DynamicConnection_unix
		case "file":
			t = DynamicConnection_file
		default:
			t = DynamicConnection_unknown
		}
		dynamic.Connections = append(dynamic.Connections, &DynamicConnection{
			Fd:     conn.Fd,
			Pid:    uint32(conn.Pid),
			Type:   t,
			Local:  conn.Local,
			Remote: conn.Remote,
			Status: conn.Status,
		})
	}
	jBuf, _ := json.Marshal(data)
	pBuf, _ := proto.Marshal(&dynamic)
	logging.Info("dyanmic data from [%s], json=%s, proto=%s, saved=%.02f%%",
		agentID,
		humanize.IBytes(uint64(len(jBuf))),
		humanize.IBytes(uint64(len(pBuf))),
		percent(len(jBuf)-len(pBuf), len(jBuf)))
}
