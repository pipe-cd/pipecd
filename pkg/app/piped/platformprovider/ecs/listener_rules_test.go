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
	"testing"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadListenerRules(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		listenerRules config.ECSListenerRules
		wantRules     []string
		wantErr       error
	}{
		{
			name: "empty listener rules",
			listenerRules: config.ECSListenerRules{
				Rules: []string{},
			},
			wantRules: nil,
			wantErr:   ErrNoListenerRule,
		},
		{
			name: "single listener rule",
			listenerRules: config.ECSListenerRules{
				Rules: []string{"rule1"},
			},
			wantRules: []string{"rule1"},
			wantErr:   nil,
		},
		{
			name: "multiple listener rules",
			listenerRules: config.ECSListenerRules{
				Rules: []string{"rule1", "rule2", "rule3"},
			},
			wantRules: []string{"rule1", "rule2", "rule3"},
			wantErr:   nil,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rules, err := loadListenerRules(tc.listenerRules)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRules, rules)
		})
	}
}
