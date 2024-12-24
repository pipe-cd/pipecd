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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type TemplatingMethod string

const (
	TemplatingMethodHelm      TemplatingMethod = "helm"
	TemplatingMethodKustomize TemplatingMethod = "kustomize"
	TemplatingMethodNone      TemplatingMethod = "none"
)

type LoaderInput struct {
	// for annotations to manage the application live state.
	PipedID    string
	CommitHash string
	AppID      string

	// for templating manifests
	AppName        string
	AppDir         string
	ConfigFilename string
	Manifests      []string

	Namespace        string
	TemplatingMethod TemplatingMethod

	KustomizeVersion string
	KustomizeOptions map[string]string

	HelmVersion string
	HelmChart   *config.InputHelmChart
	HelmOptions *config.InputHelmOptions

	Logger *zap.Logger

	// TODO: define fields for LoaderInput.
}

type Loader struct {
	toolRegistry ToolRegistry
}

type ToolRegistry interface {
	Kustomize(ctx context.Context, version string) (string, error)
	Helm(ctx context.Context, version string) (string, error)
}

func NewLoader(registry ToolRegistry) *Loader {
	return &Loader{
		toolRegistry: registry,
	}
}

func (l *Loader) LoadManifests(ctx context.Context, input LoaderInput) (manifests []Manifest, err error) {
	defer func() {
		// Add builtin annotations for tracking application live state.
		for i := range manifests {
			manifests[i].AddAnnotations(map[string]string{
				LabelManagedBy:          ManagedByPiped,
				LabelPiped:              input.PipedID,
				LabelApplication:        input.AppID,
				LabelOriginalAPIVersion: manifests[i].Key().APIVersion(),
				LabelResourceKey:        manifests[i].Key().String(),
				LabelCommitHash:         input.CommitHash,
			})
		}

		sortManifests(manifests)
	}()

	switch input.TemplatingMethod {
	case TemplatingMethodHelm:
		data, err := l.templateHelmChart(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to template helm chart: %w", err)
		}

		return ParseManifests(data)
	case TemplatingMethodKustomize:
		kustomizePath, err := l.toolRegistry.Kustomize(ctx, input.KustomizeVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to get kustomize tool: %w", err)
		}

		k := NewKustomize(kustomizePath, input.Logger)
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

func sortManifests(manifests []Manifest) {
	if len(manifests) < 2 {
		return
	}

	slices.SortFunc(manifests, func(a, b Manifest) int {
		iAns := a.body.GetAnnotations()
		// Ignore the converting error since it is not so much important.
		iIndex, _ := strconv.Atoi(iAns[AnnotationOrder])

		jAns := b.body.GetAnnotations()
		// Ignore the converting error since it is not so much important.
		jIndex, _ := strconv.Atoi(jAns[AnnotationOrder])

		return iIndex - jIndex
	})
}

func (l *Loader) templateHelmChart(ctx context.Context, input LoaderInput) (string, error) {
	helmPath, err := l.toolRegistry.Helm(ctx, input.HelmVersion)
	if err != nil {
		return "", fmt.Errorf("failed to get helm tool: %w", err)
	}

	h := NewHelm(helmPath, input.Logger)

	switch {
	case input.HelmChart.GitRemote != "":
		return "", errors.New("not implemented yet")

	case input.HelmChart.Repository != "":
		return "", errors.New("not implemented yet")

	default:
		return h.TemplateLocalChart(ctx, input.AppName, input.AppDir, input.Namespace, input.HelmChart.Path, input.HelmOptions)
	}
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
			body: &obj,
		})
	}
	return manifests, nil
}
