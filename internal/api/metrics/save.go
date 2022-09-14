package metrics

import (
	"strings"

	"github.com/Shopify/sarama"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) saveStaticData(agentID string, data *anet.HMStaticPayload) {
	var static StaticData
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
	sendKafka(h.cli, h.topic, &Data{
		AgentId: agentID,
		Payload: &Data_StaticData{
			StaticData: &static,
		},
	})
}

func (h *Handler) saveUsage(agentID string, data *anet.HMDynamicRep) {
	var dynamic DynamicData
	dynamic.Begin = timestamppb.New(data.Begin)
	dynamic.End = timestamppb.New(data.End)
	var usage DynamicUsage
	usage.CpuUsage = float32(data.Usage.Cpu.Usage)
	usage.MemoryUsed = data.Usage.Memory.Used
	usage.MemoryFree = data.Usage.Memory.Free
	usage.MemoryAvailable = data.Usage.Memory.Available
	usage.MemoryUsage = float32(data.Usage.Memory.Usage)
	usage.SwapUsed = data.Usage.Memory.SwapUsed
	usage.SwapFree = data.Usage.Memory.SwapFree
	for _, parts := range data.Usage.Partitions {
		usage.Partitions = append(usage.Partitions, &DynamicPartition{
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
		usage.Interfaces = append(usage.Interfaces, &DynamicInterface{
			Name:        intf.Name,
			BytesSent:   intf.BytesSent,
			BytesRecv:   intf.BytesRecv,
			PacketsSent: intf.PacketsSent,
			PacketsRecv: intf.PacketsRecv,
		})
	}
	dynamic.Payload = &DynamicData_Usage{Usage: &usage}
	sendKafka(h.cli, h.topic, &Data{
		AgentId: agentID,
		Payload: &Data_DynamicData{
			DynamicData: &dynamic,
		},
	})
}

func (h *Handler) saveProcess(agentID string, data *anet.HMDynamicRep) {
	var dynamic DynamicData
	dynamic.Begin = timestamppb.New(data.Begin)
	dynamic.End = timestamppb.New(data.End)
	var processes DynamicProcesses
	for _, process := range data.Process {
		processes.Data = append(processes.Data, &DynamicProcess{
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
	dynamic.Payload = &DynamicData_Processes{Processes: &processes}
	sendKafka(h.cli, h.topic, &Data{
		AgentId: agentID,
		Payload: &Data_DynamicData{
			DynamicData: &dynamic,
		},
	})
}

func (h *Handler) saveConnections(agentID string, data *anet.HMDynamicRep) {
	var dynamic DynamicData
	dynamic.Begin = timestamppb.New(data.Begin)
	dynamic.End = timestamppb.New(data.End)
	var connections DynamicConnections
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
		connections.Data = append(connections.Data, &DynamicConnection{
			Fd:     conn.Fd,
			Pid:    uint32(conn.Pid),
			Type:   t,
			Local:  conn.Local,
			Remote: conn.Remote,
			Status: conn.Status,
		})
	}
	dynamic.Payload = &DynamicData_Connections{Connections: &connections}
	sendKafka(h.cli, h.topic, &Data{
		AgentId: agentID,
		Payload: &Data_DynamicData{
			DynamicData: &dynamic,
		},
	})
}

func (h *Handler) saveSensorsTemperatures(agentID string, data *anet.HMDynamicRep) {
	var dynamic DynamicData
	dynamic.Begin = timestamppb.New(data.Begin)
	dynamic.End = timestamppb.New(data.End)
	var temps DynamicSensorsTemperatures
	for _, temp := range data.SensorsTemperatures {
		temps.Data = append(temps.Data, &DynamicSensorTemperature{
			Name: temp.Name,
			Temp: float32(temp.Temperature),
		})
	}
	dynamic.Payload = &DynamicData_Temps{Temps: &temps}
	sendKafka(h.cli, h.topic, &Data{
		AgentId: agentID,
		Payload: &Data_DynamicData{
			DynamicData: &dynamic,
		},
	})
}

func (h *Handler) saveDynamicData(agentID string, data *anet.HMDynamicRep) {
	if data.Usage != nil {
		h.saveUsage(agentID, data)
	} else if len(data.Process) > 0 {
		h.saveProcess(agentID, data)
	} else if len(data.Connections) > 0 {
		h.saveConnections(agentID, data)
	} else {
		h.saveSensorsTemperatures(agentID, data)
	}
}

func sendKafka(cli sarama.AsyncProducer, topic string, data *Data) {
	bytes, err := proto.Marshal(data)
	if err != nil {
		logging.Error("proto marshal for [%s]: %v", data.AgentId, err)
		return
	}
	if cli == nil {
		return
	}
	cli.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(data.AgentId),
		Value: sarama.ByteEncoder(bytes),
	}
}
