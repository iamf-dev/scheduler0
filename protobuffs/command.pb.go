// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: command.proto

package protobuffs

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

type Command_Type int32

const (
	Command_COMMAND_TYPE_DB_EXECUTE         Command_Type = 0
	Command_COMMAND_TYPE_JOB_QUEUE          Command_Type = 1
	Command_COMMAND_TYPE_JOB_EXECUTION_LOGS Command_Type = 2
	Command_COMMAND_TYPE_STOP_JOBS          Command_Type = 3
	Command_COMMAND_TYPE_RECOVER_JOBS       Command_Type = 4
)

// Enum value maps for Command_Type.
var (
	Command_Type_name = map[int32]string{
		0: "COMMAND_TYPE_DB_EXECUTE",
		1: "COMMAND_TYPE_JOB_QUEUE",
		2: "COMMAND_TYPE_JOB_EXECUTION_LOGS",
		3: "COMMAND_TYPE_STOP_JOBS",
		4: "COMMAND_TYPE_RECOVER_JOBS",
	}
	Command_Type_value = map[string]int32{
		"COMMAND_TYPE_DB_EXECUTE":         0,
		"COMMAND_TYPE_JOB_QUEUE":          1,
		"COMMAND_TYPE_JOB_EXECUTION_LOGS": 2,
		"COMMAND_TYPE_STOP_JOBS":          3,
		"COMMAND_TYPE_RECOVER_JOBS":       4,
	}
)

func (x Command_Type) Enum() *Command_Type {
	p := new(Command_Type)
	*p = x
	return p
}

func (x Command_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Command_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_command_proto_enumTypes[0].Descriptor()
}

func (Command_Type) Type() protoreflect.EnumType {
	return &file_command_proto_enumTypes[0]
}

func (x Command_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Command_Type.Descriptor instead.
func (Command_Type) EnumDescriptor() ([]byte, []int) {
	return file_command_proto_rawDescGZIP(), []int{0, 0}
}

type Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type         Command_Type `protobuf:"varint,1,opt,name=type,proto3,enum=command.Command_Type" json:"type,omitempty"`
	Sql          string       `protobuf:"bytes,2,opt,name=sql,proto3" json:"sql,omitempty"`
	Args         []byte       `protobuf:"bytes,3,opt,name=args,proto3" json:"args,omitempty"`
	Data         []byte       `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	ActionTarget string       `protobuf:"bytes,5,opt,name=action_target,json=actionTarget,proto3" json:"action_target,omitempty"`
}

func (x *Command) Reset() {
	*x = Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_command_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_command_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_command_proto_rawDescGZIP(), []int{0}
}

func (x *Command) GetType() Command_Type {
	if x != nil {
		return x.Type
	}
	return Command_COMMAND_TYPE_DB_EXECUTE
}

func (x *Command) GetSql() string {
	if x != nil {
		return x.Sql
	}
	return ""
}

func (x *Command) GetArgs() []byte {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *Command) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Command) GetActionTarget() string {
	if x != nil {
		return x.ActionTarget
	}
	return ""
}

var File_command_proto protoreflect.FileDescriptor

var file_command_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0xb5, 0x02, 0x0a, 0x07, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x12, 0x29, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x71, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x71,
	0x6c, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x61, 0x72, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x22, 0x9f,
	0x01, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x17, 0x43, 0x4f, 0x4d, 0x4d, 0x41,
	0x4e, 0x44, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x42, 0x5f, 0x45, 0x58, 0x45, 0x43, 0x55,
	0x54, 0x45, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x4a, 0x4f, 0x42, 0x5f, 0x51, 0x55, 0x45, 0x55, 0x45, 0x10, 0x01,
	0x12, 0x23, 0x0a, 0x1f, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x4a, 0x4f, 0x42, 0x5f, 0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x4c,
	0x4f, 0x47, 0x53, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x4f, 0x50, 0x5f, 0x4a, 0x4f, 0x42, 0x53, 0x10,
	0x03, 0x12, 0x1d, 0x0a, 0x19, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x52, 0x45, 0x43, 0x4f, 0x56, 0x45, 0x52, 0x5f, 0x4a, 0x4f, 0x42, 0x53, 0x10, 0x04,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_command_proto_rawDescOnce sync.Once
	file_command_proto_rawDescData = file_command_proto_rawDesc
)

func file_command_proto_rawDescGZIP() []byte {
	file_command_proto_rawDescOnce.Do(func() {
		file_command_proto_rawDescData = protoimpl.X.CompressGZIP(file_command_proto_rawDescData)
	})
	return file_command_proto_rawDescData
}

var file_command_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_command_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_command_proto_goTypes = []interface{}{
	(Command_Type)(0), // 0: command.Command.Type
	(*Command)(nil),   // 1: command.Command
}
var file_command_proto_depIdxs = []int32{
	0, // 0: command.Command.type:type_name -> command.Command.Type
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_command_proto_init() }
func file_command_proto_init() {
	if File_command_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_command_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Command); i {
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
			RawDescriptor: file_command_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_command_proto_goTypes,
		DependencyIndexes: file_command_proto_depIdxs,
		EnumInfos:         file_command_proto_enumTypes,
		MessageInfos:      file_command_proto_msgTypes,
	}.Build()
	File_command_proto = out.File
	file_command_proto_rawDesc = nil
	file_command_proto_goTypes = nil
	file_command_proto_depIdxs = nil
}
