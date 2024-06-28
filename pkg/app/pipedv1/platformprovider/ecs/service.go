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
	"os"

	"sigs.k8s.io/yaml"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func loadServiceDefinition(path string) (types.Service, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return types.Service{}, err
	}
	return parseServiceDefinition(data)
}

func parseServiceDefinition(data []byte) (types.Service, error) {
	var obj types.Service
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return types.Service{}, err
	}

	if obj.ClusterArn == nil {
		// Rename cluster field to clusterArn if exist
		clusterArn, err := parseServiceDefinitionForCluster(data)
		if err != nil {
			return types.Service{}, err
		}
		obj.ClusterArn = &clusterArn
	}

	if obj.RoleArn == nil {
		// Rename role field to roleArn if exist
		roleArn, err := parseServiceDefinitionForRole(data)
		if err != nil {
			return types.Service{}, err
		}
		obj.RoleArn = &roleArn
	}

	return obj, nil
}

func parseServiceDefinitionForCluster(data []byte) (string, error) {
	var obj struct {
		Cluster string `json:"cluster"`
	}
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return "", err
	}
	return obj.Cluster, nil
}

func parseServiceDefinitionForRole(data []byte) (string, error) {
	var obj struct {
		Role string `json:"role"`
	}
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return "", err
	}
	return obj.Role, nil
}
