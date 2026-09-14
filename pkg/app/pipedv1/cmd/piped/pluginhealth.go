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

package piped

import "sync"

// pluginHealth keeps track of plugins whose process has exited on its own,
// outside of piped's own shutdown. The admin /healthz handler reads this
// so that a plugin crashing can actually be observed from the outside,
// instead of piped silently staying "ok" forever.
type pluginHealth struct {
	mu        sync.Mutex
	unhealthy map[string]error
}

func newPluginHealth() *pluginHealth {
	return &pluginHealth{
		unhealthy: make(map[string]error),
	}
}

// MarkUnhealthy records that the named plugin's process has exited
// unexpectedly.
func (h *pluginHealth) MarkUnhealthy(name string, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.unhealthy[name] = err
}

// Healthy reports whether every known plugin is still considered alive.
func (h *pluginHealth) Healthy() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.unhealthy) == 0
}

// UnhealthyNames returns the names of plugins currently marked unhealthy,
// in no particular order.
func (h *pluginHealth) UnhealthyNames() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	names := make([]string, 0, len(h.unhealthy))
	for name := range h.unhealthy {
		names = append(names, name)
	}
	return names
}
