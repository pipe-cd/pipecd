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
// source: pkg/model/application.proto

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

// ApplicationSyncStatus represents the current state of syncing the application.
type ApplicationSyncStatus int32

const (
	ApplicationSyncStatus_UNKNOWN     ApplicationSyncStatus = 0
	ApplicationSyncStatus_SYNCED      ApplicationSyncStatus = 1
	ApplicationSyncStatus_DEPLOYING   ApplicationSyncStatus = 2
	ApplicationSyncStatus_OUT_OF_SYNC ApplicationSyncStatus = 3
)

// Enum value maps for ApplicationSyncStatus.
var (
	ApplicationSyncStatus_name = map[int32]string{
		0: "UNKNOWN",
		1: "SYNCED",
		2: "DEPLOYING",
		3: "OUT_OF_SYNC",
	}
	ApplicationSyncStatus_value = map[string]int32{
		"UNKNOWN":     0,
		"SYNCED":      1,
		"DEPLOYING":   2,
		"OUT_OF_SYNC": 3,
	}
)

func (x ApplicationSyncStatus) Enum() *ApplicationSyncStatus {
	p := new(ApplicationSyncStatus)
	*p = x
	return p
}

func (x ApplicationSyncStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ApplicationSyncStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_model_application_proto_enumTypes[0].Descriptor()
}

func (ApplicationSyncStatus) Type() protoreflect.EnumType {
	return &file_pkg_model_application_proto_enumTypes[0]
}

func (x ApplicationSyncStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ApplicationSyncStatus.Descriptor instead.
func (ApplicationSyncStatus) EnumDescriptor() ([]byte, []int) {
	return file_pkg_model_application_proto_rawDescGZIP(), []int{0}
}

type Application struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The generated unique identifier.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The name of the application.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The ID of the piped that should handle this application.
	PipedId string `protobuf:"bytes,4,opt,name=piped_id,json=pipedId,proto3" json:"piped_id,omitempty"`
	// The ID of the project this environment belongs to.
	ProjectId string `protobuf:"bytes,5,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	// What kind of this application.
	Kind ApplicationKind `protobuf:"varint,6,opt,name=kind,proto3,enum=model.ApplicationKind" json:"kind,omitempty"`
	// The path to the git directory of this application.
	GitPath *ApplicationGitPath `protobuf:"bytes,7,opt,name=git_path,json=gitPath,proto3" json:"git_path,omitempty"`
	// The name of cloud provider where to deploy this application.
	// This must be one of the provider names registered in the piped.
	CloudProvider string `protobuf:"bytes,8,opt,name=cloud_provider,json=cloudProvider,proto3" json:"cloud_provider,omitempty"`
	// Additional description about application.
	Description string `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
	// Custom attributes to identify applications.
	Labels map[string]string `protobuf:"bytes,10,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Basic information about the most recently successful deployment.
	// This also shows information about current running workloads.
	MostRecentlySuccessfulDeployment *ApplicationDeploymentReference `protobuf:"bytes,11,opt,name=most_recently_successful_deployment,json=mostRecentlySuccessfulDeployment,proto3" json:"most_recently_successful_deployment,omitempty"`
	// Basic information about the most recently trigered deployment.
	MostRecentlyTriggeredDeployment *ApplicationDeploymentReference `protobuf:"bytes,12,opt,name=most_recently_triggered_deployment,json=mostRecentlyTriggeredDeployment,proto3" json:"most_recently_triggered_deployment,omitempty"`
	// Current sync state.
	SyncState *ApplicationSyncState `protobuf:"bytes,13,opt,name=sync_state,json=syncState,proto3" json:"sync_state,omitempty"`
	// Whether the application is deploying or not.
	Deploying bool `protobuf:"varint,14,opt,name=deploying,proto3" json:"deploying,omitempty"`
	// Unix time when the application was deleted.
	DeletedAt int64 `protobuf:"varint,98,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	// Whether the application is deleted or not.
	Deleted bool `protobuf:"varint,99,opt,name=deleted,proto3" json:"deleted,omitempty"`
	// Whether the application is disabled or not.
	Disabled bool `protobuf:"varint,100,opt,name=disabled,proto3" json:"disabled,omitempty"`
	// Unix time when the application is created.
	CreatedAt int64 `protobuf:"varint,101,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Unix time of the last time when the application is updated.
	UpdatedAt int64 `protobuf:"varint,102,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Application) Reset() {
	*x = Application{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_application_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Application) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Application) ProtoMessage() {}

