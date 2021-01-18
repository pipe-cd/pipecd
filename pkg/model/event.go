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

package model

import (
	"fmt"
	"sort"
	"strings"
)

// MakeEventKey builds a unique identifier based on the given name and labels.
// It returns the exact same string as long as both are the same.
func MakeEventKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
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
	return b.String()
}
