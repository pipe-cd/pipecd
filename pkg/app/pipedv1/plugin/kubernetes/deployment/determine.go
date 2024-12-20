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
	"fmt"
	"slices"
	"strings"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/plugin/diff"
)

type containerImage struct {
	name string
	tag  string
}

// parseContainerImage splits the container image into name and tag.
// The image should be in the format of "name:tag".
// If the tag is not specified, it will be empty.
func parseContainerImage(image string) (img containerImage) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		img.tag = parts[1]
	}
	paths := strings.Split(parts[0], "/")
	img.name = paths[len(paths)-1]
	return
}

// determineVersions decides artifact versions of an application.
// It finds all container images that are being specified in the workload manifests then returns their names and tags.
func determineVersions(manifests []provider.Manifest) ([]*model.ArtifactVersion, error) {
	imageMap := map[string]struct{}{}
	for _, m := range manifests {
		// TODO: we should consider other fields like spec.jobTempate.spec.template.spec.containers because CronJob uses this format.
		containers, ok, err := unstructured.NestedSlice(m.Body.Object, "spec", "template", "spec", "containers")
		if err != nil {
			// if the containers field is not an array, it will return an error.
			// we define this as error because the 'containers' is plural form, so it should be an array.
			return nil, err
		}
		if !ok {
			continue
		}
		// Remove duplicate images on multiple manifests.
		for _, c := range containers {
			m, ok := c.(map[string]interface{})
			if !ok {
				// TODO: Add logging.
				continue
			}
			img, ok := m["image"]
			if !ok {
				continue
			}
			imgStr, ok := img.(string)
			if !ok {
				return nil, fmt.Errorf("invalid image format: %T(%v)", img, img)
			}
			imageMap[imgStr] = struct{}{}
		}
	}

	versions := make([]*model.ArtifactVersion, 0, len(imageMap))
	for i := range imageMap {
		image := parseContainerImage(i)
		versions = append(versions, &model.ArtifactVersion{
			Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
			Version: image.tag,
			Name:    image.name,
			Url:     i,
		})
	}

	return versions, nil
}

// findManifests returns the manifests that have the specified kind and name.
func findManifests(kind, name string, manifests []provider.Manifest) []provider.Manifest {
	out := make([]provider.Manifest, 0, len(manifests))
	for _, m := range manifests {
		if m.Body.GetKind() != kind {
			continue
		}
		if name != "" && m.Body.GetName() != name {
			continue
		}
		out = append(out, m)
	}
	return out
}

// findWorkloadManifests returns the manifests that have the specified references.
// the default kind is Deployment if it is not specified.
func findWorkloadManifests(manifests []provider.Manifest, refs []config.K8sResourceReference) []provider.Manifest {
	if len(refs) == 0 {
		return findManifests(provider.KindDeployment, "", manifests)
	}

	workloads := make([]provider.Manifest, 0)
	for _, ref := range refs {
		kind := provider.KindDeployment
		if ref.Kind != "" {
			kind = ref.Kind
		}
		ms := findManifests(kind, ref.Name, manifests)
		workloads = append(workloads, ms...)
	}
	return workloads
}

type workloadPair struct {
	old provider.Manifest
	new provider.Manifest
}

func checkImageChange(ns diff.Nodes) (string, bool) {
	const containerImageQuery = `^spec\.template\.spec\.containers\.\d+.image$`
	nodes, _ := ns.Find(containerImageQuery)
	if len(nodes) == 0 {
		return "", false
	}

	images := make([]string, 0, len(ns))
	for _, n := range nodes {
		beforeImg := parseContainerImage(n.StringX())
		afterImg := parseContainerImage(n.StringY())

		if beforeImg.name == afterImg.name {
			images = append(images, fmt.Sprintf("image %s from %s to %s", beforeImg.name, beforeImg.tag, afterImg.tag))
		} else {
			images = append(images, fmt.Sprintf("image %s:%s to %s:%s", beforeImg.name, beforeImg.tag, afterImg.name, afterImg.tag))
		}
	}
	desc := fmt.Sprintf("Sync progressively because of updating %s", strings.Join(images, ", "))
	return desc, true
}

