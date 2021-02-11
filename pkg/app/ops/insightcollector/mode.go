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

package insightcollector

import (
	"bytes"
)

type CollectorMode uint

// Options controlling the InsightCollector.
// The value is a sequence of zero or more of these letters:
// C: enable data collection for [C]hange failure rate .
// D: enable data collection for [D]evelopment frequency.
// A: enable data collection for [A]pplication count.
const (
	EnableChangeFailureRate CollectorMode = 1 << iota
	EnableDevelopmentFrequency
	EnableApplicationCount
)

func NewCollectorMode() CollectorMode {
	return CollectorMode(0)
}

func (m CollectorMode) String() string {
	var buf bytes.Buffer
	if m&EnableChangeFailureRate != 0 {
		buf.WriteByte('C')
	}
	if m&EnableDevelopmentFrequency != 0 {
		buf.WriteByte('D')
	}
	if m&EnableApplicationCount != 0 {
		buf.WriteByte('A')
	}
	return buf.String()
}

// Set parses the flag characters in m and updates *m.
func (m *CollectorMode) Set(s string) {
	var mode CollectorMode
	for _, c := range s {
		switch c {
		case 'C':
			mode |= EnableChangeFailureRate
		case 'D':
			mode |= EnableDevelopmentFrequency
		case 'A':
			mode |= EnableApplicationCount
		}
	}
	*m = mode
}

func (m CollectorMode) EnableChangeFailureRate() bool {
	return m&EnableChangeFailureRate != 0
}
func (m CollectorMode) EnableDevelopmentFrequency() bool {
	return m&EnableDevelopmentFrequency != 0
}
func (m CollectorMode) EnableApplicationCount() bool {
	return m&EnableApplicationCount != 0
}
