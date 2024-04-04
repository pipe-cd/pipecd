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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var builtInAPIVersions = map[string]struct{}{
	"admissionregistration.k8s.io/v1":      {},
	"admissionregistration.k8s.io/v1beta1": {},
	"apiextensions.k8s.io/v1":              {},
	"apiextensions.k8s.io/v1beta1":         {},
	"apiregistration.k8s.io/v1":            {},
	"apiregistration.k8s.io/v1beta1":       {},
	"apps/v1":                              {},
	"authentication.k8s.io/v1":             {},
	"authentication.k8s.io/v1beta1":        {},
	"authorization.k8s.io/v1":              {},
	"authorization.k8s.io/v1beta1":         {},
	"autoscaling/v1":                       {},
	"autoscaling/v2beta1":                  {},
	"autoscaling/v2beta2":                  {},
	"batch/v1":                             {},
	"batch/v1beta1":                        {},
	"certificates.k8s.io/v1beta1":          {},
	"coordination.k8s.io/v1":               {},
	"coordination.k8s.io/v1beta1":          {},
	"extensions/v1beta1":                   {},
	"internal.autoscaling.k8s.io/v1alpha1": {},
	"metrics.k8s.io/v1beta1":               {},
	"networking.k8s.io/v1":                 {},
	"networking.k8s.io/v1beta1":            {},
	"node.k8s.io/v1beta1":                  {},
	"policy/v1":                            {},
	"policy/v1beta1":                       {},
	"rbac.authorization.k8s.io/v1":         {},
	"rbac.authorization.k8s.io/v1beta1":    {},
	"scheduling.k8s.io/v1":                 {},
	"scheduling.k8s.io/v1beta1":            {},
	"storage.k8s.io/v1":                    {},
	"storage.k8s.io/v1beta1":               {},
	"v1":                                   {},
}

const (
	KindDeployment               = "Deployment"
	KindStatefulSet              = "StatefulSet"
	KindDaemonSet                = "DaemonSet"
	KindReplicaSet               = "ReplicaSet"
	KindPod                      = "Pod"
	KindJob                      = "Job"
	KindCronJob                  = "CronJob"
	KindConfigMap                = "ConfigMap"
	KindSecret                   = "Secret"
	KindPersistentVolume         = "PersistentVolume"
	KindPersistentVolumeClaim    = "PersistentVolumeClaim"
	KindService                  = "Service"
	KindIngress                  = "Ingress"
	KindServiceAccount           = "ServiceAccount"
	KindRole                     = "Role"
	KindRoleBinding              = "RoleBinding"
	KindClusterRole              = "ClusterRole"
	KindClusterRoleBinding       = "ClusterRoleBinding"
	KindNameSpace                = "NameSpace"
	KindPodDisruptionBudget      = "PodDisruptionBudget"
	KindCustomResourceDefinition = "CustomResourceDefinition"

	DefaultNamespace = "default"
)

type APIVersionKind struct {
	APIVersion string
	Kind       string
}

type ResourceKey struct {
	APIVersion string
	Kind       string
	Namespace  string
	Name       string
}

func (k ResourceKey) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", k.APIVersion, k.Kind, k.Namespace, k.Name)
}

func (k ResourceKey) ReadableString() string {
	return fmt.Sprintf("name=%q, kind=%q, namespace=%q, apiVersion=%q", k.Name, k.Kind, k.Namespace, k.APIVersion)
}

func (k ResourceKey) IsZero() bool {
	return k.APIVersion == "" &&
		k.Kind == "" &&
		k.Namespace == "" &&
		k.Name == ""
}

func (k ResourceKey) IsDeployment() bool {
	if k.Kind != KindDeployment {
		return false
	}
	if !IsKubernetesBuiltInResource(k.APIVersion) {
		return false
	}
	return true
}

func (k ResourceKey) IsReplicaSet() bool {
	if k.Kind != KindReplicaSet {
		return false
	}
	if !IsKubernetesBuiltInResource(k.APIVersion) {
		return false
	}
	return true
}

