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

package ecs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/ecs"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
	"github.com/pipe-cd/pipecd/pkg/app/piped/sourceprocesser"
	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/diff"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type applicationLister interface {
	ListByPlatformProvider(name string) []*model.Application
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type secretDecrypter interface {
	Decrypt(string) (string, error)
}

type reporter interface {
	ReportApplicationSyncState(ctx context.Context, appID string, state model.ApplicationSyncState) error
}

type Detector interface {
	Run(ctx context.Context) error
	ProviderName() string
}

type detector struct {
	provider          config.PipedPlatformProvider
	appLister         applicationLister
	gitClient         gitClient
	stateGetter       ecs.Getter
	reporter          reporter
	appManifestsCache cache.Cache
	interval          time.Duration
	config            *config.PipedSpec
	secretDecrypter   secretDecrypter
	logger            *zap.Logger

	gitRepos map[string]git.Repo
}

func NewDetector(
	cp config.PipedPlatformProvider,
	appLister applicationLister,
	gitClient gitClient,
	stateGetter ecs.Getter,
	reporter reporter,
	appManifestsCache cache.Cache,
	cfg *config.PipedSpec,
	sd secretDecrypter,
	logger *zap.Logger,
) Detector {

	logger = logger.Named("ecs-detector").With(
		zap.String("platform-provider", cp.Name),
	)
	return &detector{
		provider:          cp,
		appLister:         appLister,
		gitClient:         gitClient,
		stateGetter:       stateGetter,
		reporter:          reporter,
		appManifestsCache: appManifestsCache,
		interval:          time.Minute,
		config:            cfg,
		secretDecrypter:   sd,
		gitRepos:          make(map[string]git.Repo),
		logger:            logger,
	}
}

func (d *detector) Run(ctx context.Context) error {
	d.logger.Info("start running drift detector for ecs applications")

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			d.logger.Info("drift detector for ecs applications has been stopped")
			return nil

		case <-ticker.C:
			d.check(ctx)
		}
	}
}

func (d *detector) ProviderName() string {
	return d.provider.Name
}

