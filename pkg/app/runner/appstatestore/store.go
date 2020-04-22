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

// Package appstatestore provides a runner component
// that watches the live state of applications in the cluster
// to construct it cache data that will be used to provide
// data to another components quickly.
package appstatestore

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
)

// AppStateStore syncs the live state of applicaiton with the cluster
// and provides some functions for other components to query those states.
type AppStateStore struct {
	kubeConfig  *restclient.Config
	gracePeriod time.Duration
	logger      *zap.Logger
}

func NewStore(kubeConfig *restclient.Config, gracePeriod time.Duration, logger *zap.Logger) *AppStateStore {
	return &AppStateStore{
		kubeConfig:  kubeConfig,
		gracePeriod: gracePeriod,
		logger:      logger,
	}
}

func (s *AppStateStore) Run(ctx context.Context) error {
	// 1. Use discovery to discover APIs supported by the Kubernetes API server.
	// This should be run periodically with a low rate because the APIs are not added frequently.
	// https://godoc.org/k8s.io/client-go/discovery
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(s.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create discovery client: %v", err)
	}
	groupResources, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return fmt.Errorf("failed to fetch preferred resources: %v", err)
	}
	// s.logger.Info("successfully prefered resources", zap.Any("groupResources", groupResources))

	// 2. Filter above APIResources
	// - Current version, we may need to care about Ingress, Service, Deployment, Pods
	// - Or for checking the diff, we may need to care about all of resources
	// which support list and watch operator.
	targetResources := make([]schema.GroupVersionResource, 0)
	for _, gr := range groupResources {
		for _, r := range gr.APIResources {
			if r.Kind != "Deployment" {
				continue
			}
			gvk := schema.FromAPIVersionAndKind(gr.GroupVersion, r.Kind)
			gv := gvk.GroupVersion()
			target := gv.WithResource(r.Name)
			targetResources = append(targetResources, target)
		}
	}
	s.logger.Info("filtered target resources", zap.Any("targetResources", targetResources))

	// 3. Use dynamic to perform generic operations on arbitrary Kubernets API objects.
	// https://godoc.org/k8s.io/client-go/dynamic
	dynamicClient, err := dynamic.NewForConfig(s.kubeConfig)
	for _, r := range targetResources {
		rc := dynamicClient.Resource(r)
		list, err := rc.List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("failed to list resources %v, err: %v", r, err)
		}
		s.logger.Info("successfully listed resources",
			zap.Any("resource", r),
			zap.Any("list", list.Items[0]),
		)
	}

	<-ctx.Done()
	return nil
}