func (k ResourceKey) IsWorkload() bool {
	if !IsKubernetesBuiltInResource(k.APIVersion) {
		return false
	}

	switch k.Kind {
	case KindDeployment:
		return true
	case KindReplicaSet:
		return true
	case KindDaemonSet:
		return true
	case KindPod:
		return true
	}

	return false
}

func (k ResourceKey) IsService() bool {
	if k.Kind != KindService {
		return false
	}
	if !IsKubernetesBuiltInResource(k.APIVersion) {
		return false
	}
	return true
}

func (k ResourceKey) IsConfigMap() bool {
	if k.Kind != KindConfigMap {
		return false
	}
	if !IsKubernetesBuiltInResource(k.APIVersion) {
		return false
	}
	return true
}

func (k ResourceKey) IsSecret() bool {
	if k.Kind != KindSecret {
		return false
	}
	if !IsKubernetesBuiltInResource(k.APIVersion) {
		return false
	}
	return true
}

// IsLess reports whether the key should sort before the given key.
func (k ResourceKey) IsLess(a ResourceKey) bool {
	if k.APIVersion < a.APIVersion {
		return true
	} else if k.APIVersion > a.APIVersion {
		return false
	}

	if k.Kind < a.Kind {
		return true
	} else if k.Kind > a.Kind {
		return false
	}

	if k.Namespace < a.Namespace {
		return true
	} else if k.Namespace > a.Namespace {
		return false
	}

	if k.Name < a.Name {
		return true
	} else if k.Name > a.Name {
		return false
	}
	return false
}

// IsLessWithIgnoringNamespace reports whether the key should sort before the given key,
// but this ignores the comparation of the namesapce.
func (k ResourceKey) IsLessWithIgnoringNamespace(a ResourceKey) bool {
	if k.APIVersion < a.APIVersion {
		return true
	} else if k.APIVersion > a.APIVersion {
		return false
	}

	if k.Kind < a.Kind {
		return true
	} else if k.Kind > a.Kind {
		return false
	}

	if k.Name < a.Name {
		return true
	} else if k.Name > a.Name {
		return false
	}
	return false
}

// IsEqualWithIgnoringNamespace checks whether the key is equal to the given key,
// but this ignores the comparation of the namesapce.
func (k ResourceKey) IsEqualWithIgnoringNamespace(a ResourceKey) bool {
	if k.APIVersion != a.APIVersion {
		return false
	}
	if k.Kind != a.Kind {
		return false
	}
	if k.Name != a.Name {
		return false
	}
	return true
}

// GetResourceKeyFromActualResource extracts the ResourceKey from the given resource object.
// If the ResourceKey is not found in the annotations, it will create a new one.
func GetResourceKeyFromActualResource(obj *unstructured.Unstructured) ResourceKey {
	annotation, ok := obj.GetAnnotations()[LabelResourceKey]
	if !ok {
		return MakeResourceKey(obj)
	}

	key, err := DecodeResourceKey(annotation)
	if err != nil {
		return MakeResourceKey(obj)
	}
	return key
}

func MakeResourceKey(obj *unstructured.Unstructured) ResourceKey {
	k := ResourceKey{
		APIVersion: obj.GetAPIVersion(),
		Kind:       obj.GetKind(),
		Namespace:  obj.GetNamespace(),
		Name:       obj.GetName(),
	}
	if k.Namespace == "" {
		k.Namespace = DefaultNamespace
	}
	return k
}

func DecodeResourceKey(key string) (ResourceKey, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 4 {
		return ResourceKey{}, fmt.Errorf("malformed key")
	}
	return ResourceKey{
		APIVersion: parts[0],
		Kind:       parts[1],
		Namespace:  parts[2],
		Name:       parts[3],
	}, nil
}

func IsKubernetesBuiltInResource(apiVersion string) bool {
	_, ok := builtInAPIVersions[apiVersion]
	// TODO: Change the way to detect whether an APIVersion is built-in or not
	// rather than depending on this fixed list.
	return ok
}
