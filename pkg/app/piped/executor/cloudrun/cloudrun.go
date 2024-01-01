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

package cloudrun

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/cloudrun"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type registerer interface {
	Register(stage model.Stage, f executor.Factory) error
	RegisterRollback(kind model.RollbackKind, f executor.Factory) error
}

func Register(r registerer) {
	f := func(in executor.Input) executor.Executor {
		return &deployExecutor{
			Input: in,
		}
	}
	r.Register(model.StageCloudRunSync, f)
	r.Register(model.StageCloudRunPromote, f)

	r.RegisterRollback(model.RollbackKind_Rollback_CLOUDRUN, func(in executor.Input) executor.Executor {
		return &rollbackExecutor{
			Input: in,
		}
	})
}

func loadServiceManifest(in *executor.Input, serviceManifestFile string, ds *deploysource.DeploySource) (provider.ServiceManifest, bool) {
	in.LogPersister.Infof("Loading service manifest at commit %s", ds.Revision)

	sm, err := provider.LoadServiceManifest(ds.AppDir, serviceManifestFile)
	if err != nil {
		in.LogPersister.Errorf("Failed to load service manifest (%v)", err)
		return provider.ServiceManifest{}, false
	}

	in.LogPersister.Infof("Successfully loaded the service manifest at commit %s", ds.Revision)
	return sm, true
}

func findPlatformProvider(in *executor.Input) (name string, cfg *config.PlatformProviderCloudRunConfig, found bool) {
	name = in.Application.PlatformProvider
	if name == "" {
		in.LogPersister.Error("Missing the PlatformProvider name in the application configuration")
		return
	}

	cp, ok := in.PipedConfig.FindPlatformProvider(name, model.ApplicationKind_CLOUDRUN)
	if !ok {
		in.LogPersister.Errorf("The specified platform provider %q was not found in piped configuration", name)
		return
	}

	cfg = cp.CloudRunConfig
	found = true
	return
}

func decideRevisionName(sm provider.ServiceManifest, commit string, lp executor.LogPersister) (revision string, ok bool) {
	var err error
	revision, err = provider.DecideRevisionName(sm, commit)
	if err != nil {
		lp.Errorf("Unable to decide revision name for the commit %s (%v)", commit, err)
		return
	}

	ok = true
	return
}

func configureServiceManifest(sm provider.ServiceManifest, revision string, traffics []provider.RevisionTraffic, lp executor.LogPersister) bool {
	if revision != "" {
		if err := sm.SetRevision(revision); err != nil {
			lp.Errorf("Unable to set revision name to service manifest (%v)", err)
			return false
		}
	}

	if err := sm.UpdateTraffic(traffics); err != nil {
		lp.Errorf("Unable to configure traffic percentages to service manifest (%v)", err)
		return false
	}

	lp.Info("Successfully prepared service manifest with traffic percentages as below:")
	for _, t := range traffics {
		lp.Infof("  %s: %d", t.RevisionName, t.Percent)
	}

	return true
}

func apply(ctx context.Context, client provider.Client, sm provider.ServiceManifest, lp executor.LogPersister) bool {
	lp.Info("Start applying the service manifest")

	_, err := client.Update(ctx, sm)
	if err == nil {
		lp.Infof("Successfully updated the service %s", sm.Name)
		return true
	}

	if err != provider.ErrServiceNotFound {
		lp.Errorf("Failed to update the service %s (%v)", sm.Name, err)
		return false
	}

	lp.Infof("Service %s was not found, a new service will be created", sm.Name)

	if _, err := client.Create(ctx, sm); err != nil {
		lp.Errorf("Failed to create the service %s (%v)", sm.Name, err)
		return false
	}

	lp.Infof("Successfully created the service %s", sm.Name)
	return true
}

func waitRevisionReady(ctx context.Context, client provider.Client, revisionName string, retryDuration, retryTimeout time.Duration, lp executor.LogPersister) error {
	shouldCheckConditions := map[string]struct{}{
		"Active":              struct{}{},
		"Ready":               struct{}{},
		"ConfigurationsReady": struct{}{},
		"RoutesReady":         struct{}{},
		"ContainerHealthy":    struct{}{},
		"ResourcesAvailable":  struct{}{},
	}
	mustPassConditions := map[string]struct{}{
		"Ready":  struct{}{},
		"Active": struct{}{},
	}

	doCheck := func() (bool, error) {
		rvs, err := client.GetRevision(ctx, revisionName)
		// NotFound should be a retriable error.
		if err == provider.ErrRevisionNotFound {
			return true, err
		}
		if err != nil {
			return false, err
		}

		var (
			trueConds    = make(map[string]struct{}, 0)
			falseConds   = make([]string, 0, len(shouldCheckConditions))
			unknownConds = make([]string, 0, len(shouldCheckConditions))
		)
		if rvs.Status != nil {
			for _, cond := range rvs.Status.Conditions {
				if _, ok := shouldCheckConditions[cond.Type]; !ok {
					continue
				}
				switch cond.Status {
				case "True":
					trueConds[cond.Type] = struct{}{}
				case "False":
					falseConds = append(falseConds, cond.Message)
				default:
					unknownConds = append(unknownConds, cond.Message)
				}
			}
		}

		if len(falseConds) > 0 {
			return false, fmt.Errorf("%s", strings.Join(falseConds, "\n"))
		}
		if len(unknownConds) > 0 {
			return true, fmt.Errorf("%s", strings.Join(unknownConds, "\n"))
		}
		for k := range mustPassConditions {
			if _, ok := trueConds[k]; !ok {
				return true, fmt.Errorf("could not check status field %q", k)
			}
		}
		return false, nil
	}

	start := time.Now()
	for {
		retry, err := doCheck()
		if !retry {
			if err != nil {
				lp.Errorf("Revision %s was not ready: %v", revisionName, err)
				return err
			}
			lp.Infof("Revision %s is ready to receive traffic", revisionName)
			return nil
		}

		if time.Since(start) > retryTimeout {
			lp.Errorf("Revision %s was not ready: %v", revisionName, err)
			return err
		}

		lp.Infof("Revision %s is still not ready (%v), will retry after %v", revisionName, err, retryDuration)
		time.Sleep(retryDuration)
	}
}

func revisionExists(ctx context.Context, client provider.Client, revisionName string, lp executor.LogPersister) (bool, error) {
	_, err := client.GetRevision(ctx, revisionName)
	if err == nil {
		return true, nil
	}

	if err == provider.ErrRevisionNotFound {
		return false, nil
	}

	lp.Errorf("Failed while checking the existence of revision %s (%v)", revisionName, err)
	return false, err
}

func addBuiltinLabels(sm provider.ServiceManifest, hash, pipedID, appID, revisionName string, lp executor.LogPersister) bool {
	labels := map[string]string{
		provider.LabelManagedBy:   provider.ManagedByPiped,
		provider.LabelPiped:       pipedID,
		provider.LabelApplication: appID,
		provider.LabelCommitHash:  hash,
	}
	// Set builtinLabels for Service.
	sm.AddLabels(labels)

	if revisionName == "" {
		return true
	}
	// Set buildinLabels for Revision.
	labels[provider.LabelRevisionName] = revisionName
	if err := sm.AddRevisionLabels(labels); err != nil {
		lp.Errorf("Unable to add revision labels for the service manifest %s (%v)", sm.Name, err)
		return false
	}
	return true
}
