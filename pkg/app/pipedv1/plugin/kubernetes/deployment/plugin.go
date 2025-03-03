package deployment

import (
	"context"
	"errors"

	kubeconfig "github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	config "github.com/pipe-cd/pipecd/pkg/configv1"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

// Plugin implements the sdk.DeploymentPlugin interface.
type Plugin struct {
	loader loader
}

type loader interface {
	// LoadManifests renders and loads all manifests for application.
	LoadManifests(ctx context.Context, input provider.LoaderInput) ([]provider.Manifest, error)
}

var _ sdk.DeploymentPlugin[sdk.ConfigNone, kubeconfig.KubernetesDeployTargetConfig] = (*Plugin)(nil)

func (p *Plugin) Name() string {
	return "kubernetes"
}

func (p *Plugin) Version() string {
	return "0.0.1" // TODO
}

func (p *Plugin) FetchDefinedStages() []string {
	stages := make([]string, 0, len(AllStages))
	for _, s := range AllStages {
		stages = append(stages, string(s))
	}

	return stages
}

// FIXME
func (p *Plugin) BuildPipelineSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildPipelineSyncStagesInput) (*sdk.BuildPipelineSyncStagesResponse, error) {
	return nil, nil
}

// FIXME
func (p *Plugin) ExecuteStage(ctx context.Context, _ *sdk.ConfigNone, dts []*sdk.DeployTarget[kubeconfig.KubernetesDeployTargetConfig], input *sdk.ExecuteStageInput) (*sdk.ExecuteStageResponse, error) {
	switch input.Request.StageName {
	case StageK8sSync.String():
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sSyncStage(ctx, input),
		}, nil
	case StageK8sRollback.String():
		return &sdk.ExecuteStageResponse{
			Status: p.executeK8sRollbackStage(ctx, input),
		}, nil
	default:
		return nil, errors.New("unimplemented or unsupported stage")
	}
}

// FIXME
func (p *Plugin) executeK8sSyncStage(ctx context.Context, input *sdk.ExecuteStageInput) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start syncing the deployment")

	cfg, err := config.DecodeYAML[*kubeconfig.KubernetesApplicationSpec](input.Request.TargetDeploymentSource.ApplicationConfig)
	if err != nil {
		lp.Errorf("Failed while decoding application config (%v)", err)
		return sdk.StageStatusFailure
	}

	lp.Infof("Loading manifests at commit %s for handling", input.Request.TargetDeploymentSource.CommitHash)
	manifests, err := p.loadManifests(ctx, &input.Request.Deployment, cfg.Spec, &input.Request.TargetDeploymentSource)
	if err != nil {
		lp.Errorf("Failed while loading manifests (%v)", err)
		return sdk.StageStatusFailure
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
	// TODO: treat the stage options specified under "with"
	if cfg.Spec.QuickSync.AddVariantLabelToSelector {
		workloads := findWorkloadManifests(manifests, cfg.Spec.Workloads)
		for _, m := range workloads {
			if err := ensureVariantSelectorInWorkload(m, variantLabel, primaryVariant); err != nil {
				lp.Errorf("Unable to check/set %q in selector of workload %s (%v)", variantLabel+": "+primaryVariant, m.Key().ReadableString(), err)
				return sdk.StageStatusFailure
			}
		}
	}

	// Add variant annotations to all manifests.
	for i := range manifests {
		manifests[i].AddLabels(map[string]string{
			variantLabel: primaryVariant,
		})
		manifests[i].AddAnnotations(map[string]string{
			variantLabel: primaryVariant,
		})
	}

	return sdk.StageStatusFailure
}

func (p *Plugin) loadManifests(ctx context.Context, deploy *sdk.Deployment, spec *kubeconfig.KubernetesApplicationSpec, deploymentSource *sdk.DeploymentSource) ([]provider.Manifest, error) {
	manifests, err := p.loader.LoadManifests(ctx, provider.LoaderInput{
		PipedID:          deploy.PipedID,
		AppID:            deploy.ApplicationID,
		CommitHash:       deploymentSource.CommitHash,
		AppName:          deploy.ApplicationName,
		AppDir:           deploymentSource.ApplicationDirectory,
		ConfigFilename:   deploymentSource.ApplicationConfigFilename,
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

// FIXME
func (p *Plugin) executeK8sRollbackStage(ctx context.Context, input *sdk.ExecuteStageInput) sdk.StageStatus {
	return sdk.StageStatusFailure
}

// FIXME
func (p *Plugin) DetermineVersions(context.Context, *sdk.ConfigNone, *sdk.Client, sdk.TODO) (sdk.TODO, error) {
	return sdk.TODO{}, nil
}

// FIXME
func (p *Plugin) DetermineStrategy(context.Context, *sdk.ConfigNone, *sdk.Client, sdk.TODO) (sdk.TODO, error) {
	return sdk.TODO{}, nil
}

func (p *Plugin) BuildQuickSyncStages(ctx context.Context, _ *sdk.ConfigNone, input *sdk.BuildQuickSyncStagesInput) (*sdk.BuildQuickSyncStagesResponse, error) {
	return &sdk.BuildQuickSyncStagesResponse{
		Stages: BuildQuickSyncPipeline(input.Request.Rollback),
	}, nil
}
