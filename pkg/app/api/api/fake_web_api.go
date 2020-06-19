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

package api

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/app/api/service/webservice"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	fakeProjectID = "debug-project"
)

// FakeWebAPI implements the fake behaviors for the gRPC definitions of WebAPI.
type FakeWebAPI struct {
}

// NewFakeWebAPI creates a new FakeWebAPI instance.
func NewFakeWebAPI() *FakeWebAPI {
	return &FakeWebAPI{}
}

// Register registers all handling of this service into the specified gRPC server.
func (a *FakeWebAPI) Register(server *grpc.Server) {
	webservice.RegisterWebServiceServer(server, a)
}

func (a *FakeWebAPI) AddEnvironment(ctx context.Context, req *webservice.AddEnvironmentRequest) (*webservice.AddEnvironmentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *FakeWebAPI) UpdateEnvironmentDesc(ctx context.Context, req *webservice.UpdateEnvironmentDescRequest) (*webservice.UpdateEnvironmentDescResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *FakeWebAPI) ListEnvironments(ctx context.Context, req *webservice.ListEnvironmentsRequest) (*webservice.ListEnvironmentsResponse, error) {
	now := time.Now()
	envs := []*model.Environment{
		{
			Id:        fmt.Sprintf("%s:%s", fakeProjectID, "development"),
			Name:      "development",
			Desc:      "For development",
			ProjectId: fakeProjectID,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		},
		{
			Id:        fmt.Sprintf("%s:%s", fakeProjectID, "staging"),
			Name:      "staging",
			Desc:      "For staging",
			ProjectId: fakeProjectID,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		},
		{
			Id:        fmt.Sprintf("%s:%s", fakeProjectID, "production"),
			Name:      "production",
			Desc:      "For production",
			ProjectId: fakeProjectID,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		},
	}

	return &webservice.ListEnvironmentsResponse{
		Environments: envs,
	}, nil
}

func (a *FakeWebAPI) RegisterPiped(ctx context.Context, req *webservice.RegisterPipedRequest) (*webservice.RegisterPipedResponse, error) {
	return &webservice.RegisterPipedResponse{
		Id:  "e357d99f-0f83-4ce0-8c8b-27f11f432ef9",
		Key: "9bf9752a-54a2-451a-a541-444add56f96b",
	}, nil
}

func (a *FakeWebAPI) DisablePiped(ctx context.Context, req *webservice.DisablePipedRequest) (*webservice.DisablePipedResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *FakeWebAPI) ListPipeds(ctx context.Context, req *webservice.ListPipedsRequest) (*webservice.ListPipedsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *FakeWebAPI) AddApplication(ctx context.Context, req *webservice.AddApplicationRequest) (*webservice.AddApplicationResponse, error) {
	return &webservice.AddApplicationResponse{}, nil
}

func (a *FakeWebAPI) DisableApplication(ctx context.Context, req *webservice.DisableApplicationRequest) (*webservice.DisableApplicationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *FakeWebAPI) ListApplications(ctx context.Context, req *webservice.ListApplicationsRequest) (*webservice.ListApplicationsResponse, error) {
	now := time.Now()
	fakeApplications := []*model.Application{
		{
			Id:        fmt.Sprintf("%s:%s:%s", fakeProjectID, "development", "debug-app"),
			Name:      "debug-app",
			EnvId:     fmt.Sprintf("%s:%s", fakeProjectID, "development"),
			PipedId:   "debug-piped",
			ProjectId: fakeProjectID,
			Kind:      model.ApplicationKind_KUBERNETES,
			GitPath: &model.ApplicationGitPath{
				RepoId: "debug",
				Path:   "k8s",
			},
			CloudProvider: "kubernetes-default",
			MostRecentSuccessfulDeployment: &model.ApplicationCompletedDeployment{
				DeploymentId: "debug-deployment-id-01",
				CommitHash:   "3808585b46f1e90196d7ffe8dd04c807a251febc",
				Version:      "v0.1.0",
				StartedAt:    now.Add(-3 * 24 * time.Hour).Unix(),
				CompletedAt:  now.Add(-3 * 24 * time.Hour).Unix(),
			},
			SyncState: &model.ApplicationSyncState{
				Status:           model.ApplicationSyncStatus_SYNCED,
				ShortReason:      "Short resson",
				Reason:           "Reason",
				HeadDeploymentId: "debug-deployment-id-01",
				Timestamp:        now.Unix(),
			},
			Disabled:  false,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
		},
	}
	return &webservice.ListApplicationsResponse{
		Applications: fakeApplications,
	}, nil
}

