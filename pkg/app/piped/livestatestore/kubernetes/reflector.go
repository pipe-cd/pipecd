// Copyright 2023 The PipeCD Authors.
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

package kubernetes

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

	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/kubernetes/kubernetesmetrics"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/config"
)

var (
	// This is the default whitelist of resources that should be watched.
	// User can add/remove other resources to be watched in piped config at cloud provider part.
	groupWhitelist = map[string]struct{}{
		"":                          {},
		"apps":                      {},
		"extensions":                {},
		"batch":                     {},
		"storage.k8s.io":            {},
		"autoscaling":               {},
		"networking.k8s.io":         {},
		"apiextensions.k8s.io":      {},
		"rbac.authorization.k8s.io": {},
		"policy":                    {},
		"apiregistration.k8s.io":    {},
		"authorization.k8s.io":      {},
	}
	versionWhitelist = map[string]struct{}{
		"v1":      {},
		"v1beta1": {},
		"v1beta2": {},
	}
	kindWhitelist = map[string]struct{}{
		"Service":                  {},
		"Endpoints":                {},
		"Deployment":               {},
		"DaemonSet":                {},
		"StatefulSet":              {},
		"ReplicationController":    {},
		"ReplicaSet":               {},
		"Pod":                      {},
		"Job":                      {},
		"CronJob":                  {},
		"ConfigMap":                {},
		"Secret":                   {},
		"Ingress":                  {},
		"NetworkPolicy":            {},
		"StorageClass":             {},
		"PersistentVolume":         {},
		"PersistentVolumeClaim":    {},
		"HorizontalPodAutoscaler":  {},
		"ServiceAccount":           {},
		"Role":                     {},
		"RoleBinding":              {},
		"ClusterRole":              {},
		"ClusterRoleBinding":       {},
		"CustomResourceDefinition": {},
		"PodDisruptionBudget":      {},
		"PodSecurityPolicy":        {},
		"APIService":               {},
		"LocalSubjectAccessReview": {},
		"SelfSubjectAccessReview":  {},
		"SelfSubjectRulesReview":   {},
		"SubjectAccessReview":      {},
		"ResourceQuota":            {},
		"PodTemplate":              {},
		"IngressClass":             {},
		"Namespace":                {},
	}
	ignoreResourceKeys = map[string]struct{}{
		"v1:Service:default:kubernetes":               {},
		"v1:Service:kube-system:heapster":             {},
		"v1:Service:kube-system:metrics-server":       {},
		"v1:Service:kube-system:kube-dns":             {},
		"v1:Service:kube-system:kubernetes-dashboard": {},
		"v1:Service:kube-system:default-http-backend": {},

		"apps/v1:Deployment:kube-system:kube-dns":                                 {},
		"apps/v1:Deployment:kube-system:kube-dns-autoscaler":                      {},
		"apps/v1:Deployment:kube-system:fluentd-gcp-scaler":                       {},
		"apps/v1:Deployment:kube-system:kubernetes-dashboard":                     {},
		"apps/v1:Deployment:kube-system:l7-default-backend":                       {},
		"apps/v1:Deployment:kube-system:heapster-gke":                             {},
		"apps/v1:Deployment:kube-system:stackdriver-metadata-agent-cluster-level": {},

		"extensions/v1beta1:Deployment:kube-system:kube-dns":                                 {},
		"extensions/v1beta1:Deployment:kube-system:kube-dns-autoscaler":                      {},
		"extensions/v1beta1:Deployment:kube-system:fluentd-gcp-scaler":                       {},
		"extensions/v1beta1:Deployment:kube-system:kubernetes-dashboard":                     {},
		"extensions/v1beta1:Deployment:kube-system:l7-default-backend":                       {},
		"extensions/v1beta1:Deployment:kube-system:heapster-gke":                             {},
		"extensions/v1beta1:Deployment:kube-system:stackdriver-metadata-agent-cluster-level": {},

		"v1:Endpoints:kube-system:kube-controller-manager":        {},
		"v1:Endpoints:kube-system:kube-scheduler":                 {},
		"v1:Endpoints:kube-system:vpa-recommender":                {},
		"v1:Endpoints:kube-system:gcp-controller-manager":         {},
		"v1:Endpoints:kube-system:managed-certificate-controller": {},
		"v1:Endpoints:kube-system:cluster-autoscaler":             {},

		"v1:ConfigMap:kube-system:cluster-kubestore":         {},
		"v1:ConfigMap:kube-system:ingress-gce-lock":          {},
		"v1:ConfigMap:kube-system:gke-common-webhook-lock":   {},
		"v1:ConfigMap:kube-system:cluster-autoscaler-status": {},

		"rbac.authorization.k8s.io/v1:ClusterRole::system:managed-certificate-controller":        {},
		"rbac.authorization.k8s.io/v1:ClusterRoleBinding::system:managed-certificate-controller": {},
	}
)

