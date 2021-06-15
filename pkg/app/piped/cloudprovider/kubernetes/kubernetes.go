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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/git"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	LabelManagedBy            = "pipecd.dev/managed-by"             // Always be piped.
	LabelPiped                = "pipecd.dev/piped"                  // The id of piped handling this application.
	LabelApplication          = "pipecd.dev/application"            // The application this resource belongs to.
	LabelCommitHash           = "pipecd.dev/commit-hash"            // Hash value of the deployed commit.
	LabelResourceKey          = "pipecd.dev/resource-key"           // The resource key generated by apiVersion, namespace and name. e.g. apps/v1/Deployment/namespace/demo-app
	LabelOriginalAPIVersion   = "pipecd.dev/original-api-version"   // The api version defined in git configuration. e.g. apps/v1
	LabelIgnoreDriftDirection = "pipecd.dev/ignore-drift-detection" // Whether the drift detection should ignore this resource.
	AnnotationConfigHash      = "pipecd.dev/config-hash"            // The hash value of all mouting config resources.
	ManagedByPiped            = "piped"
	IgnoreDriftDetectionTrue  = "true"

	kustomizationFileName = "kustomization.yaml"
)

type TemplatingMethod string

const (
	TemplatingMethodHelm      TemplatingMethod = "helm"
	TemplatingMethodKustomize TemplatingMethod = "kustomize"
	TemplatingMethodNone      TemplatingMethod = "none"
)

type Provider interface {
	ManifestLoader
	Applier
}

type ManifestLoader interface {
	// LoadManifests renders and loads all manifests for application.
	LoadManifests(ctx context.Context) ([]Manifest, error)
}

type Applier interface {
	// Apply does applying application manifests by using the tool specified in Input.
	Apply(ctx context.Context) error
	// ApplyManifest does applying the given manifest.
	ApplyManifest(ctx context.Context, manifest Manifest) error
	// Delete deletes the given resource from Kubernetes cluster.
	Delete(ctx context.Context, key ResourceKey) error
}

type gitClient interface {
	Clone(ctx context.Context, repoID, remote, branch, destination string) (git.Repo, error)
}

var (
	// shared gitClient used inside this package for downloading dependencies.
	sharedGitClient         gitClient
	initSharedGitClientOnce sync.Once
)

type provider struct {
	appName        string
	appDir         string
	repoDir        string
	configFileName string
	input          config.KubernetesDeploymentInput
	logger         *zap.Logger

	kubectl          *Kubectl
	kustomize        *Kustomize
	helm             *Helm
	templatingMethod TemplatingMethod
	initOnce         sync.Once
	initErr          error
}

func init() {
	registerMetrics()
}

func initSharedGitClient(logger *zap.Logger) error {
	var err error
	initSharedGitClientOnce.Do(func() {
		sharedGitClient, err = git.NewClient("", "", logger)
	})
	return err
}

func NewProvider(appName, appDir, repoDir, configFileName string, input config.KubernetesDeploymentInput, logger *zap.Logger) Provider {
	return &provider{
		appName:        appName,
		appDir:         appDir,
		repoDir:        repoDir,
		configFileName: configFileName,
		input:          input,
		logger:         logger.Named("kubernetes-provider"),
	}
}

func NewManifestLoader(appName, appDir, repoDir, configFileName string, input config.KubernetesDeploymentInput, logger *zap.Logger) ManifestLoader {
	return NewProvider(appName, appDir, repoDir, configFileName, input, logger)
}

func (p *provider) init(ctx context.Context) {
	if err := initSharedGitClient(p.logger); err != nil {
		p.initErr = err
		return
	}

	p.templatingMethod = determineTemplatingMethod(p.input, p.appDir)

	// We need kubectl for all templating methods.
	p.kubectl, p.initErr = p.findKubectl(ctx, p.input.KubectlVersion)
	if p.initErr != nil {
		return
	}

	switch p.templatingMethod {
	case TemplatingMethodHelm:
		p.helm, p.initErr = p.findHelm(ctx, p.input.HelmVersion)

	case TemplatingMethodKustomize:
		p.kustomize, p.initErr = p.findKustomize(ctx, p.input.KustomizeVersion)
	}
}

