// Copyright 2026 The PipeCD Authors.
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
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
)

type containerImage struct {
	name   string
	tag    string
	digest string
}

// parseContainerImage parses an ECS container image reference into its components.
//
// Supported formats: [registry/]name[:tag|@digest]
func parseContainerImage(image string) (img containerImage) {
	ref := image

	if idx := strings.Index(ref, "@"); idx != -1 {
		img.digest = ref[idx+1:]
		ref = ref[:idx]
	}

	parts := strings.Split(ref, "/")
	last := parts[len(parts)-1]

	// Extract tag from the last segment only when there is no digest
	if img.digest == "" {
		if idx := strings.LastIndex(last, ":"); idx != -1 {
			img.tag = last[idx+1:]
			last = last[:idx]
		}
	}

	img.name = last
	return
}

// determineVersions extracts artifact versions from an ECS task definition.
//
// It finds all container images defined in the task definition's ContainerDefinitions and returns their names and tags.
//
// Duplicate image references are deduplicated.
func determineVersions(taskDef types.TaskDefinition) []sdk.ArtifactVersion {
	imageMap := map[string]struct{}{}
	for _, c := range taskDef.ContainerDefinitions {
		if c.Image == nil || *c.Image == "" {
			continue
		}
		imageMap[*c.Image] = struct{}{}
	}

	versions := make([]sdk.ArtifactVersion, 0, len(imageMap))
	for i := range imageMap {
		image := parseContainerImage(i)
		version := image.tag
		if version == "" {
			version = image.digest
		}
		versions = append(versions, sdk.ArtifactVersion{
			Version: version,
			Name:    image.name,
			URL:     i,
		})
	}
	return versions
}
