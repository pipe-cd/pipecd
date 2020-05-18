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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/kapetaniosci/pipe/pkg/config"
)

type Manifest struct {
	APIVersion string
	Kind       string
	Namespace  string
	Name       string

	originalData []byte
	unstructured *unstructured.Unstructured
}

func (m Manifest) YamlBytes() ([]byte, error) {
	return yaml.Marshal(m.unstructured)
}

func (m Manifest) AddAnnotations(annotations map[string]string) {
	if len(annotations) == 0 {
		return
	}

	annos := m.unstructured.GetAnnotations()
	if annos != nil {
		for k, v := range annotations {
			annos[k] = v
		}
	} else {
		annos = annotations
	}
	m.unstructured.SetAnnotations(annos)
}

func (m Manifest) ResourceKey() string {
	return fmt.Sprintf("%s/%s/%s/%s", m.APIVersion, m.Kind, m.Namespace, m.Name)
}

func (e *Executor) loadManifests(ctx context.Context) ([]Manifest, error) {
	switch e.templatingMethod {
	case TemplatingMethodHelm:
		return nil, nil
	case TemplatingMethodKustomize:
		return nil, nil
	case TemplatingMethodNone:
		return loadPlainYAMLMannifests(ctx, e.appDirPath, e.config.Input.Manifests)
	}
	return nil, nil
}

func loadPlainYAMLMannifests(ctx context.Context, dir string, names []string) ([]Manifest, error) {
	// If no name was specified we have to walk the app directory to collect the manifest list.
	if len(names) == 0 {
		err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			fmt.Printf("looking file %s at %s\n", f.Name(), path)
			if err != nil {
				return err
			}
			if path == dir {
				return nil
			}
			if f.IsDir() {
				return filepath.SkipDir
			}
			ext := filepath.Ext(f.Name())
			if ext != ".yaml" && ext != ".yml" {
				return nil
			}
			if f.Name() == config.DeploymentConfigurationFileName {
				return nil
			}
			fmt.Printf("found a manifest file %s at %s\n", f.Name(), path)
			names = append(names, f.Name())
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	manifests := make([]Manifest, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		ms, err := loadManifestsFromYAMLFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load maninifest at %s (%v)", path, err)
		}
		manifests = append(manifests, ms...)
	}

	return manifests, nil
}

func loadManifestsFromYAMLFile(path string) ([]Manifest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	const seperator = "\n---"
	var (
		parts     = strings.Split(string(data), seperator)
		manifests = make([]Manifest, 0, len(parts))
	)

	for _, part := range parts {
		//	Ignore all the cases where no content between separator.
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			continue
		}
		var obj unstructured.Unstructured
		if err := yaml.Unmarshal([]byte(part), &obj); err != nil {
			return nil, err
		}
		manifests = append(manifests, Manifest{
			APIVersion:   obj.GetAPIVersion(),
			Kind:         obj.GetKind(),
			Namespace:    obj.GetNamespace(),
			Name:         obj.GetName(),
			originalData: []byte(part),
			unstructured: &obj,
		})
	}
	return manifests, nil
}
