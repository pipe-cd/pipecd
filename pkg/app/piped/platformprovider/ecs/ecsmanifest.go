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

package ecs

import (
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

// ECSManifest is defined just for wrapping TaskDefinition and Service.
type ECSManifest struct {
	TaskDefinition    *types.TaskDefinition
	ServiceDefinition *types.Service
}

// Unstructured returns the unstructured representation of the ECSManifest.
func (m ECSManifest) Unstructured() (unstructured.Unstructured, error) {
	manifest := struct {
		TaskDefinition    *types.TaskDefinition
		ServiceDefinition *types.Service
		ApiVersion        string `json:"apiVersion"`
		Kind              string `json:"kind"`
	}{
		TaskDefinition:    m.TaskDefinition,
		ServiceDefinition: m.ServiceDefinition,
		// Append ApiVersion and Kind here just to avoid error
		// when unmarshalling the yaml to unstructured.Unstructured.
		// These do not affect the manifest's content.
		ApiVersion: "__pipecd.dev/v1beta1__",
		Kind:       "__PipeCDECSManifes__",
	}

	y, err := yaml.Marshal(manifest)
	if err != nil {
		return unstructured.Unstructured{}, err
	}

	var unstructuredManifest unstructured.Unstructured
	err = yaml.Unmarshal(y, &unstructuredManifest)
	if err != nil {
		return unstructured.Unstructured{}, err
	}

	return unstructuredManifest, nil
}

// LoadECSManifest loads the taskDefinition and serviceDefinition from the given files.
func LoadECSManifest(appDir, taskDefFilename, serviceDefFilename string) (ECSManifest, error) {
	service, err := LoadServiceDefinition(appDir, serviceDefFilename)
	if err != nil {
		return ECSManifest{}, err
	}
	task, err := LoadTaskDefinition(appDir, taskDefFilename)
	if err != nil {
		return ECSManifest{}, err
	}

	return ECSManifest{
		TaskDefinition:    &task,
		ServiceDefinition: &service,
	}, nil
}
