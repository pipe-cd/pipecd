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
// 	protoc        v3.19.4
// source: pkg/model/common.proto

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

type ApplicationKind int32

const (
	ApplicationKind_KUBERNETES ApplicationKind = 0
	ApplicationKind_TERRAFORM  ApplicationKind = 1
	// TODO: Uncomment when support CrossPlane application kind.
	// CROSSPLANE = 2;
	ApplicationKind_LAMBDA   ApplicationKind = 3
	ApplicationKind_CLOUDRUN ApplicationKind = 4
	ApplicationKind_ECS      ApplicationKind = 5
)

// Enum value maps for ApplicationKind.
var (
	ApplicationKind_name = map[int32]string{
		0: "KUBERNETES",
		1: "TERRAFORM",
		3: "LAMBDA",
		4: "CLOUDRUN",
		5: "ECS",
	}
	ApplicationKind_value = map[string]int32{
		"KUBERNETES": 0,
		"TERRAFORM":  1,
		"LAMBDA":     3,
		"CLOUDRUN":   4,
		"ECS":        5,
	}
)

func (x ApplicationKind) Enum() *ApplicationKind {
	p := new(ApplicationKind)
	*p = x
	return p
}

func (x ApplicationKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ApplicationKind) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_model_common_proto_enumTypes[0].Descriptor()
}

func (ApplicationKind) Type() protoreflect.EnumType {
	return &file_pkg_model_common_proto_enumTypes[0]
}

