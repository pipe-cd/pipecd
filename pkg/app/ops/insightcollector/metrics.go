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

import "bytes"

type CollectorMetrics uint

// Options controlling the InsightCollector.
const (
	ChangeFailureRate CollectorMetrics = 1 << iota
	DevelopmentFrequency
	ApplicationCount
)

func NewCollectorMetrics() CollectorMetrics {
	return CollectorMetrics(0)
}

// CollectorMetrics is represented as a sequence of zero or more of these letters:
// C: [C]hange failure rate.
// D: [D]evelopment frequency.
// A: [A]pplication count.
func (m CollectorMetrics) String() string {
	var buf bytes.Buffer
	if m.IsEnabled(ChangeFailureRate) {
		buf.WriteByte('C')
	}
	if m.IsEnabled(DevelopmentFrequency) {
		buf.WriteByte('D')
	}
	if m.IsEnabled(ApplicationCount) {
		buf.WriteByte('A')
	}
	return buf.String()
}

func (m *CollectorMetrics) Enable(a CollectorMetrics) {
	*m |= a
}

func (m CollectorMetrics) IsEnabled(a CollectorMetrics) bool {
	return m&a != 0
}
