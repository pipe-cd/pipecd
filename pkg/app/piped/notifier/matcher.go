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

package notifier

import (
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type matcher struct {
	events       map[string]struct{}
	ignoreEvents map[string]struct{}
	groups       map[string]struct{}
	ignoreGroups map[string]struct{}
	apps         map[string]struct{}
	ignoreApps   map[string]struct{}
	labels       map[string]string
	ignoreLabels map[string]string
}

func newMatcher(cfg config.NotificationRoute) *matcher {
	return &matcher{
		events:       makeStringMap(cfg.Events, "EVENT"),
		ignoreEvents: makeStringMap(cfg.IgnoreEvents, "EVENT"),
		groups:       makeStringMap(cfg.Groups, "EVENT"),
		ignoreGroups: makeStringMap(cfg.IgnoreGroups, "EVENT"),
		apps:         makeStringMap(cfg.Apps, ""),
		ignoreApps:   makeStringMap(cfg.IgnoreApps, ""),
		labels:       cfg.Labels,
		ignoreLabels: cfg.IgnoreLabels,
	}
}

type appNameMetadata interface {
	GetAppName() string
}

type labelsMetadata interface {
	GetLabels() map[string]string
}

func (m *matcher) Match(event model.NotificationEvent) bool {
	if _, ok := m.ignoreEvents[event.Type.String()]; ok {
		return false
	}
	if _, ok := m.ignoreGroups[event.Group().String()]; ok {
		return false
	}

	var appName string
	if md, ok := event.Metadata.(appNameMetadata); ok {
		appName = md.GetAppName()
	}
	if _, ok := m.ignoreApps[appName]; ok && appName != "" {
		return false
	}

	var (
		labels        = make(map[string]string)
		supportLabels = false
	)
	if md, ok := event.Metadata.(labelsMetadata); ok {
		labels = md.GetLabels()
		supportLabels = true
	}
	for k, v := range labels {
		if m.ignoreLabels[k] == v {
			return false
		}
	}

	if len(m.events) > 0 {
		if _, ok := m.events[event.Type.String()]; !ok {
			return false
		}
	}
	if len(m.groups) > 0 {
		if _, ok := m.groups[event.Group().String()]; !ok {
			return false
		}
	}
	if len(m.apps) > 0 && appName != "" {
		if _, ok := m.apps[appName]; !ok {
			return false
		}
	}
	if len(m.labels) > 0 && supportLabels {
		if len(labels) < len(m.labels) {
			return false
		}

		for k, v := range m.labels {
			if labels[k] != v {
				return false
			}
		}
	}

	return true
}

func makeStringMap(keys []string, prefix string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		if prefix != "" {
			k = prefix + "_" + k
		}
		m[k] = struct{}{}
	}
	return m
}
