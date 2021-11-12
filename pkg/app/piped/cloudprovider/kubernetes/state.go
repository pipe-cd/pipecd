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

// Copyright Istio Authors
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
	"fmt"
	"sort"
	"strings"
	"time"

	istio "istio.io/api/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/pipe-cd/pipe/pkg/model"
)

const (
	Resolution_NONE            = "NONE"
	Resolution_STATIC          = "STATIC"
	Resolution_DNS             = "DNS"
	Resolution_DNS_ROUND_ROBIN = "DNS_ROUND_ROBIN"
)

func MakeKubernetesResourceState(uid string, key ResourceKey, obj *unstructured.Unstructured, now time.Time) model.KubernetesResourceState {
	var (
		owners       = obj.GetOwnerReferences()
		ownerIDs     = make([]string, 0, len(owners))
		creationTime = obj.GetCreationTimestamp()
		status, desc = determineResourceHealth(key, obj)
	)

	for _, owner := range owners {
		ownerIDs = append(ownerIDs, string(owner.UID))
	}
	sort.Strings(ownerIDs)

	state := model.KubernetesResourceState{
		Id:       uid,
		OwnerIds: ownerIDs,
		// TODO: Think about adding more parents by using label selectors
		ParentIds:  ownerIDs,
		Name:       key.Name,
		ApiVersion: key.APIVersion,
		Kind:       key.Kind,
		Namespace:  obj.GetNamespace(),

		HealthStatus:      status,
		HealthDescription: desc,

		CreatedAt: creationTime.Unix(),
		UpdatedAt: now.Unix(),
	}

	return state
}

func determineResourceHealth(key ResourceKey, obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	if !IsKubernetesBuiltInResource(key.APIVersion) {
		desc = fmt.Sprintf("\"%s/%s\" was applied successfully but its health status couldn't be determined exactly. (Because tracking status for this kind of resource is not supported yet.)", key.APIVersion, key.Kind)
		return
	}

	switch key.Kind {
	case KindDeployment:
		return determineDeploymentHealth(obj)
	case KindStatefulSet:
		return determineStatefulSetHealth(obj)
	case KindDaemonSet:
		return determineDaemonSetHealth(obj)
	case KindReplicaSet:
		return determineReplicaSetHealth(obj)
	case KindPod:
		return determinePodHealth(obj)
	case KindJob:
		return determineJobHealth(obj)
	case KindCronJob:
		return determineCronJobHealth(obj)
	case KindService:
		return determineServiceHealth(obj)
	case KindIngress:
		return determineIngressHealth(obj)
	case KindConfigMap:
		return determineConfigMapHealth(obj)
	case KindPersistentVolume:
		return determinePersistentVolumeHealth(obj)
	case KindPersistentVolumeClaim:
		return determinePVCHealth(obj)
	case KindSecret:
		return determineSecretHealth(obj)
	case KindServiceAccount:
		return determineServiceAccountHealth(obj)
	case KindRole:
		return determineRoleHealth(obj)
	case KindRoleBinding:
		return determineRoleBindingHealth(obj)
	case KindClusterRole:
		return determineClusterRoleHealth(obj)
	case KindClusterRoleBinding:
		return determineClusterRoleBindingHealth(obj)
	case KindVirtualService:
		return determineVirtualService(obj)
	case KindDestinationRule:
		return determineDestinationRule(obj)
	case KindGateway:
		return determineGateway(obj)
	case KindServiceEntry:
		return determineServiceEntry(obj)
	default:
		desc = "Unimplemented or unknown resource"
		return
	}
}

func determineRoleHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = fmt.Sprintf("%q was applied successfully", obj.GetName())
	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineRoleBindingHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = fmt.Sprintf("%q was applied successfully", obj.GetName())
	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineClusterRoleHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = fmt.Sprintf("%q was applied successfully", obj.GetName())
	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineClusterRoleBindingHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = fmt.Sprintf("%q was applied successfully", obj.GetName())
	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineDeploymentHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	d := &appsv1.Deployment{}
	err := scheme.Scheme.Convert(obj, d, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, d, err)
		return
	}

	status = model.KubernetesResourceState_OTHER
	if d.Spec.Paused {
		desc = "Deployment is paused"
		return
	}

	// Referred to:
	//   https://github.com/kubernetes/kubernetes/blob/7942dca975b7be9386540df3c17e309c3cb2de60/staging/src/k8s.io/kubectl/pkg/polymorphichelpers/rollout_status.go#L75
	if d.Generation > d.Status.ObservedGeneration {
		desc = "Waiting for rollout to finish because observed deployment generation less than desired generation"
		return
	}
	// TimedOutReason is added in a deployment when its newest replica set fails to show any progress
	// within the given deadline (progressDeadlineSeconds).
	const timedOutReason = "ProgressDeadlineExceeded"
	var cond *appsv1.DeploymentCondition
	for i := range d.Status.Conditions {
		c := d.Status.Conditions[i]
		if c.Type == appsv1.DeploymentProgressing {
			cond = &c
			break
		}
	}
	if cond != nil && cond.Reason == timedOutReason {
		desc = fmt.Sprintf("Deployment %q exceeded its progress deadline", obj.GetName())
	}

	if d.Spec.Replicas == nil {
		desc = "The number of desired replicas is unspecified"
		return
	}
	if d.Status.UpdatedReplicas < *d.Spec.Replicas {
		desc = fmt.Sprintf("Waiting for remaining %d/%d replicas to be updated", d.Status.UpdatedReplicas, *d.Spec.Replicas)
		return
	}
	if d.Status.UpdatedReplicas < d.Status.Replicas {
		desc = fmt.Sprintf("%d old replicas are pending termination", d.Status.Replicas-d.Status.UpdatedReplicas)
		return
	}
	if d.Status.AvailableReplicas < d.Status.Replicas {
		desc = fmt.Sprintf("Waiting for remaining %d/%d replicas to be available", d.Status.Replicas-d.Status.AvailableReplicas, d.Status.Replicas)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineStatefulSetHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	s := &appsv1.StatefulSet{}
	err := scheme.Scheme.Convert(obj, s, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, s, err)
		return
	}

	// Referred to:
	//   https://github.com/kubernetes/kubernetes/blob/7942dca975b7be9386540df3c17e309c3cb2de60/staging/src/k8s.io/kubectl/pkg/polymorphichelpers/rollout_status.go#L130-L149
	status = model.KubernetesResourceState_OTHER
	if s.Status.ObservedGeneration == 0 || s.Generation > s.Status.ObservedGeneration {
		desc = "Waiting for statefulset spec update to be observed"
		return
	}

	if s.Spec.Replicas == nil {
		desc = "The number of desired replicas is unspecified"
		return
	}
	if *s.Spec.Replicas != s.Status.ReadyReplicas {
		desc = fmt.Sprintf("The number of ready replicas (%d) is different from the desired number (%d)", s.Status.ReadyReplicas, *s.Spec.Replicas)
		return
	}

	// Check if the partitioned roll out is in progress.
	if s.Spec.UpdateStrategy.Type == appsv1.RollingUpdateStatefulSetStrategyType && s.Spec.UpdateStrategy.RollingUpdate != nil {
		if s.Spec.Replicas != nil && s.Spec.UpdateStrategy.RollingUpdate.Partition != nil {
			if s.Status.UpdatedReplicas < (*s.Spec.Replicas - *s.Spec.UpdateStrategy.RollingUpdate.Partition) {
				desc = fmt.Sprintf("Waiting for partitioned roll out to finish because %d out of %d new pods have been updated",
					s.Status.UpdatedReplicas, (*s.Spec.Replicas - *s.Spec.UpdateStrategy.RollingUpdate.Partition))
				return
			}
		}
		status = model.KubernetesResourceState_HEALTHY
		return
	}

	if s.Status.UpdateRevision != s.Status.CurrentRevision {
		desc = fmt.Sprintf("Waiting for statefulset rolling update to complete %d pods at revision %s", s.Status.UpdatedReplicas, s.Status.UpdateRevision)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineDaemonSetHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	d := &appsv1.DaemonSet{}
	err := scheme.Scheme.Convert(obj, d, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, d, err)
		return
	}

	// Referred to:
	//   https://github.com/kubernetes/kubernetes/blob/7942dca975b7be9386540df3c17e309c3cb2de60/staging/src/k8s.io/kubectl/pkg/polymorphichelpers/rollout_status.go#L107-L115
	status = model.KubernetesResourceState_OTHER
	if d.Status.ObservedGeneration == 0 || d.Generation > d.Status.ObservedGeneration {
		desc = "Waiting for rollout to finish because observed daemon set generation less than desired generation"
		return
	}
	if d.Status.UpdatedNumberScheduled < d.Status.DesiredNumberScheduled {
		desc = fmt.Sprintf("Waiting for daemon set %q rollout to finish because %d out of %d new pods have been updated", d.Name, d.Status.UpdatedNumberScheduled, d.Status.DesiredNumberScheduled)
		return
	}
	if d.Status.NumberAvailable < d.Status.DesiredNumberScheduled {
		desc = fmt.Sprintf("Waiting for daemon set %q rollout to finish because %d of %d updated pods are available", d.Name, d.Status.NumberAvailable, d.Status.DesiredNumberScheduled)
		return
	}

	if d.Status.NumberMisscheduled > 0 {
		desc = fmt.Sprintf("%d nodes that are running the daemon pod, but are not supposed to run the daemon pod", d.Status.NumberMisscheduled)
		return
	}
	if d.Status.NumberUnavailable > 0 {
		desc = fmt.Sprintf("%d nodes that should be running the daemon pod and have none of the daemon pod running and available", d.Status.NumberUnavailable)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineReplicaSetHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	r := &appsv1.ReplicaSet{}
	err := scheme.Scheme.Convert(obj, r, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, r, err)
		return
	}

	status = model.KubernetesResourceState_OTHER
	if r.Status.ObservedGeneration == 0 || r.Generation > r.Status.ObservedGeneration {
		desc = "Waiting for rollout to finish because observed replica set generation less than desired generation"
		return
	}

	var cond *appsv1.ReplicaSetCondition
	for i := range r.Status.Conditions {
		c := r.Status.Conditions[i]
		if c.Type == appsv1.ReplicaSetReplicaFailure {
			cond = &c
			break
		}
	}
	switch {
	case cond != nil && cond.Status == corev1.ConditionTrue:
		desc = cond.Message
		return
	case r.Spec.Replicas == nil:
		desc = "The number of desired replicas is unspecified"
		return
	case r.Status.AvailableReplicas < *r.Spec.Replicas:
		desc = fmt.Sprintf("Waiting for rollout to finish because only %d/%d replicas are available", r.Status.AvailableReplicas, *r.Spec.Replicas)
		return
	case *r.Spec.Replicas != r.Status.ReadyReplicas:
		desc = fmt.Sprintf("The number of ready replicas (%d) is different from the desired number (%d)", r.Status.ReadyReplicas, *r.Spec.Replicas)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineCronJobHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = fmt.Sprintf("%q was applied successfully", obj.GetName())
	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineJobHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	job := &batchv1.Job{}
	err := scheme.Scheme.Convert(obj, job, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, job, err)
		return
	}

	var (
		failed    bool
		completed bool
		message   string
	)
	for _, condition := range job.Status.Conditions {
		switch condition.Type {
		case batchv1.JobFailed:
			failed = true
			completed = true
			message = condition.Message
		case batchv1.JobComplete:
			completed = true
			message = condition.Message
		}
	}

	switch {
	case !completed:
		status = model.KubernetesResourceState_HEALTHY
		desc = "Job is in progress"
	case failed:
		status = model.KubernetesResourceState_OTHER
		desc = message
	default:
		status = model.KubernetesResourceState_HEALTHY
		desc = message
	}

	return
}

func determinePodHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	p := &corev1.Pod{}
	err := scheme.Scheme.Convert(obj, p, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, p, err)
		return
	}

	// Determine based on its container statuses.
	if p.Spec.RestartPolicy == corev1.RestartPolicyAlways {
		var messages []string
		for _, s := range p.Status.ContainerStatuses {
			waiting := s.State.Waiting
			if waiting == nil {
				continue
			}
			if strings.HasPrefix(waiting.Reason, "Err") || strings.HasSuffix(waiting.Reason, "Error") || strings.HasSuffix(waiting.Reason, "BackOff") {
				status = model.KubernetesResourceState_OTHER
				messages = append(messages, waiting.Message)
			}
		}

		if status == model.KubernetesResourceState_OTHER {
			desc = strings.Join(messages, ", ")
			return
		}
	}

	// Determine based on its phase.
	switch p.Status.Phase {
	case corev1.PodRunning, corev1.PodSucceeded:
		status = model.KubernetesResourceState_HEALTHY
		desc = p.Status.Message
	default:
		status = model.KubernetesResourceState_OTHER
		desc = p.Status.Message
	}
	return
}

func determineIngressHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	check := func(ingressList []corev1.LoadBalancerIngress) {
		if len(ingressList) == 0 {
			status = model.KubernetesResourceState_OTHER
			desc = "Ingress points for the load-balancer are in progress"
			return
		}
		status = model.KubernetesResourceState_HEALTHY
		return
	}

	v1Ingress := &networkingv1.Ingress{}
	err := scheme.Scheme.Convert(obj, v1Ingress, nil)
	if err == nil {
		check(v1Ingress.Status.LoadBalancer.Ingress)
		return
	}

	// PipeCD keeps supporting Kubernetes < v1.22 for the meantime so checks deprecated versions as well.

	betaIngress := &networkingv1beta1.Ingress{}
	err = scheme.Scheme.Convert(obj, betaIngress, nil)
	if err == nil {
		check(betaIngress.Status.LoadBalancer.Ingress)
		return
	}

	extensionIngress := &extensionsv1beta1.Ingress{}
	err = scheme.Scheme.Convert(obj, extensionIngress, nil)
	if err == nil {
		check(extensionIngress.Status.LoadBalancer.Ingress)
		return
	}

	status = model.KubernetesResourceState_OTHER
	desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to neither %T, %T, nor %T: %v", obj, v1Ingress, betaIngress, extensionIngress, err)
	return

}

func determineServiceHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	s := &corev1.Service{}
	err := scheme.Scheme.Convert(obj, s, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, s, err)
		return
	}

	status = model.KubernetesResourceState_HEALTHY
	if s.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return
	}
	if len(s.Status.LoadBalancer.Ingress) == 0 {
		status = model.KubernetesResourceState_OTHER
		desc = "Ingress points for the load-balancer are in progress"
		return
	}
	return
}

func determineConfigMapHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = fmt.Sprintf("%q was applied successfully", obj.GetName())
	status = model.KubernetesResourceState_HEALTHY
	return
}

func determineSecretHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = fmt.Sprintf("%q was applied successfully", obj.GetName())
	status = model.KubernetesResourceState_HEALTHY
	return
}

func determinePersistentVolumeHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	pv := &corev1.PersistentVolume{}
	err := scheme.Scheme.Convert(obj, pv, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, pv, err)
		return
	}

	switch pv.Status.Phase {
	case corev1.VolumeBound, corev1.VolumeAvailable:
		status = model.KubernetesResourceState_HEALTHY
		desc = pv.Status.Message
		return
	default:
		status = model.KubernetesResourceState_OTHER
		desc = pv.Status.Message
		return
	}
}

func determinePVCHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := scheme.Scheme.Convert(obj, pvc, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, pvc, err)
		return
	}
	switch pvc.Status.Phase {
	case corev1.ClaimLost:
		status = model.KubernetesResourceState_OTHER
		desc = "Lost its underlying PersistentVolume"
	case corev1.ClaimPending:
		status = model.KubernetesResourceState_OTHER
		desc = "Being not yet bound"
	case corev1.ClaimBound:
		status = model.KubernetesResourceState_HEALTHY
	default:
		status = model.KubernetesResourceState_OTHER
		desc = "The current phase of PersistentVolumeClaim is unexpected"
	}
	return
}

