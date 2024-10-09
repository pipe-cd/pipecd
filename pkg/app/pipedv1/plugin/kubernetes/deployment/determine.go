// Copyright 2024 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deployment

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
	"github.com/pipe-cd/pipecd/pkg/model"
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
