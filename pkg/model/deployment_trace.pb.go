// Copyright 2025 The PipeCD Authors.
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
// 	protoc        v3.21.12
// source: pkg/model/deployment_trace.proto

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

type DeploymentTrace struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title           string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	CommitHash      string `protobuf:"bytes,3,opt,name=commit_hash,json=commitHash,proto3" json:"commit_hash,omitempty"`
	CommitUrl       string `protobuf:"bytes,4,opt,name=commit_url,json=commitUrl,proto3" json:"commit_url,omitempty"`
	CommitMessage   string `protobuf:"bytes,5,opt,name=commit_message,json=commitMessage,proto3" json:"commit_message,omitempty"`
	CommitTimestamp int64  `protobuf:"varint,6,opt,name=commit_timestamp,json=commitTimestamp,proto3" json:"commit_timestamp,omitempty"`
	Author          string `protobuf:"bytes,7,opt,name=author,proto3" json:"author,omitempty"`
	CompletedAt     int64  `protobuf:"varint,100,opt,name=completed_at,json=completedAt,proto3" json:"completed_at,omitempty"`
	CreatedAt       int64  `protobuf:"varint,101,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt       int64  `protobuf:"varint,102,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *DeploymentTrace) Reset() {
	*x = DeploymentTrace{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_deployment_trace_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentTrace) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentTrace) ProtoMessage() {}

func (x *DeploymentTrace) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_deployment_trace_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentTrace.ProtoReflect.Descriptor instead.
func (*DeploymentTrace) Descriptor() ([]byte, []int) {
	return file_pkg_model_deployment_trace_proto_rawDescGZIP(), []int{0}
}

func (x *DeploymentTrace) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DeploymentTrace) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *DeploymentTrace) GetCommitHash() string {
	if x != nil {
		return x.CommitHash
	}
	return ""
}

func (x *DeploymentTrace) GetCommitUrl() string {
	if x != nil {
		return x.CommitUrl
	}
	return ""
}

func (x *DeploymentTrace) GetCommitMessage() string {
	if x != nil {
		return x.CommitMessage
	}
	return ""
}

func (x *DeploymentTrace) GetCommitTimestamp() int64 {
	if x != nil {
		return x.CommitTimestamp
	}
	return 0
}

func (x *DeploymentTrace) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

func (x *DeploymentTrace) GetCompletedAt() int64 {
	if x != nil {
		return x.CompletedAt
	}
	return 0
}

func (x *DeploymentTrace) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *DeploymentTrace) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

var File_pkg_model_deployment_trace_proto protoreflect.FileDescriptor

var file_pkg_model_deployment_trace_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x64, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x93, 0x03, 0x0a, 0x0f, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x54, 0x72, 0x61, 0x63, 0x65, 0x12, 0x17, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1d, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x28,
	0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x63, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x26, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x55, 0x72, 0x6c,
	0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x32, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x03, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x22, 0x02, 0x28, 0x00, 0x52, 0x0f, 0x63, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1f, 0x0a, 0x06, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x2a, 0x0a, 0x0c,
	0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x64, 0x20, 0x01,
	0x28, 0x03, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x22, 0x02, 0x28, 0x00, 0x52, 0x0b, 0x63, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x26, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x65, 0x20, 0x01, 0x28, 0x03, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x22, 0x02, 0x28, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x26, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x66,
	0x20, 0x01, 0x28, 0x03, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x22, 0x02, 0x28, 0x00, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x2d, 0x63, 0x64, 0x2f, 0x70,
	0x69, 0x70, 0x65, 0x63, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_model_deployment_trace_proto_rawDescOnce sync.Once
	file_pkg_model_deployment_trace_proto_rawDescData = file_pkg_model_deployment_trace_proto_rawDesc
)

func file_pkg_model_deployment_trace_proto_rawDescGZIP() []byte {
	file_pkg_model_deployment_trace_proto_rawDescOnce.Do(func() {
		file_pkg_model_deployment_trace_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_model_deployment_trace_proto_rawDescData)
	})
	return file_pkg_model_deployment_trace_proto_rawDescData
}

var file_pkg_model_deployment_trace_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_model_deployment_trace_proto_goTypes = []interface{}{
	(*DeploymentTrace)(nil), // 0: model.DeploymentTrace
}
var file_pkg_model_deployment_trace_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_model_deployment_trace_proto_init() }
func file_pkg_model_deployment_trace_proto_init() {
	if File_pkg_model_deployment_trace_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_model_deployment_trace_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentTrace); i {
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
			RawDescriptor: file_pkg_model_deployment_trace_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_model_deployment_trace_proto_goTypes,
		DependencyIndexes: file_pkg_model_deployment_trace_proto_depIdxs,
		MessageInfos:      file_pkg_model_deployment_trace_proto_msgTypes,
	}.Build()
	File_pkg_model_deployment_trace_proto = out.File
	file_pkg_model_deployment_trace_proto_rawDesc = nil
	file_pkg_model_deployment_trace_proto_goTypes = nil
	file_pkg_model_deployment_trace_proto_depIdxs = nil
}
