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

package pipedsdk

import "github.com/pipe-cd/pipecd/pkg/plugin/pipedapi"

// Client is a toolkit for interacting with the piped service.
// It provides methods to call the piped service APIs.
// It's a wrapper around the raw piped service client.
type Client struct {
	base *pipedapi.PipedServiceClient

	// applicationID is used to identify the application that the client is working with.
	applicationID string
	// deploymentID is used to identify the deployment that the client is working with.
	// This field exists only when the client is working with a specific deployment; for example, when this client is passed as the deployoment plugin's argument.
	deploymentID string
	// stageID is used to identify the stage that the client is working with.
	// This field exists only when the client is working with a specific stage; for example, when this client is passed as the ExecuteStage method's argument.
	stageID string
}
