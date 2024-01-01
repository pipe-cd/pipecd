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

package planpreviewmetrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	statusKey = "status"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
)

var (
	commandReceivedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "plan_preview_command_received_total",
			Help: "Total number of plan-preview commands received at piped.",
		},
	)
	commandHandledTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "plan_preview_command_handled_total",
			Help: "Total number of plan-preview commands handled at piped.",
		},
		[]string{statusKey},
	)

	commandHandlingSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "plan_preview_command_handling_seconds",
			Help:    "Histogram of handling seconds of plan-preview commands.",
			Buckets: []float64{1, 10, 30, 60, 120, 300, 600},
		},
		[]string{statusKey},
	)
)

func ReceivedCommands(n int) {
	commandReceivedTotal.Add(float64(n))
}

func HandledCommand(s Status, d time.Duration) {
	commandHandledTotal.With(prometheus.Labels{
		statusKey: string(s),
	}).Inc()

	commandHandlingSeconds.With(prometheus.Labels{
		statusKey: string(s),
	}).Observe(d.Seconds())
}

func Register(r prometheus.Registerer) {
	r.MustRegister(
		commandReceivedTotal,
		commandHandledTotal,
		commandHandlingSeconds,
	)
}
