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

package kubernetesmetrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/tools/metrics"
)

const (
	hostKey         = "host"
	methodKey       = "method"
	codeKey         = "code"
	eventKey        = "event"
	eventHandledKey = "handled"
)

type EventKind string

const (
	LabelEventAdd    EventKind = "add"
	LabelEventUpdate EventKind = "update"
	LabelEventDelete EventKind = "delete"
)

type EventHandledVal string

const (
	LabelEventHandled       EventHandledVal = "true"
	LabelEventNotYetHandled EventHandledVal = "false"
)

var (
	apiRequestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "livestatestore_kubernetes_api_requests_total",
			Help: "Number of requests sent to kubernetes api server.",
		},
		[]string{
			hostKey,
			methodKey,
			codeKey,
		},
	)
	resourceEventsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "livestatestore_kubernetes_resource_events_total",
			Help: "Number of resource events received from kubernetes server.",
		},
		[]string{
			eventKey,
			eventHandledKey,
		},
	)
)

func Register(r prometheus.Registerer) {
	r.MustRegister(
		apiRequestsCounter,
		resourceEventsCounter,
	)

	opts := metrics.RegisterOpts{
		RequestResult: requestResultCollector{},
	}
	metrics.Register(opts)
}

type requestResultCollector struct {
}

func (c requestResultCollector) Increment(ctx context.Context, code string, method string, host string) {
	apiRequestsCounter.With(prometheus.Labels{
		hostKey:   host,
		methodKey: method,
		codeKey:   code,
	}).Inc()
}

func IncResourceEventsCounter(event EventKind, handled EventHandledVal) {
	resourceEventsCounter.With(prometheus.Labels{
		eventKey:        string(event),
		eventHandledKey: string(handled),
	}).Inc()
}
