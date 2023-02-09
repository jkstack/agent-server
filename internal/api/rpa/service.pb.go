// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.2
// source: service.proto

package network

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ControlArgsStatus int32

const (
	ControlArgs_pause    ControlArgsStatus = 0 // 暂停
	ControlArgs_stop     ControlArgsStatus = 1 // 停止
	ControlArgs_continue ControlArgsStatus = 2 // 继续运行
)

// Enum value maps for ControlArgsStatus.
var (
	ControlArgsStatus_name = map[int32]string{
		0: "pause",
		1: "stop",
		2: "continue",
	}
	ControlArgsStatus_value = map[string]int32{
		"pause":    0,
		"stop":     1,
		"continue": 2,
	}
)

func (x ControlArgsStatus) Enum() *ControlArgsStatus {
	p := new(ControlArgsStatus)
	*p = x
	return p
}

func (x ControlArgsStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ControlArgsStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_service_proto_enumTypes[0].Descriptor()
}

func (ControlArgsStatus) Type() protoreflect.EnumType {
	return &file_service_proto_enumTypes[0]
}

func (x ControlArgsStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ControlArgsStatus.Descriptor instead.
func (ControlArgsStatus) EnumDescriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2, 0}
}

type RunArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url     string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`                         // zip包下载地址
	IsDebug bool   `protobuf:"varint,2,opt,name=is_debug,json=isDebug,proto3" json:"is_debug,omitempty"` // 是否调试运行
}

func (x *RunArgs) Reset() {
	*x = RunArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunArgs) ProtoMessage() {}

func (x *RunArgs) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunArgs.ProtoReflect.Descriptor instead.
func (*RunArgs) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

func (x *RunArgs) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *RunArgs) GetIsDebug() bool {
	if x != nil {
		return x.IsDebug
	}
	return false
}

type Log struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"` // 日志内容
}

func (x *Log) Reset() {
	*x = Log{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Log) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log) ProtoMessage() {}

func (x *Log) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log.ProtoReflect.Descriptor instead.
func (*Log) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

func (x *Log) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type ControlArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	St ControlArgsStatus `protobuf:"varint,1,opt,name=st,proto3,enum=rpa.ControlArgsStatus" json:"st,omitempty"`
}

func (x *ControlArgs) Reset() {
	*x = ControlArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ControlArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ControlArgs) ProtoMessage() {}

func (x *ControlArgs) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ControlArgs.ProtoReflect.Descriptor instead.
func (*ControlArgs) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2}
}

func (x *ControlArgs) GetSt() ControlArgsStatus {
	if x != nil {
		return x.St
	}
	return ControlArgs_pause
}

type ControlResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok  bool   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`  // 操作是否成功
	Msg string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"` // 失败原因
}

func (x *ControlResponse) Reset() {
	*x = ControlResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ControlResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ControlResponse) ProtoMessage() {}

func (x *ControlResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ControlResponse.ProtoReflect.Descriptor instead.
func (*ControlResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3}
}

func (x *ControlResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *ControlResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_service_proto protoreflect.FileDescriptor

var file_service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x72, 0x70, 0x61, 0x22, 0x37, 0x0a, 0x08, 0x72, 0x75, 0x6e, 0x5f, 0x61, 0x72, 0x67, 0x73,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x72, 0x6c, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x62, 0x75, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x44, 0x65, 0x62, 0x75, 0x67, 0x22, 0x19, 0x0a,
	0x03, 0x6c, 0x6f, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x65, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x5f, 0x61, 0x72, 0x67, 0x73, 0x12, 0x28, 0x0a, 0x02, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x72, 0x70, 0x61, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x5f, 0x61, 0x72, 0x67, 0x73, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x02,
	0x73, 0x74, 0x22, 0x2b, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x09, 0x0a, 0x05,
	0x70, 0x61, 0x75, 0x73, 0x65, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x73, 0x74, 0x6f, 0x70, 0x10,
	0x01, 0x12, 0x0c, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x10, 0x02, 0x22,
	0x34, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x02, 0x6f, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6d, 0x73, 0x67, 0x32, 0x5c, 0x0a, 0x03, 0x72, 0x70, 0x61, 0x12, 0x20, 0x0a, 0x03,
	0x72, 0x75, 0x6e, 0x12, 0x0d, 0x2e, 0x72, 0x70, 0x61, 0x2e, 0x72, 0x75, 0x6e, 0x5f, 0x61, 0x72,
	0x67, 0x73, 0x1a, 0x08, 0x2e, 0x72, 0x70, 0x61, 0x2e, 0x6c, 0x6f, 0x67, 0x30, 0x01, 0x12, 0x33,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x11, 0x2e, 0x72, 0x70, 0x61, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x61, 0x72, 0x67, 0x73, 0x1a, 0x15, 0x2e, 0x72,
	0x70, 0x61, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x3b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_service_proto_rawDescOnce sync.Once
	file_service_proto_rawDescData = file_service_proto_rawDesc
)

func file_service_proto_rawDescGZIP() []byte {
	file_service_proto_rawDescOnce.Do(func() {
		file_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_proto_rawDescData)
	})
	return file_service_proto_rawDescData
}

var file_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_service_proto_goTypes = []interface{}{
	(ControlArgsStatus)(0),  // 0: rpa.control_args.status
	(*RunArgs)(nil),         // 1: rpa.run_args
	(*Log)(nil),             // 2: rpa.log
	(*ControlArgs)(nil),     // 3: rpa.control_args
	(*ControlResponse)(nil), // 4: rpa.control_response
}
var file_service_proto_depIdxs = []int32{
	0, // 0: rpa.control_args.st:type_name -> rpa.control_args.status
	1, // 1: rpa.rpa.run:input_type -> rpa.run_args
	3, // 2: rpa.rpa.control:input_type -> rpa.control_args
	2, // 3: rpa.rpa.run:output_type -> rpa.log
	4, // 4: rpa.rpa.control:output_type -> rpa.control_response
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Log); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ControlArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ControlResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
		EnumInfos:         file_service_proto_enumTypes,
		MessageInfos:      file_service_proto_msgTypes,
	}.Build()
	File_service_proto = out.File
	file_service_proto_rawDesc = nil
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
