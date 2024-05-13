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

package kubernetes

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
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
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

type loader struct {
	appName               string
	appDir                string
	repoDir               string
	configFileName        string
	input                 config.KubernetesDeploymentInput
	isNamespacedResources map[schema.GroupVersionKind]bool
	gc                    gitClient
	logger                *zap.Logger

	templatingMethod TemplatingMethod
	kustomize        *Kustomize
	helm             *Helm
	initOnce         sync.Once
	initErr          error
}

func NewLoader(
	appName, appDir, repoDir, configFileName string,
	input config.KubernetesDeploymentInput,
	isNamespacedResources map[schema.GroupVersionKind]bool,
	gc gitClient,
	logger *zap.Logger,
) Loader {

	return &loader{
		appName:               appName,
		appDir:                appDir,
		repoDir:               repoDir,
		configFileName:        configFileName,
		input:                 input,
		isNamespacedResources: isNamespacedResources,
		gc:                    gc,
		logger:                logger.Named("kubernetes-loader"),
	}
}

// LoadManifests renders and loads all manifests for application.
func (l *loader) LoadManifests(ctx context.Context) (manifests []Manifest, err error) {
	defer func() {
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

	for i := range manifests {
		namespace, err := l.refineNamespace(manifests[i])
		if err != nil {
			return nil, err
		}
		manifests[i].Key.Namespace = namespace
		manifests[i].u.SetNamespace(namespace)
	}

	return
}

// refineNamespace returns the namespace to use for the given manifest.
// The priority is as follows:
// 1. The namespace set in the application configuration.
// 2. The namespace set in the manifest.
// 3. The default namespace.
// If the resource is cluster-scoped, it returns an empty string.
func (l *loader) refineNamespace(m Manifest) (string, error) {
	namespaced, ok := l.isNamespacedResources[m.u.GroupVersionKind()]
	if !ok {
		return "", fmt.Errorf("unknown resource kind %s", m.u.GroupVersionKind().String())
	}

	// cluster-scoped resource
	if !namespaced {
		return "", nil
	}

	// namespace-scoped resource from here
	if l.input.Namespace != "" {
		return l.input.Namespace, nil
	}

	ns := m.u.GetNamespace()
	if ns != "" {
		return ns, nil
	}

	return "default", nil
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
	path, installed, err := toolregistry.DefaultRegistry().Kustomize(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("no kustomize %s (%v)", version, err)
	}
	if installed {
		l.logger.Info(fmt.Sprintf("kustomize %s has just been installed because of no pre-installed binary for that version", version))
	}
	return NewKustomize(version, path, l.logger), nil
}

func (l *loader) findHelm(ctx context.Context, version string) (*Helm, error) {
	path, installed, err := toolregistry.DefaultRegistry().Helm(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("no helm %s (%v)", version, err)
	}
	if installed {
		l.logger.Info(fmt.Sprintf("helm %s has just been installed because of no pre-installed binary for that version", version))
	}
	return NewHelm(version, path, l.logger), nil
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
