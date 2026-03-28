// Copyright 2026 The PipeCD Authors.
package deployment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/pipe-cd/piped-plugin-sdk-go"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/cloudrunservice/config"
)

func TestParseContainerImage(t *testing.T) {
	tests := []struct {
		name         string
		image        string
		expectedName string
		expectedTag  string
	}{
		{
			name:         "typical image with tag",
			image:        "gcr.io/project/app:v1.0.0",
			expectedName: "app",
			expectedTag:  "v1.0.0",
		},
		{
			name:         "image with registry port and tag",
			image:        "localhost:5000/app:v1.0.0",
			expectedName: "app",
			expectedTag:  "v1.0.0",
		},
		{
			name:         "image without tag",
			image:        "ubuntu",
			expectedName: "ubuntu",
			expectedTag:  "latest",
		},
		{
			name:         "image with registry port without tag",
			image:        "registry.local:8080/foo/bar",
			expectedName: "bar",
			expectedTag:  "latest",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			name, tag := parseContainerImage(tc.image)
			assert.Equal(t, tc.expectedName, name)
			assert.Equal(t, tc.expectedTag, tag)
		})
	}
}

func TestDetermineStrategy(t *testing.T) {
	tests := []struct {
		name                    string
		runningDeploymentSource string
		expectedStrategy        sdk.SyncStrategy
	}{
		{
			name:                    "first time deployment",
			runningDeploymentSource: "",
			expectedStrategy:        sdk.SyncStrategyQuickSync,
		},
		{
			name:                    "subsequent deployment",
			runningDeploymentSource: "/path/to/app",
			expectedStrategy:        sdk.SyncStrategyPipelineSync,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := &sdk.DetermineStrategyInput[config.CloudRunApplicationSpec]{
				Request: sdk.DetermineStrategyRequest[config.CloudRunApplicationSpec]{
					RunningDeploymentSource: sdk.DeploymentSource[config.CloudRunApplicationSpec]{
						ApplicationDirectory: tc.runningDeploymentSource,
					},
				},
			}

			resp, err := determineStrategy(context.Background(), input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStrategy, resp.Strategy)
		})
	}
}
