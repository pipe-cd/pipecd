// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.14.0
// source: pkg/model/logblock.proto

package model

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type LogSeverity int32

const (
	LogSeverity_INFO    LogSeverity = 0
	LogSeverity_SUCCESS LogSeverity = 1
	LogSeverity_ERROR   LogSeverity = 2
)

// Enum value maps for LogSeverity.
var (
	LogSeverity_name = map[int32]string{
		0: "INFO",
		1: "SUCCESS",
		2: "ERROR",
	}
	LogSeverity_value = map[string]int32{
		"INFO":    0,
		"SUCCESS": 1,
		"ERROR":   2,
	}
)

func (x LogSeverity) Enum() *LogSeverity {
	p := new(LogSeverity)
	*p = x
	return p
}

func (x LogSeverity) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LogSeverity) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_model_logblock_proto_enumTypes[0].Descriptor()
}

func (LogSeverity) Type() protoreflect.EnumType {
	return &file_pkg_model_logblock_proto_enumTypes[0]
}

func (x LogSeverity) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LogSeverity.Descriptor instead.
func (LogSeverity) EnumDescriptor() ([]byte, []int) {
	return file_pkg_model_logblock_proto_rawDescGZIP(), []int{0}
}

type LogBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The index of log block.
	Index int64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	// The log content.
	Log string `protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`
	// Severity level for this block.
	Severity LogSeverity `protobuf:"varint,3,opt,name=severity,proto3,enum=model.LogSeverity" json:"severity,omitempty"`
	// Unix time when the log block was created.
	CreatedAt int64 `protobuf:"varint,14,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *LogBlock) Reset() {
	*x = LogBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_logblock_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogBlock) ProtoMessage() {}

func (x *LogBlock) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_logblock_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogBlock.ProtoReflect.Descriptor instead.
func (*LogBlock) Descriptor() ([]byte, []int) {
	return file_pkg_model_logblock_proto_rawDescGZIP(), []int{0}
}

func (x *LogBlock) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *LogBlock) GetLog() string {
	if x != nil {
		return x.Log
	}
	return ""
}

func (x *LogBlock) GetSeverity() LogSeverity {
	if x != nil {
		return x.Severity
	}
	return LogSeverity_INFO
}

func (x *LogBlock) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

var File_pkg_model_logblock_proto protoreflect.FileDescriptor

var file_pkg_model_logblock_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x6c, 0x6f, 0x67, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9d, 0x01, 0x0a, 0x08, 0x4c,
	0x6f, 0x67, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x19, 0x0a,
	0x03, 0x6c, 0x6f, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x03, 0x6c, 0x6f, 0x67, 0x12, 0x38, 0x0a, 0x08, 0x73, 0x65, 0x76, 0x65,
	0x72, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x2e, 0x4c, 0x6f, 0x67, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69,
	0x74, 0x79, 0x12, 0x26, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x22, 0x02, 0x20, 0x00, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2a, 0x2f, 0x0a, 0x0b, 0x4c, 0x6f,
	0x67, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x4e, 0x46,
	0x4f, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x01,
	0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x02, 0x42, 0x25, 0x5a, 0x23, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x2d, 0x63,
	0x64, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x63, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_model_logblock_proto_rawDescOnce sync.Once
	file_pkg_model_logblock_proto_rawDescData = file_pkg_model_logblock_proto_rawDesc
)

func file_pkg_model_logblock_proto_rawDescGZIP() []byte {
	file_pkg_model_logblock_proto_rawDescOnce.Do(func() {
		file_pkg_model_logblock_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_model_logblock_proto_rawDescData)
	})
	return file_pkg_model_logblock_proto_rawDescData
}

var file_pkg_model_logblock_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_model_logblock_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_model_logblock_proto_goTypes = []interface{}{
	(LogSeverity)(0), // 0: model.LogSeverity
	(*LogBlock)(nil), // 1: model.LogBlock
}
var file_pkg_model_logblock_proto_depIdxs = []int32{
	0, // 0: model.LogBlock.severity:type_name -> model.LogSeverity
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_model_logblock_proto_init() }
func file_pkg_model_logblock_proto_init() {
	if File_pkg_model_logblock_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_model_logblock_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogBlock); i {
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
			RawDescriptor: file_pkg_model_logblock_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_model_logblock_proto_goTypes,
		DependencyIndexes: file_pkg_model_logblock_proto_depIdxs,
		EnumInfos:         file_pkg_model_logblock_proto_enumTypes,
		MessageInfos:      file_pkg_model_logblock_proto_msgTypes,
	}.Build()
	File_pkg_model_logblock_proto = out.File
	file_pkg_model_logblock_proto_rawDesc = nil
	file_pkg_model_logblock_proto_goTypes = nil
	file_pkg_model_logblock_proto_depIdxs = nil
}
