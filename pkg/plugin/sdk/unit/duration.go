// Copyright 2025 The PipeCD Authors.
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

package unit

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration represents a time duration that can be marshaled/unmarshaled as a string in JSON.
// It supports both numeric values (interpreted as nanoseconds) and string values (e.g., "5m", "1h30m").
type Duration time.Duration

// Duration returns the time.Duration value.
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// MarshalJSON marshals the Duration to a JSON string representation (e.g., "5m30s").
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON unmarshals a JSON value to Duration.
// It accepts both numeric values (interpreted as nanoseconds) and string values (e.g., "5m", "1h30m").
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch raw := v.(type) {
	case float64:
		*d = Duration(time.Duration(raw))
		return nil
	case string:
		value, err := time.ParseDuration(raw)
		if err != nil {
			return err
		}
		*d = Duration(value)
		return nil
	default:
		return fmt.Errorf("invalid duration: %v", string(b))
	}
}
