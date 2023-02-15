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

package lambda

import (
	"encoding/json"
)

// TrafficConfigKeyName represents key for lambda service config map.
type TrafficConfigKeyName string

const (
	// TrafficPrimaryVersionKeyName represents the key points to primary version config on traffic routing map.
	TrafficPrimaryVersionKeyName TrafficConfigKeyName = "primary"
	// TrafficSecondaryVersionKeyName represents the key points to primary version config on traffic routing map.
	TrafficSecondaryVersionKeyName TrafficConfigKeyName = "secondary"
)

// RoutingTrafficConfig presents a map of primary and secondary version traffic for lambda function alias.
type RoutingTrafficConfig map[TrafficConfigKeyName]VersionTraffic

func (c *RoutingTrafficConfig) Encode() (string, error) {
	out, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (c *RoutingTrafficConfig) Decode(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}

// VersionTraffic presents the version, and the percent of traffic that's routed to it.
type VersionTraffic struct {
	Version string  `json:"version"`
	Percent float64 `json:"percent"`
}
