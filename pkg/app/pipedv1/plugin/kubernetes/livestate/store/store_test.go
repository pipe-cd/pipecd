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

package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetApplicationResources(t *testing.T) {
	dtr := newDeployTargetResources("target-1")

	app1 := dtr.getApplicationResources("app-1")
	require.NotNil(t, app1)
	assert.Equal(t, "target-1", app1.deployTarget)

	// Fetching the same appID returns the same instance
	app1Again := dtr.getApplicationResources("app-1")
	assert.Same(t, app1, app1Again)

	// Fetching a different appID returns a distinct instance
	app2 := dtr.getApplicationResources("app-2")
	require.NotNil(t, app2)
	assert.NotSame(t, app1, app2)
}

func TestGetApplicationResources_Concurrent(t *testing.T) {
	dtr := newDeployTargetResources("target-concurrent")

	var wg sync.WaitGroup
	workers := 100

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			appID := fmt.Sprintf("app-%d", workerID%10)
			app := dtr.getApplicationResources(appID)
			assert.NotNil(t, app)
		}(i)
	}

	wg.Wait()

	// Verify all 10 application entries exist in the map
	dtr.mu.RLock()
	defer dtr.mu.RUnlock()
	assert.Len(t, dtr.applications, 10)
}
