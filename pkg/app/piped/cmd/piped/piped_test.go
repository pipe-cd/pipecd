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

package piped

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasTooManyConfigFlags(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		title     string
		p         *piped
		expectErr bool
	}{
		{
			title: "no config",
			p: &piped{
				configFile:            "",
				configGCPSecret:       "",
				configAWSSecret:       "",
				configAWSSsmParameter: "",
			},
			expectErr: false,
		},
		{
			title: "only one config is set",
			p: &piped{
				configFile:            "config.yaml",
				configGCPSecret:       "",
				configAWSSecret:       "",
				configAWSSsmParameter: "",
			},
			expectErr: false,
		},
		{
			title: "two configs are set",
			p: &piped{
				configFile:            "config.yaml",
				configGCPSecret:       "xxx",
				configAWSSecret:       "",
				configAWSSsmParameter: "",
			},
			expectErr: true,
		},
		{
			title: "all configs are set",
			p: &piped{
				configFile:            "config.yaml",
				configGCPSecret:       "xxx",
				configAWSSecret:       "yyy",
				configAWSSsmParameter: "zzz",
			},
			expectErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			err := tc.p.hasTooManyConfigFlags()
			assert.Equal(t, tc.expectErr, err != nil)
		})
	}
}
