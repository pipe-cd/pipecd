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
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/git"
)

type TemplatingMethod string

const (
	TemplatingMethodHelm      TemplatingMethod = "helm"
	TemplatingMethodKustomize TemplatingMethod = "kustomize"
	TemplatingMethodNone      TemplatingMethod = "none"
)

type Loader interface {
	// LoadManifests renders and loads all manifests for application.
	LoadManifests(ctx context.Context) ([]Manifest, error)
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

type registry interface {
	Kubectl(ctx context.Context, version string) (string, error)
	Kustomize(ctx context.Context, version string) (string, error)
	Helm(ctx context.Context, version string) (string, error)
}

type loader struct {
	appName        string
	appDir         string
	repoDir        string
	configFileName string
	input          config.KubernetesDeploymentInput
	gc             gitClient
	logger         *zap.Logger
	toolregistry   registry

	templatingMethod TemplatingMethod
	kustomize        *Kustomize
	helm             *Helm
	initOnce         sync.Once
	initErr          error
}

func NewLoader(
	appName, appDir, repoDir, configFileName string,
	input config.KubernetesDeploymentInput,
	gc gitClient,
	logger *zap.Logger,
	toolregistry registry,
) Loader {

	return &loader{
		appName:        appName,
		appDir:         appDir,
		repoDir:        repoDir,
		configFileName: configFileName,
		input:          input,
		gc:             gc,
		logger:         logger.Named("kubernetes-loader"),
		toolregistry:   toolregistry,
	}
}

// LoadManifests renders and loads all manifests for application.
func (l *loader) LoadManifests(ctx context.Context) (manifests []Manifest, err error) {
	defer func() {
		// Override namespace if set because ParseManifests does not parse it
		// if namespace is not explicitly specified in the manifests.
		setNamespace(manifests, l.input.Namespace)
		sortManifests(manifests)
	}()
	l.initOnce.Do(func() {
		var initErrorHelm, initErrorKustomize error
		l.templatingMethod = determineTemplatingMethod(l.input, l.appDir)
		if l.templatingMethod != TemplatingMethodNone {
			l.helm, initErrorHelm = l.findHelm(ctx, l.input.HelmVersion)
			l.kustomize, initErrorKustomize = l.findKustomize(ctx, l.input.KustomizeVersion)
			l.initErr = errors.Join(initErrorHelm, initErrorKustomize)
		}
	})
	if l.initErr != nil {
		return nil, l.initErr
	}

	switch l.templatingMethod {
	case TemplatingMethodHelm:
		var data string
		switch {
		case l.input.HelmChart.GitRemote != "":
			chart := helmRemoteGitChart{
				GitRemote: l.input.HelmChart.GitRemote,
				Ref:       l.input.HelmChart.Ref,
				Path:      l.input.HelmChart.Path,
			}
			data, err = l.helm.TemplateRemoteGitChart(ctx,
				l.appName,
				l.appDir,
				l.input.Namespace,
				chart,
				l.gc,
				l.input.HelmOptions)

		case l.input.HelmChart.Repository != "":
			chart := helmRemoteChart{
				Repository: l.input.HelmChart.Repository,
				Name:       l.input.HelmChart.Name,
				Version:    l.input.HelmChart.Version,
				Insecure:   l.input.HelmChart.Insecure,
			}
			data, err = l.helm.TemplateRemoteChart(ctx,
				l.appName,
				l.appDir,
				l.input.Namespace,
				chart,
				l.input.HelmOptions)

		default:
			data, err = l.helm.TemplateLocalChart(ctx,
				l.appName,
				l.appDir,
				l.input.Namespace,
				l.input.HelmChart.Path,
				l.input.HelmOptions)
		}

		if err != nil {
			err = fmt.Errorf("unable to run helm template: %w", err)
			return
		}
		manifests, err = ParseManifests(data)

	case TemplatingMethodKustomize:
		var data string
		data, err = l.kustomize.Template(ctx, l.appName, l.appDir, l.input.KustomizeOptions)
		if err != nil {
			err = fmt.Errorf("unable to run kustomize template: %w", err)
			return
		}
		manifests, err = ParseManifests(data)

	case TemplatingMethodNone:
		manifests, err = LoadPlainYAMLManifests(l.appDir, l.input.Manifests, l.configFileName)

	default:
		err = fmt.Errorf("unsupport templating method %v", l.templatingMethod)
	}

	return
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
	sort.Slice(manifests, func(i, j int) bool {
		iAns := manifests[i].GetAnnotations()
		// Ignore the converting error since it is not so much important.
		iIndex, _ := strconv.Atoi(iAns[AnnotationOrder])

		jAns := manifests[j].GetAnnotations()
		// Ignore the converting error since it is not so much important.
		jIndex, _ := strconv.Atoi(jAns[AnnotationOrder])

		return iIndex < jIndex
	})
}

func (l *loader) findKustomize(ctx context.Context, version string) (*Kustomize, error) {
	path, err := l.toolregistry.Kustomize(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("no kustomize %s (%v)", version, err)
	}
	return NewKustomize(version, path, l.logger), nil
}

func (l *loader) findHelm(ctx context.Context, version string) (*Helm, error) {
	path, err := l.toolregistry.Helm(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("no helm %s (%v)", version, err)
	}
	return NewHelm(version, path, l.logger, l.toolregistry), nil
}

func determineTemplatingMethod(input config.KubernetesDeploymentInput, appDirPath string) TemplatingMethod {
	if input.HelmChart != nil {
		return TemplatingMethodHelm
	}
	if _, err := os.Stat(filepath.Join(appDirPath, kustomizationFileName)); err == nil {
		return TemplatingMethodKustomize
	}
	return TemplatingMethodNone
}
