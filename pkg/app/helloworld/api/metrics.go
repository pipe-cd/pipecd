// Copyright 2020 The Dianomi Authors.
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

package api

import (
	"sync"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
)

var (
	mCalls = stats.Int64("hello_calls", "The distribution of hello calls", "1")

	keyGender, _ = tag.NewKey("gender")

	callCountView = &view.View{
		Name:        "calls",
		Measure:     mCalls,
		Description: "The number of hello calls",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyGender},
	}

	registerOnce sync.Once
)

func regsiterMetrics(logger *zap.Logger) {
	registerOnce.Do(func() {
		err := view.Register(callCountView)
		if err != nil {
			logger.Fatal("failed to register metrics", zap.Error(err))
		}
	})
}