func determineServiceAccountHealth(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	desc = fmt.Sprintf("%q was applied successfully", obj.GetName())
	status = model.KubernetesResourceState_HEALTHY
	return
}

// This function based on [istio project](https://github.com/istio/istio/blob/8ce0defcc905873dadf3fa1d8c3f3629cd39895c/pkg/config/validation/validation.go#L1930-L2041).
// See the file headers for detail information.
func determineVirtualService(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	vs := &istio.VirtualService{}
	err := scheme.Scheme.Convert(obj, vs, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, vs, err)
		return
	}
	if len(vs.Hosts) == 0 {
		desc = "Virtual service must have at least one host"
		return
	}

	if len(vs.Http) == 0 && len(vs.Tcp) == 0 && len(vs.Tls) == 0 {
		desc = "Http, tcp or tls must be provided in virtual service"
		return
	}
	for _, httpRoute := range vs.Http {
		if httpRoute == nil {
			desc = "Http route may not be null"
			continue
		}
	}
	return
}

// This function based on [istio project](https://github.com/istio/istio/blob/8ce0defcc905873dadf3fa1d8c3f3629cd39895c/pkg/config/validation/validation.go#L442-L484).
// See the file headers for detail information.
func determineGateway(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	g := &istio.Gateway{}
	err := scheme.Scheme.Convert(obj, g, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, g, err)
		return
	}

	if len(g.Servers) == 0 {
		status = model.KubernetesResourceState_OTHER
		desc = "Gateway must have at least one server"
		return
	}

	portNames := make(map[string]bool)
	for _, s := range g.Servers {
		if s == nil {
			desc = "Server may not be nil"
			return
		}
		if s.Port != nil {
			if portNames[s.Port.Name] {
				desc = fmt.Sprintf("Port names in servers must be unique: duplicate name: %s", s.Port.Name)
				return
			}
			portNames[s.Port.Name] = true
			switch strings.ToLower(s.Port.Protocol) {
			case "http", "http2", "http-proxy", "grpc", "grpc-web":
				if s.GetTls().GetHttpsRedirect() {
					desc = fmt.Sprintf("Tls.httpsRedirect should only be used with http servers: %T", s.Tls)
					return
				}
			}
		}
	}
	return
}

// This function based on [istio project](https://github.com/istio/istio/blob/8ce0defcc905873dadf3fa1d8c3f3629cd39895c/pkg/config/validation/validation.go#L611-L653).
// See the file headers for detail information.
func determineDestinationRule(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	dr := &istio.DestinationRule{}
	err := scheme.Scheme.Convert(obj, dr, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, dr, err)
		return
	}

	for _, subset := range dr.Subsets {
		if subset == nil {
			desc = "The subset may not be nil"
			return
		}
	}
	return
}

