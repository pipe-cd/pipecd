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

package lambda

import (
	"context"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/lambda"
	"github.com/pipe-cd/pipe/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.ApplicationKind, f executor.Factory) error
}

func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &deployExecutor{
			Input: in,
		}
	}
	r.Register(model.StageLambdaSync, f)
	r.Register(model.StageLambdaPromote, f)
}

func findCloudProvider(in *executor.Input) (name string, cfg *config.CloudProviderLambdaConfig, found bool) {
	name = in.Application.CloudProvider
	if name == "" {
		in.LogPersister.Errorf("Missing the CloudProvider name in the application configuration")
		return
	}

	cp, ok := in.PipedConfig.FindCloudProvider(name, model.CloudProviderLambda)
	if !ok {
		in.LogPersister.Errorf("The specified cloud provider %q was not found in piped configuration", name)
		return
	}

	cfg = cp.LambdaConfig
	found = true
	return
}

func loadFunctionManifest(in *executor.Input, functionManifestFile string, ds *deploysource.DeploySource) (provider.FunctionManifest, bool) {
	in.LogPersister.Infof("Loading service manifest at the %s commit (%s)", ds.RevisionName, ds.RevisionName)

	fm, err := provider.LoadFunctionManifest(ds.AppDir, functionManifestFile)
	if err != nil {
		in.LogPersister.Errorf("Failed to load lambda function manifest (%v)", err)
		return provider.FunctionManifest{}, false
	}

	in.LogPersister.Infof("Successfully loaded the lambda function manifest at the %s commit", ds.RevisionName)
	return fm, true
}

func decideRevisionName(in *executor.Input, fm provider.FunctionManifest, commit string) (revision string, ok bool) {
	var err error
	revision, err = provider.DecideRevisionName(fm, commit)
	if err != nil {
		in.LogPersister.Errorf("Unable to decide revision name for the commit %s (%v)", commit, err)
		return
	}

	ok = true
	return
}

func apply(ctx context.Context, in *executor.Input, cloudProviderName string, cloudProviderCfg *config.CloudProviderLambdaConfig, fm provider.FunctionManifest) bool {
	in.LogPersister.Infof("Start applying the lambda function manifest")
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create Lambda client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	if err := client.Apply(ctx, fm, cloudProviderCfg.Role); err != nil {
		in.LogPersister.Errorf("Failed to apply the lambda function manifest (%v)", err)
		return false
	}

	in.LogPersister.Infof("Successfully applied the lambda function manifest")
	return true
}
