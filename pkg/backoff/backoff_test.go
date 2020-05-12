// Copyright 2020 The PipeCD Authors.
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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitNext(t *testing.T) {
	var (
		bo          = NewConstant(time.Millisecond)
		r           = NewRetry(10, bo)
		ctx, cancel = context.WithCancel(context.TODO())
	)
	ok := r.WaitNext(ctx)
	assert.Equal(t, true, ok)

	cancel()
	ok = r.WaitNext(ctx)
	assert.Equal(t, false, ok)
}
