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

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes/resource"
	"github.com/pipe-cd/pipe/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipe/pkg/app/piped/diff"
	"github.com/pipe-cd/pipe/pkg/app/piped/planner"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	versionUnknown = "unknown"
)

// Planner plans the deployment pipeline for kubernetes application.
type Planner struct {
}

type registerer interface {
	Register(k model.ApplicationKind, p planner.Planner) error
}

// Register registers this planner into the given registerer.
func Register(r registerer) {
	r.Register(model.ApplicationKind_KUBERNETES, &Planner{})
}

// Plan decides which pipeline should be used for the given input.
func (p *Planner) Plan(ctx context.Context, in planner.Input) (out planner.Output, err error) {
	ds, err := in.TargetDSP.Get(ctx, ioutil.Discard)
	if err != nil {
		err = fmt.Errorf("error while preparing deploy source data (%v)", err)
		return
	}
	cfg := ds.DeploymentConfig.KubernetesDeploymentSpec
	if cfg == nil {
		err = fmt.Errorf("missing KubernetesDeploymentSpec in deployment configuration")
		return
	}

	manifestCache := provider.AppManifestsCache{
		AppID:  in.Deployment.ApplicationId,
		Cache:  in.AppManifestsCache,
		Logger: in.Logger,
	}

	// Load previous deployed manifests and new manifests to compare.
	newManifests, ok := manifestCache.Get(in.Deployment.Trigger.Commit.Hash)
	if !ok {
		// When the manifests were not in the cache we have to load them.
		loader := provider.NewManifestLoader(in.Deployment.ApplicationName, ds.AppDir, ds.RepoDir, in.Deployment.GitPath.ConfigFilename, cfg.Input, in.Logger)
		newManifests, err = loader.LoadManifests(ctx)
		if err != nil {
			return
		}
		manifestCache.Put(in.Deployment.Trigger.Commit.Hash, newManifests)
	}

	// Determine application version from the manifests.
	if version, e := determineVersion(newManifests); e != nil {
		in.Logger.Error("unable to determine version", zap.Error(e))
		out.Version = versionUnknown
	} else {
		out.Version = version
	}

	// If the progressive pipeline was not configured
	// we have only one choise to do is applying all manifestt.
	if cfg.Pipeline == nil || len(cfg.Pipeline.Stages) == 0 {
		out.Stages = buildPipeline(cfg.Input.AutoRollback, time.Now())
		out.Summary = "Quick sync by applying all manifests (no pipeline was configured)"
		return
	}

	// This deployment is triggered by a commit with the intent to perform pipeline.
	// Commit Matcher will be ignored when triggered by a command.
	if p := cfg.CommitMatcher.Pipeline; p != "" && in.Deployment.Trigger.Commander == "" {
		pipelineRegex, err := in.RegexPool.Get(p)
		if err != nil {
			err = fmt.Errorf("failed to compile commitMatcher.pipeline(%s): %w", p, err)
			return out, err
		}
		if pipelineRegex.MatchString(in.Deployment.Trigger.Commit.Message) {
			out.Stages = buildProgressivePipeline(cfg.Pipeline, cfg.Input.AutoRollback, time.Now())
			out.Summary = fmt.Sprintf("Sync progressively because the commit message was matching %q", p)
			return out, err
		}
	}

	// This deployment is triggered by a commit with the intent to synchronize.
	// Commit Matcher will be ignored when triggered by a command.
	if s := cfg.CommitMatcher.QuickSync; s != "" && in.Deployment.Trigger.Commander == "" {
		syncRegex, err := in.RegexPool.Get(s)
		if err != nil {
			err = fmt.Errorf("failed to compile commitMatcher.sync(%s): %w", s, err)
			return out, err
		}
		if syncRegex.MatchString(in.Deployment.Trigger.Commit.Message) {
			out.Stages = buildPipeline(cfg.Input.AutoRollback, time.Now())
			out.Summary = fmt.Sprintf("Quick sync by applying all manifests because the commit message was matching %q", s)
			return out, err
		}
	}

	// This is the first time to deploy this application
	// or it was unable to retrieve that value.
	// We just apply all manifests.
	if in.MostRecentSuccessfulCommitHash == "" {
		out.Stages = buildPipeline(cfg.Input.AutoRollback, time.Now())
		out.Summary = "Quick sync by applying all manifests because it seems this is the first deployment"
		return
	}

	// Load manifests of the previously applied commit.
	oldManifests, ok := manifestCache.Get(in.MostRecentSuccessfulCommitHash)
	if !ok {
		// When the manifests were not in the cache we have to load them.
		var runningDs *deploysource.DeploySource
		runningDs, err = in.RunningDSP.Get(ctx, ioutil.Discard)
		if err != nil {
			err = fmt.Errorf("failed to prepare the running deploy source data (%v)", err)
			return
		}

		loader := provider.NewManifestLoader(in.Deployment.ApplicationName, runningDs.AppDir, runningDs.RepoDir, in.Deployment.GitPath.ConfigFilename, cfg.Input, in.Logger)
		oldManifests, err = loader.LoadManifests(ctx)
		if err != nil {
			err = fmt.Errorf("failed to load previously deployed manifests: %w", err)
			return
		}
		manifestCache.Put(in.MostRecentSuccessfulCommitHash, oldManifests)
	}

	progressive, desc := decideStrategy(oldManifests, newManifests)
	out.Summary = desc

	if progressive {
		out.Stages = buildProgressivePipeline(cfg.Pipeline, cfg.Input.AutoRollback, time.Now())
		return
	}

	out.Stages = buildPipeline(cfg.Input.AutoRollback, time.Now())
	return
}

