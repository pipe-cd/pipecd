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

package planpreview

import (
	"testing"

	"github.com/pipe-cd/pipecd/pkg/model"

	"github.com/stretchr/testify/require"
)

func TestSortResults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []*model.ApplicationPlanPreviewResult
		expected []*model.ApplicationPlanPreviewResult
	}{
		{
			name: "sort by name, kind",
			input: []*model.ApplicationPlanPreviewResult{
				{ApplicationName: "appB", ApplicationKind: model.ApplicationKind_CLOUDRUN},
				{ApplicationName: "appB", ApplicationKind: model.ApplicationKind_KUBERNETES},
				{ApplicationName: "appA", ApplicationKind: model.ApplicationKind_CLOUDRUN},
				{ApplicationName: "appA", ApplicationKind: model.ApplicationKind_KUBERNETES},
			},
			expected: []*model.ApplicationPlanPreviewResult{
				{ApplicationName: "appA", ApplicationKind: model.ApplicationKind_KUBERNETES},
				{ApplicationName: "appA", ApplicationKind: model.ApplicationKind_CLOUDRUN},
				{ApplicationName: "appB", ApplicationKind: model.ApplicationKind_KUBERNETES},
				{ApplicationName: "appB", ApplicationKind: model.ApplicationKind_CLOUDRUN},
			},
		},
		{
			name: "sort by name, env, kind",
			input: []*model.ApplicationPlanPreviewResult{
				{ApplicationName: "appB", Labels: map[string]string{labelEnvKey: "dev"}, ApplicationKind: model.ApplicationKind_CLOUDRUN},
				{ApplicationName: "appB", Labels: map[string]string{labelEnvKey: "dev"}, ApplicationKind: model.ApplicationKind_KUBERNETES},
				{ApplicationName: "appA", Labels: map[string]string{labelEnvKey: "prd"}, ApplicationKind: model.ApplicationKind_CLOUDRUN},
				{ApplicationName: "appA", Labels: map[string]string{labelEnvKey: "dev"}, ApplicationKind: model.ApplicationKind_KUBERNETES},
			},
			expected: []*model.ApplicationPlanPreviewResult{
				{ApplicationName: "appA", Labels: map[string]string{labelEnvKey: "dev"}, ApplicationKind: model.ApplicationKind_KUBERNETES},
				{ApplicationName: "appA", Labels: map[string]string{labelEnvKey: "prd"}, ApplicationKind: model.ApplicationKind_CLOUDRUN},
				{ApplicationName: "appB", Labels: map[string]string{labelEnvKey: "dev"}, ApplicationKind: model.ApplicationKind_KUBERNETES},
				{ApplicationName: "appB", Labels: map[string]string{labelEnvKey: "dev"}, ApplicationKind: model.ApplicationKind_CLOUDRUN},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sortResults(tt.input)
			require.Equal(t, tt.expected, tt.input)
		})
	}
}
