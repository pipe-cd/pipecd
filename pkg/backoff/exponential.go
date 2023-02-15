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
	"math"
	"math/rand"
	"time"
)

type exponential struct {
	base  time.Duration
	max   time.Duration
	calls int
	rand  *rand.Rand
}

func NewExponential(base, max time.Duration) Backoff {
	return &exponential{
		base: base,
		max:  max,
		rand: rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
	}
}

// Implemented FullJitter algorithm
// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
func (b *exponential) Next() time.Duration {
	defer func() {
		b.calls++
	}()
	if b.calls == 0 {
		return 0
	}
	d := math.Min(float64(b.max), float64(b.base)*math.Pow(2, float64(b.calls-1)))
	d *= b.rand.Float64()
	return time.Duration(d)
}

func (b *exponential) Calls() int {
	return b.calls
}

func (b *exponential) Reset() {
	b.calls = 0
}

func (b *exponential) Clone() Backoff {
	return &exponential{
		base: b.base,
		max:  b.max,
		rand: b.rand,
	}
}
