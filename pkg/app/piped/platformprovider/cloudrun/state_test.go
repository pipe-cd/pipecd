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

package cloudrun

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestMakeResourceStates(t *testing.T) {
	t.Parallel()

	sm, err := ParseServiceManifest([]byte(serviceManifest))
	require.NoError(t, err)

	svc, err := sm.RunService()
	require.NoError(t, err)

	s := (*Service)(svc)

	rm, err := ParseRevisionManifest([]byte(revisionManifest))
	require.NoError(t, err)

	rev, err := rm.RunRevision()
	require.NoError(t, err)

	r := (*Revision)(rev)

	// MakeResourceStates
	rs := []*Revision{r}
	states := MakeResourceStates(s, rs, time.Now())
	require.Len(t, states, 2)
	assert.Equal(t, model.CloudRunResourceState_OTHER, states[0].HealthStatus)
	assert.Equal(t, model.CloudRunResourceState_HEALTHY, states[1].HealthStatus)
}
