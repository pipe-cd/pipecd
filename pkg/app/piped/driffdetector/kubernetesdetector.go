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

package driffdetector

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/kapetaniosci/pipe/pkg/app/piped/livestatestore/kubernetes"
	"github.com/kapetaniosci/pipe/pkg/config"
)

type kubernetesDetector struct {
	provider    config.PipedCloudProvider
	appLister   applicationLister
	stateGetter kubernetes.Getter
	apiClient   apiClient
	interval    time.Duration
	logger      *zap.Logger
}

func newKubernetesDetector(cp config.PipedCloudProvider, appLister applicationLister, stateGetter kubernetes.Getter, apiClient apiClient, logger *zap.Logger) *kubernetesDetector {
	logger = logger.Named("kubernetes-detector").With(
		zap.String("cloud-provider", cp.Name),
	)
	return &kubernetesDetector{
		provider:    cp,
		appLister:   appLister,
		stateGetter: stateGetter,
		apiClient:   apiClient,
		interval:    time.Minute,
		logger:      logger,
	}
}

func (d *kubernetesDetector) Run(ctx context.Context) error {
	d.logger.Info("start running driff detector for kubernetes applications")

	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

L:
	for {
		select {

		case <-ticker.C:
			d.check(ctx)

		case <-ctx.Done():
			break L
		}
	}

	d.logger.Info("driff detector for kubernetes applications has been stopped")
	return nil
}

func (d *kubernetesDetector) check(ctx context.Context) error {
	apps := d.appLister.ListByCloudProvider(d.provider.Name)
	for _, app := range apps {
		liveManifests := d.stateGetter.GetAppLiveManifests(app.Id)
		d.logger.Info(fmt.Sprintf("application %s has %d live manifests", app.Id, len(liveManifests)))
	}
	return nil
}

func (d *kubernetesDetector) ProviderName() string {
	return d.provider.Name
}
