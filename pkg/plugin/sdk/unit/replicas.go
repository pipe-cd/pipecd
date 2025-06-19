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
	"math"
	"strconv"
	"strings"
)

// Replicas represents a replica count that can be either an absolute number or a percentage.
// When IsPercentage is true, Number represents a percentage value.
type Replicas struct {
	Number       int
	IsPercentage bool
}

// String returns the string representation of the replicas.
// If IsPercentage is true, it includes the % suffix.
func (r Replicas) String() string {
	s := strconv.FormatInt(int64(r.Number), 10)
	if r.IsPercentage {
		return s + "%"
	}
	return s
}

// Calculate returns the actual replica count based on the total replicas.
// If the replica is a percentage, it calculates the ceiling of (Number * total / 100).
// If Number is 0, it returns the defaultValue.
// If IsPercentage is false, it returns the Number directly.
func (r Replicas) Calculate(total, defaultValue int) int {
	if r.Number == 0 {
		return defaultValue
	}
	if !r.IsPercentage {
		return r.Number
	}
	num := float64(r.Number*total) / 100.0
	return int(math.Ceil(num))
}

// MarshalJSON marshals the Replicas to a JSON string representation.
func (r Replicas) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON unmarshals a JSON value to Replicas.
// It accepts numeric values, string numbers (e.g., "5"), and percentage strings (e.g., "50%").
func (r *Replicas) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch raw := v.(type) {
	case float64:
		*r = Replicas{
			Number:       int(raw),
			IsPercentage: false,
		}
		return nil
	case string:
		replicas := Replicas{
			IsPercentage: false,
		}
		if strings.HasSuffix(raw, "%") {
			replicas.IsPercentage = true
			raw = strings.TrimSuffix(raw, "%")
		}
		value, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("invalid replicas: %v", err)
		}
		replicas.Number = value
		*r = replicas
		return nil
	default:
		return fmt.Errorf("invalid replicas: %v", string(b))
	}
}
