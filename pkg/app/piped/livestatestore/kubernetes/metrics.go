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

package kubernetes

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/tools/metrics"
)

const (
	metricsLabelHost         = "host"
	metricsLabelMethod       = "method"
	metricsLabelCode         = "code"
	metricsLabelEvent        = "event"
	metricsLabelEventHandled = "handled"
)

var (
	metricsAPIRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "livestatestore_kubernetes_api_requests_total",
			Help: "Number of requests sent to kubernetes api server.",
		},
		[]string{
			metricsLabelHost,
			metricsLabelMethod,
			metricsLabelCode,
		},
	)
	metricsResourceEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "livestatestore_kubernetes_resource_events_total",
			Help: "Number of resource events received from kubernetes server.",
		},
		[]string{
			metricsLabelEvent,
			metricsLabelEventHandled,
		},
	)
)

func registerMetrics() {
	prometheus.MustRegister(
		metricsAPIRequests,
		metricsResourceEvents,
	)

	opts := metrics.RegisterOpts{
		RequestResult: requestResultCollector{},
	}
	metrics.Register(opts)
}

type requestResultCollector struct {
}

func (c requestResultCollector) Increment(code string, method string, host string) {
	metricsAPIRequests.With(prometheus.Labels{
		metricsLabelHost:   host,
		metricsLabelMethod: method,
		metricsLabelCode:   code,
	}).Inc()
}

func incrementResourceEventCounter(event string, handled bool) {
	metricsResourceEvents.With(prometheus.Labels{
		metricsLabelEvent:        event,
		metricsLabelEventHandled: strconv.FormatBool(handled),
	}).Inc()
}
