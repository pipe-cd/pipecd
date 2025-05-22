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
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"

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
		manifests = append(manifests, manifest)
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
