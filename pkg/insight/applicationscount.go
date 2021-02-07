// Copyright 2020 The PipeCD Authors.
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

type ApplicationsCount struct {
	AccumulatedTo int64 `json:"accumulated_to"`
	// Use model.ApplicationKind as the key
	// and key "all" means data of all kind of applications.
	Counts map[string]Count `json:"counts"`
}

type Count struct {
	Deploying int64 `json:"deploying"`
	Disabled  int64 `json:"disabled"`
	Deleted   int64 `json:"deleted"`
}