func (x *Application) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_application_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Application.ProtoReflect.Descriptor instead.
func (*Application) Descriptor() ([]byte, []int) {
	return file_pkg_model_application_proto_rawDescGZIP(), []int{0}
}

func (x *Application) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Application) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Application) GetPipedId() string {
	if x != nil {
		return x.PipedId
	}
	return ""
}

func (x *Application) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *Application) GetKind() ApplicationKind {
	if x != nil {
		return x.Kind
	}
	return ApplicationKind_KUBERNETES
}

func (x *Application) GetGitPath() *ApplicationGitPath {
	if x != nil {
		return x.GitPath
	}
	return nil
}

func (x *Application) GetCloudProvider() string {
	if x != nil {
		return x.CloudProvider
	}
	return ""
}

func (x *Application) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Application) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Application) GetMostRecentlySuccessfulDeployment() *ApplicationDeploymentReference {
	if x != nil {
		return x.MostRecentlySuccessfulDeployment
	}
	return nil
}

func (x *Application) GetMostRecentlyTriggeredDeployment() *ApplicationDeploymentReference {
	if x != nil {
		return x.MostRecentlyTriggeredDeployment
	}
	return nil
}

func (x *Application) GetSyncState() *ApplicationSyncState {
	if x != nil {
		return x.SyncState
	}
	return nil
}

func (x *Application) GetDeploying() bool {
	if x != nil {
		return x.Deploying
	}
	return false
}

func (x *Application) GetDeletedAt() int64 {
	if x != nil {
		return x.DeletedAt
	}
	return 0
}

func (x *Application) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

func (x *Application) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *Application) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Application) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

// Current sync state of a specific application.
// This part is determined by drift detector component of piped.
type ApplicationSyncState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status           ApplicationSyncStatus `protobuf:"varint,1,opt,name=status,proto3,enum=model.ApplicationSyncStatus" json:"status,omitempty"`
	ShortReason      string                `protobuf:"bytes,2,opt,name=short_reason,json=shortReason,proto3" json:"short_reason,omitempty"`
	Reason           string                `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	HeadDeploymentId string                `protobuf:"bytes,4,opt,name=head_deployment_id,json=headDeploymentId,proto3" json:"head_deployment_id,omitempty"`
	Timestamp        int64                 `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *ApplicationSyncState) Reset() {
	*x = ApplicationSyncState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_application_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplicationSyncState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplicationSyncState) ProtoMessage() {}

func (x *ApplicationSyncState) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_application_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplicationSyncState.ProtoReflect.Descriptor instead.
func (*ApplicationSyncState) Descriptor() ([]byte, []int) {
	return file_pkg_model_application_proto_rawDescGZIP(), []int{1}
}

func (x *ApplicationSyncState) GetStatus() ApplicationSyncStatus {
	if x != nil {
		return x.Status
	}
	return ApplicationSyncStatus_UNKNOWN
}

func (x *ApplicationSyncState) GetShortReason() string {
	if x != nil {
		return x.ShortReason
	}
	return ""
}

func (x *ApplicationSyncState) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *ApplicationSyncState) GetHeadDeploymentId() string {
	if x != nil {
		return x.HeadDeploymentId
	}
	return ""
}

