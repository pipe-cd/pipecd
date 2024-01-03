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

package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

var hashFunc = sha256.New

func (e *Event) IsHandled() bool {
	return e.Status == EventStatus_EVENT_SUCCESS || e.Status == EventStatus_EVENT_FAILURE
}

// ContainLabels checks if it has all the given labels.
func (e *Event) ContainLabels(labels map[string]string) bool {
	if len(e.Labels) < len(labels) {
		return false
	}

	for k, v := range labels {
		value, ok := e.Labels[k]
		if !ok {
			return false
		}
		if value != v {
			return false
		}
	}
	return true
}

func (e *Event) SetUpdatedAt(t int64) {
	e.UpdatedAt = t
}

// MakeEventKey builds a fixed-length identifier based on the given name
// and labels. It returns the exact same string as long as both are the same.
func MakeEventKey(name string, labels map[string]string) string {
	h := hashFunc()
	if len(labels) == 0 {
		h.Write([]byte(name))
		return hex.EncodeToString(h.Sum(nil))
	}

	var b strings.Builder
	b.WriteString(name)

	// Guarantee uniqueness by sorting by keys.
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		b.WriteString(fmt.Sprintf("/%s:%s", key, labels[key]))
	}

	// Make a hash to be fixed length regardless of the length of the event key.
	h.Write([]byte(b.String()))
	return hex.EncodeToString(h.Sum(nil))
}
