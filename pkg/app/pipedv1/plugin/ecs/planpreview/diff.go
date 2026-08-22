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

package planpreview

import (
	"fmt"
	"strings"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pmezard/go-difflib/difflib"
	"sigs.k8s.io/yaml"
)

// diffDefinitions marshals old and new to YAML and returns a unified diff string
func diffDefinitions[T any](old, new *T, name string) (string, error) {
	var oldYAML string
	if old != nil {
		data, err := yaml.Marshal(old)
		if err != nil {
			return "", fmt.Errorf("failed to marshal old %s: %w", name, err)
		}
		oldYAML = string(data)
	}

	newData, err := yaml.Marshal(new)
	if err != nil {
		return "", fmt.Errorf("failed to marshal new %s: %w", name, err)
	}

	ud := difflib.UnifiedDiff{
		A:        difflib.SplitLines(oldYAML),
		B:        difflib.SplitLines(string(newData)),
		FromFile: fmt.Sprintf("%s (running)", name),
		ToFile:   fmt.Sprintf("%s (target)", name),
		Context:  3,
	}
	return difflib.GetUnifiedDiffString(ud)
}

func toResponse(deployTarget, taskDefDiff, serviceDiff string) *sdk.GetPlanPreviewResponse {
	details := buildDetails(taskDefDiff, serviceDiff)
	noChange := taskDefDiff == "" && serviceDiff == ""

	var summary string
	if noChange {
		summary = "No changes were detected"
	} else {
		summary = buildSummary(taskDefDiff, serviceDiff)
	}

	return &sdk.GetPlanPreviewResponse{
		Results: []sdk.PlanPreviewResult{
			{
				DeployTarget: deployTarget,
				NoChange:     noChange,
				Summary:      summary,
				Details:      details,
				DiffLanguage: "diff",
			},
		},
	}
}

func buildSummary(taskDefDiff, serviceDiff string) string {
	var parts []string
	if taskDefDiff != "" {
		parts = append(parts, "task definition changed")
	}
	if serviceDiff != "" {
		parts = append(parts, "service definition changed")
	}
	return strings.Join(parts, ", ")
}

func buildDetails(taskDefDiff, serviceDiff string) []byte {
	var sb strings.Builder
	if taskDefDiff != "" {
		sb.WriteString(taskDefDiff)
	}
	if serviceDiff != "" {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(serviceDiff)
	}
	if sb.Len() == 0 {
		return nil
	}
	return []byte(sb.String())
}