func (a *FakeWebAPI) SyncApplication(ctx context.Context, req *webservice.SyncApplicationRequest) (*webservice.SyncApplicationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *FakeWebAPI) GetApplication(ctx context.Context, req *webservice.GetApplicationRequest) (*webservice.GetApplicationResponse, error) {
	now := time.Now()
	application := model.Application{
		Id:        fmt.Sprintf("%s:%s:%s", fakeProjectID, "development", "debug-app"),
		Name:      "debug-app",
		EnvId:     fmt.Sprintf("%s:%s", fakeProjectID, "development"),
		PipedId:   "debug-piped",
		ProjectId: fakeProjectID,
		Kind:      model.ApplicationKind_KUBERNETES,
		GitPath: &model.ApplicationGitPath{
			RepoId: "debug",
			Path:   "k8s",
		},
		CloudProvider: "kubernetes-default",
		MostRecentSuccessfulDeployment: &model.ApplicationCompletedDeployment{
			DeploymentId: "debug-deployment-id-01",
			CommitHash:   "3808585b46f1e90196d7ffe8dd04c807a251febc",
			Version:      "v0.1.0",
			StartedAt:    now.Add(-3 * 24 * time.Hour).Unix(),
			CompletedAt:  now.Add(-3 * 24 * time.Hour).Unix(),
		},
		SyncState: &model.ApplicationSyncState{
			Status:           model.ApplicationSyncStatus_SYNCED,
			ShortReason:      "Short resson",
			Reason:           "Reason",
			HeadDeploymentId: "debug-deployment-id-01",
			Timestamp:        now.Unix(),
		},
		Disabled:  false,
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}
	return &webservice.GetApplicationResponse{
		Application: &application,
	}, nil
}

func (a *FakeWebAPI) ListDeployments(ctx context.Context, req *webservice.ListDeploymentsRequest) (*webservice.ListDeploymentsResponse, error) {
	now := time.Now()
	deploymentTime := now
	fakeDeployments := make([]*model.Deployment, 15)
	for i := 0; i < 15; i++ {
		// 5 hour intervals
		deploymentTime := deploymentTime.Add(time.Duration(-5*i) * time.Hour)
		fakeDeployments[i] = &model.Deployment{
			Id:            fmt.Sprintf("debug-deployment-id-%02d", i),
			ApplicationId: fmt.Sprintf("%s:%s:%s", fakeProjectID, "development", "debug-app"),
			EnvId:         fmt.Sprintf("%s:%s", fakeProjectID, "development"),
			PipedId:       "debug-piped",
			ProjectId:     fakeProjectID,
			GitPath: &model.ApplicationGitPath{
				RepoId: "debug",
				Path:   "k8s",
			},
			Trigger: &model.DeploymentTrigger{
				Commit: &model.Commit{
					Hash:      "3808585b46f1e90196d7ffe8dd04c807a251febc",
					Message:   "Add web page routing (#133)",
					Author:    "cakecatz",
					Branch:    "master",
					CreatedAt: deploymentTime.Unix(),
				},
				Commander: "",
				Timestamp: deploymentTime.Unix(),
			},
			RunningCommitHash: "3808585b46f1e90196d7ffe8dd04c807a251febc",
			Description:       fmt.Sprintf("This deployment is debug-%02d", i),
			Status:            model.DeploymentStatus_DEPLOYMENT_SUCCESS,
			Stages: []*model.PipelineStage{
				{
					Id:           "fake-stage-id-0-0",
					Name:         model.StageK8sCanaryRollout.String(),
					Index:        0,
					Predefined:   true,
					Status:       model.StageStatus_STAGE_SUCCESS,
					RetriedCount: 0,
					CompletedAt:  now.Unix(),
					CreatedAt:    now.Unix(),
					UpdatedAt:    now.Unix(),
				},
				{
					Id:         "fake-stage-id-1-0",
					Name:       model.StageK8sCanaryRollout.String(),
					Index:      0,
					Predefined: true,
					Requires: []string{
						"fake-stage-id-0-0",
					},
					Status:       model.StageStatus_STAGE_RUNNING,
					RetriedCount: 0,
					CompletedAt:  0,
					CreatedAt:    now.Unix(),
					UpdatedAt:    now.Unix(),
				},
				{
					Id:         "fake-stage-id-1-1",
					Name:       model.StageK8sPrimaryUpdate.String(),
					Index:      1,
					Predefined: true,
					Requires: []string{
						"fake-stage-id-0-0",
					},
					Status:       model.StageStatus_STAGE_SUCCESS,
					RetriedCount: 0,
					CompletedAt:  now.Unix(),
					CreatedAt:    now.Unix(),
					UpdatedAt:    now.Unix(),
				},
				{
					Id:         "fake-stage-id-1-2",
					Name:       model.StageK8sCanaryRollout.String(),
					Index:      2,
					Predefined: true,
					Requires: []string{
						"fake-stage-id-0-0",
					},
					Status:       model.StageStatus_STAGE_FAILURE,
					RetriedCount: 0,
					CompletedAt:  now.Unix(),
					CreatedAt:    now.Unix(),
					UpdatedAt:    now.Unix(),
				},
				{
					Id:         "fake-stage-id-2-0",
					Name:       model.StageK8sCanaryClean.String(),
					Desc:       "waiting approval",
					Index:      0,
					Predefined: true,
					Requires: []string{
						"fake-stage-id-1-0",
						"fake-stage-id-1-1",
						"fake-stage-id-1-2",
					},
					Status:       model.StageStatus_STAGE_NOT_STARTED_YET,
					RetriedCount: 0,
					CompletedAt:  0,
					CreatedAt:    now.Unix(),
					UpdatedAt:    now.Unix(),
				},
				{
					Id:         "fake-stage-id-2-1",
					Name:       model.StageK8sCanaryClean.String(),
					Desc:       "approved by cakecatz",
					Index:      1,
					Predefined: true,
					Requires: []string{
						"fake-stage-id-1-0",
						"fake-stage-id-1-1",
						"fake-stage-id-1-2",
					},
					Status:       model.StageStatus_STAGE_NOT_STARTED_YET,
					RetriedCount: 0,
					CompletedAt:  0,
					CreatedAt:    now.Unix(),
					UpdatedAt:    now.Unix(),
				},
				{
					Id:         "fake-stage-id-3-0",
					Name:       model.StageK8sCanaryRollout.String(),
					Index:      0,
					Predefined: true,
					Requires: []string{
						"fake-stage-id-2-0",
						"fake-stage-id-2-1",
					},
					Status:       model.StageStatus_STAGE_NOT_STARTED_YET,
					RetriedCount: 0,
					CompletedAt:  0,
					CreatedAt:    now.Unix(),
					UpdatedAt:    now.Unix(),
				},
			},
			CreatedAt: deploymentTime.Unix(),
			UpdatedAt: deploymentTime.Unix(),
		}
	}
	return &webservice.ListDeploymentsResponse{
		Deployments: fakeDeployments,
	}, nil
}

