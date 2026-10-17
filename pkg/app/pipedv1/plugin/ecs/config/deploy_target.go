// Copyright 2026 The PipeCD Authors.
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

// ECSDeployTargetConfig represents the deployment target configuration for ECS plugin.
type ECSDeployTargetConfig struct {
	// Region is the AWS region where the ECS cluster is located
	// (e.g., "us-west-2").
	Region string `json:"region"`

	// Profile is the AWS profile to use from the credentials file
	// If empty, uses the default profile or "default" if AWS_PROFILE env var is not set
	Profile string `json:"profile,omitempty"`

	// CredentialsFile is the path to the AWS shared credentials file
	// (e.g., "~/.aws/credentials")
	// If empty, uses the default location
	CredentialsFile string `json:"credentialsFile,omitempty"`

	// RoleARN is the IAM role ARN to assume when accessing AWS resources
	// (e.g., "arn:aws:iam::123456789:role/ecs-deployment-role").
	// Required when assuming a role across accounts
	RoleARN string `json:"roleARN,omitempty"`

	// TokenFile is the path to the OIDC token file for web identity federation.
	// Required when RoleARN is set for OIDC-based authentication
	TokenFile string `json:"tokenFile,omitempty"`
}