// LoadManifests renders and loads all manifests for application.
func (p *provider) LoadManifests(ctx context.Context) (manifests []Manifest, err error) {
	p.initOnce.Do(func() { p.init(ctx) })
	if p.initErr != nil {
		return nil, p.initErr
	}

	switch p.templatingMethod {
	case TemplatingMethodHelm:
		var data string
		switch {
		case p.input.HelmChart.GitRemote != "":
			chart := helmRemoteGitChart{
				GitRemote: p.input.HelmChart.GitRemote,
				Ref:       p.input.HelmChart.Ref,
				Path:      p.input.HelmChart.Path,
			}
			data, err = p.helm.TemplateRemoteGitChart(ctx,
				p.appName,
				p.appDir,
				p.input.Namespace,
				chart,
				sharedGitClient,
				p.input.HelmOptions)

		case p.input.HelmChart.Repository != "":
			chart := helmRemoteChart{
				Repository: p.input.HelmChart.Repository,
				Name:       p.input.HelmChart.Name,
				Version:    p.input.HelmChart.Version,
				Insecure:   p.input.HelmChart.Insecure,
			}
			data, err = p.helm.TemplateRemoteChart(ctx,
				p.appName,
				p.appDir,
				p.input.Namespace,
				chart,
				p.input.HelmOptions)

		default:
			data, err = p.helm.TemplateLocalChart(ctx,
				p.appName,
				p.appDir,
				p.input.Namespace,
				p.input.HelmChart.Path,
				p.input.HelmOptions)
		}

		if err != nil {
			err = fmt.Errorf("unable to run helm template: %w", err)
			return
		}
		manifests, err = ParseManifests(data)

	case TemplatingMethodKustomize:
		var data string
		data, err = p.kustomize.Template(ctx, p.appName, p.appDir, p.input.KustomizeOptions)
		if err != nil {
			err = fmt.Errorf("unable to run kustomize template: %w", err)
			return
		}
		manifests, err = ParseManifests(data)

	case TemplatingMethodNone:
		manifests, err = LoadPlainYAMLManifests(p.appDir, p.input.Manifests, p.configFileName)

	default:
		err = fmt.Errorf("unsupport templating method %v", p.templatingMethod)
	}

	return
}

// Apply does applying application manifests by using the tool specified in Input.
func (p *provider) Apply(ctx context.Context) error {
	return nil
}

// ApplyManifest does applying the given manifest.
func (p *provider) ApplyManifest(ctx context.Context, manifest Manifest) error {
	p.initOnce.Do(func() { p.init(ctx) })
	if p.initErr != nil {
		return p.initErr
	}

	return p.kubectl.Apply(ctx, p.getNamespaceToRun(manifest.Key), manifest)
}

// Delete deletes the given resource from Kubernetes cluster.
func (p *provider) Delete(ctx context.Context, k ResourceKey) (err error) {
	p.initOnce.Do(func() { p.init(ctx) })
	if p.initErr != nil {
		return p.initErr
	}

	return p.kubectl.Delete(ctx, p.getNamespaceToRun(k), k)
}

// getNamespaceToRun returns namespace used on kubectl apply/delete commands.
// priority: config.KubernetesDeploymentInput > kubernetes.ResourceKey
func (p *provider) getNamespaceToRun(k ResourceKey) string {
	if p.input.Namespace != "" {
		return p.input.Namespace
	}
	return k.Namespace
}

func (p *provider) findKubectl(ctx context.Context, version string) (*Kubectl, error) {
	path, installed, err := toolregistry.DefaultRegistry().Kubectl(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("no kubectl %s (%v)", version, err)
	}
	if installed {
		p.logger.Info(fmt.Sprintf("kubectl %s has just been installed because of no pre-installed binary for that version", version))
	}
	return NewKubectl(version, path), nil
}

func (p *provider) findKustomize(ctx context.Context, version string) (*Kustomize, error) {
	path, installed, err := toolregistry.DefaultRegistry().Kustomize(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("no kustomize %s (%v)", version, err)
	}
	if installed {
		p.logger.Info(fmt.Sprintf("kustomize %s has just been installed because of no pre-installed binary for that version", version))
	}
	return NewKustomize(version, path, p.logger), nil
}

func (p *provider) findHelm(ctx context.Context, version string) (*Helm, error) {
	path, installed, err := toolregistry.DefaultRegistry().Helm(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("no helm %s (%v)", version, err)
	}
	if installed {
		p.logger.Info(fmt.Sprintf("helm %s has just been installed because of no pre-installed binary for that version", version))
	}
	return NewHelm(version, path, p.logger), nil
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
