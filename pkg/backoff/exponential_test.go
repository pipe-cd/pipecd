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

package backoff

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExponential(t *testing.T) {
	eb := NewExponential(time.Millisecond, time.Second)
	assert.Equal(t, 0, eb.Calls())
	assert.Equal(t, time.Duration(0), eb.Next())
	assert.Equal(t, 1, eb.Calls())

	for i := 2; i < 100; i++ {
		d := eb.Next()
		des := fmt.Sprintf("i = %d duration: %v", i, d)
		require.True(t, d >= 0, des)
		require.True(t, d <= time.Second, des)
		require.True(t, eb.Calls() == i, des)
	}

	eb.Reset()
	assert.Equal(t, 0, eb.Calls())
}
