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

package config

import (
	"encoding/json"
	"testing"
)

// Each body previously panicked in UnmarshalJSON, which decodes into
// map[string]interface{} and asserts on the lookups.
func TestSharedSSOConfigMalformed(t *testing.T) {
	for name, body := range map[string]string{
		"provider missing":  `{"name":"github"}`,
		"provider a number": `{"name":"github","provider":42}`,
		"name a number":     `{"name":42,"provider":"GITHUB"}`,
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked on user config: %v", r)
				}
			}()
			var c SharedSSOConfig
			if err := json.Unmarshal([]byte(body), &c); err == nil {
				t.Fatal("want an error, got nil")
			}
		})
	}
}
