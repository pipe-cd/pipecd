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
	"github.com/pipe-cd/pipecd/pkg/config"
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

		// TODO: Tag application-id to Lambda funcs on create/update
		// Use the application ID tag as the key.
		appId := ""
		// for _, tag := range service.Tags {
		// 	if *tag.Key == provider.LabelApplication {
		// 		appId = *tag.Value
		// 		break
		// 	}
		// }
		if appId == "" {
			// Skip a service which is not managed by PipeCD.
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
	architectures := make([]provider.Architecture, 0, len(f.Configuration.Architectures))
	for _, arch := range f.Configuration.Architectures {
		architectures = append(architectures, provider.Architecture{Name: string(arch)})
	}

	layerArns := make([]string, 0, len(f.Configuration.Layers))
	for _, layer := range f.Configuration.Layers {
		layerArns = append(layerArns, *layer.Arn)
	}

	return provider.FunctionManifest{
		Kind:       "LambdaFunction",
		APIVersion: config.VersionV1Beta1,
		Spec: provider.FunctionManifestSpec{
			Name: *f.Configuration.FunctionName,
			Role: *f.Configuration.Role,

			ImageURI: *f.Code.ImageUri,

			// TODO:
			// S3Bucket         : *f.Code.,
			// S3Key            : *f.,
			// S3ObjectVersion  : *f.,
			// SourceCode       : *f. , Maybe not exist in Lambda's param

			Handler:       *f.Configuration.Handler,
			Architectures: architectures,
			EphemeralStorage: &provider.EphemeralStorage{
				Size: *f.Configuration.EphemeralStorage.Size,
			},
			Runtime:      string(f.Configuration.Runtime),
			Memory:       *f.Configuration.MemorySize,
			Timeout:      *f.Configuration.Timeout,
			Tags:         f.Tags,
			Environments: f.Configuration.Environment.Variables,
			VPCConfig: &provider.VPCConfig{
				SecurityGroupIDs: f.Configuration.VpcConfig.SecurityGroupIds,
				SubnetIDs:        f.Configuration.VpcConfig.SubnetIds,
			},
			Layers: layerArns,
		},
	}
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
