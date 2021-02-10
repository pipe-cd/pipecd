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

package insight

type ApplicationCount struct {
	Counts          []ApplicationCountByLabelSet `json:"counts"`
	AccumulatedFrom int64                        `json:"accumulated_from"`
	AccumulatedTo   int64                        `json:"accumulated_to"`
}

type ApplicationCountByLabelSet struct {
	LabelSet ApplicationCountLabelSet `json:"label_set"`
	Count    int                      `json:"count"`
}

type ApplicationCountLabelSet struct {
	// KUBERNETES, TERRAFORM, CLOUDRUN...
	Kind string `json:"kind"`
	// enable, disable or deleted
	Status string `json:"status"`
}