// First up, checks to see if the workload's `spec.template` has been changed,
// and then checks if the configmap/secret's data.
func decideStrategy(olds, news []provider.Manifest) (progressive bool, desc string) {
	oldWorkload, ok := findWorkload(olds)
	if !ok {
		desc = "Quick sync by applying all manifests because it was unable to find the currently running workloads"
		return
	}
	newWorkload, ok := findWorkload(news)
	if !ok {
		desc = "Quick sync by applying all manifests because it was unable to find workloads in the new manifests"
		return
	}

	// If the workload's pod template was touched
	// do progressive deployment with the specified pipeline.
	diffResult, err := provider.Diff(oldWorkload, newWorkload)
	if err != nil {
		progressive = true
		desc = fmt.Sprintf("Sync progressively due to an error while calculating the diff (%v)", err)
		return
	}
	diffNodes := diffResult.Nodes()

	templateDiffs := diffNodes.FindByPrefix("spec.template")
	if len(templateDiffs) > 0 {
		progressive = true

		if msg, changed := checkImageChange(templateDiffs); changed {
			desc = msg
			return
		}

		desc = fmt.Sprintf("Sync progressively because pod template of workload %s was changed", newWorkload.Key.Name)
		return
	}

	// If the config/secret was touched, we also need to do progressive
	// deployment to check run with the new config/secret content.
	oldConfigs := findConfigs(olds)
	newConfigs := findConfigs(news)
	if len(oldConfigs) > len(newConfigs) {
		progressive = true
		desc = fmt.Sprintf("Sync progressively because %d configmap/secret deleted", len(oldConfigs)-len(newConfigs))
		return
	}
	if len(oldConfigs) < len(newConfigs) {
		progressive = true
		desc = fmt.Sprintf("Sync progressively because new %d configmap/secret added", len(newConfigs)-len(oldConfigs))
		return
	}
	for k, oc := range oldConfigs {
		nc, ok := newConfigs[k]
		if !ok {
			progressive = true
			desc = fmt.Sprintf("Sync progressively because %s %s was deleted", oc.Key.Kind, oc.Key.Name)
			return
		}
		result, err := provider.Diff(oc, nc)
		if err != nil {
			progressive = true
			desc = fmt.Sprintf("Sync progressively due to an error while calculating the diff (%v)", err)
			return
		}
		if result.HasDiff() {
			progressive = true
			desc = fmt.Sprintf("Sync progressively because %s %s was updated", oc.Key.Kind, oc.Key.Name)
			return
		}
	}

	// Check if this is a scaling commit.
	if msg, changed := checkReplicasChange(diffNodes); changed {
		desc = msg
		return
	}

	desc = "Quick sync by applying all manifests"
	return
}

// The assumption that an application has only one workload.
func findWorkload(manifests []provider.Manifest) (provider.Manifest, bool) {
	for _, m := range manifests {
		if !m.Key.IsDeployment() {
			continue
		}
		return m, true
	}
	return provider.Manifest{}, false
}

func findConfigs(manifests []provider.Manifest) map[provider.ResourceKey]provider.Manifest {
	configs := make(map[provider.ResourceKey]provider.Manifest)
	for _, m := range manifests {
		if m.Key.IsConfigMap() {
			configs[m.Key] = m
		}
		if m.Key.IsSecret() {
			configs[m.Key] = m
		}
	}
	return configs
}

func checkImageChange(ns diff.Nodes) (string, bool) {
	const containerImageQuery = `^spec\.template\.spec\.containers\.\d+.image$`
	nodes, _ := ns.Find(containerImageQuery)
	if len(nodes) == 0 {
		return "", false
	}

	images := make([]string, 0, len(ns))
	for _, n := range ns {
		beforeName, beforeTag := parseContainerImage(n.StringX())
		afterName, afterTag := parseContainerImage(n.StringY())

		if beforeName == afterName {
			images = append(images, fmt.Sprintf("image %s from %s to %s", beforeName, beforeTag, afterTag))
		} else {
			images = append(images, fmt.Sprintf("image %s:%s to %s:%s", beforeName, beforeTag, afterName, afterTag))
		}
	}
	desc := fmt.Sprintf("Sync progressively because of updating %s", strings.Join(images, ", "))
	return desc, true
}

func checkReplicasChange(ns diff.Nodes) (string, bool) {
	const replicasQuery = `^spec\.replicas$`
	node, err := ns.FindOne(replicasQuery)
	if err != nil {
		return "", false
	}

	desc := fmt.Sprintf("Scale workload from %s to %s.", node.StringX(), node.StringY())
	return desc, true
}

func parseContainerImage(image string) (name, tag string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		tag = parts[1]
	}
	paths := strings.Split(parts[0], "/")
	name = paths[len(paths)-1]
	return
}

// TODO: Add ability to configure how to determine application version.
func determineVersion(manifests []provider.Manifest) (string, error) {
	for _, m := range manifests {
		if !m.Key.IsDeployment() {
			continue
		}
		data, err := m.MarshalJSON()
		if err != nil {
			return "", err
		}
		var d resource.Deployment
		if err := json.Unmarshal(data, &d); err != nil {
			return "", err
		}

		containers := d.Spec.Template.Spec.Containers
		if len(containers) == 0 {
			return versionUnknown, nil
		}
		_, tag := parseContainerImage(containers[0].Image)
		return tag, nil
	}
	return versionUnknown, nil
}
