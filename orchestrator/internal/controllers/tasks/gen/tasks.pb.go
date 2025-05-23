// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: tasks.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Operator int32

const (
	Operator_ADDICTION      Operator = 0
	Operator_SUBTRACTION    Operator = 1
	Operator_MULTIPLICATION Operator = 2
	Operator_DIVISION       Operator = 3
)

// Enum value maps for Operator.
var (
	Operator_name = map[int32]string{
		0: "ADDICTION",
		1: "SUBTRACTION",
		2: "MULTIPLICATION",
		3: "DIVISION",
	}
	Operator_value = map[string]int32{
		"ADDICTION":      0,
		"SUBTRACTION":    1,
		"MULTIPLICATION": 2,
		"DIVISION":       3,
	}
)

func (x Operator) Enum() *Operator {
	p := new(Operator)
	*p = x
	return p
}

func (x Operator) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Operator) Descriptor() protoreflect.EnumDescriptor {
	return file_tasks_proto_enumTypes[0].Descriptor()
}

func (Operator) Type() protoreflect.EnumType {
	return &file_tasks_proto_enumTypes[0]
}

func (x Operator) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Operator.Descriptor instead.
func (Operator) EnumDescriptor() ([]byte, []int) {
	return file_tasks_proto_rawDescGZIP(), []int{0}
}

type UpdateTaskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Result        string                 `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateTaskRequest) Reset() {
	*x = UpdateTaskRequest{}
	mi := &file_tasks_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTaskRequest) ProtoMessage() {}

func (x *UpdateTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tasks_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTaskRequest.ProtoReflect.Descriptor instead.
func (*UpdateTaskRequest) Descriptor() ([]byte, []int) {
	return file_tasks_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateTaskRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateTaskRequest) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type TaskResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Arg1          string                 `protobuf:"bytes,2,opt,name=arg1,proto3" json:"arg1,omitempty"`
	Arg2          string                 `protobuf:"bytes,3,opt,name=arg2,proto3" json:"arg2,omitempty"`
	Time          int32                  `protobuf:"varint,4,opt,name=time,proto3" json:"time,omitempty"`
	Operator      Operator               `protobuf:"varint,5,opt,name=operator,proto3,enum=gen.Operator" json:"operator,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskResponse) Reset() {
	*x = TaskResponse{}
	mi := &file_tasks_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskResponse) ProtoMessage() {}

func (x *TaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tasks_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskResponse.ProtoReflect.Descriptor instead.
func (*TaskResponse) Descriptor() ([]byte, []int) {
	return file_tasks_proto_rawDescGZIP(), []int{1}
}

func (x *TaskResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TaskResponse) GetArg1() string {
	if x != nil {
		return x.Arg1
	}
	return ""
}

func (x *TaskResponse) GetArg2() string {
	if x != nil {
		return x.Arg2
	}
	return ""
}

func (x *TaskResponse) GetTime() int32 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *TaskResponse) GetOperator() Operator {
	if x != nil {
		return x.Operator
	}
	return Operator_ADDICTION
}

var File_tasks_proto protoreflect.FileDescriptor

const file_tasks_proto_rawDesc = "" +
	"\n" +
	"\vtasks.proto\x12\x03gen\x1a\x1bgoogle/protobuf/empty.proto\";\n" +
	"\x11UpdateTaskRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x16\n" +
	"\x06result\x18\x02 \x01(\tR\x06result\"\x85\x01\n" +
	"\fTaskResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04arg1\x18\x02 \x01(\tR\x04arg1\x12\x12\n" +
	"\x04arg2\x18\x03 \x01(\tR\x04arg2\x12\x12\n" +
	"\x04time\x18\x04 \x01(\x05R\x04time\x12)\n" +
	"\boperator\x18\x05 \x01(\x0e2\r.gen.OperatorR\boperator*L\n" +
	"\bOperator\x12\r\n" +
	"\tADDICTION\x10\x00\x12\x0f\n" +
	"\vSUBTRACTION\x10\x01\x12\x12\n" +
	"\x0eMULTIPLICATION\x10\x02\x12\f\n" +
	"\bDIVISION\x10\x032\x84\x01\n" +
	"\fOrchestrator\x126\n" +
	"\aGetTask\x12\x16.google.protobuf.Empty\x1a\x11.gen.TaskResponse0\x01\x12<\n" +
	"\n" +
	"UpdateTask\x12\x16.gen.UpdateTaskRequest\x1a\x16.google.protobuf.EmptyB-Z+orchestrator/internal/controllers/tasks/genb\x06proto3"

var (
	file_tasks_proto_rawDescOnce sync.Once
	file_tasks_proto_rawDescData []byte
)

func file_tasks_proto_rawDescGZIP() []byte {
	file_tasks_proto_rawDescOnce.Do(func() {
		file_tasks_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_tasks_proto_rawDesc), len(file_tasks_proto_rawDesc)))
	})
	return file_tasks_proto_rawDescData
}

var file_tasks_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_tasks_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_tasks_proto_goTypes = []any{
	(Operator)(0),             // 0: gen.Operator
	(*UpdateTaskRequest)(nil), // 1: gen.UpdateTaskRequest
	(*TaskResponse)(nil),      // 2: gen.TaskResponse
	(*emptypb.Empty)(nil),     // 3: google.protobuf.Empty
}
var file_tasks_proto_depIdxs = []int32{
	0, // 0: gen.TaskResponse.operator:type_name -> gen.Operator
	3, // 1: gen.Orchestrator.GetTask:input_type -> google.protobuf.Empty
	1, // 2: gen.Orchestrator.UpdateTask:input_type -> gen.UpdateTaskRequest
	2, // 3: gen.Orchestrator.GetTask:output_type -> gen.TaskResponse
	3, // 4: gen.Orchestrator.UpdateTask:output_type -> google.protobuf.Empty
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tasks_proto_init() }
func file_tasks_proto_init() {
	if File_tasks_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_tasks_proto_rawDesc), len(file_tasks_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tasks_proto_goTypes,
		DependencyIndexes: file_tasks_proto_depIdxs,
		EnumInfos:         file_tasks_proto_enumTypes,
		MessageInfos:      file_tasks_proto_msgTypes,
	}.Build()
	File_tasks_proto = out.File
	file_tasks_proto_goTypes = nil
	file_tasks_proto_depIdxs = nil
}
