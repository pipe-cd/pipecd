// Copyright 2023 The PipeCD Authors.
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
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func loadTaskDefinition(path string) (types.TaskDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return types.TaskDefinition{}, err
	}
	return parseTaskDefinition(data)
}

func parseTaskDefinition(data []byte) (types.TaskDefinition, error) {
	var obj types.TaskDefinition
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return types.TaskDefinition{}, err
	}
	return obj, nil
}

// FindImageTag parses image tag from given ECS task definition.
func FindImageTag(taskDefinition types.TaskDefinition) (string, error) {
	if len(taskDefinition.ContainerDefinitions) == 0 {
		return "", fmt.Errorf("container definition could not be empty")
	}
	name, tag := parseContainerImage(*taskDefinition.ContainerDefinitions[0].Image)
	if name == "" {
		return "", fmt.Errorf("image name could not be empty")
	}
	return tag, nil
}

func parseContainerImage(image string) (name, tag string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		tag = parts[1]
	}
	paths := strings.Split(parts[0], "/")
	name = paths[len(paths)-1]
	return
}

// FindArtifactVersions parses artifact versions from ECS task definition.
func FindArtifactVersions(taskDefinition types.TaskDefinition) ([]*model.ArtifactVersion, error) {
	if len(taskDefinition.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("container definition could not be empty")
	}

	// Remove duplicate images.
	imageMap := map[string]struct{}{}
	for _, cd := range taskDefinition.ContainerDefinitions {
		imageMap[*cd.Image] = struct{}{}
	}

	versions := make([]*model.ArtifactVersion, 0, len(imageMap))
	for i := range imageMap {
		name, tag := parseContainerImage(i)
		if name == "" {
			return nil, fmt.Errorf("image name could not be empty")
		}

		versions = append(versions, &model.ArtifactVersion{
			Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
			Version: tag,
			Name:    name,
			Url:     i,
		})
	}

	return versions, nil
}
