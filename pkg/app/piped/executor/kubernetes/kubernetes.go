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
	"path/filepath"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/kubernetes"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	variantLabel = "pipecd.dev/variant" // Variant name: primary, stage, baseline
)

type Executor struct {
	executor.Input

	provider provider.Provider
	config   *config.KubernetesDeploymentSpec
}

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.ApplicationKind, f executor.Factory) error
}

// Register registers this executor factory into a given registerer.
func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &Executor{
			Input: in,
		}
	}
	r.Register(model.StageK8sPrimaryRollout, f)
	r.Register(model.StageK8sCanaryRollout, f)
	r.Register(model.StageK8sCanaryClean, f)
	r.Register(model.StageK8sBaselineRollout, f)
	r.Register(model.StageK8sBaselineClean, f)
	r.Register(model.StageK8sTrafficRouting, f)

	r.RegisterRollback(model.ApplicationKind_KUBERNETES, f)
}

func (e *Executor) Execute(sig executor.StopSignal) model.StageStatus {
	e.config = e.DeploymentConfig.KubernetesDeploymentSpec
	if e.config == nil {
		e.LogPersister.AppendError("Malformed deployment configuration: missing KubernetesDeploymentSpec")
		return model.StageStatus_STAGE_FAILURE
	}

	var (
		ctx    = sig.Context()
		appDir = filepath.Join(e.RepoDir, e.Deployment.GitPath.Path)
	)
	e.provider = provider.NewProvider(e.Deployment.ApplicationName, appDir, e.RepoDir, e.config.Input, e.Logger)

	e.Logger.Info("start executing kubernetes stage",
		zap.String("stage-name", e.Stage.Name),
		zap.String("app-dir", appDir),
	)

	var (
		originalStatus = e.Stage.Status
		status         model.StageStatus
	)

	switch model.Stage(e.Stage.Name) {
	case model.StageK8sSync:
		status = e.ensureSync(ctx)

	case model.StageK8sPrimaryRollout:
		status = e.ensurePrimaryRollout(ctx)

	case model.StageK8sCanaryRollout:
		status = e.ensureCanaryRollout(ctx)

	case model.StageK8sCanaryClean:
		status = e.ensureCanaryClean(ctx)

	case model.StageK8sBaselineRollout:
		status = e.ensureBaselineRollout(ctx)

	case model.StageK8sBaselineClean:
		status = e.ensureBaselineClean(ctx)

	case model.StageK8sTrafficRouting:
		status = e.ensureTrafficRouting(ctx)

	case model.StageRollback:
		status = e.ensureRollback(ctx)

	default:
		e.LogPersister.AppendError(fmt.Sprintf("Unsupported stage %s for kubernetes application", e.Stage.Name))
		return model.StageStatus_STAGE_FAILURE
	}

	return executor.DetermineStageStatus(sig.Signal(), originalStatus, status)
}

func (e *Executor) loadManifests(ctx context.Context) ([]provider.Manifest, error) {
	cache := provider.AppManifestsCache{
		AppID:  e.Deployment.ApplicationId,
		Cache:  e.AppManifestsCache,
		Logger: e.Logger,
	}
	manifests, ok := cache.Get(e.Deployment.Trigger.Commit.Hash)
	if ok {
		return manifests, nil
	}

	// When the manifests were not in the cache we have to load them.
	manifests, err := e.provider.LoadManifests(ctx)
	if err != nil {
		return nil, err
	}
	cache.Put(e.Deployment.Trigger.Commit.Hash, manifests)

	return manifests, nil
}

func (e *Executor) loadRunningManifests(ctx context.Context) (manifests []provider.Manifest, err error) {
	runningCommit := e.Deployment.RunningCommitHash
	if runningCommit == "" {
		return nil, fmt.Errorf("Unable to determine running commit")
	}

	cache := provider.AppManifestsCache{
		AppID:  e.Deployment.ApplicationId,
		Cache:  e.AppManifestsCache,
		Logger: e.Logger,
	}
	manifests, ok := cache.Get(runningCommit)
	if ok {
		return manifests, nil
	}

	// When the manifests were not in the cache we have to load them.
	var (
		appDir = filepath.Join(e.RepoDir, e.Deployment.GitPath.Path)
		p      = provider.NewProvider(e.Deployment.ApplicationName, appDir, e.RunningRepoDir, e.config.Input, e.Logger)
	)
	manifests, err = p.LoadManifests(ctx)
	if err != nil {
		return nil, err
	}
	cache.Put(runningCommit, manifests)

	return manifests, nil
}

