// Copyright 2021 The PipeCD Authors.
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
	"io"

	"github.com/pipe-cd/pipe/pkg/cache"
)

type PipedStatsBuilder struct {
}

func NewPipedStatsBuilder(_ cache.Cache) *PipedStatsBuilder {
	return &PipedStatsBuilder{}
}

// TODO: Implement builder to gather piped stats metrics in redis.
func (b *PipedStatsBuilder) Build() (io.ReadCloser, error) {
	return nil, nil
}
