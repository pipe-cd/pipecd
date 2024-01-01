#!/usr/bin/env bash

# Copyright 2024 The PipeCD Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

CLUSTER=$1
REG_NAME='kind-registry'
REG_PORT='5001'

# Create registry container unless it already exists
echo "Creating local registry container..."
running="$(docker inspect -f '{{.State.Running}}' "${REG_NAME}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -e REGISTRY_HTTP_ADDR=0.0.0.0:5001 \
    -d --restart=always -p "127.0.0.1:${REG_PORT}:5001" --name "${REG_NAME}" \
    registry:2
fi

# Create a cluster with the local registry enabled in containerd
REG_CONFIG_DIR="/etc/containerd/certs.d"
cat <<EOF | kind create cluster --name ${CLUSTER} --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry]
    config_path = "${REG_CONFIG_DIR}"
EOF

# Connect the registry to the cluster network
# (the network may already be connected)
docker network connect "kind" "${REG_NAME}" || true

# Create containerd config files in cluster
NODE=$(kind get nodes --name ${CLUSTER})
docker exec "$NODE" /bin/bash -c "
set -o errexit
set -o nounset
set -o pipefail

mkdir -p  \"${REG_CONFIG_DIR}/localhost:${REG_PORT}\"
cat <<EOF >> \"${REG_CONFIG_DIR}/localhost:${REG_PORT}/hosts.toml\"
server = \"https://localhost:${REG_PORT}\"

[host.\"http://${REG_NAME}:${REG_PORT}\"]
  capabilities = [\"pull\", \"resolve\", \"push\"]
  skip_verify = true
  plain-http = true
EOF
systemctl restart containerd
"

# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${REG_PORT}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF
