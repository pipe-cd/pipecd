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

package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAnalysisTemplate(t *testing.T) {
	testcases := []struct {
		name          string
		repoDir       string
		expectedSpec  interface{}
		expectedError error
	}{
		{
			name:    "Load analysis template successfully",
			repoDir: "testdata",
			expectedSpec: &AnalysisTemplateSpec{
				Metrics: map[string]AnalysisMetrics{
					"app_http_error_percentage": {
						Strategy:  AnalysisStrategyThreshold,
						Query:     "http_error_percentage{env={{ .App.Env }}, app={{ .App.Name }}}",
						Expected:  AnalysisExpected{Max: floatPointer(0.1)},
						Interval:  Duration(time.Minute),
						Timeout:   Duration(30 * time.Second),
						Provider:  "datadog-dev",
						Deviation: AnalysisDeviationEither,
					},
					"container_cpu_usage_seconds_total": {
						Strategy:     AnalysisStrategyThreshold,
						Query:        "sum(\n  max(kube_pod_labels{label_app=~\"{{ .App.Name }}\", label_pipecd_dev_variant=~\"canary\"}) by (label_app, label_pipecd_dev_variant, pod)\n  *\n  on(pod)\n  group_right(label_app, label_pipecd_dev_variant)\n  label_replace(\n    sum by(pod_name) (\n      rate(container_cpu_usage_seconds_total{namespace=\"default\"}[5m])\n    ), \"pod\", \"$1\", \"pod_name\", \"(.+)\"\n  )\n) by (label_app, label_pipecd_dev_variant)\n",
						Expected:     AnalysisExpected{Max: floatPointer(0.0001)},
						FailureLimit: 2,
						Interval:     Duration(10 * time.Second),
						Timeout:      Duration(30 * time.Second),
						Provider:     "prometheus-dev",
						Deviation:    AnalysisDeviationEither,
					},
					"grpc_error_rate-percentage": {
						Strategy:     AnalysisStrategyThreshold,
						Query:        "100 - sum(\n    rate(\n        grpc_server_handled_total{\n          grpc_code!=\"OK\",\n          kubernetes_namespace=\"{{ .Args.namespace }}\",\n          kubernetes_pod_name=~\"{{ .App.Name }}-[0-9a-zA-Z]+(-[0-9a-zA-Z]+)\"\n        }[{{ .Args.interval }}]\n    )\n)\n/\nsum(\n    rate(\n        grpc_server_started_total{\n          kubernetes_namespace=\"{{ .Args.namespace }}\",\n          kubernetes_pod_name=~\"{{ .App.Name }}-[0-9a-zA-Z]+(-[0-9a-zA-Z]+)\"\n        }[{{ .Args.interval }}]\n    )\n) * 100\n",
						Expected:     AnalysisExpected{Max: floatPointer(10)},
						FailureLimit: 1,
						Interval:     Duration(time.Minute),
						Timeout:      Duration(30 * time.Second),
						Provider:     "prometheus-dev",
						Deviation:    AnalysisDeviationEither,
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			spec, err := LoadAnalysisTemplate(tc.repoDir)
			require.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedSpec, spec)
			}
		})
	}
}
