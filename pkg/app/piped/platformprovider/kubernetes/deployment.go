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
	"sort"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func FindReferencingConfigMapsInDeployment(d *appsv1.Deployment) []string {
	m := make(map[string]struct{}, 0)

	// Find all configmaps specified in Volumes.
	for _, v := range d.Spec.Template.Spec.Volumes {
		if cm := v.ConfigMap; cm != nil {
			m[cm.Name] = struct{}{}
		}
	}

	findInContainers := func(containers []corev1.Container) {
		for _, c := range containers {
			for _, env := range c.Env {
				if source := env.ValueFrom; source != nil {
					if ref := source.ConfigMapKeyRef; ref != nil {
						m[ref.Name] = struct{}{}
					}
				}
			}
			for _, env := range c.EnvFrom {
				if ref := env.ConfigMapRef; ref != nil {
					m[ref.Name] = struct{}{}
				}
			}
		}
	}

	// Find all configmaps specified in Env.
	findInContainers(d.Spec.Template.Spec.Containers)
	findInContainers(d.Spec.Template.Spec.InitContainers)

	if len(m) == 0 {
		return nil
	}

	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)

	return out
}

func FindReferencingSecretsInDeployment(d *appsv1.Deployment) []string {
	m := make(map[string]struct{}, 0)

	// Find all secrets specified in Volumes.
	for _, v := range d.Spec.Template.Spec.Volumes {
		if s := v.Secret; s != nil {
			m[s.SecretName] = struct{}{}
		}
	}

	findInContainers := func(containers []corev1.Container) {
		for _, c := range containers {
			for _, env := range c.Env {
				if source := env.ValueFrom; source != nil {
					if ref := source.SecretKeyRef; ref != nil {
						m[ref.Name] = struct{}{}
					}
				}
			}
			for _, env := range c.EnvFrom {
				if ref := env.SecretRef; ref != nil {
					m[ref.Name] = struct{}{}
				}
			}
		}
	}

	// Find all secrets specified in Env.
	findInContainers(d.Spec.Template.Spec.Containers)
	findInContainers(d.Spec.Template.Spec.InitContainers)

	if len(m) == 0 {
		return nil
	}

	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)

	return out
}