// checkReplicasChange checks if the replicas field is changed.
func checkReplicasChange(ns diff.Nodes) (before, after string, changed bool) {
	const replicasQuery = `^spec\.replicas$`
	node, err := ns.FindOne(replicasQuery)
	if err != nil {
		// no difference between the before and after manifests, or unknown error occurred.
		return "", "", false
	}

	if node.TypeX == nil {
		// The replicas field is not found in the before manifest.
		// There is difference between the before and after manifests, So it means the replicas field is added in the after manifest.
		// So the replicas field in the before manifest is nil, we should return "<nil>" as the before value.
		return "<nil>", node.StringY(), true
	}

	if node.TypeY == nil {
		// The replicas field is not found in the after manifest.
		// There is difference between the before and after manifests, So it means the replicas field is removed in the after manifest.
		// So the replicas field in the after manifest is nil, we should return "<nil>" as the after value.
		return node.StringX(), "<nil>", true
	}

	return node.StringX(), node.StringY(), true
}

// First up, checks to see if the workload's `spec.template` has been changed,
// and then checks if the configmap/secret's data.
func determineStrategy(olds, news []provider.Manifest, workloadRefs []config.K8sResourceReference, logger *zap.Logger) (strategy model.SyncStrategy, summary string) {
	oldWorkloads := findWorkloadManifests(olds, workloadRefs)
	if len(oldWorkloads) == 0 {
		return model.SyncStrategy_QUICK_SYNC, "Quick sync by applying all manifests because it was unable to find the currently running workloads"
	}
	newWorkloads := findWorkloadManifests(news, workloadRefs)
	if len(newWorkloads) == 0 {
		return model.SyncStrategy_QUICK_SYNC, "Quick sync by applying all manifests because it was unable to find workloads in the new manifests"
	}

	workloads := provider.FindUpdatedWorkloads(oldWorkloads, newWorkloads)
	diffs := make(map[provider.ResourceKey]diff.Nodes, len(workloads))

	for _, w := range workloads {
		// If the workload's pod template was touched
		// do progressive deployment with the specified pipeline.
		diffResult, err := provider.Diff(w.Old, w.New, logger)
		if err != nil {
			return model.SyncStrategy_PIPELINE, fmt.Sprintf("Sync progressively due to an error while calculating the diff (%v)", err)
		}
		diffNodes := diffResult.Nodes()
		diffs[w.New.Key] = diffNodes

		templateDiffs := diffNodes.FindByPrefix("spec.template")
		if len(templateDiffs) > 0 {
			if msg, changed := checkImageChange(templateDiffs); changed {
				return model.SyncStrategy_PIPELINE, msg
			}
			return model.SyncStrategy_PIPELINE, fmt.Sprintf("Sync progressively because pod template of workload %s was changed", w.New.Key.Name)
		}
	}

	// If the config/secret was touched, we also need to do progressive
	// deployment to check run with the new config/secret content.
	oldConfigs := provider.FindConfigsAndSecrets(olds)
	newConfigs := provider.FindConfigsAndSecrets(news)
	if len(oldConfigs) > len(newConfigs) {
		return model.SyncStrategy_PIPELINE, fmt.Sprintf("Sync progressively because %d configmap/secret deleted", len(oldConfigs)-len(newConfigs))
	}
	if len(oldConfigs) < len(newConfigs) {
		return model.SyncStrategy_PIPELINE, fmt.Sprintf("Sync progressively because new %d configmap/secret added", len(newConfigs)-len(oldConfigs))
	}
	for k, oc := range oldConfigs {
		nc, ok := newConfigs[k]
		if !ok {
			return model.SyncStrategy_PIPELINE, fmt.Sprintf("Sync progressively because %s %s was deleted", oc.Key.Kind, oc.Key.Name)
		}
		result, err := provider.Diff(oc, nc, logger)
		if err != nil {
			return model.SyncStrategy_PIPELINE, fmt.Sprintf("Sync progressively due to an error while calculating the diff (%v)", err)
		}
		if result.HasDiff() {
			return model.SyncStrategy_PIPELINE, fmt.Sprintf("Sync progressively because %s %s was updated", oc.Key.Kind, oc.Key.Name)
		}
	}

	// Check if this is a scaling commit.
	scales := make([]string, 0, len(diffs))
	for k, d := range diffs {
		if before, after, changed := checkReplicasChange(d); changed {
			scales = append(scales, fmt.Sprintf("%s/%s from %s to %s", k.Kind, k.Name, before, after))
		}

	}
	if len(scales) > 0 {
		slices.Sort(scales)
		return model.SyncStrategy_QUICK_SYNC, fmt.Sprintf("Quick sync to scale %s", strings.Join(scales, ", "))
	}

	return model.SyncStrategy_QUICK_SYNC, "Quick sync by applying all manifests"
}
