// Copyright 2023 The PipeCD Authors.
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
	"sync"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/config"
)

type Applier interface {
	// ApplyManifest does applying the given manifest.
	ApplyManifest(ctx context.Context, manifest Manifest) error
	// CreateManifest does creating resource from given manifest.
	CreateManifest(ctx context.Context, manifest Manifest) error
	// ReplaceManifest does replacing resource from given manifest.
	ReplaceManifest(ctx context.Context, manifest Manifest) error
	// Delete deletes the given resource from Kubernetes cluster.
	Delete(ctx context.Context, key ResourceKey) error
}

type applier struct {
	input            config.KubernetesDeploymentInput
	platformProvider config.PlatformProviderKubernetesConfig
	logger           *zap.Logger

	kubectl  *Kubectl
	initOnce sync.Once
	initErr  error
}

func NewApplier(input config.KubernetesDeploymentInput, cp config.PlatformProviderKubernetesConfig, logger *zap.Logger) Applier {
	return &applier{
		input:            input,
		platformProvider: cp,
		logger:           logger.Named("kubernetes-applier"),
	}
}

// ApplyManifest does applying the given manifest.
func (a *applier) ApplyManifest(ctx context.Context, manifest Manifest) error {
	a.initOnce.Do(func() {
		a.kubectl, a.initErr = a.findKubectl(ctx, a.getToolVersionToRun())
	})
	if a.initErr != nil {
		return a.initErr
	}
	if a.input.AutoCreateNamespace == true {
		if err := a.kubectl.CreateNamespace(
			ctx,
			a.platformProvider.KubeConfigPath,
			a.getNamespaceToRun(manifest.Key),
		); err != nil {
			return err
		}
	}

	return a.kubectl.Apply(
		ctx,
		a.platformProvider.KubeConfigPath,
		a.getNamespaceToRun(manifest.Key),
		manifest,
	)
}

// CreateManifest uses kubectl to create the given manifests.
func (a *applier) CreateManifest(ctx context.Context, manifest Manifest) error {
	a.initOnce.Do(func() {
		a.kubectl, a.initErr = a.findKubectl(ctx, a.getToolVersionToRun())
	})
	if a.initErr != nil {
		return a.initErr
	}
	if a.input.AutoCreateNamespace == true {
		if err := a.kubectl.CreateNamespace(
			ctx,
			a.platformProvider.KubeConfigPath,
			a.getNamespaceToRun(manifest.Key),
		); err != nil {
			return err
		}
	}

	return a.kubectl.Create(
		ctx,
		a.platformProvider.KubeConfigPath,
		a.getNamespaceToRun(manifest.Key),
		manifest,
	)
}

// ReplaceManifest uses kubectl to replace the given manifests.
func (a *applier) ReplaceManifest(ctx context.Context, manifest Manifest) error {
	a.initOnce.Do(func() {
		a.kubectl, a.initErr = a.findKubectl(ctx, a.getToolVersionToRun())
	})
	if a.initErr != nil {
		return a.initErr
	}

	err := a.kubectl.Replace(
		ctx,
		a.platformProvider.KubeConfigPath,
		a.getNamespaceToRun(manifest.Key),
		manifest,
	)
	if err == nil {
		return nil
	}

	if errors.Is(err, errorReplaceNotFound) {
		return ErrNotFound
	}

	return err
}

// Delete deletes the given resource from Kubernetes cluster.
// If the resource key is different, this returns ErrNotFound.
func (a *applier) Delete(ctx context.Context, k ResourceKey) (err error) {
	a.initOnce.Do(func() {
		a.kubectl, a.initErr = a.findKubectl(ctx, a.getToolVersionToRun())
	})
	if a.initErr != nil {
		return a.initErr
	}

	m, err := a.kubectl.Get(
		ctx,
		a.platformProvider.KubeConfigPath,
		a.getNamespaceToRun(k),
		k,
	)

	if err != nil {
		return err
	}

	if k.String() != m.GetAnnotations()[LabelResourceKey] {
		return ErrNotFound
	}

	return a.kubectl.Delete(
		ctx,
		a.platformProvider.KubeConfigPath,
		a.getNamespaceToRun(k),
		k,
	)
}

// getNamespaceToRun returns namespace used on kubectl apply/delete commands.
// priority: config.KubernetesDeploymentInput > kubernetes.ResourceKey
func (a *applier) getNamespaceToRun(k ResourceKey) string {
	if a.input.Namespace != "" {
		return a.input.Namespace
	}
	return k.Namespace
}

// getToolVersionToRun returns version of kubectl which should be used for commands.
// priority: applicationConfig.KubectlVersion > pipedConfig.KubectlVersion
func (a *applier) getToolVersionToRun() string {
	if a.input.KubectlVersion != "" {
		return a.input.KubectlVersion
	}
	return a.platformProvider.KubectlVersion
}

func (a *applier) findKubectl(ctx context.Context, version string) (*Kubectl, error) {
	path, installed, err := toolregistry.DefaultRegistry().Kubectl(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("no kubectl %s (%v)", version, err)
	}
	if installed {
		a.logger.Info(fmt.Sprintf("kubectl %s has just been installed because of no pre-installed binary for that version", version))
	}
	return NewKubectl(version, path), nil
}

type multiApplier struct {
	appliers []Applier
}

// NewMultiApplier creates an applier that duplicates its operations to all the provided appliers.
func NewMultiApplier(appliers ...Applier) Applier {
	return &multiApplier{
		appliers: appliers,
	}
}

func (a *multiApplier) ApplyManifest(ctx context.Context, manifest Manifest) error {
	for _, a := range a.appliers {
		if err := a.ApplyManifest(ctx, manifest); err != nil {
			return err
		}
	}
	return nil
}

func (a *multiApplier) CreateManifest(ctx context.Context, manifest Manifest) error {
	for _, a := range a.appliers {
		if err := a.CreateManifest(ctx, manifest); err != nil {
			return err
		}
	}
	return nil
}

func (a *multiApplier) ReplaceManifest(ctx context.Context, manifest Manifest) error {
	for _, a := range a.appliers {
		if err := a.ReplaceManifest(ctx, manifest); err != nil {
			return err
		}
	}
	return nil
}

func (a *multiApplier) Delete(ctx context.Context, key ResourceKey) error {
	for _, a := range a.appliers {
		if err := a.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}
