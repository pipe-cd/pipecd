// Copyright 2020 The Pipe Authors.
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

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	testcases := []struct {
		fileName      string
		expected      *Config
		expectedError error
	}{
		// {
		// 	fileName: "testdata/k8s-plain-apply.yaml",
		// 	expected: &Config{
		// 		Version: "v1",
		// 		Kind:    "K8sApp",
		// 		Name:    "account",
		// 		Stages: []*Stage{
		// 			&Stage{
		// 				Name: "K8S_APPLY",
		// 				Desc: "Rolling Update",
		// 			},
		// 			&Stage{
		// 				Name: "VERIFICATION",
		// 				Desc: "Smoke Test",
		// 			},
		// 		},
		// 	},
		// 	expectedError: nil,
		// },
	}
	for _, tc := range testcases {
		t.Run(tc.fileName, func(t *testing.T) {
			cfg, err := LoadFromYAML(tc.fileName)
			assert.Equal(t, tc.expected, cfg)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
