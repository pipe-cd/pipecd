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
	"errors"
	"time"

	provider "github.com/pipe-cd/pipe/pkg/app/piped/cloudprovider/lambda"
	"github.com/pipe-cd/pipe/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/backoff"
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

func sync(ctx context.Context, in *executor.Input, cloudProviderName string, cloudProviderCfg *config.CloudProviderLambdaConfig, fm provider.FunctionManifest) bool {
	in.LogPersister.Infof("Start applying the lambda function manifest")
	client, err := provider.DefaultRegistry().Client(cloudProviderName, cloudProviderCfg, in.Logger)
	if err != nil {
		in.LogPersister.Errorf("Unable to create Lambda client for the provider %s: %v", cloudProviderName, err)
		return false
	}

	found, err := client.IsFunctionExist(ctx, fm.Spec.Name)
	if err != nil {
		in.LogPersister.Errorf("Unable to validate function name %s: %v", fm.Spec.Name, err)
		return false
	}
	if found {
		if err := client.UpdateFunction(ctx, fm); err != nil {
			in.LogPersister.Errorf("Failed to update lambda function %s: %v", fm.Spec.Name, err)
			return false
		}
	} else {
		if err := client.CreateFunction(ctx, fm); err != nil {
			in.LogPersister.Errorf("Failed to create lambda function %s: %v", fm.Spec.Name, err)
			return false
		}
	}

	// TODO: Using backoff instead of time sleep waiting for a specific duration of time.
	// Wait before ready to commit change.
	in.LogPersister.Info("Waiting to update lambda function in progress...")
	// time.Sleep(3 * time.Minute)

	retry := backoff.NewRetry(3, backoff.NewConstant(time.Duration(1)*time.Minute))
	updateFunctionSucceed := false
	startWaitingStamp := time.Now()
	var version string
	for retry.WaitNext(ctx) {
		if !updateFunctionSucceed {
			// Commit version for applied Lambda function.
			// Note: via the current docs of [Lambda.PublishVersion](https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#Lambda.PublishVersion)
			// AWS Lambda doesn't publish a version if the function's configuration and code haven't changed since the last version.
			// But currently, unchanged revision is able to make publish (versionId++) as usual.
			version, err = client.PublishFunction(ctx, fm)
			if err != nil {
				in.LogPersister.Errorf("Failed to commit new version for Lambda function %s: %v", fm.Spec.Name, err)
				return false
			}
			in.LogPersister.Errorf("Commit new version for Lambda function %s after duration %v", fm.Spec.Name, time.Since(startWaitingStamp))
			updateFunctionSucceed = true

		}

		if updateFunctionSucceed {
			break
		}
	}

	_, err = client.GetTrafficConfig(ctx, fm)
	// Create Alias on not yet existed.
	if errors.Is(err, provider.ErrNotFound) {
		if err := client.CreateTrafficConfig(ctx, fm, version); err != nil {
			in.LogPersister.Errorf("Failed to create traffic routing for Lambda function %s (version: %s): %v", fm.Spec.Name, version, err)
			return false
		}
		in.LogPersister.Infof("Successfully applied the lambda function manifest")
		return true
	}
	if err != nil {
		in.LogPersister.Errorf("Failed to prepare traffic routing for Lambda function %s: %v", fm.Spec.Name, err)
		return false
	}

	// Update 100% traffic to the new lambda version.
	routingCfg := []provider.VersionTraffic{
		{
			Version: version,
			Percent: 100,
		},
	}
	if err = client.UpdateTrafficConfig(ctx, fm, routingCfg); err != nil {
		in.LogPersister.Errorf("Failed to update traffic routing for Lambda function %s (version: %s): %v", fm.Spec.Name, version, err)
		return false
	}

	in.LogPersister.Infof("Successfully applied the lambda function manifest")
	return true
}
