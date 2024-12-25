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
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	KindDeployment = "Deployment"
	KindSecret     = "Secret"
	KindConfigMap  = "ConfigMap"

	DefaultNamespace = "default"
)

type ResourceKey struct {
	groupKind schema.GroupKind
	namespace string
	name      string
}

func (k ResourceKey) Kind() string {
	return k.groupKind.Kind
}

func (k ResourceKey) Name() string {
	return k.name
}

func (k ResourceKey) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", k.groupKind.Group, k.groupKind.Kind, k.namespace, k.name)
}

func (k ResourceKey) ReadableString() string {
	return fmt.Sprintf("name=%q, kind=%q, namespace=%q, apiGroup=%q", k.name, k.groupKind.Kind, k.namespace, k.groupKind.Group)
}

func makeResourceKey(obj *unstructured.Unstructured) ResourceKey {
	k := ResourceKey{
		groupKind: obj.GroupVersionKind().GroupKind(),
		namespace: obj.GetNamespace(),
		name:      obj.GetName(),
	}
	return k
}