// reflector watches the live state of application with the cluster
// and triggers the specified callbacks.
type reflector struct {
	config      *config.PlatformProviderKubernetesConfig
	kubeConfig  *restclient.Config
	pipedConfig *config.PipedSpec

	onAdd    func(obj *unstructured.Unstructured)
	onUpdate func(oldObj, obj *unstructured.Unstructured)
	onDelete func(obj *unstructured.Unstructured)

	watchingResourceKinds []provider.APIVersionKind
	stopCh                chan struct{}
	logger                *zap.Logger
}

func (r *reflector) start(ctx context.Context) error {
	matcher := newResourceMatcher(r.config.AppStateInformer)

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
	r.logger.Info(fmt.Sprintf("successfully preferred resources that contains for %d groups", len(groupResources)))

	// Filter above APIResources.
	var (
		targetResources           = make([]schema.GroupVersionResource, 0)
		namespacedTargetResources = make([]schema.GroupVersionResource, 0)
	)
	for _, gr := range groupResources {
		for _, resource := range gr.APIResources {
			gvk := schema.FromAPIVersionAndKind(gr.GroupVersion, resource.Kind)
			if !matcher.Match(gvk) {
				r.logger.Info(fmt.Sprintf("skip watching %v because of not matching the configured list", gvk))
				continue
			}

			if !isSupportedList(resource) || !isSupportedWatch(resource) {
				r.logger.Info(fmt.Sprintf("skip watching %v because of not supporting watch or list verb", gvk))
				continue
			}

			gv := gvk.GroupVersion()
			r.watchingResourceKinds = append(r.watchingResourceKinds, provider.APIVersionKind{
				APIVersion: gv.String(),
				Kind:       gvk.Kind,
			})
			target := gv.WithResource(resource.Name)
			if resource.Namespaced {
				namespacedTargetResources = append(namespacedTargetResources, target)
			} else {
				targetResources = append(targetResources, target)
			}
		}
	}
	r.logger.Info("filtered target resources",
		zap.Any("targetResources", targetResources),
		zap.Any("namespacedTargetResources", namespacedTargetResources),
	)

	// Use dynamic to perform generic operations on arbitrary Kubernets API objects.
	// https://godoc.org/k8s.io/client-go/dynamic
	dynamicClient, err := dynamic.NewForConfig(r.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	stopCh := make(chan struct{})

	startInformer := func(namespace string, resources []schema.GroupVersionResource) {
		factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 30*time.Minute, namespace, nil)
		for _, tr := range resources {
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
	}

	ns := r.config.AppStateInformer.Namespace
	if ns == "" {
		ns = metav1.NamespaceAll
	}
	r.logger.Info(fmt.Sprintf("start running %d namespaced-resource informers", len(namespacedTargetResources)))
	startInformer(ns, namespacedTargetResources)

	if ns == metav1.NamespaceAll {
		r.logger.Info(fmt.Sprintf("start running %d non-namespaced-resource informers", len(targetResources)))
		startInformer(metav1.NamespaceAll, targetResources)
	}

	r.logger.Info("all informer caches have been synced")
	return nil
}

func (r *reflector) onObjectAdd(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	key := provider.MakeResourceKeyFromCluster(u)

	// Ignore all predefined ones.
	if _, ok := ignoreResourceKeys[key.String()]; ok {
		kubernetesmetrics.IncResourceEventsCounter(
			kubernetesmetrics.LabelEventAdd,
			kubernetesmetrics.LabelEventNotYetHandled,
		)
		return
	}

	// Ignore all objects that are not handled by this piped.
	pipedID := u.GetAnnotations()[provider.LabelPiped]
	if pipedID != "" && pipedID != r.pipedConfig.PipedID {
		kubernetesmetrics.IncResourceEventsCounter(
			kubernetesmetrics.LabelEventAdd,
			kubernetesmetrics.LabelEventNotYetHandled,
		)
		return
	}

	r.logger.Debug(fmt.Sprintf("received add event for %s", key.String()))
	r.onAdd(u)
	kubernetesmetrics.IncResourceEventsCounter(
		kubernetesmetrics.LabelEventAdd,
		kubernetesmetrics.LabelEventHandled,
	)
}

func (r *reflector) onObjectUpdate(oldObj, obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	oldU := oldObj.(*unstructured.Unstructured)

	// Ignore all predefined ones.
	key := provider.MakeResourceKeyFromCluster(u)
	if _, ok := ignoreResourceKeys[key.String()]; ok {
		kubernetesmetrics.IncResourceEventsCounter(
			kubernetesmetrics.LabelEventUpdate,
			kubernetesmetrics.LabelEventNotYetHandled,
		)
		return
	}

	// Ignore all objects that are not handled by this piped.
	pipedID := u.GetAnnotations()[provider.LabelPiped]
	if pipedID != "" && pipedID != r.pipedConfig.PipedID {
		kubernetesmetrics.IncResourceEventsCounter(
			kubernetesmetrics.LabelEventUpdate,
			kubernetesmetrics.LabelEventNotYetHandled,
		)
		return
	}

	r.logger.Debug(fmt.Sprintf("received update event for %s", key.String()))
	r.onUpdate(oldU, u)
	kubernetesmetrics.IncResourceEventsCounter(
		kubernetesmetrics.LabelEventUpdate,
		kubernetesmetrics.LabelEventHandled,
	)
}

func (r *reflector) onObjectDelete(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	key := provider.MakeResourceKeyFromCluster(u)

	// Ignore all predefined ones.
	if _, ok := ignoreResourceKeys[key.String()]; ok {
		kubernetesmetrics.IncResourceEventsCounter(
			kubernetesmetrics.LabelEventDelete,
			kubernetesmetrics.LabelEventNotYetHandled,
		)
		return
	}

	// Ignore all objects that are not handled by this piped.
	pipedID := u.GetAnnotations()[provider.LabelPiped]
	if pipedID != "" && pipedID != r.pipedConfig.PipedID {
		kubernetesmetrics.IncResourceEventsCounter(
			kubernetesmetrics.LabelEventDelete,
			kubernetesmetrics.LabelEventNotYetHandled,
		)
		return
	}

	r.logger.Debug(fmt.Sprintf("received delete event for %s", key.String()))
	r.onDelete(u)
	kubernetesmetrics.IncResourceEventsCounter(
		kubernetesmetrics.LabelEventDelete,
		kubernetesmetrics.LabelEventHandled,
	)
}

func isSupportedWatch(r metav1.APIResource) bool {
	for _, v := range r.Verbs {
		if v == "watch" {
			return true
		}
	}
	return false
}

func isSupportedList(r metav1.APIResource) bool {
	for _, v := range r.Verbs {
		if v == "list" {
			return true
		}
	}
	return false
}

type resourceMatcher struct {
	includes map[string]struct{}
	excludes map[string]struct{}
}

func newResourceMatcher(cfg config.KubernetesAppStateInformer) *resourceMatcher {
	r := &resourceMatcher{
		includes: make(map[string]struct{}, len(cfg.IncludeResources)),
		excludes: make(map[string]struct{}, len(cfg.ExcludeResources)),
	}

	for _, m := range cfg.IncludeResources {
		if m.Kind == "" {
			r.includes[m.APIVersion] = struct{}{}
		} else {
			r.includes[m.APIVersion+":"+m.Kind] = struct{}{}
		}
	}
	for _, m := range cfg.ExcludeResources {
		if m.Kind == "" {
			r.excludes[m.APIVersion] = struct{}{}
		} else {
			r.excludes[m.APIVersion+":"+m.Kind] = struct{}{}
		}
	}
	return r
}

func (m *resourceMatcher) Match(gvk schema.GroupVersionKind) bool {
	var (
		gv         = gvk.GroupVersion()
		apiVersion = gv.String()
		key        = apiVersion + ":" + gvk.Kind
	)

	// Any resource matches the specified ExcludeResources will be ignored.
	if _, ok := m.excludes[apiVersion]; ok {
		return false
	}
	if _, ok := m.excludes[key]; ok {
		return false
	}

	// Any resources matches the specified IncludeResources will be included.
	if _, ok := m.includes[apiVersion]; ok {
		return true
	}
	if _, ok := m.includes[key]; ok {
		return true
	}

	// Check the predefined list.
	if _, ok := kindWhitelist[gvk.Kind]; !ok {
		return false
	}
	if _, ok := groupWhitelist[gv.Group]; !ok {
		return false
	}
	if _, ok := versionWhitelist[gv.Version]; !ok {
		return false
	}

	return true
}
