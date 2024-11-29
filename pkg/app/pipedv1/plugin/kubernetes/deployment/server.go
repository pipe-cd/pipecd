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

package deployment

import (
	"context"
	"time"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/deployment"
	"github.com/pipe-cd/pipecd/pkg/plugin/logpersister"
	"github.com/pipe-cd/pipecd/pkg/regexpool"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type toolRegistry interface {
	InstallTool(ctx context.Context, name, version string) (path string, err error)
}

type loader interface {
	// LoadManifests renders and loads all manifests for application.
	LoadManifests(ctx context.Context, input provider.LoaderInput) ([]provider.Manifest, error)
}

type applier interface {
	// ApplyManifest does applying the given manifest.
	ApplyManifest(ctx context.Context, manifest provider.Manifest) error
	// CreateManifest does creating resource from given manifest.
	CreateManifest(ctx context.Context, manifest provider.Manifest) error
	// ReplaceManifest does replacing resource from given manifest.
	ReplaceManifest(ctx context.Context, manifest provider.Manifest) error
	// ForceReplaceManifest does force replacing resource from given manifest.
	ForceReplaceManifest(ctx context.Context, manifest provider.Manifest) error
}

type logPersister interface {
	StageLogPersister(deploymentID, stageID string) logpersister.StageLogPersister
}

type DeploymentService struct {
	deployment.UnimplementedDeploymentServiceServer

	RegexPool    *regexpool.Pool
	Logger       *zap.Logger
	ToolRegistry toolRegistry
	Loader       loader
	LogPersister logPersister
}

// NewDeploymentService creates a new planService.
func NewDeploymentService(
	logger *zap.Logger,
) *DeploymentService {
	return &DeploymentService{
		RegexPool:    regexpool.DefaultPool(),
		Logger:       logger.Named("planner"),
		ToolRegistry: nil, // TODO: set the tool registry
	}
}

// Register registers all handling of this service into the specified gRPC server.
func (a *DeploymentService) Register(server *grpc.Server) {
	deployment.RegisterDeploymentServiceServer(server, a)
}

// DetermineStrategy implements deployment.DeploymentServiceServer.
func (a *DeploymentService) DetermineStrategy(ctx context.Context, request *deployment.DetermineStrategyRequest) (*deployment.DetermineStrategyResponse, error) {
	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](request.GetInput().GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	runnings, err := a.loadManifests(ctx, request.GetInput().GetDeployment(), cfg.Spec, request.GetInput().GetRunningDeploymentSource())

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	targets, err := a.loadManifests(ctx, request.GetInput().GetDeployment(), cfg.Spec, request.GetInput().GetTargetDeploymentSource())

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	strategy, summary := determineStrategy(runnings, targets, cfg.Spec.Workloads, a.Logger)

	return &deployment.DetermineStrategyResponse{
		SyncStrategy: strategy,
		Summary:      summary,
	}, nil

}

// DetermineVersions implements deployment.DeploymentServiceServer.
func (a *DeploymentService) DetermineVersions(ctx context.Context, request *deployment.DetermineVersionsRequest) (*deployment.DetermineVersionsResponse, error) {
	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](request.GetInput().GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	manifests, err := a.loadManifests(ctx, request.GetInput().GetDeployment(), cfg.Spec, request.GetInput().GetTargetDeploymentSource())

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	versions, err := determineVersions(manifests)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &deployment.DetermineVersionsResponse{
		Versions: versions,
	}, nil
}

// BuildPipelineSyncStages implements deployment.DeploymentServiceServer.
func (a *DeploymentService) BuildPipelineSyncStages(ctx context.Context, request *deployment.BuildPipelineSyncStagesRequest) (*deployment.BuildPipelineSyncStagesResponse, error) {
	now := time.Now()
	stages := buildPipelineStages(request.GetStages(), request.GetRollback(), now)
	return &deployment.BuildPipelineSyncStagesResponse{
		Stages: stages,
	}, nil
}

// BuildQuickSyncStages implements deployment.DeploymentServiceServer.
func (a *DeploymentService) BuildQuickSyncStages(ctx context.Context, request *deployment.BuildQuickSyncStagesRequest) (*deployment.BuildQuickSyncStagesResponse, error) {
	now := time.Now()
	stages := buildQuickSyncPipeline(request.GetRollback(), now)
	return &deployment.BuildQuickSyncStagesResponse{
		Stages: stages,
	}, nil
}

// FetchDefinedStages implements deployment.DeploymentServiceServer.
func (a *DeploymentService) FetchDefinedStages(context.Context, *deployment.FetchDefinedStagesRequest) (*deployment.FetchDefinedStagesResponse, error) {
	stages := make([]string, 0, len(AllStages))
	for _, s := range AllStages {
		stages = append(stages, string(s))
	}

	return &deployment.FetchDefinedStagesResponse{
		Stages: stages,
	}, nil
}

func (a *DeploymentService) loadManifests(ctx context.Context, deploy *model.Deployment, spec *kubeconfig.KubernetesApplicationSpec, deploymentSource *deployment.DeploymentSource) ([]provider.Manifest, error) {
	manifests, err := a.Loader.LoadManifests(ctx, provider.LoaderInput{
		PipedID:          deploy.GetPipedId(),
		AppID:            deploy.GetApplicationId(),
		CommitHash:       deploy.GetTrigger().GetCommit().GetHash(),
		AppName:          deploy.GetApplicationName(),
		AppDir:           deploymentSource.GetApplicationDirectory(),
		ConfigFilename:   deploymentSource.GetApplicationConfigFilename(),
		Manifests:        spec.Input.Manifests,
		Namespace:        spec.Input.Namespace,
		TemplatingMethod: provider.TemplatingMethodNone, // TODO: Implement detection of templating method or add it to the config spec.

		// TODO: Define other fields for LoaderInput
	})

	if err != nil {
		return nil, err
	}

	return manifests, nil
}

func (a *DeploymentService) ExecuteStage(ctx context.Context, request *deployment.ExecuteStageRequest) (*deployment.ExecuteStageResponse, error) {
	switch request.GetInput().GetStage().GetName() {
	case StageK8sSync.String():
		return a.executeK8sSyncStage(ctx, request.GetInput())
	case StageK8sRollback.String():
		return a.executeK8sRollbackStage(ctx, request.GetInput())
	default:
		return nil, status.Error(codes.InvalidArgument, "unimplemented or unsupported stage")
	}
}

func (a *DeploymentService) executeK8sSyncStage(ctx context.Context, input *deployment.ExecutePluginInput) (*deployment.ExecuteStageResponse, error) {
	lp := a.LogPersister.StageLogPersister(input.GetDeployment().GetId(), input.GetStage().GetId())
	lp.Infof("Start syncing the deployment")

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.GetTargetDeploymentSource().GetApplicationConfig())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	lp.Infof("Loading manifests at commit %s for handling", input.GetDeployment().GetTrigger().GetCommit().GetHash())
	manifests, err := a.loadManifests(ctx, input.GetDeployment(), cfg.Spec, input.GetTargetDeploymentSource())
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	lp.Successf("Successfully loaded %d manifests", len(manifests))

	// Because the loaded manifests are read-only
	// we duplicate them to avoid updating the shared manifests data in cache.
	// TODO: implement duplicateManifests function

	// When addVariantLabelToSelector is true, ensure that all workloads
	// have the variant label in their selector.
	var (
		variantLabel   = cfg.Spec.VariantLabel.Key
		primaryVariant = cfg.Spec.VariantLabel.PrimaryValue
	)
	// TODO: handle cfg.Spec.QuickSync.AddVariantLabelToSelector

	// Add variant annotations to all manifests.
	for i := range manifests {
		manifests[i].AddAnnotations(map[string]string{
			variantLabel: primaryVariant,
		})
	}

	// TODO: implement annotateConfigHash to ensure restart of workloads when config changes

	// Get the applier for the target cluster.
	var applier applier // TODO: build applier from the plugin config

	// Start applying all manifests to add or update running resources.
	if err := applyManifests(ctx, applier, manifests, cfg.Spec.Input.Namespace, lp); err != nil {
		lp.Errorf("Failed while applying manifests (%v)", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: implement prune resources

	return &deployment.ExecuteStageResponse{
		Status: model.StageStatus_STAGE_SUCCESS,
	}, nil
}

func (a *DeploymentService) executeK8sRollbackStage(ctx context.Context, input *deployment.ExecutePluginInput) (*deployment.ExecuteStageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
