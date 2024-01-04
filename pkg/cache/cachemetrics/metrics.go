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

package cachemetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	statusKey = "status"
	sourceKey = "source"
)

type StatusLabel string

const (
	LabelStatusHit  StatusLabel = "hit"
	LabelStatusMiss StatusLabel = "miss"
)

type SourceLabel string

const (
	LabelSourceRedis    SourceLabel = "redis"
	LabelSourceInmemory SourceLabel = "inmemory"
)

var (
	getCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_get_operation_total",
			Help: "Number of cache get operation while processing",
		},
		[]string{
			statusKey,
			sourceKey,
		},
	)
)

func Register(r prometheus.Registerer) {
	r.MustRegister(getCounter)
}

func IncGetOperationCounter(source SourceLabel, status StatusLabel) {
	getCounter.With(prometheus.Labels{
		statusKey: string(status),
		sourceKey: string(source),
	}).Inc()
}
