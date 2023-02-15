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
	"time"
)

type constant struct {
	calls    int
	interval time.Duration
}

func NewConstant(interval time.Duration) Backoff {
	return &constant{
		interval: interval,
	}
}

func (b *constant) Next() time.Duration {
	defer func() {
		b.calls++
	}()
	if b.calls == 0 {
		return 0
	}
	return b.interval
}

func (b *constant) Calls() int {
	return b.calls
}

func (b *constant) Reset() {
	b.calls = 0
}

func (b *constant) Clone() Backoff {
	return &constant{
		interval: b.interval,
	}
}