func (x ApplicationKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ApplicationKind.Descriptor instead.
func (ApplicationKind) EnumDescriptor() ([]byte, []int) {
	return file_pkg_model_common_proto_rawDescGZIP(), []int{0}
}

type ApplicationActiveStatus int32

const (
	ApplicationActiveStatus_ENABLED  ApplicationActiveStatus = 0
	ApplicationActiveStatus_DISABLED ApplicationActiveStatus = 1
	ApplicationActiveStatus_DELETED  ApplicationActiveStatus = 2
)

// Enum value maps for ApplicationActiveStatus.
var (
	ApplicationActiveStatus_name = map[int32]string{
		0: "ENABLED",
		1: "DISABLED",
		2: "DELETED",
	}
	ApplicationActiveStatus_value = map[string]int32{
		"ENABLED":  0,
		"DISABLED": 1,
		"DELETED":  2,
	}
)

func (x ApplicationActiveStatus) Enum() *ApplicationActiveStatus {
	p := new(ApplicationActiveStatus)
	*p = x
	return p
}

func (x ApplicationActiveStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ApplicationActiveStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_model_common_proto_enumTypes[1].Descriptor()
}

func (ApplicationActiveStatus) Type() protoreflect.EnumType {
	return &file_pkg_model_common_proto_enumTypes[1]
}

func (x ApplicationActiveStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ApplicationActiveStatus.Descriptor instead.
func (ApplicationActiveStatus) EnumDescriptor() ([]byte, []int) {
	return file_pkg_model_common_proto_rawDescGZIP(), []int{1}
}

type SyncStrategy int32

const (
	SyncStrategy_AUTO       SyncStrategy = 0
	SyncStrategy_QUICK_SYNC SyncStrategy = 1
	SyncStrategy_PIPELINE   SyncStrategy = 2
)

// Enum value maps for SyncStrategy.
var (
	SyncStrategy_name = map[int32]string{
		0: "AUTO",
		1: "QUICK_SYNC",
		2: "PIPELINE",
	}
	SyncStrategy_value = map[string]int32{
		"AUTO":       0,
		"QUICK_SYNC": 1,
		"PIPELINE":   2,
	}
)

func (x SyncStrategy) Enum() *SyncStrategy {
	p := new(SyncStrategy)
	*p = x
	return p
}

func (x SyncStrategy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SyncStrategy) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_model_common_proto_enumTypes[2].Descriptor()
}

func (SyncStrategy) Type() protoreflect.EnumType {
	return &file_pkg_model_common_proto_enumTypes[2]
}

func (x SyncStrategy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SyncStrategy.Descriptor instead.
func (SyncStrategy) EnumDescriptor() ([]byte, []int) {
	return file_pkg_model_common_proto_rawDescGZIP(), []int{2}
}

type ArtifactVersion_Kind int32

const (
	ArtifactVersion_UNKNOWN          ArtifactVersion_Kind = 0
	ArtifactVersion_CONTAINER_IMAGE  ArtifactVersion_Kind = 1
	ArtifactVersion_S3_OBJECT        ArtifactVersion_Kind = 2
	ArtifactVersion_GIT_SOURCE       ArtifactVersion_Kind = 3
	ArtifactVersion_TERRAFORM_MODULE ArtifactVersion_Kind = 4
)

// Enum value maps for ArtifactVersion_Kind.
var (
	ArtifactVersion_Kind_name = map[int32]string{
		0: "UNKNOWN",
		1: "CONTAINER_IMAGE",
		2: "S3_OBJECT",
		3: "GIT_SOURCE",
		4: "TERRAFORM_MODULE",
	}
	ArtifactVersion_Kind_value = map[string]int32{
		"UNKNOWN":          0,
		"CONTAINER_IMAGE":  1,
		"S3_OBJECT":        2,
		"GIT_SOURCE":       3,
		"TERRAFORM_MODULE": 4,
	}
)

func (x ArtifactVersion_Kind) Enum() *ArtifactVersion_Kind {
	p := new(ArtifactVersion_Kind)
	*p = x
	return p
}

func (x ArtifactVersion_Kind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ArtifactVersion_Kind) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_model_common_proto_enumTypes[3].Descriptor()
}

func (ArtifactVersion_Kind) Type() protoreflect.EnumType {
	return &file_pkg_model_common_proto_enumTypes[3]
}

func (x ArtifactVersion_Kind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ArtifactVersion_Kind.Descriptor instead.
func (ArtifactVersion_Kind) EnumDescriptor() ([]byte, []int) {
	return file_pkg_model_common_proto_rawDescGZIP(), []int{3, 0}
}

type ApplicationGitPath struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The repository that was configured at piped.
	Repo *ApplicationGitRepository `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
	// TODO: Make this field immutable.
	Path string `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	// Deprecated: Do not use.
	ConfigPath     string `protobuf:"bytes,3,opt,name=config_path,json=configPath,proto3" json:"config_path,omitempty"`
	ConfigFilename string `protobuf:"bytes,4,opt,name=config_filename,json=configFilename,proto3" json:"config_filename,omitempty"`
	Url            string `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *ApplicationGitPath) Reset() {
	*x = ApplicationGitPath{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplicationGitPath) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplicationGitPath) ProtoMessage() {}

func (x *ApplicationGitPath) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplicationGitPath.ProtoReflect.Descriptor instead.
func (*ApplicationGitPath) Descriptor() ([]byte, []int) {
	return file_pkg_model_common_proto_rawDescGZIP(), []int{0}
}

func (x *ApplicationGitPath) GetRepo() *ApplicationGitRepository {
	if x != nil {
		return x.Repo
	}
	return nil
}

func (x *ApplicationGitPath) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

// Deprecated: Do not use.
func (x *ApplicationGitPath) GetConfigPath() string {
	if x != nil {
		return x.ConfigPath
	}
	return ""
}

func (x *ApplicationGitPath) GetConfigFilename() string {
	if x != nil {
		return x.ConfigFilename
	}
	return ""
}

func (x *ApplicationGitPath) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type ApplicationGitRepository struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Remote string `protobuf:"bytes,2,opt,name=remote,proto3" json:"remote,omitempty"`
	Branch string `protobuf:"bytes,3,opt,name=branch,proto3" json:"branch,omitempty"`
}

func (x *ApplicationGitRepository) Reset() {
	*x = ApplicationGitRepository{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplicationGitRepository) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplicationGitRepository) ProtoMessage() {}

func (x *ApplicationGitRepository) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplicationGitRepository.ProtoReflect.Descriptor instead.
func (*ApplicationGitRepository) Descriptor() ([]byte, []int) {
	return file_pkg_model_common_proto_rawDescGZIP(), []int{1}
}

func (x *ApplicationGitRepository) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ApplicationGitRepository) GetRemote() string {
	if x != nil {
		return x.Remote
	}
	return ""
}

func (x *ApplicationGitRepository) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

type ApplicationInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// This field is not allowed to be changed.
	Kind   ApplicationKind   `protobuf:"varint,3,opt,name=kind,proto3,enum=model.ApplicationKind" json:"kind,omitempty"`
	Labels map[string]string `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// This field is not allowed to be changed.
	RepoId string `protobuf:"bytes,5,opt,name=repo_id,json=repoId,proto3" json:"repo_id,omitempty"`
	// This field is not allowed to be changed.
	Path string `protobuf:"bytes,6,opt,name=path,proto3" json:"path,omitempty"`
	// This field is not allowed to be changed.
	ConfigFilename string `protobuf:"bytes,7,opt,name=config_filename,json=configFilename,proto3" json:"config_filename,omitempty"`
	PipedId        string `protobuf:"bytes,8,opt,name=piped_id,json=pipedId,proto3" json:"piped_id,omitempty"`
	Description    string `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *ApplicationInfo) Reset() {
	*x = ApplicationInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplicationInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplicationInfo) ProtoMessage() {}

func (x *ApplicationInfo) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplicationInfo.ProtoReflect.Descriptor instead.
func (*ApplicationInfo) Descriptor() ([]byte, []int) {
	return file_pkg_model_common_proto_rawDescGZIP(), []int{2}
}

func (x *ApplicationInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ApplicationInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ApplicationInfo) GetKind() ApplicationKind {
	if x != nil {
		return x.Kind
	}
	return ApplicationKind_KUBERNETES
}

func (x *ApplicationInfo) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *ApplicationInfo) GetRepoId() string {
	if x != nil {
		return x.RepoId
	}
	return ""
}

func (x *ApplicationInfo) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ApplicationInfo) GetConfigFilename() string {
	if x != nil {
		return x.ConfigFilename
	}
	return ""
}

func (x *ApplicationInfo) GetPipedId() string {
	if x != nil {
		return x.PipedId
	}
	return ""
}

func (x *ApplicationInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type ArtifactVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind    ArtifactVersion_Kind `protobuf:"varint,1,opt,name=kind,proto3,enum=model.ArtifactVersion_Kind" json:"kind,omitempty"`
	Version string               `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Name    string               `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Url     string               `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *ArtifactVersion) Reset() {
	*x = ArtifactVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArtifactVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArtifactVersion) ProtoMessage() {}

func (x *ArtifactVersion) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArtifactVersion.ProtoReflect.Descriptor instead.
func (*ArtifactVersion) Descriptor() ([]byte, []int) {
	return file_pkg_model_common_proto_rawDescGZIP(), []int{3}
}

func (x *ArtifactVersion) GetKind() ArtifactVersion_Kind {
	if x != nil {
		return x.Kind
	}
	return ArtifactVersion_UNKNOWN
}

func (x *ArtifactVersion) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *ArtifactVersion) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ArtifactVersion) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_pkg_model_common_proto protoreflect.FileDescriptor

