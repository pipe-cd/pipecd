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
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginHealth(t *testing.T) {
	h := newPluginHealth()

	assert.True(t, h.Healthy())
	assert.Empty(t, h.UnhealthyNames())

	h.MarkUnhealthy("kubernetes", errors.New("exit status 1"))

	assert.False(t, h.Healthy())
	assert.Equal(t, []string{"kubernetes"}, h.UnhealthyNames())

	// Marking the same plugin again must not create a duplicate entry.
	h.MarkUnhealthy("kubernetes", errors.New("exit status 2"))
	assert.Equal(t, []string{"kubernetes"}, h.UnhealthyNames())

	h.MarkUnhealthy("terraform", errors.New("killed"))
	assert.False(t, h.Healthy())
	assert.ElementsMatch(t, []string{"kubernetes", "terraform"}, h.UnhealthyNames())
}

func TestPluginHealthConcurrentAccess(t *testing.T) {
	h := newPluginHealth()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.MarkUnhealthy("plugin", errors.New("boom"))
			h.Healthy()
			h.UnhealthyNames()
		}()
	}
	wg.Wait()

	assert.False(t, h.Healthy())
	assert.Equal(t, []string{"plugin"}, h.UnhealthyNames())
}
