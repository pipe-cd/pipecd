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

package provider

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type TemplatingMethod string

const (
	TemplatingMethodHelm      TemplatingMethod = "helm"
	TemplatingMethodKustomize TemplatingMethod = "kustomize"
	TemplatingMethodNone      TemplatingMethod = "none"
)

type LoaderInput struct {
	AppName        string
	AppDir         string
	ConfigFilename string
	Manifests      []string

	Namespace        string
	TemplatingMethod TemplatingMethod

	KustomizeVersion string
	KustomizeOptions map[string]string

	// TODO: define fields for LoaderInput.
}

type Loader struct {
	toolRegistry ToolRegistry
}

type ToolRegistry interface {
	Kustomize(ctx context.Context, version string) (string, error)
	Helm(ctx context.Context, version string) (string, error)
}

func (l *Loader) LoadManifests(ctx context.Context, input LoaderInput) (manifests []Manifest, err error) {
	defer func() {
		// Override namespace if set because ParseManifests does not parse it
		// if namespace is not explicitly specified in the manifests.
		setNamespace(manifests, input.Namespace)
		sortManifests(manifests)
	}()

	switch input.TemplatingMethod {
	case TemplatingMethodHelm:
		return nil, errors.New("not implemented yet")
	case TemplatingMethodKustomize:
		kustomizePath, err := l.toolRegistry.Kustomize(ctx, input.KustomizeVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to get kustomize tool: %w", err)
		}

		k := NewKustomize(kustomizePath, zap.NewNop()) // TODO: pass logger
		data, err := k.Template(ctx, input.AppName, input.AppDir, input.KustomizeOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to template kustomize manifests: %w", err)
		}

		return ParseManifests(data)
	case TemplatingMethodNone:
		return LoadPlainYAMLManifests(input.AppDir, input.Manifests, input.ConfigFilename)
	default:
		return nil, fmt.Errorf("unsupported templating method %s", input.TemplatingMethod)
	}
}

func setNamespace(manifests []Manifest, namespace string) {
	if namespace == "" {
		return
	}
	for i := range manifests {
		manifests[i].Key.Namespace = namespace
	}
}

func sortManifests(manifests []Manifest) {
	if len(manifests) < 2 {
		return
	}

	slices.SortFunc(manifests, func(a, b Manifest) int {
		iAns := a.Body.GetAnnotations()
		// Ignore the converting error since it is not so much important.
		iIndex, _ := strconv.Atoi(iAns[AnnotationOrder])

		jAns := b.Body.GetAnnotations()
		// Ignore the converting error since it is not so much important.
		jIndex, _ := strconv.Atoi(jAns[AnnotationOrder])

		return iIndex - jIndex
	})
}

func LoadPlainYAMLManifests(dir string, names []string, configFilename string) ([]Manifest, error) {
	// If no name was specified we have to walk the app directory to collect the manifest list.
	if len(names) == 0 {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path == dir {
				return nil
			}
			if d.IsDir() {
				return fs.SkipDir
			}
			if ext := filepath.Ext(d.Name()); ext != ".yaml" && ext != ".yml" && ext != ".json" {
				return nil
			}
			if model.IsApplicationConfigFile(d.Name()) {
				// MEMO: can we remove this check because we have configFilename?
				return nil
			}
			if d.Name() == configFilename {
				return nil
			}
			names = append(names, d.Name())
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	manifests := make([]Manifest, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		ms, err := LoadManifestsFromYAMLFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load manifest at %s (%w)", path, err)
		}
		manifests = append(manifests, ms...)
	}

	return manifests, nil
}

// LoadManifestsFromYAMLFile loads the manifests from the given file.
func LoadManifestsFromYAMLFile(path string) ([]Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseManifests(string(data))
}

// ParseManifests parses the given data and returns a list of Manifest.
func ParseManifests(data string) ([]Manifest, error) {
	const separator = "\n---"
	var (
		parts     = strings.Split(data, separator)
		manifests = make([]Manifest, 0, len(parts))
	)

	for i, part := range parts {
		// Ignore all the cases where no content between separator.
		if len(strings.TrimSpace(part)) == 0 {
			continue
		}
		// Append new line which trim by document separator.
		if i != len(parts)-1 {
			part += "\n"
		}
		var obj unstructured.Unstructured
		if err := yaml.Unmarshal([]byte(part), &obj); err != nil {
			return nil, err
		}
		if len(obj.Object) == 0 {
			continue
		}
		manifests = append(manifests, Manifest{
			Key:  MakeResourceKey(&obj),
			Body: &obj,
		})
	}
	return manifests, nil
}
