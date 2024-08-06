// Copyright 2024 The PipeCD Authors.
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
// source: pkg/model/deployment_source.proto

package model

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

type DeploymentSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The application directory where the source code is located.
	ApplicationDirectory string `protobuf:"bytes,1,opt,name=application_directory,json=applicationDirectory,proto3" json:"application_directory,omitempty"`
	// The git commit revision of the source code.
	Revision string `protobuf:"bytes,2,opt,name=revision,proto3" json:"revision,omitempty"`
	// The configuration of the application which is independent for plugins.
	GenericApplicationConfig *GenericApplicationSpec `protobuf:"bytes,3,opt,name=generic_application_config,json=genericApplicationConfig,proto3" json:"generic_application_config,omitempty"`
	// The configuration of the application which is specific for plugins.
	ApplicationConfig *PluginApplicationSpec `protobuf:"bytes,4,opt,name=application_config,json=applicationConfig,proto3" json:"application_config,omitempty"`
}

func (x *DeploymentSource) Reset() {
	*x = DeploymentSource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_deployment_source_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentSource) ProtoMessage() {}

func (x *DeploymentSource) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_deployment_source_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentSource.ProtoReflect.Descriptor instead.
func (*DeploymentSource) Descriptor() ([]byte, []int) {
	return file_pkg_model_deployment_source_proto_rawDescGZIP(), []int{0}
}

func (x *DeploymentSource) GetApplicationDirectory() string {
	if x != nil {
		return x.ApplicationDirectory
	}
	return ""
}

func (x *DeploymentSource) GetRevision() string {
	if x != nil {
		return x.Revision
	}
	return ""
}

func (x *DeploymentSource) GetGenericApplicationConfig() *GenericApplicationSpec {
	if x != nil {
		return x.GenericApplicationConfig
	}
	return nil
}

func (x *DeploymentSource) GetApplicationConfig() *PluginApplicationSpec {
	if x != nil {
		return x.ApplicationConfig
	}
	return nil
}

type GenericApplicationSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GenericApplicationSpec) Reset() {
	*x = GenericApplicationSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_deployment_source_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenericApplicationSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenericApplicationSpec) ProtoMessage() {}

func (x *GenericApplicationSpec) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_deployment_source_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenericApplicationSpec.ProtoReflect.Descriptor instead.
func (*GenericApplicationSpec) Descriptor() ([]byte, []int) {
	return file_pkg_model_deployment_source_proto_rawDescGZIP(), []int{1}
}

type PluginApplicationSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The kind of the application spec.
	Kind string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	// The version of the application spec.
	ApiVersion string `protobuf:"bytes,2,opt,name=api_version,json=apiVersion,proto3" json:"api_version,omitempty"`
	// The raw data of the application spec.
	// We store this spec in bytes because we don't want to define the structure of the plugin-specific config.
	// The serialization format is JSON.
	Spec []byte `protobuf:"bytes,3,opt,name=spec,proto3" json:"spec,omitempty"`
}

func (x *PluginApplicationSpec) Reset() {
	*x = PluginApplicationSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_deployment_source_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginApplicationSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginApplicationSpec) ProtoMessage() {}

func (x *PluginApplicationSpec) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_deployment_source_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginApplicationSpec.ProtoReflect.Descriptor instead.
func (*PluginApplicationSpec) Descriptor() ([]byte, []int) {
	return file_pkg_model_deployment_source_proto_rawDescGZIP(), []int{2}
}

func (x *PluginApplicationSpec) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *PluginApplicationSpec) GetApiVersion() string {
	if x != nil {
		return x.ApiVersion
	}
	return ""
}

func (x *PluginApplicationSpec) GetSpec() []byte {
	if x != nil {
		return x.Spec
	}
	return nil
}

var File_pkg_model_deployment_source_proto protoreflect.FileDescriptor

var file_pkg_model_deployment_source_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x64, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x22, 0x8d, 0x02, 0x0a, 0x10, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x33, 0x0a, 0x15, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x5b, 0x0a, 0x1a, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x5f, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x47, 0x65, 0x6e,
	0x65, 0x72, 0x69, 0x63, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x18, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x41, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x4b, 0x0a,
	0x12, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x11, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x18, 0x0a, 0x16, 0x47, 0x65,
	0x6e, 0x65, 0x72, 0x69, 0x63, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x70, 0x65, 0x63, 0x22, 0x60, 0x0a, 0x15, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x41, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x70, 0x65, 0x63, 0x12, 0x12, 0x0a,
	0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x70, 0x69, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x2d, 0x63, 0x64, 0x2f, 0x70, 0x69, 0x70,
	0x65, 0x63, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_model_deployment_source_proto_rawDescOnce sync.Once
	file_pkg_model_deployment_source_proto_rawDescData = file_pkg_model_deployment_source_proto_rawDesc
)

func file_pkg_model_deployment_source_proto_rawDescGZIP() []byte {
	file_pkg_model_deployment_source_proto_rawDescOnce.Do(func() {
		file_pkg_model_deployment_source_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_model_deployment_source_proto_rawDescData)
	})
	return file_pkg_model_deployment_source_proto_rawDescData
}

var file_pkg_model_deployment_source_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_model_deployment_source_proto_goTypes = []interface{}{
	(*DeploymentSource)(nil),       // 0: model.DeploymentSource
	(*GenericApplicationSpec)(nil), // 1: model.GenericApplicationSpec
	(*PluginApplicationSpec)(nil),  // 2: model.PluginApplicationSpec
}
var file_pkg_model_deployment_source_proto_depIdxs = []int32{
	1, // 0: model.DeploymentSource.generic_application_config:type_name -> model.GenericApplicationSpec
	2, // 1: model.DeploymentSource.application_config:type_name -> model.PluginApplicationSpec
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_model_deployment_source_proto_init() }
func file_pkg_model_deployment_source_proto_init() {
	if File_pkg_model_deployment_source_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_model_deployment_source_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentSource); i {
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
		file_pkg_model_deployment_source_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenericApplicationSpec); i {
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
		file_pkg_model_deployment_source_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginApplicationSpec); i {
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
			RawDescriptor: file_pkg_model_deployment_source_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_model_deployment_source_proto_goTypes,
		DependencyIndexes: file_pkg_model_deployment_source_proto_depIdxs,
		MessageInfos:      file_pkg_model_deployment_source_proto_msgTypes,
	}.Build()
	File_pkg_model_deployment_source_proto = out.File
	file_pkg_model_deployment_source_proto_rawDesc = nil
	file_pkg_model_deployment_source_proto_goTypes = nil
	file_pkg_model_deployment_source_proto_depIdxs = nil
}
