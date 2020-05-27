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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ResourceKey struct {
	APIVersion string
	Kind       string
	Namespace  string
	Name       string
}

func (k ResourceKey) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", k.APIVersion, k.Kind, k.Namespace, k.Name)
}

func (k ResourceKey) IsZero() bool {
	return k.APIVersion == "" &&
		k.Kind == "" &&
		k.Namespace == "" &&
		k.Name == ""
}

func MakeResourceKey(obj *unstructured.Unstructured) ResourceKey {
	return ResourceKey{
		APIVersion: obj.GetAPIVersion(),
		Kind:       obj.GetKind(),
		Namespace:  obj.GetNamespace(),
		Name:       obj.GetName(),
	}
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
