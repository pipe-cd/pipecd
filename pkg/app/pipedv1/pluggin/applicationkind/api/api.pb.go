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
// 	protoc        v3.18.1
// source: pkg/app/pipedv1/pluggin/applicationkind/api/api.proto

package api

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	model "github.com/pipe-cd/pipecd/pkg/model"
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

type BuildPlanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkingDir string `protobuf:"bytes,1,opt,name=working_dir,json=workingDir,proto3" json:"working_dir,omitempty"`
	// Last successful commit hash and config file name.
	// Use to build deployment source object for last successful deployment.
	LastSuccessfulCommitHash     string `protobuf:"bytes,2,opt,name=last_successful_commit_hash,json=lastSuccessfulCommitHash,proto3" json:"last_successful_commit_hash,omitempty"`
	LastSuccessfulConfigFileName string `protobuf:"bytes,3,opt,name=last_successful_config_file_name,json=lastSuccessfulConfigFileName,proto3" json:"last_successful_config_file_name,omitempty"`
	// The configuration of the piped that handles the deployment.
	PipedConfig []byte `protobuf:"bytes,4,opt,name=piped_config,json=pipedConfig,proto3" json:"piped_config,omitempty"`
	// The deployment to build a plan for.
	Deployment *model.Deployment `protobuf:"bytes,5,opt,name=deployment,proto3" json:"deployment,omitempty"`
}

func (x *BuildPlanRequest) Reset() {
	*x = BuildPlanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildPlanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildPlanRequest) ProtoMessage() {}

func (x *BuildPlanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildPlanRequest.ProtoReflect.Descriptor instead.
func (*BuildPlanRequest) Descriptor() ([]byte, []int) {
	return file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescGZIP(), []int{0}
}

func (x *BuildPlanRequest) GetWorkingDir() string {
	if x != nil {
		return x.WorkingDir
	}
	return ""
}

func (x *BuildPlanRequest) GetLastSuccessfulCommitHash() string {
	if x != nil {
		return x.LastSuccessfulCommitHash
	}
	return ""
}

func (x *BuildPlanRequest) GetLastSuccessfulConfigFileName() string {
	if x != nil {
		return x.LastSuccessfulConfigFileName
	}
	return ""
}

func (x *BuildPlanRequest) GetPipedConfig() []byte {
	if x != nil {
		return x.PipedConfig
	}
	return nil
}

func (x *BuildPlanRequest) GetDeployment() *model.Deployment {
	if x != nil {
		return x.Deployment
	}
	return nil
}

type BuildPlanResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The built deployment plan.
	Plan *DeploymentPlan `protobuf:"bytes,1,opt,name=plan,proto3" json:"plan,omitempty"`
}

func (x *BuildPlanResponse) Reset() {
	*x = BuildPlanResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildPlanResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildPlanResponse) ProtoMessage() {}

func (x *BuildPlanResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildPlanResponse.ProtoReflect.Descriptor instead.
func (*BuildPlanResponse) Descriptor() ([]byte, []int) {
	return file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescGZIP(), []int{1}
}

func (x *BuildPlanResponse) GetPlan() *DeploymentPlan {
	if x != nil {
		return x.Plan
	}
	return nil
}

type DeploymentPlan struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SyncStrategy model.SyncStrategy `protobuf:"varint,1,opt,name=sync_strategy,json=syncStrategy,proto3,enum=model.SyncStrategy" json:"sync_strategy,omitempty"`
	// Text summary of planned deployment.
	Summary  string                   `protobuf:"bytes,2,opt,name=summary,proto3" json:"summary,omitempty"`
	Versions []*model.ArtifactVersion `protobuf:"bytes,3,rep,name=versions,proto3" json:"versions,omitempty"`
	Stages   []*model.PipelineStage   `protobuf:"bytes,4,rep,name=stages,proto3" json:"stages,omitempty"`
}

func (x *DeploymentPlan) Reset() {
	*x = DeploymentPlan{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentPlan) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentPlan) ProtoMessage() {}

func (x *DeploymentPlan) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentPlan.ProtoReflect.Descriptor instead.
func (*DeploymentPlan) Descriptor() ([]byte, []int) {
	return file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescGZIP(), []int{2}
}

func (x *DeploymentPlan) GetSyncStrategy() model.SyncStrategy {
	if x != nil {
		return x.SyncStrategy
	}
	return model.SyncStrategy(0)
}

func (x *DeploymentPlan) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *DeploymentPlan) GetVersions() []*model.ArtifactVersion {
	if x != nil {
		return x.Versions
	}
	return nil
}

func (x *DeploymentPlan) GetStages() []*model.PipelineStage {
	if x != nil {
		return x.Stages
	}
	return nil
}

type ExecutePipelineRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The deployment plan to execute.
	Plan *DeploymentPlan `protobuf:"bytes,1,opt,name=plan,proto3" json:"plan,omitempty"`
}

func (x *ExecutePipelineRequest) Reset() {
	*x = ExecutePipelineRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecutePipelineRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutePipelineRequest) ProtoMessage() {}

func (x *ExecutePipelineRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecutePipelineRequest.ProtoReflect.Descriptor instead.
func (*ExecutePipelineRequest) Descriptor() ([]byte, []int) {
	return file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescGZIP(), []int{3}
}

func (x *ExecutePipelineRequest) GetPlan() *DeploymentPlan {
	if x != nil {
		return x.Plan
	}
	return nil
}

type ExecutePipelineResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status model.DeploymentStatus `protobuf:"varint,1,opt,name=status,proto3,enum=model.DeploymentStatus" json:"status,omitempty"`
	Log    string                 `protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`
}

func (x *ExecutePipelineResponse) Reset() {
	*x = ExecutePipelineResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecutePipelineResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutePipelineResponse) ProtoMessage() {}

func (x *ExecutePipelineResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecutePipelineResponse.ProtoReflect.Descriptor instead.
func (*ExecutePipelineResponse) Descriptor() ([]byte, []int) {
	return file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescGZIP(), []int{4}
}

func (x *ExecutePipelineResponse) GetStatus() model.DeploymentStatus {
	if x != nil {
		return x.Status
	}
	return model.DeploymentStatus(0)
}

func (x *ExecutePipelineResponse) GetLog() string {
	if x != nil {
		return x.Log
	}
	return ""
}

var File_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto protoreflect.FileDescriptor

var file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDesc = []byte{
	0x0a, 0x35, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x64, 0x76,
	0x31, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x67, 0x69, 0x6e, 0x2f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6b, 0x69, 0x6e, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70,
	0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x67, 0x69, 0x6e, 0x2e, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x6b, 0x69, 0x6e, 0x64, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16,
	0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x2f, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xac, 0x02, 0x0a, 0x10, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x50, 0x6c, 0x61, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x69,
	0x6e, 0x67, 0x5f, 0x64, 0x69, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x44, 0x69,
	0x72, 0x12, 0x3d, 0x0a, 0x1b, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x66, 0x75, 0x6c, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x18, 0x6c, 0x61, 0x73, 0x74, 0x53, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x48, 0x61, 0x73, 0x68,
	0x12, 0x46, 0x0a, 0x20, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x66, 0x75, 0x6c, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1c, 0x6c, 0x61, 0x73, 0x74,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x0c, 0x70, 0x69, 0x70, 0x65,
	0x64, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x7a, 0x02, 0x10, 0x01, 0x52, 0x0b, 0x70, 0x69, 0x70, 0x65, 0x64, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x3b, 0x0a, 0x0a, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x08, 0xfa, 0x42, 0x05,
	0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0a, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x22, 0x55, 0x0a, 0x11, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x50, 0x6c, 0x61, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x04, 0x70, 0x6c, 0x61, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x6c, 0x75, 0x67,
	0x67, 0x69, 0x6e, 0x2e, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6b,
	0x69, 0x6e, 0x64, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x6c,
	0x61, 0x6e, 0x52, 0x04, 0x70, 0x6c, 0x61, 0x6e, 0x22, 0xc6, 0x01, 0x0a, 0x0e, 0x44, 0x65, 0x70,
	0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x6c, 0x61, 0x6e, 0x12, 0x38, 0x0a, 0x0d, 0x73,
	0x79, 0x6e, 0x63, 0x5f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x53,
	0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x52, 0x0c, 0x73, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x72,
	0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12,
	0x32, 0x0a, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61,
	0x63, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x50, 0x69, 0x70, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x53, 0x74, 0x61, 0x67, 0x65, 0x52, 0x06, 0x73, 0x74, 0x61, 0x67, 0x65,
	0x73, 0x22, 0x64, 0x0a, 0x16, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x50, 0x69, 0x70, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x4a, 0x0a, 0x04, 0x70,
	0x6c, 0x61, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x70, 0x6c, 0x75, 0x67, 0x67, 0x69, 0x6e, 0x2e, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x6b, 0x69, 0x6e, 0x64, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x50, 0x6c, 0x61, 0x6e, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10,
	0x01, 0x52, 0x04, 0x70, 0x6c, 0x61, 0x6e, 0x22, 0x5c, 0x0a, 0x17, 0x45, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x65, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x17, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6c, 0x6f, 0x67, 0x32, 0x7e, 0x0a, 0x0e, 0x50, 0x6c, 0x61, 0x6e, 0x6e, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6c, 0x0a, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64,
	0x50, 0x6c, 0x61, 0x6e, 0x12, 0x2e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x6c, 0x75, 0x67,
	0x67, 0x69, 0x6e, 0x2e, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6b,
	0x69, 0x6e, 0x64, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x50, 0x6c, 0x61, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x6c, 0x75, 0x67,
	0x67, 0x69, 0x6e, 0x2e, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6b,
	0x69, 0x6e, 0x64, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x50, 0x6c, 0x61, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x94, 0x01, 0x0a, 0x0f, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x80, 0x01, 0x0a, 0x0f, 0x45, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x65, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x34, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x67, 0x69, 0x6e, 0x2e, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6b, 0x69, 0x6e, 0x64, 0x2e, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x65, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x35, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x67,
	0x69, 0x6e, 0x2e, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6b, 0x69,
	0x6e, 0x64, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69,
	0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x47, 0x5a, 0x45,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x2d,
	0x63, 0x64, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x63, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70,
	0x70, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x64, 0x76, 0x31, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x67, 0x69,
	0x6e, 0x2f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6b, 0x69, 0x6e,
	0x64, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescOnce sync.Once
	file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescData = file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDesc
)

func file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescGZIP() []byte {
	file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescOnce.Do(func() {
		file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescData)
	})
	return file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDescData
}

