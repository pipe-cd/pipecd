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

package planner

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin"
	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/platform"
	"github.com/pipe-cd/pipecd/pkg/regexpool"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type PlannerService struct {
	platform.UnimplementedPlannerServiceServer

	Decrypter secretDecrypter
	RegexPool *regexpool.Pool
	Logger    *zap.Logger
}

// Register registers all handling of this service into the specified gRPC server.
func (a *PlannerService) Register(server *grpc.Server) {
	platform.RegisterPlannerServiceServer(server, a)
}

// NewPlannerService creates a new planService.
func NewPlannerService(
	decrypter secretDecrypter,
	logger *zap.Logger,
) *PlannerService {
	return &PlannerService{
		Decrypter: decrypter,
		RegexPool: regexpool.DefaultPool(),
		Logger:    logger.Named("planner"),
	}
}

// type gitClient interface {
// 	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
// 	Clean() error
// }

// type secretDecrypter interface {
// 	Decrypt(string) (string, error)
// }

// const (
// 	versionUnknown = "unknown"
// )

func (ps *PlannerService) DetermineStrategy(ctx context.Context, in *platform.DetermineStrategyRequest) (*platform.DetermineStrategyResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (ps *PlannerService) QuickSyncPlan(ctx context.Context, in *platform.QuickSyncPlanRequest) (*platform.QuickSyncPlanResponse, error) {
	now := time.Now()

	cloner, err := plugin.GetPlanSourceCloner(in.GetInput())
	if err != nil {
		return nil, err
	}

	d, err := os.MkdirTemp("", "") // TODO
	if err != nil {
		return nil, fmt.Errorf("failed to prepare temporary directory (%w)", err)
	}
	defer os.RemoveAll(d)

	p := deploysource.NewProvider(
		d,
		cloner,
		*in.GetInput().GetDeployment().GetGitPath(),
		ps.Decrypter,
	)

	ds, err := p.Get(ctx, io.Discard /* TODO */)
	if err != nil {
		return nil, err
	}

	cfg := ds.ApplicationConfig.KubernetesApplicationSpec
	if cfg == nil {
		return nil, fmt.Errorf("missing KubernetesApplicationSpec in application configuration")
	}

	return &platform.QuickSyncPlanResponse{
		Stages: buildQuickSyncPipeline(*cfg.Input.AutoRollback, now),
	}, nil
}

func (ps *PlannerService) PipelineSyncPlan(ctx context.Context, in *platform.PipelineSyncPlanRequest) (*platform.PipelineSyncPlanResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// func (ps *PlannerService) BuildPlan(ctx context.Context, in *platform.BuildPlanRequest) (*platform.BuildPlanResponse, error) {
// 	var (
// 		pipedConfig     *config.PipedSpec
// 		gitClient       gitClient
// 		secretDecrypter secretDecrypter

// 		repoCfg = config.PipedRepository{
// 			RepoID: in.Deployment.GitPath.Repo.Id,
// 			Remote: in.Deployment.GitPath.Repo.Remote,
// 			Branch: in.Deployment.GitPath.Repo.Branch,
// 		}
// 		targetDSP  deploysource.Provider
// 		runningDSP deploysource.Provider
// 		out        = &platform.DeploymentPlan{}
// 	)

// 	rawCfg, err := config.DecodeYAML(in.PipedConfig)
// 	if err != nil {
// 		err = fmt.Errorf("failed to decode piped configuration (%v)", err)
// 		return nil, err
// 	}

// 	pipedConfig = rawCfg.PipedSpec

// 	// Initialize git client.
// 	gitOptions := []git.Option{
// 		git.WithUserName(pipedConfig.Git.Username),
// 		git.WithEmail(pipedConfig.Git.Email),
// 		git.WithLogger(ps.Logger),
// 	}
// 	for _, repo := range pipedConfig.GitHelmChartRepositories() {
// 		if f := repo.SSHKeyFile; f != "" {
// 			// Configure git client to use the specified SSH key while fetching private Helm charts.
// 			env := fmt.Sprintf("GIT_SSH_COMMAND=ssh -i %s -o StrictHostKeyChecking=no -F /dev/null", f)
// 			gitOptions = append(gitOptions, git.WithGitEnvForRepo(repo.GitRemote, env))
// 		}
// 	}
// 	gitClient, err = git.NewClient(gitOptions...)
// 	if err != nil {
// 		err = fmt.Errorf("failed to create git client (%v)", err)
// 		return nil, err
// 	}
// 	defer func() {
// 		if err := gitClient.Clean(); err != nil {
// 			ps.Logger.Error("had an error while cleaning gitClient", zap.Error(err))
// 			return
// 		}
// 		ps.Logger.Info("successfully cleaned gitClient")
// 	}()

// 	// Initialize secret decrypter.
// 	secretDecrypter, err = initializeSecretDecrypter(pipedConfig)
// 	if err != nil {
// 		err = fmt.Errorf("failed to initialize secret decrypter (%v)", err)
// 		return nil, err
// 	}

// 	targetDSP = deploysource.NewProvider(
// 		filepath.Join(in.WorkingDir, "target-deploysource"),
// 		deploysource.NewGitSourceCloner(gitClient, repoCfg, "target", in.Deployment.Trigger.Commit.Hash),
// 		*in.Deployment.GitPath,
// 		secretDecrypter,
// 	)

// 	if in.LastSuccessfulCommitHash != "" {
// 		gp := *in.Deployment.GitPath
// 		gp.ConfigFilename = in.LastSuccessfulConfigFileName

// 		runningDSP = deploysource.NewProvider(
// 			filepath.Join(in.WorkingDir, "running-deploysource"),
// 			deploysource.NewGitSourceCloner(gitClient, repoCfg, "running", in.LastSuccessfulCommitHash),
// 			gp,
// 			secretDecrypter,
// 		)
// 	}

// 	ds, err := targetDSP.Get(ctx, io.Discard)
// 	if err != nil {
// 		err = fmt.Errorf("error while preparing deploy source data (%v)", err)
// 		return nil, err
// 	}
// 	cfg := ds.ApplicationConfig.KubernetesApplicationSpec
// 	if cfg == nil {
// 		err = fmt.Errorf("missing KubernetesApplicationSpec in application configuration")
// 		return nil, err
// 	}

// 	// TODO: get from request parameter
// 	if cfg.Input.HelmChart != nil {
// 		chartRepoName := cfg.Input.HelmChart.Repository
// 		if chartRepoName != "" {
// 			cfg.Input.HelmChart.Insecure = isInsecureChartRepository(pipedConfig, chartRepoName)
// 		}
// 	}

// 	manifestCache := provider.AppManifestsCache{
// 		AppID:  in.Deployment.ApplicationId,
// 		Cache:  memorycache.NewTTLCache(ctx, time.Hour, time.Minute),
// 		Logger: ps.Logger,
// 	}

// 	// Load previous deployed manifests and new manifests to compare.
// 	newManifests, ok := manifestCache.Get(in.Deployment.Trigger.Commit.Hash)
// 	if !ok {
// 		// When the manifests were not in the cache we have to load them.
// 		loader := provider.NewLoader(in.Deployment.ApplicationName, ds.AppDir, ds.RepoDir, in.Deployment.GitPath.ConfigFilename, cfg.Input, gitClient, ps.Logger)
// 		newManifests, err = loader.LoadManifests(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		manifestCache.Put(in.Deployment.Trigger.Commit.Hash, newManifests)
// 	}

// 	if versions, e := determineVersions(newManifests); e != nil || len(versions) == 0 {
// 		ps.Logger.Warn("unable to determine versions", zap.Error(e))
// 		out.Versions = []*model.ArtifactVersion{
// 			{
// 				Kind:    model.ArtifactVersion_UNKNOWN,
// 				Version: versionUnknown,
// 			},
// 		}
// 	} else {
// 		out.Versions = versions
// 	}

// 	autoRollback := *cfg.Input.AutoRollback

// 	// In case the strategy has been decided by trigger.
// 	// For example: user triggered the deployment via web console.
// 	switch in.Deployment.Trigger.SyncStrategy {
// 	case model.SyncStrategy_QUICK_SYNC:
// 		out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
// 		out.Stages = buildQuickSyncPipeline(autoRollback, time.Now())
// 		out.Summary = in.Deployment.Trigger.StrategySummary
// 		return &platform.BuildPlanResponse{Plan: out}, nil
// 	case model.SyncStrategy_PIPELINE:
// 		if cfg.Pipeline == nil {
// 			err = fmt.Errorf("unable to force sync with pipeline because no pipeline was specified")
// 			return nil, err
// 		}
// 		out.SyncStrategy = model.SyncStrategy_PIPELINE
// 		out.Stages = buildProgressivePipeline(cfg.Pipeline, autoRollback, time.Now())
// 		out.Summary = in.Deployment.Trigger.StrategySummary
// 		return &platform.BuildPlanResponse{Plan: out}, nil
// 	}

// 	// If the progressive pipeline was not configured
// 	// we have only one choise to do is applying all manifestt.
// 	if cfg.Pipeline == nil || len(cfg.Pipeline.Stages) == 0 {
// 		out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
// 		out.Stages = buildQuickSyncPipeline(autoRollback, time.Now())
// 		out.Summary = "Quick sync by applying all manifests (no pipeline was configured)"
// 		return &platform.BuildPlanResponse{Plan: out}, nil
// 	}

// 	// Force to use pipeline when the alwaysUsePipeline field was configured.
// 	if cfg.Planner.AlwaysUsePipeline {
// 		out.SyncStrategy = model.SyncStrategy_PIPELINE
// 		out.Stages = buildProgressivePipeline(cfg.Pipeline, autoRollback, time.Now())
// 		out.Summary = "Sync with the specified pipeline (alwaysUsePipeline was set)"
// 		return &platform.BuildPlanResponse{Plan: out}, nil
// 	}

// 	// This deployment is triggered by a commit with the intent to perform pipeline.
// 	// Commit Matcher will be ignored when triggered by a command.
// 	if p := cfg.CommitMatcher.Pipeline; p != "" && in.Deployment.Trigger.Commander == "" {
// 		pipelineRegex, err := ps.RegexPool.Get(p)
// 		if err != nil {
// 			err = fmt.Errorf("failed to compile commitMatcher.pipeline(%s): %w", p, err)
// 			return &platform.BuildPlanResponse{Plan: out}, err
// 		}
// 		if pipelineRegex.MatchString(in.Deployment.Trigger.Commit.Message) {
// 			out.SyncStrategy = model.SyncStrategy_PIPELINE
// 			out.Stages = buildProgressivePipeline(cfg.Pipeline, autoRollback, time.Now())
// 			out.Summary = fmt.Sprintf("Sync progressively because the commit message was matching %q", p)
// 			return &platform.BuildPlanResponse{Plan: out}, err
// 		}
// 	}

// 	// This deployment is triggered by a commit with the intent to synchronize.
// 	// Commit Matcher will be ignored when triggered by a command.
// 	if s := cfg.CommitMatcher.QuickSync; s != "" && in.Deployment.Trigger.Commander == "" {
// 		syncRegex, err := ps.RegexPool.Get(s)
// 		if err != nil {
// 			err = fmt.Errorf("failed to compile commitMatcher.sync(%s): %w", s, err)
// 			return &platform.BuildPlanResponse{Plan: out}, err
// 		}
// 		if syncRegex.MatchString(in.Deployment.Trigger.Commit.Message) {
// 			out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
// 			out.Stages = buildQuickSyncPipeline(autoRollback, time.Now())
// 			out.Summary = fmt.Sprintf("Quick sync by applying all manifests because the commit message was matching %q", s)
// 			return &platform.BuildPlanResponse{Plan: out}, nil
// 		}
// 	}

// 	// This is the first time to deploy this application
// 	// or it was unable to retrieve that value.
// 	// We just apply all manifests.
// 	if in.LastSuccessfulCommitHash == "" {
// 		out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
// 		out.Stages = buildQuickSyncPipeline(autoRollback, time.Now())
// 		out.Summary = "Quick sync by applying all manifests because it seems this is the first deployment"
// 		return &platform.BuildPlanResponse{Plan: out}, nil
// 	}

// 	// Load manifests of the previously applied commit.
// 	oldManifests, ok := manifestCache.Get(in.LastSuccessfulCommitHash)
// 	if !ok {
// 		// When the manifests were not in the cache we have to load them.
// 		var runningDs *deploysource.DeploySource
// 		runningDs, err = runningDSP.Get(ctx, io.Discard)
// 		if err != nil {
// 			err = fmt.Errorf("failed to prepare the running deploy source data (%v)", err)
// 			return &platform.BuildPlanResponse{Plan: out}, err
// 		}

// 		runningCfg := runningDs.ApplicationConfig.KubernetesApplicationSpec
// 		if runningCfg == nil {
// 			err = fmt.Errorf("unable to find the running configuration (%v)", err)
// 			return &platform.BuildPlanResponse{Plan: out}, err
// 		}
// 		loader := provider.NewLoader(in.Deployment.ApplicationName, runningDs.AppDir, runningDs.RepoDir, in.Deployment.GitPath.ConfigFilename, runningCfg.Input, gitClient, ps.Logger)
// 		oldManifests, err = loader.LoadManifests(ctx)
// 		if err != nil {
// 			err = fmt.Errorf("failed to load previously deployed manifests: %w", err)
// 			return &platform.BuildPlanResponse{Plan: out}, err
// 		}
// 		manifestCache.Put(in.LastSuccessfulCommitHash, oldManifests)
// 	}

// 	progressive, desc := decideStrategy(oldManifests, newManifests, cfg.Workloads, ps.Logger)
// 	out.Summary = desc

// 	if progressive {
// 		out.SyncStrategy = model.SyncStrategy_PIPELINE
// 		out.Stages = buildProgressivePipeline(cfg.Pipeline, autoRollback, time.Now())
// 		return &platform.BuildPlanResponse{Plan: out}, err
// 	}

// 	out.SyncStrategy = model.SyncStrategy_QUICK_SYNC
// 	out.Stages = buildQuickSyncPipeline(autoRollback, time.Now())
// 	return &platform.BuildPlanResponse{Plan: out}, err
// }

// // First up, checks to see if the workload's `spec.template` has been changed,
// // and then checks if the configmap/secret's data.
// func decideStrategy(olds, news []provider.Manifest, workloadRefs []config.K8sResourceReference, logger *zap.Logger) (progressive bool, desc string) {
// 	oldWorkloads := findWorkloadManifests(olds, workloadRefs)
// 	if len(oldWorkloads) == 0 {
// 		desc = "Quick sync by applying all manifests because it was unable to find the currently running workloads"
// 		return
// 	}
// 	newWorkloads := findWorkloadManifests(news, workloadRefs)
// 	if len(newWorkloads) == 0 {
// 		desc = "Quick sync by applying all manifests because it was unable to find workloads in the new manifests"
// 		return
// 	}

// 	workloads := findUpdatedWorkloads(oldWorkloads, newWorkloads)
// 	diffs := make(map[provider.ResourceKey]diff.Nodes, len(workloads))

// 	for _, w := range workloads {
// 		// If the workload's pod template was touched
// 		// do progressive deployment with the specified pipeline.
// 		diffResult, err := provider.Diff(w.old, w.new, logger)
// 		if err != nil {
// 			progressive = true
// 			desc = fmt.Sprintf("Sync progressively due to an error while calculating the diff (%v)", err)
// 			return
// 		}
// 		diffNodes := diffResult.Nodes()
// 		diffs[w.new.Key] = diffNodes

// 		templateDiffs := diffNodes.FindByPrefix("spec.template")
// 		if len(templateDiffs) > 0 {
// 			progressive = true

// 			if msg, changed := checkImageChange(templateDiffs); changed {
// 				desc = msg
// 				return
// 			}

// 			desc = fmt.Sprintf("Sync progressively because pod template of workload %s was changed", w.new.Key.Name)
// 			return
// 		}
// 	}

// 	// If the config/secret was touched, we also need to do progressive
// 	// deployment to check run with the new config/secret content.
// 	oldConfigs := findConfigs(olds)
// 	newConfigs := findConfigs(news)
// 	if len(oldConfigs) > len(newConfigs) {
// 		progressive = true
// 		desc = fmt.Sprintf("Sync progressively because %d configmap/secret deleted", len(oldConfigs)-len(newConfigs))
// 		return
// 	}
// 	if len(oldConfigs) < len(newConfigs) {
// 		progressive = true
// 		desc = fmt.Sprintf("Sync progressively because new %d configmap/secret added", len(newConfigs)-len(oldConfigs))
// 		return
// 	}
// 	for k, oc := range oldConfigs {
// 		nc, ok := newConfigs[k]
// 		if !ok {
// 			progressive = true
// 			desc = fmt.Sprintf("Sync progressively because %s %s was deleted", oc.Key.Kind, oc.Key.Name)
// 			return
// 		}
// 		result, err := provider.Diff(oc, nc, logger)
// 		if err != nil {
// 			progressive = true
// 			desc = fmt.Sprintf("Sync progressively due to an error while calculating the diff (%v)", err)
// 			return
// 		}
// 		if result.HasDiff() {
// 			progressive = true
// 			desc = fmt.Sprintf("Sync progressively because %s %s was updated", oc.Key.Kind, oc.Key.Name)
// 			return
// 		}
// 	}

// 	// Check if this is a scaling commit.
// 	scales := make([]string, 0, len(diffs))
// 	for k, d := range diffs {
// 		if before, after, changed := checkReplicasChange(d); changed {
// 			scales = append(scales, fmt.Sprintf("%s/%s from %s to %s", k.Kind, k.Name, before, after))
// 		}

// 	}
// 	sort.Strings(scales)
// 	if len(scales) > 0 {
// 		desc = fmt.Sprintf("Quick sync to scale %s", strings.Join(scales, ", "))
// 		return
// 	}

// 	desc = "Quick sync by applying all manifests"
// 	return
// }

// func initializeSecretDecrypter(cfg *config.PipedSpec) (crypto.Decrypter, error) {
// 	sm := cfg.SecretManagement
// 	if sm == nil {
// 		return nil, nil
// 	}

// 	switch sm.Type {
// 	case model.SecretManagementTypeNone:
// 		return nil, nil

// 	case model.SecretManagementTypeKeyPair:
// 		key, err := sm.KeyPair.LoadPrivateKey()
// 		if err != nil {
// 			return nil, err
// 		}
// 		decrypter, err := crypto.NewHybridDecrypter(key)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to initialize decrypter (%w)", err)
// 		}
// 		return decrypter, nil

// 	case model.SecretManagementTypeGCPKMS:
// 		return nil, fmt.Errorf("type %q is not implemented yet", sm.Type.String())

// 	case model.SecretManagementTypeAWSKMS:
// 		return nil, fmt.Errorf("type %q is not implemented yet", sm.Type.String())

// 	default:
// 		return nil, fmt.Errorf("unsupported secret management type: %s", sm.Type.String())
// 	}
// }

// func isInsecureChartRepository(cfg *config.PipedSpec, name string) bool {
// 	for _, cr := range cfg.ChartRepositories {
// 		if cr.Name == name {
// 			return cr.Insecure
// 		}
// 	}
// 	return false
// }

// func findWorkloadManifests(manifests []provider.Manifest, refs []config.K8sResourceReference) []provider.Manifest {
// 	if len(refs) == 0 {
// 		return findManifests(provider.KindDeployment, "", manifests)
// 	}

// 	workloads := make([]provider.Manifest, 0)
// 	for _, ref := range refs {
// 		kind := provider.KindDeployment
// 		if ref.Kind != "" {
// 			kind = ref.Kind
// 		}
// 		ms := findManifests(kind, ref.Name, manifests)
// 		workloads = append(workloads, ms...)
// 	}
// 	return workloads
// }

// func findManifests(kind, name string, manifests []provider.Manifest) []provider.Manifest {
// 	out := make([]provider.Manifest, 0, len(manifests))
// 	for _, m := range manifests {
// 		if m.Key.Kind != kind {
// 			continue
// 		}
// 		if name != "" && m.Key.Name != name {
// 			continue
// 		}
// 		out = append(out, m)
// 	}
// 	return out
// }

// type workloadPair struct {
// 	old provider.Manifest
// 	new provider.Manifest
// }

// func findUpdatedWorkloads(olds, news []provider.Manifest) []workloadPair {
// 	pairs := make([]workloadPair, 0)
// 	oldMap := make(map[provider.ResourceKey]provider.Manifest, len(olds))
// 	nomalizeKey := func(k provider.ResourceKey) provider.ResourceKey {
// 		// Ignoring APIVersion because user can upgrade to the new APIVersion for the same workload.
// 		k.APIVersion = ""
// 		if k.Namespace == provider.DefaultNamespace {
// 			k.Namespace = ""
// 		}
// 		return k
// 	}
// 	for _, m := range olds {
// 		key := nomalizeKey(m.Key)
// 		oldMap[key] = m
// 	}
// 	for _, n := range news {
// 		key := nomalizeKey(n.Key)
// 		if o, ok := oldMap[key]; ok {
// 			pairs = append(pairs, workloadPair{
// 				old: o,
// 				new: n,
// 			})
// 		}
// 	}
// 	return pairs
// }

// func findConfigs(manifests []provider.Manifest) map[provider.ResourceKey]provider.Manifest {
// 	configs := make(map[provider.ResourceKey]provider.Manifest)
// 	for _, m := range manifests {
// 		if m.Key.IsConfigMap() {
// 			configs[m.Key] = m
// 		}
// 		if m.Key.IsSecret() {
// 			configs[m.Key] = m
// 		}
// 	}
// 	return configs
// }

// func checkImageChange(ns diff.Nodes) (string, bool) {
// 	const containerImageQuery = `^spec\.template\.spec\.containers\.\d+.image$`
// 	nodes, _ := ns.Find(containerImageQuery)
// 	if len(nodes) == 0 {
// 		return "", false
// 	}

// 	images := make([]string, 0, len(ns))
// 	for _, n := range nodes {
// 		beforeImg := parseContainerImage(n.StringX())
// 		afterImg := parseContainerImage(n.StringY())

// 		if beforeImg.name == afterImg.name {
// 			images = append(images, fmt.Sprintf("image %s from %s to %s", beforeImg.name, beforeImg.tag, afterImg.tag))
// 		} else {
// 			images = append(images, fmt.Sprintf("image %s:%s to %s:%s", beforeImg.name, beforeImg.tag, afterImg.name, afterImg.tag))
// 		}
// 	}
// 	desc := fmt.Sprintf("Sync progressively because of updating %s", strings.Join(images, ", "))
// 	return desc, true
// }

// func checkReplicasChange(ns diff.Nodes) (before, after string, changed bool) {
// 	const replicasQuery = `^spec\.replicas$`
// 	node, err := ns.FindOne(replicasQuery)
// 	if err != nil {
// 		return
// 	}

// 	before = node.StringX()
// 	after = node.StringY()
// 	changed = true
// 	return
// }

// type containerImage struct {
// 	name string
// 	tag  string
// }

// func parseContainerImage(image string) (img containerImage) {
// 	parts := strings.Split(image, ":")
// 	if len(parts) == 2 {
// 		img.tag = parts[1]
// 	}
// 	paths := strings.Split(parts[0], "/")
// 	img.name = paths[len(paths)-1]
// 	return
// }

// // determineVersions decides artifact versions of an application.
// // It finds all container images that are being specified in the workload manifests then returns their names, version numbers, and urls.
// func determineVersions(manifests []provider.Manifest) ([]*model.ArtifactVersion, error) {
// 	imageMap := map[string]struct{}{}
// 	for _, m := range manifests {
// 		// TODO: Determine container image version from other workload kinds such as StatefulSet, Pod, Daemon, CronJob...
// 		if !m.Key.IsDeployment() {
// 			continue
// 		}
// 		data, err := m.MarshalJSON()
// 		if err != nil {
// 			return nil, err
// 		}
// 		var d resource.Deployment
// 		if err := json.Unmarshal(data, &d); err != nil {
// 			return nil, err
// 		}

// 		containers := d.Spec.Template.Spec.Containers
// 		// Remove duplicate images on multiple manifests.
// 		for _, c := range containers {
// 			imageMap[c.Image] = struct{}{}
// 		}
// 	}

// 	versions := make([]*model.ArtifactVersion, 0, len(imageMap))
// 	for i := range imageMap {
// 		image := parseContainerImage(i)
// 		versions = append(versions, &model.ArtifactVersion{
// 			Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
// 			Version: image.tag,
// 			Name:    image.name,
// 			Url:     i,
// 		})
// 	}

// 	return versions, nil
// }