var file_pkg_model_common_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x1a,
	0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd8, 0x01, 0x0a, 0x12, 0x41, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x47, 0x69, 0x74, 0x50, 0x61, 0x74, 0x68, 0x12,
	0x3d, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x47, 0x69, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x23,
	0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0f, 0xfa, 0x42,
	0x0c, 0x72, 0x0a, 0x32, 0x08, 0x5e, 0x5b, 0x5e, 0x2f, 0x5d, 0x2e, 0x2b, 0x24, 0x52, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x12, 0x23, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x70, 0x61,
	0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52, 0x0a, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x50, 0x61, 0x74, 0x68, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x22, 0x63, 0x0a, 0x18, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x47, 0x69, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12,
	0x17, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x22, 0xa7, 0x03, 0x0a, 0x0f, 0x41, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e,
	0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x69, 0x6e, 0x64, 0x42,
	0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12,
	0x3a, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x20, 0x0a, 0x07, 0x72,
	0x65, 0x70, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x49, 0x64, 0x12, 0x23, 0x0a,
	0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0f, 0xfa, 0x42, 0x0c,
	0x72, 0x0a, 0x32, 0x08, 0x5e, 0x5b, 0x5e, 0x2f, 0x5d, 0x2e, 0x2b, 0x24, 0x52, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x08, 0x70,
	0x69, 0x70, 0x65, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x07, 0x70, 0x69, 0x70, 0x65, 0x64, 0x49, 0x64, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x4a, 0x04, 0x08, 0x0e,
	0x10, 0x0f, 0x22, 0xeb, 0x01, 0x0a, 0x0f, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x39, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x72, 0x74,
	0x69, 0x66, 0x61, 0x63, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x4b, 0x69, 0x6e,
	0x64, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x22, 0x5d, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x4f, 0x4e, 0x54, 0x41, 0x49,
	0x4e, 0x45, 0x52, 0x5f, 0x49, 0x4d, 0x41, 0x47, 0x45, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x53,
	0x33, 0x5f, 0x4f, 0x42, 0x4a, 0x45, 0x43, 0x54, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x47, 0x49,
	0x54, 0x5f, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x10, 0x03, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x45,
	0x52, 0x52, 0x41, 0x46, 0x4f, 0x52, 0x4d, 0x5f, 0x4d, 0x4f, 0x44, 0x55, 0x4c, 0x45, 0x10, 0x04,
	0x2a, 0x53, 0x0a, 0x0f, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4b,
	0x69, 0x6e, 0x64, 0x12, 0x0e, 0x0a, 0x0a, 0x4b, 0x55, 0x42, 0x45, 0x52, 0x4e, 0x45, 0x54, 0x45,
	0x53, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x45, 0x52, 0x52, 0x41, 0x46, 0x4f, 0x52, 0x4d,
	0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x41, 0x4d, 0x42, 0x44, 0x41, 0x10, 0x03, 0x12, 0x0c,
	0x0a, 0x08, 0x43, 0x4c, 0x4f, 0x55, 0x44, 0x52, 0x55, 0x4e, 0x10, 0x04, 0x12, 0x07, 0x0a, 0x03,
	0x45, 0x43, 0x53, 0x10, 0x05, 0x2a, 0x41, 0x0a, 0x17, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x0b, 0x0a, 0x07, 0x45, 0x4e, 0x41, 0x42, 0x4c, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a,
	0x08, 0x44, 0x49, 0x53, 0x41, 0x42, 0x4c, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x44,
	0x45, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x02, 0x2a, 0x36, 0x0a, 0x0c, 0x53, 0x79, 0x6e, 0x63,
	0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x08, 0x0a, 0x04, 0x41, 0x55, 0x54, 0x4f,
	0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x51, 0x55, 0x49, 0x43, 0x4b, 0x5f, 0x53, 0x59, 0x4e, 0x43,
	0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x49, 0x50, 0x45, 0x4c, 0x49, 0x4e, 0x45, 0x10, 0x02,
	0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70,
	0x69, 0x70, 0x65, 0x2d, 0x63, 0x64, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x63, 0x64, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_model_common_proto_rawDescOnce sync.Once
	file_pkg_model_common_proto_rawDescData = file_pkg_model_common_proto_rawDesc
)