func (a *FakeWebAPI) GetDeployment(ctx context.Context, req *webservice.GetDeploymentRequest) (*webservice.GetDeploymentResponse, error) {
	now := time.Now()
	resp := &model.Deployment{
		Id:            "debug-deployment-id-01",
		ApplicationId: fmt.Sprintf("%s:%s:%s", fakeProjectID, "development", "debug-app"),
		EnvId:         fmt.Sprintf("%s:%s", fakeProjectID, "development"),
		PipedId:       "debug-piped",
		ProjectId:     fakeProjectID,
		Kind:          model.ApplicationKind_KUBERNETES,
		GitPath: &model.ApplicationGitPath{
			RepoId: "debug",
			Path:   "k8s",
		},
		Trigger: &model.DeploymentTrigger{
			Commit: &model.Commit{
				Hash:      "3808585b46f1e90196d7ffe8dd04c807a251febc",
				Message:   "Add web page routing (#133)",
				Author:    "cakecatz",
				Branch:    "master",
				CreatedAt: now.Add(-30 * time.Minute).Unix(),
			},
			Commander: "cakecatz",
			Timestamp: now.Add(-30 * time.Minute).Unix(),
		},
		RunningCommitHash: "3808585b46f1e90196d7ffe8dd04c807a251febc",
		Description:       "This deployment is debug",
		Status:            model.DeploymentStatus_DEPLOYMENT_RUNNING,
		Stages: []*model.PipelineStage{
			{
				Id:           "fake-stage-id-0-0",
				Name:         model.StageK8sCanaryRollout.String(),
				Index:        0,
				Predefined:   true,
				Status:       model.StageStatus_STAGE_SUCCESS,
				RetriedCount: 0,
				CompletedAt:  now.Unix(),
				CreatedAt:    now.Unix(),
				UpdatedAt:    now.Unix(),
			},
			{
				Id:         "fake-stage-id-1-0",
				Name:       model.StageK8sCanaryRollout.String(),
				Index:      0,
				Predefined: true,
				Requires: []string{
					"fake-stage-id-0-0",
				},
				Status:       model.StageStatus_STAGE_RUNNING,
				RetriedCount: 0,
				CompletedAt:  0,
				CreatedAt:    now.Unix(),
				UpdatedAt:    now.Unix(),
			},
			{
				Id:         "fake-stage-id-1-1",
				Name:       model.StageK8sPrimaryUpdate.String(),
				Index:      1,
				Predefined: true,
				Requires: []string{
					"fake-stage-id-0-0",
				},
				Status:       model.StageStatus_STAGE_SUCCESS,
				RetriedCount: 0,
				CompletedAt:  now.Unix(),
				CreatedAt:    now.Unix(),
				UpdatedAt:    now.Unix(),
			},
			{
				Id:         "fake-stage-id-1-2",
				Name:       model.StageK8sCanaryRollout.String(),
				Index:      2,
				Predefined: true,
				Requires: []string{
					"fake-stage-id-0-0",
				},
				Status:       model.StageStatus_STAGE_FAILURE,
				RetriedCount: 0,
				CompletedAt:  now.Unix(),
				CreatedAt:    now.Unix(),
				UpdatedAt:    now.Unix(),
			},
			{
				Id:         "fake-stage-id-2-0",
				Name:       model.StageK8sCanaryClean.String(),
				Desc:       "waiting approval",
				Index:      0,
				Predefined: true,
				Requires: []string{
					"fake-stage-id-1-0",
					"fake-stage-id-1-1",
					"fake-stage-id-1-2",
				},
				Status:       model.StageStatus_STAGE_NOT_STARTED_YET,
				RetriedCount: 0,
				CompletedAt:  0,
				CreatedAt:    now.Unix(),
				UpdatedAt:    now.Unix(),
			},
			{
				Id:         "fake-stage-id-2-1",
				Name:       model.StageK8sCanaryClean.String(),
				Desc:       "approved by cakecatz",
				Index:      1,
				Predefined: true,
				Requires: []string{
					"fake-stage-id-1-0",
					"fake-stage-id-1-1",
					"fake-stage-id-1-2",
				},
				Status:       model.StageStatus_STAGE_NOT_STARTED_YET,
				RetriedCount: 0,
				CompletedAt:  0,
				CreatedAt:    now.Unix(),
				UpdatedAt:    now.Unix(),
			},
			{
				Id:         "fake-stage-id-3-0",
				Name:       model.StageK8sCanaryRollout.String(),
				Index:      0,
				Predefined: true,
				Requires: []string{
					"fake-stage-id-2-0",
					"fake-stage-id-2-1",
				},
				Status:       model.StageStatus_STAGE_NOT_STARTED_YET,
				RetriedCount: 0,
				CompletedAt:  0,
				CreatedAt:    now.Unix(),
				UpdatedAt:    now.Unix(),
			},
		},
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	}

	return &webservice.GetDeploymentResponse{
		Deployment: resp,
	}, nil
}

func (a *FakeWebAPI) GetStageLog(ctx context.Context, req *webservice.GetStageLogRequest) (*webservice.GetStageLogResponse, error) {
	startTime := time.Now().Add(-10 * time.Minute)
	resp := []*model.LogBlock{
		{
			Index:     1,
			Log:       "+ make build",
			Severity:  model.LogSeverity_INFO,
			CreatedAt: startTime.Unix(),
		},
		{
			Index:     2,
			Log:       "bazelisk  --output_base=/workspace/bazel_out build  --config=ci -- //...",
			Severity:  model.LogSeverity_INFO,
			CreatedAt: startTime.Add(5 * time.Second).Unix(),
		},
		{
			Index:     3,
			Log:       "2020/06/01 08:52:07 Downloading https://releases.bazel.build/3.1.0/release/bazel-3.1.0-linux-x86_64...",
			Severity:  model.LogSeverity_INFO,
			CreatedAt: startTime.Add(10 * time.Second).Unix(),
		},
		{
			Index:     4,
			Log:       "Extracting Bazel installation...",
			Severity:  model.LogSeverity_INFO,
			CreatedAt: startTime.Add(15 * time.Second).Unix(),
		},
		{
			Index:     5,
			Log:       "Starting local Bazel server and connecting to it...",
			Severity:  model.LogSeverity_INFO,
			CreatedAt: startTime.Add(20 * time.Second).Unix(),
		},
		{
			Index:     6,
			Log:       "(08:52:14) Loading: 0 packages loaded",
			Severity:  model.LogSeverity_SUCCESS,
			CreatedAt: startTime.Add(30 * time.Second).Unix(),
		},
		{
			Index:     7,
			Log:       "(08:53:21) Analyzing: 157 targets (88 packages loaded, 0 targets configured)",
			Severity:  model.LogSeverity_SUCCESS,
			CreatedAt: startTime.Add(35 * time.Second).Unix(),
		},
		{
			Index:     8,
			Log:       "Error: Error building: logged 2 error(s)",
			Severity:  model.LogSeverity_ERROR,
			CreatedAt: startTime.Add(45 * time.Second).Unix(),
		},
	}

	return &webservice.GetStageLogResponse{
		Blocks: resp,
	}, nil
}