func (x *ApplicationSyncState) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type ApplicationDeploymentReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeploymentId   string             `protobuf:"bytes,1,opt,name=deployment_id,json=deploymentId,proto3" json:"deployment_id,omitempty"`
	Trigger        *DeploymentTrigger `protobuf:"bytes,2,opt,name=trigger,proto3" json:"trigger,omitempty"`
	Summary        string             `protobuf:"bytes,3,opt,name=summary,proto3" json:"summary,omitempty"`
	Version        string             `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
	ConfigFilename string             `protobuf:"bytes,5,opt,name=config_filename,json=configFilename,proto3" json:"config_filename,omitempty"`
	Versions       []*ArtifactVersion `protobuf:"bytes,6,rep,name=versions,proto3" json:"versions,omitempty"`
	StartedAt      int64              `protobuf:"varint,14,opt,name=started_at,json=startedAt,proto3" json:"started_at,omitempty"`
	CompletedAt    int64              `protobuf:"varint,15,opt,name=completed_at,json=completedAt,proto3" json:"completed_at,omitempty"`
}

func (x *ApplicationDeploymentReference) Reset() {
	*x = ApplicationDeploymentReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_model_application_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplicationDeploymentReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplicationDeploymentReference) ProtoMessage() {}

func (x *ApplicationDeploymentReference) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_model_application_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplicationDeploymentReference.ProtoReflect.Descriptor instead.
func (*ApplicationDeploymentReference) Descriptor() ([]byte, []int) {
	return file_pkg_model_application_proto_rawDescGZIP(), []int{2}
}

func (x *ApplicationDeploymentReference) GetDeploymentId() string {
	if x != nil {
		return x.DeploymentId
	}
	return ""
}

func (x *ApplicationDeploymentReference) GetTrigger() *DeploymentTrigger {
	if x != nil {
		return x.Trigger
	}
	return nil
}

func (x *ApplicationDeploymentReference) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *ApplicationDeploymentReference) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *ApplicationDeploymentReference) GetConfigFilename() string {
	if x != nil {
		return x.ConfigFilename
	}
	return ""
}

func (x *ApplicationDeploymentReference) GetVersions() []*ArtifactVersion {
	if x != nil {
		return x.Versions
	}
	return nil
}

func (x *ApplicationDeploymentReference) GetStartedAt() int64 {
	if x != nil {
		return x.StartedAt
	}
	return 0
}

func (x *ApplicationDeploymentReference) GetCompletedAt() int64 {
	if x != nil {
		return x.CompletedAt
	}
	return 0
}

var File_pkg_model_application_proto protoreflect.FileDescriptor

var file_pkg_model_application_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x70,
	0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2f, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xc2, 0x07, 0x0a, 0x0b, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x17, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10,
	0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x08, 0x70, 0x69, 0x70, 0x65, 0x64,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02,
	0x10, 0x01, 0x52, 0x07, 0x70, 0x69, 0x70, 0x65, 0x64, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0a, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x16, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x69, 0x6e, 0x64, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x82, 0x01,
	0x02, 0x10, 0x01, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x3e, 0x0a, 0x08, 0x67, 0x69, 0x74,
	0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x47,
	0x69, 0x74, 0x50, 0x61, 0x74, 0x68, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01,
	0x52, 0x07, 0x67, 0x69, 0x74, 0x50, 0x61, 0x74, 0x68, 0x12, 0x2e, 0x0a, 0x0e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x0d, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x36, 0x0a, 0x06, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x12, 0x74, 0x0a, 0x23, 0x6d, 0x6f, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x63, 0x65,
	0x6e, 0x74, 0x6c, 0x79, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x5f,
	0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x25, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x20, 0x6d, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x63,
	0x65, 0x6e, 0x74, 0x6c, 0x79, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x72, 0x0a, 0x22, 0x6d, 0x6f, 0x73,
	0x74, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x6e, 0x74, 0x6c, 0x79, 0x5f, 0x74, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x65, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x1f, 0x6d, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x6e, 0x74, 0x6c, 0x79, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x65, 0x64, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x3a, 0x0a,
	0x0a, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x09,
	0x73, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x65, 0x70,
	0x6c, 0x6f, 0x79, 0x69, 0x6e, 0x67, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x64, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x69, 0x6e, 0x67, 0x12, 0x26, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x62, 0x20, 0x01, 0x28, 0x03, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x22, 0x02, 0x28, 0x00, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x63, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73,
	0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x64, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x73,
	0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x26, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x65, 0x20, 0x01, 0x28, 0x03, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x22, 0x02,
	0x20, 0x00, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x26, 0x0a,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x66, 0x20, 0x01, 0x28,
	0x03, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x22, 0x02, 0x20, 0x00, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x22, 0xe6, 0x01, 0x0a, 0x14, 0x41, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x3e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1c, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x08, 0xfa,
	0x42, 0x05, 0x82, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x21, 0x0a, 0x0c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x12, 0x68, 0x65,
	0x61, 0x64, 0x5f, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x68, 0x65, 0x61, 0x64, 0x44, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x22, 0x02, 0x20, 0x00, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22,
	0xf1, 0x02, 0x0a, 0x1e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x12, 0x2c, 0x0a, 0x0d, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02,
	0x10, 0x01, 0x52, 0x0c, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x3c, 0x0a, 0x07, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x42, 0x08, 0xfa, 0x42, 0x05,
	0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x07, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x26, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x03, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x22, 0x02, 0x20, 0x00, 0x52, 0x09, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2a, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x03, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x22, 0x02, 0x28, 0x00, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x2a, 0x50, 0x0a, 0x15, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x59, 0x4e,
	0x43, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x49,
	0x4e, 0x47, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x4f, 0x55, 0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x53,
	0x59, 0x4e, 0x43, 0x10, 0x03, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x2d, 0x63, 0x64, 0x2f, 0x70, 0x69, 0x70, 0x65,
	0x63, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_model_application_proto_rawDescOnce sync.Once
	file_pkg_model_application_proto_rawDescData = file_pkg_model_application_proto_rawDesc
)

func file_pkg_model_application_proto_rawDescGZIP() []byte {
	file_pkg_model_application_proto_rawDescOnce.Do(func() {
		file_pkg_model_application_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_model_application_proto_rawDescData)
	})
	return file_pkg_model_application_proto_rawDescData
}

var file_pkg_model_application_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_model_application_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_model_application_proto_goTypes = []interface{}{
	(ApplicationSyncStatus)(0),             // 0: model.ApplicationSyncStatus
	(*Application)(nil),                    // 1: model.Application
	(*ApplicationSyncState)(nil),           // 2: model.ApplicationSyncState
	(*ApplicationDeploymentReference)(nil), // 3: model.ApplicationDeploymentReference
	nil,                                    // 4: model.Application.LabelsEntry
	(ApplicationKind)(0),                   // 5: model.ApplicationKind
	(*ApplicationGitPath)(nil),             // 6: model.ApplicationGitPath
	(*DeploymentTrigger)(nil),              // 7: model.DeploymentTrigger
	(*ArtifactVersion)(nil),                // 8: model.ArtifactVersion
}
var file_pkg_model_application_proto_depIdxs = []int32{
	5, // 0: model.Application.kind:type_name -> model.ApplicationKind
	6, // 1: model.Application.git_path:type_name -> model.ApplicationGitPath
	4, // 2: model.Application.labels:type_name -> model.Application.LabelsEntry
	3, // 3: model.Application.most_recently_successful_deployment:type_name -> model.ApplicationDeploymentReference
	3, // 4: model.Application.most_recently_triggered_deployment:type_name -> model.ApplicationDeploymentReference
	2, // 5: model.Application.sync_state:type_name -> model.ApplicationSyncState
	0, // 6: model.ApplicationSyncState.status:type_name -> model.ApplicationSyncStatus
	7, // 7: model.ApplicationDeploymentReference.trigger:type_name -> model.DeploymentTrigger
	8, // 8: model.ApplicationDeploymentReference.versions:type_name -> model.ArtifactVersion
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_pkg_model_application_proto_init() }
func file_pkg_model_application_proto_init() {
	if File_pkg_model_application_proto != nil {
		return
	}
	file_pkg_model_common_proto_init()
	file_pkg_model_deployment_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pkg_model_application_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Application); i {
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
		file_pkg_model_application_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplicationSyncState); i {
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
		file_pkg_model_application_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplicationDeploymentReference); i {
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
			RawDescriptor: file_pkg_model_application_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_model_application_proto_goTypes,
		DependencyIndexes: file_pkg_model_application_proto_depIdxs,
		EnumInfos:         file_pkg_model_application_proto_enumTypes,
		MessageInfos:      file_pkg_model_application_proto_msgTypes,
	}.Build()
	File_pkg_model_application_proto = out.File
	file_pkg_model_application_proto_rawDesc = nil
	file_pkg_model_application_proto_goTypes = nil
	file_pkg_model_application_proto_depIdxs = nil
}