func (e *Executor) builtinAnnotations(m provider.Manifest, variant, hash string) map[string]string {
	return map[string]string{
		provider.LabelManagedBy:          provider.ManagedByPiped,
		provider.LabelPiped:              e.PipedConfig.PipedID,
		provider.LabelApplication:        e.Deployment.ApplicationId,
		variantLabel:                     variant,
		provider.LabelOriginalAPIVersion: m.Key.APIVersion,
		provider.LabelResourceKey:        m.Key.String(),
		provider.LabelCommitHash:         hash,
	}
}

func findManifests(kind, name string, manifests []provider.Manifest) []provider.Manifest {
	var out []provider.Manifest
	for _, m := range manifests {
		if m.Key.Kind != kind {
			continue
		}
		if name != "" && m.Key.Name != name {
			continue
		}
		out = append(out, m)
	}
	return out
}

func findConfigMapManifests(manifests []provider.Manifest) []provider.Manifest {
	var out []provider.Manifest
	for _, m := range manifests {
		if !m.Key.IsConfigMap() {
			continue
		}
		out = append(out, m)
	}
	return out
}

func findSecretManifests(manifests []provider.Manifest) []provider.Manifest {
	var out []provider.Manifest
	for _, m := range manifests {
		if !m.Key.IsSecret() {
			continue
		}
		out = append(out, m)
	}
	return out
}

func generateServiceManifests(services []provider.Manifest, variant, nameSuffix string) ([]provider.Manifest, error) {
	manifests := make([]provider.Manifest, 0, len(services))
	updateService := func(s *corev1.Service) {
		s.Name = makeSuffixedName(s.Name, nameSuffix)
		// Currently, we suppose that all generated services should be ClusterIP.
		s.Spec.Type = corev1.ServiceTypeClusterIP
		// Append the variant label to the selector
		// to ensure that the generated service is using only workloads of this variant.
		if s.Spec.Selector == nil {
			s.Spec.Selector = map[string]string{}
		}
		s.Spec.Selector[variantLabel] = variant
		// Empty all unneeded fields.
		s.Spec.ExternalIPs = nil
		s.Spec.LoadBalancerIP = ""
		s.Spec.LoadBalancerSourceRanges = nil
	}

	for _, m := range services {
		s := &corev1.Service{}
		if err := m.ConvertToStructuredObject(s); err != nil {
			return nil, err
		}
		updateService(s)
		manifest, err := provider.ParseFromStructuredObject(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Service object to Manifest: %w", err)
		}
		manifests = append(manifests, manifest)
	}
	return manifests, nil
}

func generateWorkloadManifests(workloads, configmaps, secrets []provider.Manifest, variant, nameSuffix string, replicasCalculator func(*int32) int32) ([]provider.Manifest, error) {
	manifests := make([]provider.Manifest, 0, len(workloads))

	cmNames := make(map[string]struct{}, len(configmaps))
	for i := range configmaps {
		cmNames[configmaps[i].Key.Name] = struct{}{}
	}

	secretNames := make(map[string]struct{}, len(secrets))
	for i := range secrets {
		secretNames[secrets[i].Key.Name] = struct{}{}
	}

	updatePod := func(pod *corev1.PodTemplateSpec) {
		// Add variant labels.
		if pod.Labels == nil {
			pod.Labels = map[string]string{}
		}
		pod.Labels[variantLabel] = variant

		// Update volumes to use canary's ConfigMaps and Secrets.
		for i := range pod.Spec.Volumes {
			if cm := pod.Spec.Volumes[i].ConfigMap; cm != nil {
				if _, ok := cmNames[cm.Name]; ok {
					cm.Name = makeSuffixedName(cm.Name, nameSuffix)
				}
			}
			if s := pod.Spec.Volumes[i].Secret; s != nil {
				if _, ok := secretNames[s.SecretName]; ok {
					s.SecretName = makeSuffixedName(s.SecretName, nameSuffix)
				}
			}
		}
	}

	updateDeployment := func(d *appsv1.Deployment) {
		d.Name = makeSuffixedName(d.Name, nameSuffix)
		if replicasCalculator != nil {
			replicas := replicasCalculator(d.Spec.Replicas)
			d.Spec.Replicas = &replicas
		}
		d.Spec.Selector = metav1.AddLabelToSelector(d.Spec.Selector, variantLabel, variant)
		updatePod(&d.Spec.Template)
	}

	for _, m := range workloads {
		switch m.Key.Kind {
		case provider.KindDeployment:
			d := &appsv1.Deployment{}
			if err := m.ConvertToStructuredObject(d); err != nil {
				return nil, err
			}
			updateDeployment(d)
			manifest, err := provider.ParseFromStructuredObject(d)
			if err != nil {
				return nil, err
			}
			manifests = append(manifests, manifest)

		default:
			return nil, fmt.Errorf("unsupported workload kind %s", m.Key.Kind)
		}
	}

	return manifests, nil
}

func makeSuffixedName(name, suffix string) string {
	if suffix != "" {
		return name + "-" + suffix
	}
	return name
}
