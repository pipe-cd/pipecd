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
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Replicas struct {
	Number       int
	IsPercentage bool
}

func (r Replicas) String() string {
	s := strconv.FormatInt(int64(r.Number), 10)
	if r.IsPercentage {
		return s + "%"
	}
	return s
}

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

func (r Replicas) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Replicas) UnmarshalJSON(b []byte) error {
	var v interface{}
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
