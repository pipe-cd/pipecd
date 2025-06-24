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

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes_multicluster/config"
)

type kubectl interface {
	Apply(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error
	Create(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error
	Replace(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error
	ForceReplace(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error
	Delete(ctx context.Context, kubeconfig, namespace string, key ResourceKey) error
	Get(ctx context.Context, kubeconfig, namespace string, key ResourceKey) (Manifest, error)
	CreateNamespace(ctx context.Context, kubeconfig, namespace string) error
}

type Applier struct {
	kubectl kubectl

	input        config.KubernetesDeploymentInput
	deployTarget config.KubernetesDeployTargetConfig
	logger       *zap.Logger
}

func NewApplier(kubectl kubectl, input config.KubernetesDeploymentInput, cp config.KubernetesDeployTargetConfig, logger *zap.Logger) *Applier {
	return &Applier{
		kubectl:      kubectl,
		input:        input,
		deployTarget: cp,
		logger:       logger.Named("kubernetes-applier"),
	}
}

// ApplyManifest does applying the given manifest.
func (a *Applier) ApplyManifest(ctx context.Context, manifest Manifest) error {
	if a.input.AutoCreateNamespace {
		err := a.kubectl.CreateNamespace(
			ctx,
			a.deployTarget.KubeConfigPath,
			manifest.Key().Namespace(),
		)
		if err != nil && !errors.Is(err, errResourceAlreadyExists) {
			return err
		}
	}

	return a.kubectl.Apply(
		ctx,
		a.deployTarget.KubeConfigPath,
		manifest.Key().Namespace(),
		manifest,
	)
}

// CreateManifest uses kubectl to create the given manifests.
func (a *Applier) CreateManifest(ctx context.Context, manifest Manifest) error {
	if a.input.AutoCreateNamespace {
		err := a.kubectl.CreateNamespace(
			ctx,
			a.deployTarget.KubeConfigPath,
			manifest.Key().Namespace(),
		)
		if err != nil && !errors.Is(err, errResourceAlreadyExists) {
			return err
		}
	}

	return a.kubectl.Create(
		ctx,
		a.deployTarget.KubeConfigPath,
		manifest.Key().Namespace(),
		manifest,
	)
}

// ReplaceManifest uses kubectl to replace the given manifests.
func (a *Applier) ReplaceManifest(ctx context.Context, manifest Manifest) error {
	err := a.kubectl.Replace(
		ctx,
		a.deployTarget.KubeConfigPath,
		manifest.Key().Namespace(),
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

// ForceReplaceManifest uses kubectl to forcefully replace the given manifests.
func (a *Applier) ForceReplaceManifest(ctx context.Context, manifest Manifest) error {
	err := a.kubectl.ForceReplace(
		ctx,
		a.deployTarget.KubeConfigPath,
		manifest.Key().Namespace(),
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
func (a *Applier) Delete(ctx context.Context, k ResourceKey) (err error) {
	m, err := a.kubectl.Get(
		ctx,
		a.deployTarget.KubeConfigPath,
		k.Namespace(),
		k,
	)

	if err != nil {
		return err
	}

	if k.String() != m.Key().String() {
		return ErrNotFound
	}

	return a.kubectl.Delete(
		ctx,
		a.deployTarget.KubeConfigPath,
		k.Namespace(),
		k,
	)
}

// getNamespaceToRun returns namespace used on kubectl apply/delete commands.
// Applier.input.Namespace is not used here because it is referenced when the manifest is loaded.
func (a *Applier) getNamespaceToRun(k ResourceKey) string {
	return k.namespace
}