func file_pkg_model_common_proto_rawDescGZIP() []byte {
	file_pkg_model_common_proto_rawDescOnce.Do(func() {
		file_pkg_model_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_model_common_proto_rawDescData)
	})
	return file_pkg_model_common_proto_rawDescData
}

var file_pkg_model_common_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_pkg_model_common_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_model_common_proto_goTypes = []interface{}{
	(ApplicationKind)(0),             // 0: model.ApplicationKind
	(ApplicationActiveStatus)(0),     // 1: model.ApplicationActiveStatus
	(SyncStrategy)(0),                // 2: model.SyncStrategy
	(ArtifactVersion_Kind)(0),        // 3: model.ArtifactVersion.Kind
	(*ApplicationGitPath)(nil),       // 4: model.ApplicationGitPath
	(*ApplicationGitRepository)(nil), // 5: model.ApplicationGitRepository
	(*ApplicationInfo)(nil),          // 6: model.ApplicationInfo
	(*ArtifactVersion)(nil),          // 7: model.ArtifactVersion
	nil,                              // 8: model.ApplicationInfo.LabelsEntry
}
var file_pkg_model_common_proto_depIdxs = []int32{
	5, // 0: model.ApplicationGitPath.repo:type_name -> model.ApplicationGitRepository
	0, // 1: model.ApplicationInfo.kind:type_name -> model.ApplicationKind
	8, // 2: model.ApplicationInfo.labels:type_name -> model.ApplicationInfo.LabelsEntry
	3, // 3: model.ArtifactVersion.kind:type_name -> model.ArtifactVersion.Kind
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pkg_model_common_proto_init() }
func file_pkg_model_common_proto_init() {
	if File_pkg_model_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_model_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplicationGitPath); i {
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
		file_pkg_model_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplicationGitRepository); i {
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
		file_pkg_model_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplicationInfo); i {
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
		file_pkg_model_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArtifactVersion); i {
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
			RawDescriptor: file_pkg_model_common_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_model_common_proto_goTypes,
		DependencyIndexes: file_pkg_model_common_proto_depIdxs,
		EnumInfos:         file_pkg_model_common_proto_enumTypes,
		MessageInfos:      file_pkg_model_common_proto_msgTypes,
	}.Build()
	File_pkg_model_common_proto = out.File
	file_pkg_model_common_proto_rawDesc = nil
	file_pkg_model_common_proto_goTypes = nil
	file_pkg_model_common_proto_depIdxs = nil
}
