// Copyright 2022 The PipeCD Authors.
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

package cloudrun

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type RevisionManifest struct {
	Name string
	u    *unstructured.Unstructured
}

func ParseRevisionManifest(data []byte) (RevisionManifest, error) {
	var obj unstructured.Unstructured
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return RevisionManifest{}, err
	}

	return RevisionManifest{
		Name: obj.GetName(),
		u:    &obj,
	}, nil
}

func (r RevisionManifest) YamlBytes() ([]byte, error) {
	return yaml.Marshal(r.u)
}
