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
	"strings"
	"time"

	"go.uber.org/zap"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes/resource"
	"github.com/pipe-cd/pipe/pkg/app/piped/planner"
	"github.com/pipe-cd/pipe/pkg/model"
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
	cfg := in.DeploymentConfig.KubernetesDeploymentSpec
	if cfg == nil {
		err = fmt.Errorf("malfored deployment configuration: missing KubernetesDeploymentSpec")
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
		loader := provider.NewManifestLoader(in.Deployment.ApplicationName, in.AppDir, in.RepoDir, cfg.Input, in.Logger)
		newManifests, err = loader.LoadManifests(ctx)
		if err != nil {
			err = fmt.Errorf("failed to load new manifests: %v", err)
			return
		}
		manifestCache.Put(in.Deployment.Trigger.Commit.Hash, newManifests)
	}

	// Determine application version from the manifests.
	if version, e := determineVersion(newManifests); e != nil {
		in.Logger.Error("unable to determine version", zap.Error(e))
	} else {
		out.Version = version
	}

	// This is the first time to deploy this application
	// or it was unable to retrieve that value.
	// We just apply all manifests.
	if in.MostRecentSuccessfulCommitHash == "" {
		out.Stages = buildPipeline(cfg.Input.AutoRollback, time.Now())
		out.Description = "Apply all manifests because it was unable to find the most recent successful commit."
		return
	}

	// If the commit is a revert one. Let's apply primary to rollback.
	// TODO: Find a better way to determine the revert commit.
	if strings.Contains(in.Deployment.Trigger.Commit.Message, "/pipecd rollback ") {
		out.Stages = buildPipeline(cfg.Input.AutoRollback, time.Now())
		out.Description = fmt.Sprintf("Rollback from commit %s.", in.MostRecentSuccessfulCommitHash)
		return
	}

	// Checkout to the most recent successful commit to load its manifests.
	err = in.Repo.Checkout(ctx, in.MostRecentSuccessfulCommitHash)
	if err != nil {
		err = fmt.Errorf("failed to checkout to commit %s: %v", in.MostRecentSuccessfulCommitHash, err)
		return
	}

	// Load manifests of the previously applied commit.
	oldManifests, ok := manifestCache.Get(in.MostRecentSuccessfulCommitHash)
	if !ok {
		// When the manifests were not in the cache we have to load them.
		loader := provider.NewManifestLoader(in.Deployment.ApplicationName, in.AppDir, in.RepoDir, cfg.Input, in.Logger)
		oldManifests, err = loader.LoadManifests(ctx)
		if err != nil {
			err = fmt.Errorf("failed to load previously deployed manifests: %v", err)
			return
		}
		manifestCache.Put(in.MostRecentSuccessfulCommitHash, oldManifests)
	}

	progressive, desc := decideStrategy(oldManifests, newManifests)
	out.Description = desc

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
		desc = "Apply all manifests because it was unable to find the currently running workloads."
		return
	}
	newWorkload, ok := findWorkload(news)
	if !ok {
		desc = "Apply all manifests because it was unable to find workloads in the new manifests."
		return
	}

	// If the workload's pod template was touched
	// do progressive deployment with the specified pipeline.
	var (
		workloadDiffs = provider.Diff(oldWorkload, newWorkload, provider.WithDiffPathPrefix("spec"))
		templateDiffs = workloadDiffs.FindByPrefix("spec.template")
	)
	if len(templateDiffs) > 0 {
		progressive = true

		if msg, changed := checkImageChange(templateDiffs); changed {
			desc = msg
			return
		}

		desc = fmt.Sprintf("Progressive deployment because pod template of workload %s was changed.", newWorkload.Key.Name)
		return
	}

	// If the config/secret was touched,
	// we also need to do progressive deployment to check run with the new config/secret content.
	oldConfigs := findConfigs(olds)
	newConfigs := findConfigs(news)
	if len(oldConfigs) > 0 && len(newConfigs) > 0 {
		for k, oc := range oldConfigs {
			nc, ok := newConfigs[k]
			if !ok {
				desc = fmt.Sprintf("Progressive deployment because %s %s was deleted", oc.Key.Kind, oc.Key.Name)
				return
			}
			diffs := provider.Diff(oc, nc, provider.WithDiffPathPrefix("data"))
			if len(diffs) > 0 {
				progressive = true
				desc = fmt.Sprintf("Progressive deployment because %s %s was updated", oc.Key.Kind, oc.Key.Name)
				return
			}
			delete(newConfigs, k)
		}
		for _, nc := range newConfigs {
			desc = fmt.Sprintf("Progressive deployment because %s %s was added", nc.Key.Kind, nc.Key.Name)
			return
		}
	}

	// Check if this is a scaling commit.
	if msg, changed := checkReplicasChange(workloadDiffs); changed {
		desc = msg
		return
	}

	desc = "Apply all manifests"
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
	configs := make(map[provider.ResourceKey]provider.Manifest, 0)
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

func checkImageChange(diffList provider.DiffResultList) (string, bool) {
	const containerImageQuery = `^spec.template.spec.containers.\[\d+\].image$`
	imageDiffs := diffList.FindAll(containerImageQuery)

	if len(imageDiffs) == 0 {
		return "", false
	}

	images := make([]string, 0, len(imageDiffs))
	for _, d := range imageDiffs {
		beforeName, beforeTag := parseContainerImage(d.Before)
		afterName, afterTag := parseContainerImage(d.After)

		if beforeName == afterName {
			images = append(images, fmt.Sprintf("image %s from %s to %s", beforeName, beforeTag, afterTag))
		} else {
			images = append(images, fmt.Sprintf("image %s:%s to %s:%s", beforeName, beforeTag, afterName, afterTag))
		}
	}
	desc := fmt.Sprintf("Progressive deployment because of updating %s.", strings.Join(images, ", "))
	return desc, true
}

func checkReplicasChange(diffList provider.DiffResultList) (string, bool) {
	const replicasQuery = `^spec.replicas$`
	diff, found, _ := diffList.Find(replicasQuery)
	if !found {
		return "", false
	}

	desc := fmt.Sprintf("Scale workload from %s to %s.", diff.Before, diff.After)
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
			return "", nil
		}
		_, tag := parseContainerImage(containers[0].Image)
		return tag, nil
	}
	return "", nil
}
