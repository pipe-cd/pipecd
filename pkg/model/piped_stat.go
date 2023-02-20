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

package model

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	PipedStatsRetention = 2 * time.Minute
)

func UnmarshalPipedStat(data interface{}, ps *PipedStat) error {
	value, okValue := data.([]byte)
	if !okValue {
		return errors.New("error value not a bulk of string value")
	}
	return json.Unmarshal(value, ps)
}

func (ps *PipedStat) IsStaled(d time.Duration) bool {
	return time.Since(time.Unix(ps.Timestamp, 0)) > d
}
