// Copyright 2025 The PipeCD Authors.
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

package deployment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

func ensureVariantSelectorInWorkload(m provider.Manifest, variantLabel, variant string) error {
	variantMap := map[string]string{
		variantLabel: variant,
	}
	if err := m.AddStringMapValues(variantMap, "spec", "selector", "matchLabels"); err != nil {
		return err
	}
	return m.AddStringMapValues(variantMap, "spec", "template", "metadata", "labels")
}

func checkVariantSelectorInWorkload(manifest provider.Manifest, variantLabel, variant string) error {
	var (
		matchLabelsFields = []string{"spec", "selector", "matchLabels"}
		labelsFields      = []string{"spec", "template", "metadata", "labels"}
	)

	value, ok, err := manifest.NestedString(append(matchLabelsFields, variantLabel)...)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("missing %s key in spec.selector.matchLabels", variantLabel)
	}
	if value != variant {
		return fmt.Errorf("require %s but got %s for %s key in %s", variant, value, variantLabel, strings.Join(matchLabelsFields, "."))
	}

	value, ok, err = manifest.NestedString(append(labelsFields, variantLabel)...)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("missing %s key in spec.template.metadata.labels", variantLabel)
	}
	if value != variant {
		return fmt.Errorf("require %s but got %s for %s key in %s", variant, value, variantLabel, strings.Join(labelsFields, "."))
	}

	return nil
}

// generateVariantServiceManifests generates Service manifests for the specified variant.
// It duplicates the given Service manifests, adds a name suffix, sets type to ClusterIP,
// appends the variant label to the selector, and clears unnecessary fields.
func generateVariantServiceManifests(services []provider.Manifest, variantLabel, variant, nameSuffix string) ([]provider.Manifest, error) {
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
		manifest, err := provider.FromStructuredObject(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Service object to Manifest: %w", err)
		}
		// This is because the resource key differs between variants because of the name suffix.
		// For example, The Service named "simple" has the resource key ":Service:some-namespace:simple"
		// and its baseline variant has the resource key ":Service:some-namespace:simple-baseline".
		manifest.AddAnnotations(map[string]string{
			provider.LabelResourceKey: manifest.Key().String(),
		})
		manifests = append(manifests, manifest)
	}
	return manifests, nil
}

