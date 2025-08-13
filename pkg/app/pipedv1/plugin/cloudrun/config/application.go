package config

import "github.com/pipe-cd/piped-plugin-sdk-go/unit"

type CloudRunApplicationSpec struct {
	// Input for CloudRun deployment such as docker image...
	Input CloudRunDeploymentInput `json:"input"`
	// Configuration for quick sync.
	QuickSync CloudRunSyncStageOptions `json:"quickSync"`
}

type CloudRunDeploymentInput struct {
	// The name of service manifest file placing in application directory.
	// Default is service.yaml
	ServiceManifestFile string `json:"serviceManifestFile"`
}

// CloudRunSyncStageOptions contains all configurable values for a CLOUDRUN_SYNC stage.
type CloudRunSyncStageOptions struct {
}

// CloudRunPromoteStageOptions contains all configurable values for a CLOUDRUN_PROMOTE stage.
type CloudRunPromoteStageOptions struct {
	// Percentage of traffic should be routed to the new version.
	Percent unit.Percentage `json:"percent"`
}

type CloudRunDeployTargetConfig struct {
	// The GCP project hosting the CloudRun service.
	Project string `json:"project"`
	// The region of running CloudRun service.
	Region string `json:"region"`
	// The path to the service account file for accessing CloudRun service.
	CredentialsFile string `json:"credentialsFile,omitempty"`
}
