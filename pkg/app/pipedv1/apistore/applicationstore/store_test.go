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

package applicationstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"go.uber.org/zap/zaptest"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestListByPluginName(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name       string
		storedApps []*model.Application
		plugin     string
		expected   []*model.Application
	}{
		{
			name:       "There is no stored application",
			storedApps: []*model.Application{},
			plugin:     "plugin-a",
			expected:   []*model.Application{},
		},
		{
			name: "no matching",
			storedApps: []*model.Application{
				{Id: "app-1", Plugins: []string{"plugin-b"}},
				{Id: "app-2", Plugins: []string{"plugin-c"}},
			},
			plugin:   "plugin-a",
			expected: []*model.Application{},
		},
		{
			name: "one matched application",
			storedApps: []*model.Application{
				{Id: "app-1", Plugins: []string{"plugin-a", "plugin-b"}},
				{Id: "app-2", Plugins: []string{"plugin-b"}},
			},
			plugin: "plugin-a",
			expected: []*model.Application{
				{Id: "app-1", Plugins: []string{"plugin-a", "plugin-b"}},
			},
		},
		{
			name: "matched some applications",
			storedApps: []*model.Application{
				{Id: "app-1", Plugins: []string{"plugin-a", "plugin-b"}},
				{Id: "app-2", Plugins: []string{"plugin-a"}},
				{Id: "app-3", Plugins: []string{"plugin-b"}},
			},
			plugin: "plugin-a",
			expected: []*model.Application{
				{Id: "app-1", Plugins: []string{"plugin-a", "plugin-b"}},
				{Id: "app-2", Plugins: []string{"plugin-a"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &store{
				applicationList: atomic.Value{},
				logger:          logger,
			}
			s.applicationList.Store(tt.storedApps)

			got := s.ListByPluginName(tt.plugin)
			assert.Equal(t, tt.expected, got)
		})
	}
}
