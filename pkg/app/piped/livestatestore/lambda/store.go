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

package lambda

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/service/lambda"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/lambda"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type store struct {
	apps   atomic.Value
	logger *zap.Logger
	client provider.Client
}

type app struct {
	functionManifest provider.FunctionManifest

	// States of functions
	states  []*model.LambdaResourceState
	version model.ApplicationLiveStateVersion
}

func (s *store) run(ctx context.Context) error {
	apps := map[string]app{}
	now := time.Now()
	version := model.ApplicationLiveStateVersion{
		Timestamp: now.Unix(),
	}

	funcCfgs, err := s.client.ListFunctions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list Lambda functions: %w", err)
	}

	for _, funcCfg := range funcCfgs {
		f, err := s.client.GetFunction(ctx, *funcCfg.FunctionName)
		if err != nil {
			return fmt.Errorf("failed to get Lambda function %s: %w", *funcCfg.FunctionName, err)
		}

		appId, ok := f.Tags[provider.LabelApplication]
		if !ok {
			// Skip a function not managed by PipeCD.
			continue
		}

		apps[appId] = app{
			functionManifest: convertToManifest(f),
			states: []*model.LambdaResourceState{
				provider.MakeFunctionResourceState(f.Configuration),
			},
			version: version,
		}

	}

	s.apps.Store(apps)

	return nil
}

func convertToManifest(f *lambda.GetFunctionOutput) provider.FunctionManifest {
	fc := f.Configuration

	architectures := make([]provider.Architecture, 0, len(fc.Architectures))
	for _, arch := range fc.Architectures {
		architectures = append(architectures, provider.Architecture{Name: string(arch)})
	}

	layerArns := make([]string, 0, len(fc.Layers))
	for _, layer := range fc.Layers {
		layerArns = append(layerArns, *layer.Arn)
	}

	m := provider.FunctionManifest{
		Kind:       provider.FunctionManifestKind,
		APIVersion: provider.VersionV1Beta1,
		Spec: provider.FunctionManifestSpec{
			Name: *fc.FunctionName,
			// S3Bucket, S3Key, S3ObjectVersion, and SourceCode cannot be retrieved from Lambda's response.

			Architectures: architectures,
			Runtime:       string(fc.Runtime),
			Memory:        *fc.MemorySize,
			Timeout:       *fc.Timeout,
			Tags:          f.Tags,
			Layers:        layerArns,
		},
	}

	if fc.Role != nil {
		m.Spec.Role = *fc.Role
	}
	if f.Code.ImageUri != nil {
		m.Spec.ImageURI = *f.Code.ImageUri
	}
	if fc.Handler != nil {
		m.Spec.Handler = *fc.Handler
	}
	if fc.EphemeralStorage != nil && fc.EphemeralStorage.Size != nil {
		m.Spec.EphemeralStorage = &provider.EphemeralStorage{
			Size: *fc.EphemeralStorage.Size,
		}
	}
	if fc.Environment != nil {
		m.Spec.Environments = fc.Environment.Variables
	}
	if fc.VpcConfig != nil {
		m.Spec.VPCConfig = &provider.VPCConfig{
			SecurityGroupIDs: fc.VpcConfig.SecurityGroupIds,
			SubnetIDs:        fc.VpcConfig.SubnetIds,
		}
	}

	return m
}

func (s *store) loadApps() map[string]app {
	apps := s.apps.Load()
	if apps == nil {
		return nil
	}
	return apps.(map[string]app)
}

func (s *store) getFunctionManifest(appID string) (provider.FunctionManifest, bool) {
	apps := s.loadApps()
	if apps == nil {
		return provider.FunctionManifest{}, false
	}

	app, ok := apps[appID]
	if !ok {
		return provider.FunctionManifest{}, false
	}

	return app.functionManifest, true
}

func (s *store) getState(appID string) (State, bool) {
	apps := s.loadApps()
	if apps == nil {
		return State{}, false
	}

	app, ok := apps[appID]
	if !ok {
		return State{}, false
	}

	return State{
		Resources: app.states,
		Version:   app.version,
	}, true
}
