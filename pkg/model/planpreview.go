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

import "time"

func (r *PlanPreviewCommandResult) FillURLs(baseURL string) {
	r.PipedUrl = MakePipedURL(baseURL, r.PipedId)
	for _, ar := range r.Results {
		ar.ApplicationUrl = MakeApplicationURL(baseURL, ar.ApplicationId)
	}
}

func MakeApplicationPlanPreviewResult(app Application) *ApplicationPlanPreviewResult {
	r := &ApplicationPlanPreviewResult{
		ApplicationId:        app.Id,
		ApplicationName:      app.Name,
		ApplicationKind:      app.Kind,
		ApplicationDirectory: app.GitPath.Path,
		Labels:               app.Labels,
		PipedId:              app.PipedId,
		ProjectId:            app.ProjectId,
		CreatedAt:            time.Now().Unix(),
	}
	return r
}
