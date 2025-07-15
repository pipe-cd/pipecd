package store

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

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

type reflector struct {
	namespace            string
	targetMatcher        gvkMatcher
	resourceEventHandler resourceEventHandler

	kubeConfig *restclient.Config
	logger     *zap.Logger
}

// gvkMatcher is used to filter the resources that should be watched.
type gvkMatcher interface {
	// matchGVK returns true if the resource should be watched.
	matchGVK(gvk schema.GroupVersionKind) bool
}

// resourceEventHandler is used to handle the events of the resources.
type resourceEventHandler interface {
	// matchResource returns true if the resource's event should be handled by the event handler.
	matchResource(m provider.Manifest) bool
	// onAdd is called when a resource is added.
	onAdd(m provider.Manifest)
	// onUpdate is called when a resource is updated.
	onUpdate(old, new provider.Manifest)
	// onDelete is called when a resource is deleted.
	onDelete(m provider.Manifest)
}

func (r *reflector) start(ctx context.Context) error {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(r.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create discovery client: %v", err)
	}
	_, lists, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		return fmt.Errorf("failed to fetch groups and resources: %v", err)
	}

	var (
		targetResources, namespacedTargetResources []schema.GroupVersionResource
	)

	for _, list := range lists {
		for _, resource := range list.APIResources {
			gvk := schema.FromAPIVersionAndKind(list.GroupVersion, resource.Kind)
			if !r.targetMatcher.matchGVK(gvk) {
				r.logger.Info("skipping resource because it does not match the target matcher", zap.String("group", gvk.Group), zap.String("version", gvk.Version), zap.String("kind", gvk.Kind))
				continue
			}
			if !isSupportedWatch(resource) || !isSupportedList(resource) {
				r.logger.Info("skipping resource because it does not support watch or list", zap.String("group", gvk.Group), zap.String("version", gvk.Version), zap.String("kind", gvk.Kind))
				continue
			}
			gv := gvk.GroupVersion()
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

	dynamicClient, err := dynamic.NewForConfig(r.kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	// we can use ctx.Done() to handle the stop signal
	stopCh := ctx.Done()

	startInformer := func(namespace string, resources []schema.GroupVersionResource) {
		factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 30*time.Minute, namespace, nil)
		for _, tr := range resources {
			di := factory.ForResource(tr).Informer()
			di.AddEventHandler(r)
			di.Run(stopCh)
			if cache.WaitForCacheSync(stopCh, di.HasSynced) {
				r.logger.Info(fmt.Sprintf("informer cache for %v has been synced", tr))
			} else {
				// TODO: Handle the case informer cache has not been synced correctly.
				r.logger.Info(fmt.Sprintf("informer cache for %v has not been synced correctly", tr))
			}
		}
	}

	ns := r.namespace
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

// OnAdd implements cache.ResourceEventHandler.
func (r *reflector) OnAdd(obj interface{}) {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		r.logger.Error("failed to convert object to unstructured", zap.Any("object", obj))
		return
	}

	m := provider.FromUnstructured(u)
	if !r.resourceEventHandler.matchResource(m) {
		r.logger.Info("skipping resource because it does not match the resource matcher", zap.String("manifest", m.Key().ReadableString()))
		return
	}

	r.logger.Debug("received add event", zap.String("manifest", m.Key().ReadableString()))
	r.resourceEventHandler.onAdd(m)
}

// OnUpdate implements cache.ResourceEventHandler.
func (r *reflector) OnUpdate(oldObj, newObj interface{}) {
	u, ok := newObj.(*unstructured.Unstructured)
	if !ok {
		r.logger.Error("failed to convert object to unstructured", zap.Any("object", newObj))
		return
	}
	oldU, ok := oldObj.(*unstructured.Unstructured)
	if !ok {
		r.logger.Error("failed to convert object to unstructured", zap.Any("object", oldObj))
		return
	}

	m := provider.FromUnstructured(u)
	oldM := provider.FromUnstructured(oldU)
	if !r.resourceEventHandler.matchResource(m) {
		r.logger.Info("skipping resource because it does not match the resource matcher", zap.String("manifest", m.Key().ReadableString()))
		return
	}

	r.logger.Debug("received update event", zap.String("manifest", m.Key().ReadableString()))
	r.resourceEventHandler.onUpdate(oldM, m)
}

// OnDelete implements cache.ResourceEventHandler.
func (r *reflector) OnDelete(obj interface{}) {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		r.logger.Error("failed to convert object to unstructured", zap.Any("object", obj))
		return
	}

	m := provider.FromUnstructured(u)
	if !r.resourceEventHandler.matchResource(m) {
		r.logger.Info("skipping resource because it does not match the resource matcher", zap.String("manifest", m.Key().ReadableString()))
		return
	}

	r.logger.Debug("received delete event", zap.String("manifest", m.Key().ReadableString()))
	r.resourceEventHandler.onDelete(m)
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
