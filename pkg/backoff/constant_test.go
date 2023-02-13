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

package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstant(t *testing.T) {
	bo := NewConstant(time.Millisecond)
	assert.Equal(t, 0, bo.Calls())
	assert.Equal(t, time.Duration(0), bo.Next())
	assert.Equal(t, 1, bo.Calls())

	for i := 2; i < 10; i++ {
		d := bo.Next()
		assert.Equal(t, time.Millisecond, d)
		assert.Equal(t, i, bo.Calls())
	}
}
