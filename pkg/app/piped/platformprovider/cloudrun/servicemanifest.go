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

package cloudrun

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"google.golang.org/api/run/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type ServiceManifest struct {
	Name string
	u    *unstructured.Unstructured
}

func (m ServiceManifest) SetRevision(name string) error {
	return unstructured.SetNestedField(m.u.Object, name, "spec", "template", "metadata", "name")
}

type RevisionTraffic struct {
	RevisionName string `json:"revisionName"`
	Percent      int    `json:"percent"`
}

func (m ServiceManifest) UpdateTraffic(revisions []RevisionTraffic) error {
	items := []any{}
	for i := range revisions {
		out, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&revisions[i])
		if err != nil {
			return fmt.Errorf("unable to set traffic for object: %w", err)
		}
		items = append(items, out)
	}

	return unstructured.SetNestedSlice(m.u.Object, items, "spec", "traffic")
}

func (m ServiceManifest) UpdateAllTraffic(revision string) error {
	return m.UpdateTraffic([]RevisionTraffic{
		{
			RevisionName: revision,
			Percent:      100,
		},
	})
}

func (m ServiceManifest) YamlBytes() ([]byte, error) {
	return yaml.Marshal(m.u)
}

func (m ServiceManifest) Labels() map[string]string {
	return m.u.GetLabels()
}

func (m ServiceManifest) RevisionLabels() map[string]string {
	v, _, _ := unstructured.NestedStringMap(m.u.Object, "spec", "template", "metadata", "labels")
	return v
}

func (m ServiceManifest) AppID() (string, bool) {
	v := m.Labels()
	if v == nil || v[LabelApplication] == "" {
		return "", false
	}
	return v[LabelApplication], true
}

func (m ServiceManifest) AddLabels(labels map[string]string) {
	if len(labels) == 0 {
		return
	}

	lbls := m.u.GetLabels()
	if lbls == nil {
		m.u.SetLabels(labels)
		return
	}
	maps.Copy(lbls, labels)
	m.u.SetLabels(lbls)
}

func (m ServiceManifest) AddRevisionLabels(labels map[string]string) error {
	if len(labels) == 0 {
		return nil
	}

	fields := []string{"spec", "template", "metadata", "labels"}
	lbls, ok, err := unstructured.NestedStringMap(m.u.Object, fields...)
	if err != nil {
		return err
	}
	if !ok {
		return unstructured.SetNestedStringMap(m.u.Object, labels, fields...)
	}

	maps.Copy(lbls, labels)
	return unstructured.SetNestedStringMap(m.u.Object, lbls, fields...)
}

func (m ServiceManifest) RunService() (*run.Service, error) {
	data, err := m.YamlBytes()
	if err != nil {
		return nil, err
	}

	var s run.Service
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func loadServiceManifest(path string) (ServiceManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ServiceManifest{}, err
	}
	return ParseServiceManifest(data)
}

func ParseServiceManifest(data []byte) (ServiceManifest, error) {
	var obj unstructured.Unstructured
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return ServiceManifest{}, err
	}

	return ServiceManifest{
		Name: obj.GetName(),
		u:    &obj,
	}, nil
}

func DecideRevisionName(sm ServiceManifest, commit string) (string, error) {
	tag, err := FindImageTag(sm)
	if err != nil {
		return "", err
	}
	tag = strings.ReplaceAll(tag, ".", "")

	if len(commit) > 7 {
		commit = commit[:7]
	}
	return fmt.Sprintf("%s-%s-%s", sm.Name, tag, commit), nil
}

// FindImageTag returns the tag of the service's container image.
// It is empty when the image has no tag, e.g. when it is referenced only by digest.
func FindImageTag(sm ServiceManifest) (string, error) {
	image, err := findContainerImage(sm)
	if err != nil {
		return "", err
	}
	_, tag, _ := parseContainerImage(image)
	return tag, nil
}

// FindImageVersion returns the version of the service's container image.
// It is the image digest, or the image tag when the image has no digest.
func FindImageVersion(sm ServiceManifest) (string, error) {
	image, err := findContainerImage(sm)
	if err != nil {
		return "", err
	}
	_, tag, digest := parseContainerImage(image)
	return imageVersion(tag, digest), nil
}

func findContainerImage(sm ServiceManifest) (string, error) {
	containers, ok, err := unstructured.NestedSlice(sm.u.Object, "spec", "template", "spec", "containers")
	if err != nil {
		return "", err
	}
	if !ok || len(containers) == 0 {
		return "", fmt.Errorf("spec.template.spec.containers was missing")
	}

	container, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&containers[0])
	if err != nil {
		return "", fmt.Errorf("invalid container format")
	}

	image, ok, err := unstructured.NestedString(container, "image")
	if err != nil {
		return "", err
	}
	if !ok || image == "" {
		return "", fmt.Errorf("image was missing")
	}
	return image, nil
}

// parseContainerImage parses an image reference in the format of
// [registry/]name[:tag][@digest] into its name, tag and digest.
func parseContainerImage(image string) (name, tag, digest string) {
	// Split off the digest first, so that "@sha256:<digest>" is not taken as the tag.
	if idx := strings.Index(image, "@"); idx != -1 {
		digest = image[idx+1:]
		image = image[:idx]
	}
	// Look for the tag only in the last path segment, so that a registry
	// port such as "localhost:5000" is not taken as the tag.
	paths := strings.Split(image, "/")
	name = paths[len(paths)-1]
	if idx := strings.LastIndex(name, ":"); idx != -1 {
		tag = name[idx+1:]
		name = name[:idx]
	}
	return
}

// imageVersion returns the digest, or the tag when the image has no digest.
// The digest is preferred because it decides which image is actually deployed.
func imageVersion(tag, digest string) string {
	if digest != "" {
		return digest
	}
	return tag
}

func FindArtifactVersions(sm ServiceManifest) ([]*model.ArtifactVersion, error) {
	containers, ok, err := unstructured.NestedSlice(sm.u.Object, "spec", "template", "spec", "containers")
	if err != nil {
		return nil, err
	}
	if !ok || len(containers) == 0 {
		return nil, fmt.Errorf("spec.template.spec.containers was missing")
	}

	container, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&containers[0])
	if err != nil {
		return nil, fmt.Errorf("invalid container format")
	}

	image, ok, err := unstructured.NestedString(container, "image")
	if err != nil {
		return nil, err
	}
	if !ok || image == "" {
		return nil, fmt.Errorf("image was missing")
	}
	name, tag, digest := parseContainerImage(image)

	return []*model.ArtifactVersion{
		{
			Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
			Version: imageVersion(tag, digest),
			Name:    name,
			Url:     image,
		},
	}, nil
}
