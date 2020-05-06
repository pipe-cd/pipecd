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

package appstatestore

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	// At this version, we only care about Ingress, Service, Deployment and Pod,
	// but in the next version we will expand to watch all resource kinds.
	groupWhitelist = map[string]struct{}{
		"":                  {},
		"apps":              {},
		"extensions":        {},
		"networking.k8s.io": {},
	}
	versionWhitelist = map[string]struct{}{
		"v1":      {},
		"v1beta1": {},
		"v1beta2": {},
	}
	kindWhitelist = map[string]struct{}{
		"Ingress":     {},
		"Service":     {},
		"Deployment":  {},
		"DaemonSet":   {},
		"StatefulSet": {},
		"ReplicaSet":  {},
		"Pod":         {},
		"ConfigMap":   {},
		"Secret":      {},
	}
	ignoreResourceKeys = map[string]struct{}{
		"v1/Service/default/kubernetes":               {},
		"v1/Service/kube-system/heapster":             {},
		"v1/Service/kube-system/metrics-server":       {},
		"v1/Service/kube-system/kube-dns":             {},
		"v1/Service/kube-system/kubernetes-dashboard": {},
		"v1/Service/kube-system/default-http-backend": {},

		"apps/v1/Deployment/kube-system/kube-dns":                                 {},
		"apps/v1/Deployment/kube-system/kube-dns-autoscaler":                      {},
		"apps/v1/Deployment/kube-system/fluentd-gcp-scaler":                       {},
		"apps/v1/Deployment/kube-system/kubernetes-dashboard":                     {},
		"apps/v1/Deployment/kube-system/l7-default-backend":                       {},
		"apps/v1/Deployment/kube-system/heapster-gke":                             {},
		"apps/v1/Deployment/kube-system/stackdriver-metadata-agent-cluster-level": {},

		"extensions/v1beta1/Deployment/kube-system/kube-dns":                                 {},
		"extensions/v1beta1/Deployment/kube-system/kube-dns-autoscaler":                      {},
		"extensions/v1beta1/Deployment/kube-system/fluentd-gcp-scaler":                       {},
		"extensions/v1beta1/Deployment/kube-system/kubernetes-dashboard":                     {},
		"extensions/v1beta1/Deployment/kube-system/l7-default-backend":                       {},
		"extensions/v1beta1/Deployment/kube-system/heapster-gke":                             {},
		"extensions/v1beta1/Deployment/kube-system/stackdriver-metadata-agent-cluster-level": {},

		"v1/ConfigMap/kube-system/ingress-gce-lock":          {},
		"v1/ConfigMap/kube-system/gke-common-webhook-lock":   {},
		"v1/ConfigMap/kube-system/cluster-autoscaler-status": {},
	}
)

// reflector watches the live state of applicaiton with the cluster
// and triggers the specified callbacks.
type reflector struct {
	kubeConfig *restclient.Config

	onAdd    func(obj *unstructured.Unstructured)
	onUpdate func(oldObj, obj *unstructured.Unstructured)
	onDelete func(obj *unstructured.Unstructured)

	stopCh chan struct{}
	logger *zap.Logger
}

func (r *reflector) start(ctx context.Context) error {
	// Use discovery to discover APIs supported by the Kubernetes API server.
	// This should be run periodically with a low rate because the APIs are not added frequently.
	// https://godoc.org/k8s.io/client-go/discovery
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(r.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create discovery client: %v", err)
	}
	groupResources, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return fmt.Errorf("failed to fetch preferred resources: %v", err)
	}
	r.logger.Info(fmt.Sprintf("successfully prefered resources that contains for %d groups", len(groupResources)))

	// Filter above APIResources.
	targetResources := make([]schema.GroupVersionResource, 0)
	for _, gr := range groupResources {
		for _, r := range gr.APIResources {
			if _, ok := kindWhitelist[r.Kind]; !ok {
				continue
			}
			gvk := schema.FromAPIVersionAndKind(gr.GroupVersion, r.Kind)
			gv := gvk.GroupVersion()
			if _, ok := groupWhitelist[gv.Group]; !ok {
				continue
			}
			if _, ok := versionWhitelist[gv.Version]; !ok {
				continue
			}
			target := gv.WithResource(r.Name)
			targetResources = append(targetResources, target)
		}
	}
	r.logger.Info("filtered target resources", zap.Any("targetResources", targetResources))

	// Use dynamic to perform generic operations on arbitrary Kubernets API objects.
	// https://godoc.org/k8s.io/client-go/dynamic
	dynamicClient, err := dynamic.NewForConfig(r.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 30*time.Minute, metav1.NamespaceAll, nil)
	stopCh := make(chan struct{})

	for _, tr := range targetResources {
		di := factory.ForResource(tr).Informer()
		di.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc:    r.onObjectAdd,
			UpdateFunc: r.onObjectUpdate,
			DeleteFunc: r.onObjectDelete,
		})
		go di.Run(r.stopCh)
		if cache.WaitForCacheSync(stopCh, di.HasSynced) {
			r.logger.Info(fmt.Sprintf("informer cache for %v has been synced", tr))
		} else {
			// TODO: Handle the case informer cache has not been synced correctly.
			r.logger.Info(fmt.Sprintf("informer cache for %v has not been synced correctly", tr))
		}
	}

	r.logger.Info("all informer caches have been synced")
	return nil
}

func (r *reflector) onObjectAdd(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	key := makeResourceKey(u)
	if _, ok := ignoreResourceKeys[key.String()]; ok {
		return
	}
	r.logger.Info("received add event", zap.Stringer("key", key))
	r.onAdd(u)
}

func (r *reflector) onObjectUpdate(oldObj, obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	oldU := oldObj.(*unstructured.Unstructured)
	key := makeResourceKey(u)
	if _, ok := ignoreResourceKeys[key.String()]; ok {
		return
	}
	r.logger.Info("received update event", zap.Stringer("key", key))
	r.onUpdate(oldU, u)
}

func (r *reflector) onObjectDelete(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	key := makeResourceKey(u)
	if _, ok := ignoreResourceKeys[key.String()]; ok {
		return
	}
	r.logger.Info("received delete event", zap.Stringer("key", key))
	r.onDelete(u)
}
