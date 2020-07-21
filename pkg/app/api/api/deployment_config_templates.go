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

package api

const (
	k8sCanaryDeploymentConfigTemplate = `
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  commitMatcher:
    sync: "^Revert"
  pipeline:
    stages:
      # Deploy the workloads of CANARY variant. In this case, the number of
      # workload replicas of CANARY variant is 10% of the replicas number of PRIMARY variant.
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 10%
      # Wait 10 seconds before going to the next stage.
      - name: WAIT
        with:
          duration: 10s
      # Update the workload of PRIMARY variant to the new version.
      - name: K8S_PRIMARY_ROLLOUT
      # Destroy all workloads of CANARY variant.
      - name: K8S_CANARY_CLEAN
`

	k8sBluegreenDeploymentConfigTemplate = `
apiVersion: pipecd.dev/v1beta1
kind: KubernetesApp
spec:
  pipeline:
    stages:
      # Deploy the workloads of CANARY variant. In this case, the number of
      # workload replicas of CANARY variant is the same with PRIMARY variant.
      - name: K8S_CANARY_ROLLOUT
        with:
          replicas: 100%
      # The percentage of traffic each variant should receive.
      # In this case, CANARY variant will receive all of the traffic.
      - name: K8S_TRAFFIC_ROUTING
        with:
          all: canary
      - name: WAIT_APPROVAL
      # Update the workload of PRIMARY variant to the new version.
      - name: K8S_PRIMARY_ROLLOUT
      # The percentage of traffic each variant should receive.
      # In this case, PRIMARY variant will receive all of the traffic.
      - name: K8S_TRAFFIC_ROUTING
        with:
          primary: 100
      # Destroy all workloads of CANARY variant.
      - name: K8S_CANARY_CLEAN
  # This example is not using service mesh.
  trafficSplit:
    method: pod
`
)
