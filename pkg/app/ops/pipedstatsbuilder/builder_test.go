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

package pipedstatsbuilder

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/cache"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type mockBuilderBackend struct {
	cache.Cache
	srcs []string
}

func newMockBuilderBackend() *mockBuilderBackend {
	return &mockBuilderBackend{
		srcs: []string{
			"./testdata/piped_stat_1",
			"./testdata/piped_stat_2",
		},
	}
}

func (m *mockBuilderBackend) GetAll() (map[string]interface{}, error) {
	out := make(map[string]interface{}, len(m.srcs))
	for _, file := range m.srcs {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		val, _ := json.Marshal(model.PipedStat{Metrics: data, Timestamp: time.Now().Unix()})
		out[file] = val
	}
	return out, nil
}

func TestBuildPipedStat(t *testing.T) {
	builder := NewPipedStatsBuilder(newMockBuilderBackend(), zap.NewNop())
	rc, err := builder.Build()
	require.NoError(t, err)
	require.NotNil(t, rc)

	buf := new(strings.Builder)
	io.Copy(buf, rc)
	actOutElements := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n\n")

	data, _ := os.ReadFile("./testdata/expected")
	expOutElements := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n\n")

	require.Equal(t, len(expOutElements), len(actOutElements))
	assert.ElementsMatch(t, expOutElements, actOutElements)
}
