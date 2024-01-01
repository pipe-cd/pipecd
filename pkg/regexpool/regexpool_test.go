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

package regexpool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPool(t *testing.T) {
	pool := DefaultPool()
	require.NotNil(t, pool)
}

func TestPool(t *testing.T) {
	pool, err := NewPool(2)
	require.NoError(t, err)

	regex, err := pool.Get("(gopher){2}")
	assert.NoError(t, err)
	assert.NotNil(t, regex)

	regex, err = pool.Get("(abc")
	assert.Equal(t, fmt.Errorf("unable to compile: (abc"), err)
	assert.Nil(t, regex)
}