// generateVariantWorkloadManifests generates Workload manifests for the specified variant.
// It duplicates the given Workload manifests, adds a name suffix, sets the variant label to the selector,
// and updates the ENV references in containers to use canary's ConfigMaps and Secrets.
func generateVariantWorkloadManifests(workloads, configmaps, secrets []provider.Manifest, variantLabel, variant, nameSuffix string, replicasCalculator func(*int32) int32) ([]provider.Manifest, error) {
	manifests := make([]provider.Manifest, 0, len(workloads))

	cmNames := make(map[string]struct{}, len(configmaps))
	for _, cm := range configmaps {
		cmNames[cm.Name()] = struct{}{}
	}

	secretNames := make(map[string]struct{}, len(secrets))
	for _, secret := range secrets {
		secretNames[secret.Name()] = struct{}{}
	}

	updateContainers := func(containers []corev1.Container) {
		for _, container := range containers {
			for _, env := range container.Env {
				if v := env.ValueFrom; v != nil {
					if ref := v.ConfigMapKeyRef; ref != nil {
						if _, ok := cmNames[ref.Name]; ok {
							ref.Name = makeSuffixedName(ref.Name, nameSuffix)
						}
					}
					if ref := v.SecretKeyRef; ref != nil {
						if _, ok := secretNames[ref.Name]; ok {
							ref.Name = makeSuffixedName(ref.Name, nameSuffix)
						}
					}
				}
			}
			for _, envFrom := range container.EnvFrom {
				if ref := envFrom.ConfigMapRef; ref != nil {
					if _, ok := cmNames[ref.Name]; ok {
						ref.Name = makeSuffixedName(ref.Name, nameSuffix)
					}
				}
				if ref := envFrom.SecretRef; ref != nil {
					if _, ok := secretNames[ref.Name]; ok {
						ref.Name = makeSuffixedName(ref.Name, nameSuffix)
					}
				}
			}
		}
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

		// Update ENV references in containers.
		updateContainers(pod.Spec.InitContainers)
		updateContainers(pod.Spec.Containers)
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
		switch m.Kind() {
		case provider.KindDeployment:
			d := &appsv1.Deployment{}
			if err := m.ConvertToStructuredObject(d); err != nil {
				return nil, err
			}
			updateDeployment(d)
			manifest, err := provider.FromStructuredObject(d)
			if err != nil {
				return nil, err
			}
			// This is because the resource key differs between variants because of the name suffix.
			// For example, The Deployment named "simple" has the resource key "apps:Deployment:some-namespace:simple"
			// and its baseline variant has the resource key "apps:Deployment:some-namespace:simple-baseline".
			manifest.AddAnnotations(map[string]string{
				provider.LabelResourceKey: manifest.Key().String(),
			})
			manifests = append(manifests, manifest)

		default:
			return nil, fmt.Errorf("unsupported workload kind %s", m.Kind())
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

// addVariantLabelsAndAnnotations adds the variant label and annotation to the given manifests.
func addVariantLabelsAndAnnotations(m []provider.Manifest, variantLabel, variant string) {
	for _, m := range m {
		m.AddLabels(map[string]string{
			variantLabel: variant,
		})
		m.AddAnnotations(map[string]string{
			variantLabel: variant,
		})
	}
}

// duplicateManifests duplicates the given manifests and appends a name suffix to each manifest.
func duplicateManifests(manifests []provider.Manifest, nameSuffix string) []provider.Manifest {
	copied := make([]provider.Manifest, len(manifests))
	for i, m := range manifests {
		copied[i] = m.DeepCopyWithName(makeSuffixedName(m.Name(), nameSuffix))
	}
	return copied
}

// deleteResources deletes the given resources.
// It returns the number of deleted resources.
func deleteResources(ctx context.Context, lp sdk.StageLogPersister, applier *provider.Applier, keys []provider.ResourceKey) int {
	var deletedCount int

	for _, k := range keys {
		if err := applier.Delete(ctx, k); err != nil {
			if errors.Is(err, provider.ErrNotFound) {
				lp.Infof("Specified resource does not exist, so skip deleting the resource: %s (%v)", k.ReadableString(), err)
				continue
			}
			lp.Errorf("Failed while deleting resource %s (%v)", k.ReadableString(), err)
			continue // continue to delete other resources
		}
		deletedCount++
		lp.Successf("- deleted resource: %s", k.ReadableString())
	}

	return deletedCount
}

// deleteVariantResources deletes the resources of the specified variant.
// It finds the resources of the specified variant and deletes them.
// It deletes the resources in the order of Service -> Workload -> Others -> Cluster-scoped resources.
func deleteVariantResources(ctx context.Context, lp sdk.StageLogPersister, kubectl *provider.Kubectl, kubeConfig string, applier *provider.Applier, applicationID, variantLabel, variant string) error {
	namespacedLiveResources, clusterScopedLiveResources, err := provider.GetLiveResources(ctx, kubectl, kubeConfig, applicationID, fmt.Sprintf("%s=%s", variantLabel, variant))
	if err != nil {
		return err
	}

	services := make([]provider.ResourceKey, 0, len(namespacedLiveResources))
	workloads := make([]provider.ResourceKey, 0, len(namespacedLiveResources))
	others := make([]provider.ResourceKey, 0, len(namespacedLiveResources))
	clusterScoped := make([]provider.ResourceKey, 0, len(clusterScopedLiveResources))

	for _, r := range namespacedLiveResources {
		switch {
		case r.IsService():
			services = append(services, r.Key())
		case r.IsWorkload():
			workloads = append(workloads, r.Key())
		default:
			others = append(others, r.Key())
		}
	}

	for _, r := range clusterScopedLiveResources {
		clusterScoped = append(clusterScoped, r.Key())
	}

	var deletedCount int
	deletedCount += deleteResources(ctx, lp, applier, services)
	deletedCount += deleteResources(ctx, lp, applier, workloads)
	deletedCount += deleteResources(ctx, lp, applier, others)
	deletedCount += deleteResources(ctx, lp, applier, clusterScoped)
	lp.Successf("Successfully deleted %d resources", deletedCount)

	return nil
}