var file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_goTypes = []interface{}{
	(*BuildPlanRequest)(nil),        // 0: grpc.pluggin.applicationkind.BuildPlanRequest
	(*BuildPlanResponse)(nil),       // 1: grpc.pluggin.applicationkind.BuildPlanResponse
	(*DeploymentPlan)(nil),          // 2: grpc.pluggin.applicationkind.DeploymentPlan
	(*ExecutePipelineRequest)(nil),  // 3: grpc.pluggin.applicationkind.ExecutePipelineRequest
	(*ExecutePipelineResponse)(nil), // 4: grpc.pluggin.applicationkind.ExecutePipelineResponse
	(*model.Deployment)(nil),        // 5: model.Deployment
	(model.SyncStrategy)(0),         // 6: model.SyncStrategy
	(*model.ArtifactVersion)(nil),   // 7: model.ArtifactVersion
	(*model.PipelineStage)(nil),     // 8: model.PipelineStage
	(model.DeploymentStatus)(0),     // 9: model.DeploymentStatus
}
var file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_depIdxs = []int32{
	5, // 0: grpc.pluggin.applicationkind.BuildPlanRequest.deployment:type_name -> model.Deployment
	2, // 1: grpc.pluggin.applicationkind.BuildPlanResponse.plan:type_name -> grpc.pluggin.applicationkind.DeploymentPlan
	6, // 2: grpc.pluggin.applicationkind.DeploymentPlan.sync_strategy:type_name -> model.SyncStrategy
	7, // 3: grpc.pluggin.applicationkind.DeploymentPlan.versions:type_name -> model.ArtifactVersion
	8, // 4: grpc.pluggin.applicationkind.DeploymentPlan.stages:type_name -> model.PipelineStage
	2, // 5: grpc.pluggin.applicationkind.ExecutePipelineRequest.plan:type_name -> grpc.pluggin.applicationkind.DeploymentPlan
	9, // 6: grpc.pluggin.applicationkind.ExecutePipelineResponse.status:type_name -> model.DeploymentStatus
	0, // 7: grpc.pluggin.applicationkind.PlannerService.BuildPlan:input_type -> grpc.pluggin.applicationkind.BuildPlanRequest
	3, // 8: grpc.pluggin.applicationkind.ExecutorService.ExecutePipeline:input_type -> grpc.pluggin.applicationkind.ExecutePipelineRequest
	1, // 9: grpc.pluggin.applicationkind.PlannerService.BuildPlan:output_type -> grpc.pluggin.applicationkind.BuildPlanResponse
	4, // 10: grpc.pluggin.applicationkind.ExecutorService.ExecutePipeline:output_type -> grpc.pluggin.applicationkind.ExecutePipelineResponse
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_init() }
func file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_init() {
	if File_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildPlanRequest); i {
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
		file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildPlanResponse); i {
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
		file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentPlan); i {
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
		file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecutePipelineRequest); i {
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
		file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecutePipelineResponse); i {
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
			RawDescriptor: file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_goTypes,
		DependencyIndexes: file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_depIdxs,
		MessageInfos:      file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_msgTypes,
	}.Build()
	File_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto = out.File
	file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_rawDesc = nil
	file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_goTypes = nil
	file_pkg_app_pipedv1_pluggin_applicationkind_api_api_proto_depIdxs = nil
}