func (a *FakeWebAPI) CancelDeployment(ctx context.Context, req *webservice.CancelDeploymentRequest) (*webservice.CancelDeploymentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *FakeWebAPI) ApproveStage(ctx context.Context, req *webservice.ApproveStageRequest) (*webservice.ApproveStageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *FakeWebAPI) GetApplicationLiveState(ctx context.Context, req *webservice.GetApplicationLiveStateRequest) (*webservice.GetApplicationLiveStateResponse, error) {
	now := time.Now()
	snapshot := &model.ApplicationLiveStateSnapshot{
		ApplicationId: fmt.Sprintf("%s:%s:%s", fakeProjectID, "development", "debug-app"),
		EnvId:         fmt.Sprintf("%s:%s", fakeProjectID, "development"),
		PipedId:       "debug-piped",
		ProjectId:     fakeProjectID,
		Kind:          model.ApplicationKind_KUBERNETES,
		Kubernetes: &model.KubernetesApplicationLiveState{
			Resources: []*model.KubernetesResourceState{
				{
					Id:         "f2c832a3-1f5b-4982-8f6e-72345ecb3c82",
					Name:       "demo-application",
					ApiVersion: "networking.k8s.io/v1beta1",
					Kind:       "Ingress",
					Namespace:  "default",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
				{
					Id:         "8423fb53-5170-4864-a7d2-b84f8d36cb02",
					Name:       "demo-application",
					ApiVersion: "v1",
					Kind:       "Service",
					Namespace:  "default",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
				{
					Id:         "660ecdfd-307b-4e47-becd-1fde4e0c1e7a",
					Name:       "demo-application",
					ApiVersion: "apps/v1",
					Kind:       "Deployment",
					Namespace:  "default",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
				{
					Id: "8621f186-6641-4f7a-9be4-5983eb647f8d",
					OwnerIds: []string{
						"660ecdfd-307b-4e47-becd-1fde4e0c1e7a",
					},
					ParentIds: []string{
						"660ecdfd-307b-4e47-becd-1fde4e0c1e7a",
					},
					Name:       "demo-application-9504e8601a",
					ApiVersion: "apps/v1",
					Kind:       "ReplicaSet",
					Namespace:  "default",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
				{
					Id: "ae5d0031-1f63-4396-b929-fa9987d1e6de",
					OwnerIds: []string{
						"660ecdfd-307b-4e47-becd-1fde4e0c1e7a",
					},
					ParentIds: []string{
						"8621f186-6641-4f7a-9be4-5983eb647f8d",
					},
					Name:       "demo-application-9504e8601a-7vrdw",
					ApiVersion: "v1",
					Kind:       "Pod",
					Namespace:  "default",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
				{
					Id: "f55c7891-ba25-44bb-bca4-ffbc16b0089f",
					OwnerIds: []string{
						"660ecdfd-307b-4e47-becd-1fde4e0c1e7a",
					},
					ParentIds: []string{
						"8621f186-6641-4f7a-9be4-5983eb647f8d",
					},
					Name:       "demo-application-9504e8601a-vlgd5",
					ApiVersion: "v1",
					Kind:       "Pod",
					Namespace:  "default",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
				{
					Id: "c2a81415-5bbf-44e8-9101-98bbd636bbeb",
					OwnerIds: []string{
						"660ecdfd-307b-4e47-becd-1fde4e0c1e7a",
					},
					ParentIds: []string{
						"8621f186-6641-4f7a-9be4-5983eb647f8d",
					},
					Name:       "demo-application-9504e8601a-tmwp5",
					ApiVersion: "v1",
					Kind:       "Pod",
					Namespace:  "default",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				},
			},
		},
		Version: &model.ApplicationLiveStateVersion{
			Index:     1,
			Timestamp: now.Unix(),
		},
	}
	return &webservice.GetApplicationLiveStateResponse{
		Snapshot: snapshot,
	}, nil
}

func (a *FakeWebAPI) GetProject(ctx context.Context, req *webservice.GetProjectRequest) (*webservice.GetProjectResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (a *FakeWebAPI) GetMe(ctx context.Context, req *webservice.GetMeRequest) (*webservice.GetMeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