func (d *detector) check(ctx context.Context) {
	appsByRepo := d.listGroupedApplication()

	for repoID, apps := range appsByRepo {
		gitRepo, ok := d.gitRepos[repoID]
		if !ok {
			// Clone repository for the first time.
			gr, err := d.cloneGitRepository(ctx, repoID)
			if err != nil {
				d.logger.Error("failed to clone git repository",
					zap.String("repo-id", repoID),
					zap.Error(err),
				)
				continue
			}
			gitRepo = gr
			d.gitRepos[repoID] = gitRepo
		}

		// Fetch the latest commit to compare the states.
		branch := gitRepo.GetClonedBranch()
		if err := gitRepo.Pull(ctx, branch); err != nil {
			d.logger.Error("failed to pull repository branch",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
			continue
		}

		// Get the head commit of the repository.
		headCommit, err := gitRepo.GetLatestCommit(ctx)
		if err != nil {
			d.logger.Error("failed to get head commit hash",
				zap.String("repo-id", repoID),
				zap.Error(err),
			)
			continue
		}

		// Start checking all applications in this repository.
		for _, app := range apps {
			if err := d.checkApplication(ctx, app, gitRepo, headCommit); err != nil {
				d.logger.Error(fmt.Sprintf("failed to check application: %s", app.Id), zap.Error(err))
			}
		}
	}
}

func (d *detector) cloneGitRepository(ctx context.Context, repoID string) (git.Repo, error) {
	repoCfg, ok := d.config.GetRepository(repoID)
	if !ok {
		return nil, fmt.Errorf("repository %s was not found in piped configuration", repoID)
	}
	return d.gitClient.Clone(ctx, repoID, repoCfg.Remote, repoCfg.Branch, "")
}

// listGroupedApplication retrieves all applications those should be handled by this director
// and then groups them by repoID.
func (d *detector) listGroupedApplication() map[string][]*model.Application {
	var (
		apps = d.appLister.ListByPlatformProvider(d.provider.Name)
		m    = make(map[string][]*model.Application)
	)
	for _, app := range apps {
		repoID := app.GitPath.Repo.Id
		m[repoID] = append(m[repoID], app)
	}
	return m
}

func (d *detector) checkApplication(ctx context.Context, app *model.Application, repo git.Repo, headCommit git.Commit) error {
	headManifests, err := d.loadConfigs(app, repo, headCommit)
	if err != nil {
		return err
	}
	d.logger.Info(fmt.Sprintf("application %s has ecs definition files at commit %s", app.Id, headCommit.Hash))

	liveManifests, ok := d.stateGetter.GetECSManifests(app.Id)
	if !ok {
		return fmt.Errorf("failed to get live ecs definition files")
	}
	d.logger.Info(fmt.Sprintf("application %s has live ecs definition files", app.Id))

	// Ignore some fields whech are not necessary or unable to detect diff.
	live, head := ignoreParameters(liveManifests, headManifests)

	result, err := provider.Diff(
		live,
		head,
		diff.WithEquateEmpty(),
		diff.WithIgnoreAddingMapKeys(),
		diff.WithCompareNumberAndNumericString(),
	)
	if err != nil {
		return err
	}

	state := makeSyncState(result, headCommit.Hash)

	return d.reporter.ReportApplicationSyncState(ctx, app.Id, state)
}

// ignoreParameters adjusts the fields to ignore unnecessary diff.
//
// TODO: We should check diff of following fields. Currently they are ignored:
//   - service.LoadBalancers
//   - service.Tags
//   - taskDefinition.ContainerDefinitions[].PortMappings[].HostPort
//
// TODO: Maybe we should check diff of following fields when not set in the head manifests in some way. Currently they are ignored:
//   - service.DeploymentConfiguration
//   - service.PlatformVersion
//   - service.RoleArn
func ignoreParameters(liveManifests provider.ECSManifests, headManifests provider.ECSManifests) (live, head provider.ECSManifests) {
	liveService := *liveManifests.ServiceDefinition
	liveService.CreatedAt = nil
	liveService.CreatedBy = nil
	liveService.Events = nil
	liveService.LoadBalancers = nil // TODO: We should set values in headService from the head manifests .
	liveService.PendingCount = 0
	liveService.PlatformFamily = nil // Users cannot specify PlatformFamily in a service definition file. It is automatically set by ECS.
	liveService.RunningCount = 0
	liveService.ServiceArn = nil
	liveService.Status = nil         // Service's Status is shown on the WebUI as healthStatus. So we don't need Status in the driftdetection.
	liveService.TaskDefinition = nil // TODO: Find a way to compare the task definition if possible.
	liveService.TaskSets = nil

	var liveTask types.TaskDefinition
	if liveManifests.TaskDefinition != nil {
		// When liveTask does not exist, e.g. right after the service is created.
		liveTask = *liveManifests.TaskDefinition
		liveTask.RegisteredAt = nil
		liveTask.RegisteredBy = nil
		liveTask.RequiresAttributes = nil
		liveTask.Revision = 0 // TODO: Find a way to compare the revision if possible.
		liveTask.TaskDefinitionArn = nil
		for i := range liveTask.ContainerDefinitions {
			liveTask.ContainerDefinitions[i].Environment = sortKeyPairs(liveTask.ContainerDefinitions[i].Environment)

			for j := range liveTask.ContainerDefinitions[i].PortMappings {
				// We ignore diff of HostPort because it has several default values. See https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_ContainerDefinition.html#ECS-Type-ContainerDefinition-portMappings.
				liveTask.ContainerDefinitions[i].PortMappings[j].HostPort = nil
			}
		}
	}

	headService := *headManifests.ServiceDefinition
	if headService.PlatformVersion == nil {
		// The LATEST platform version is used by default if PlatformVersion is not specified.
		// See https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateService.html#ECS-CreateService-request-platformVersion.
		headService.PlatformVersion = liveService.PlatformVersion
	}
	if headService.RoleArn == nil || len(*headService.RoleArn) == 0 {
		// User's ECS service-linked role is used by default if RoleArn is not specified.
		// See https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateService.html#ECS-CreateService-request-role.
		headService.RoleArn = liveService.RoleArn
	}
	if headService.NetworkConfiguration != nil && headService.NetworkConfiguration.AwsvpcConfiguration != nil {
		awsvpcCfg := *headService.NetworkConfiguration.AwsvpcConfiguration
		awsvpcCfg.Subnets = slices.Clone(awsvpcCfg.Subnets)
		slices.Sort(awsvpcCfg.Subnets)
		if len(awsvpcCfg.AssignPublicIp) == 0 {
			// AssignPublicIp is DISABLED by default.
			// See https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_AwsVpcConfiguration.html#ECS-Type-AwsVpcConfiguration-assignPublicIp.
			awsvpcCfg.AssignPublicIp = types.AssignPublicIpDisabled
		}
		headService.NetworkConfiguration = &types.NetworkConfiguration{AwsvpcConfiguration: &awsvpcCfg}
	}

	// Sort the subnets of the live service as well
	if liveService.NetworkConfiguration != nil && liveService.NetworkConfiguration.AwsvpcConfiguration != nil {
		awsvpcCfg := *liveService.NetworkConfiguration.AwsvpcConfiguration
		awsvpcCfg.Subnets = slices.Clone(awsvpcCfg.Subnets)
		slices.Sort(awsvpcCfg.Subnets)
		liveService.NetworkConfiguration = &types.NetworkConfiguration{AwsvpcConfiguration: &awsvpcCfg}
	}

	if headService.DeploymentConfiguration == nil {
		liveService.DeploymentConfiguration = nil
	}

	// TODO: In order to check diff of the tags, we need to add pipecd-managed tags and sort.
	liveService.Tags = nil
	headService.Tags = nil

	headTask := *headManifests.TaskDefinition
	headTask.Status = types.TaskDefinitionStatusActive // If livestate's status is not ACTIVE, we should re-deploy a new task definition.
	if liveManifests.TaskDefinition != nil {
		headTask.Compatibilities = liveTask.Compatibilities // Users can specify Compatibilities in a task definition file, but it is not used when registering a task definition.
	}

	headTask.ContainerDefinitions = slices.Clone(headManifests.TaskDefinition.ContainerDefinitions)
	for i := range headTask.ContainerDefinitions {
		cd := &headTask.ContainerDefinitions[i]
		if cd.Essential == nil {
			// Essential is true by default. See https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_ContainerDefinition.html#ECS-Type-ContainerDefinition-es.
			cd.Essential = aws.Bool(true)
		}

		cd.Environment = sortKeyPairs(cd.Environment)

		cd.PortMappings = slices.Clone(cd.PortMappings)
		for j := range cd.PortMappings {
			pm := &cd.PortMappings[j]
			if len(pm.Protocol) == 0 {
				// Protocol is tcp by default. See https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_PortMapping.html#ECS-Type-PortMapping-protocol.
				pm.Protocol = types.TransportProtocolTcp
			}
			pm.HostPort = nil // We ignore diff of HostPort because it has several default values.
		}
	}

	live = provider.ECSManifests{ServiceDefinition: &liveService, TaskDefinition: &liveTask}
	head = provider.ECSManifests{ServiceDefinition: &headService, TaskDefinition: &headTask}
	return live, head
}

func (d *detector) loadConfigs(app *model.Application, repo git.Worktree, headCommit git.Commit) (provider.ECSManifests, error) {
	var (
		manifestCache = provider.ECSManifestsCache{
			AppID:  app.Id,
			Cache:  d.appManifestsCache,
			Logger: d.logger,
		}
		repoDir = repo.GetPath()
		appDir  = filepath.Join(repoDir, app.GitPath.Path)
	)

	manifests, ok := manifestCache.Get(headCommit.Hash)
	if ok {
		return manifests, nil
	}
	// When the manifests were not in the cache we have to load them.
	cfg, err := d.loadApplicationConfiguration(repoDir, app)
	if err != nil {
		return provider.ECSManifests{}, fmt.Errorf("failed to load application configuration: %w", err)
	}

	gds, ok := cfg.GetGenericApplication()
	if !ok {
		return provider.ECSManifests{}, fmt.Errorf("unsupported application kind %s", cfg.Kind)
	}

	var (
		encryptionUsed = d.secretDecrypter != nil && gds.Encryption != nil
		attachmentUsed = gds.Attachment != nil
	)

	// We have to copy repository into another directory because
	// decrypting the sealed secrets or attaching files might change the git repository.
	if attachmentUsed || encryptionUsed {
		dir, err := os.MkdirTemp("", "detector-git-processing")
		if err != nil {
			return provider.ECSManifests{}, fmt.Errorf("failed to prepare a temporary directory for git repository (%w)", err)
		}
		defer os.RemoveAll(dir)

		repo, err = repo.Copy(filepath.Join(dir, "repo"))
		if err != nil {
			return provider.ECSManifests{}, fmt.Errorf("failed to copy the cloned git repository (%w)", err)
		}
		defer repo.Clean()

		repoDir := repo.GetPath()
		appDir = filepath.Join(repoDir, app.GitPath.Path)
	}

	var templProcessors []sourceprocesser.SourceTemplateProcessor
	// Decrypting secrets to manifests.
	if encryptionUsed {
		templProcessors = append(templProcessors, sourceprocesser.NewSecretDecrypterProcessor(gds.Encryption, d.secretDecrypter))
	}
	// Then attaching configurated files to manifests.
	if attachmentUsed {
		templProcessors = append(templProcessors, sourceprocesser.NewAttachmentProcessor(gds.Attachment))
	}
	if len(templProcessors) > 0 {
		sp := sourceprocesser.NewSourceProcessor(appDir, templProcessors...)
		if err := sp.Process(); err != nil {
			return provider.ECSManifests{}, fmt.Errorf("failed to process source files: %w", err)
		}
	}

	var serviceDefFile string
	var taskDefFile string
	if cfg.ECSApplicationSpec != nil {
		serviceDefFile = cfg.ECSApplicationSpec.Input.ServiceDefinitionFile
		taskDefFile = cfg.ECSApplicationSpec.Input.TaskDefinitionFile
	}
	serviceDef, err := provider.LoadServiceDefinition(appDir, serviceDefFile)
	if err != nil {
		return provider.ECSManifests{}, fmt.Errorf("failed to load new service definition: %w", err)
	}
	taskDef, err := provider.LoadTaskDefinition(appDir, taskDefFile)
	if err != nil {
		return provider.ECSManifests{}, fmt.Errorf("failed to load new task definition: %w", err)
	}

	manifests = provider.ECSManifests{
		ServiceDefinition: &serviceDef,
		TaskDefinition:    &taskDef,
	}
	manifestCache.Put(headCommit.Hash, manifests)

	return manifests, nil
}

func (d *detector) loadApplicationConfiguration(repoPath string, app *model.Application) (*config.Config, error) {
	path := filepath.Join(repoPath, app.GitPath.GetApplicationConfigFilePath())
	cfg, err := config.LoadFromYAML(path)
	if err != nil {
		return nil, err
	}
	if appKind, ok := cfg.Kind.ToApplicationKind(); !ok || appKind != app.Kind {
		return nil, fmt.Errorf("application in application configuration file is not match, got: %s, expected: %s", appKind, app.Kind)
	}
	return cfg, nil
}

func makeSyncState(r *provider.DiffResult, commit string) model.ApplicationSyncState {
	if r.NoChange() {
		return model.ApplicationSyncState{
			Status:    model.ApplicationSyncStatus_SYNCED,
			Timestamp: time.Now().Unix(),
		}
	}

	if ignoreAutoScalingDiff(r) {
		return model.ApplicationSyncState{
			Status:      model.ApplicationSyncStatus_SYNCED,
			ShortReason: "Ignore diff of `desiredCount`.",
			Reason:      "`desiredCount` is 0 or not defined in your config (which means ignoring updating desiredCount) and only `desiredCount` is changed.",
			Timestamp:   time.Now().Unix(),
		}
	}

	shortReason := "The ecs definition files are not synced"
	if len(commit) >= 7 {
		commit = commit[:7]
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Diff between the defined state in Git at commit %s and actual live state:\n\n", commit))
	b.WriteString("--- Actual   (LiveState)\n+++ Expected (Git)\n\n")

	details := r.Render(provider.DiffRenderOptions{
		// Currently, we do not use the diff command to render the result
		// because ECS adds a large number of default values to the
		// running manifest that causes a wrong diff text.
		UseDiffCommand: false,
	})
	b.WriteString(details)

	return model.ApplicationSyncState{
		Status:      model.ApplicationSyncStatus_OUT_OF_SYNC,
		ShortReason: shortReason,
		Reason:      b.String(),
		Timestamp:   time.Now().Unix(),
	}
}

// ignoreAutoScalingDiff returns true if the diff contains only autoscaled desiredCount.
func ignoreAutoScalingDiff(r *provider.DiffResult) bool {
	return r.Diff.NumNodes() == 1 &&
		r.New.ServiceDefinition.DesiredCount == 0 && // When desiredCount is 0 or not defined in the head manifest, autoscaling may be enabled.
		r.Old.ServiceDefinition.DesiredCount != r.New.ServiceDefinition.DesiredCount
}

func sortKeyPairs(kps []types.KeyValuePair) []types.KeyValuePair {
	sorted := slices.Clone(kps)
	sort.Slice(sorted, func(i, j int) bool {
		return *sorted[i].Name < *sorted[j].Name
	})

	return sorted
}