// This function based on [istio project](https://github.com/istio/istio/blob/8ce0defcc905873dadf3fa1d8c3f3629cd39895c/pkg/config/validation/validation.go#L2941-L3102).
// See the file headers for detail information.
func determineServiceEntry(obj *unstructured.Unstructured) (status model.KubernetesResourceState_HealthStatus, desc string) {
	se := &istio.ServiceEntry{}
	err := scheme.Scheme.Convert(obj, se, nil)
	if err != nil {
		status = model.KubernetesResourceState_OTHER
		desc = fmt.Sprintf("Unexpected error while calculating: unable to convert %T to %T: %v", obj, se, err)
		return
	}

	if se.WorkloadSelector != nil && se.Endpoints != nil {
		status = model.KubernetesResourceState_OTHER
		desc = "Only one of WorkloadSelector or Endpoints is allowed in Service Entry"
		return
	}
	if len(se.Hosts) == 0 {
		status = model.KubernetesResourceState_OTHER
		desc = "Service entry must have at least one host"
		return
	}
	for _, hostname := range se.Hosts {
		if hostname == "*" {
			desc = fmt.Sprintf("Invalid host %s", hostname)
			return
		}
	}
	cidrFound := false
	for _, address := range se.Addresses {
		cidrFound = cidrFound || strings.Contains(address, "/")
	}
	if cidrFound {
		if se.Resolution.String() != Resolution_NONE && se.Resolution.String() != Resolution_STATIC {
			desc = "CIDR addresses are allowed only for NONE/STATIC resolution types"
			return
		}
	}
	servicePortNumbers := make(map[uint32]bool)
	servicePorts := make(map[string]bool, len(se.Ports))
	for _, port := range se.Ports {
		if port == nil {
			desc = "Service entry port may not be null"
			return
		}
		if servicePorts[port.Name] {
			desc = fmt.Sprintf("Service entry port name %q already defined", port.Name)
			return
		}
		servicePorts[port.Name] = true
		if servicePortNumbers[port.Number] {
			desc = fmt.Sprintf("Service entry port %d already defined", port.Number)
			return
		}
		servicePortNumbers[port.Number] = true
		if port.TargetPort == 0 {
			desc = "Service entry port may not be 0"
			return
		}
	}
	switch se.Resolution.String() {
	case Resolution_NONE:
		if len(se.Endpoints) != 0 {
			desc = "No endpoints should be provided for resolution type none"
			return
		}
	case Resolution_STATIC:
		unixEndpoint := false
		for _, endpoint := range se.Endpoints {
			addr := endpoint.GetAddress()
			if strings.HasPrefix(addr, "unix://") {
				unixEndpoint = true
				if len(endpoint.Ports) != 0 {
					desc = fmt.Sprintf("Unix endpoint %s must not include ports", addr)
					return
				}
			} else {
				for name, port := range endpoint.Ports {
					if !servicePorts[name] {
						desc = fmt.Sprintf("Endpoint port %v is not defined by the service entry", port)
						return
					}
				}
			}
		}
		if unixEndpoint && len(se.Ports) != 1 {
			desc = "Exactly 1 service port required for unix endpoints"
			return
		}
	case Resolution_DNS, Resolution_DNS_ROUND_ROBIN:
		if len(se.Addresses) > 0 {
			for _, port := range se.Ports {
				switch strings.ToLower(port.Protocol) {
				case "tcp", "https", "tls", "mongo", "redis", "mysql", "thrift":
					if len(se.Hosts) > 1 {
						// TODO: prevent this invalid setting, maybe when istio version is 1.11+
						desc = "Service entry can not have more than one host specified simultaneously with address and tcp port"
					}
					return
				}
			}
		}
	default:
		desc = fmt.Sprintf("Unsupported resolution type %s", se.Resolution.String())
		return
	}
	if se.Resolution.String() != Resolution_NONE && len(se.Hosts) > 1 {
		for _, port := range se.Ports {
			switch strings.ToLower(port.Protocol) {
			case "http", "http2", "http-proxy", "grpc", "grpc-web", "https", "tls":
				desc = "Multiple hosts provided with non-HTTP, non-TLS ports"
				return
			}
		}
	}
	return
}
