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

package provider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	errorReplaceNotFound     = errors.New("specified resource is not found")
	errorNotFoundLiteral     = "Error from server (NotFound)"
	errResourceAlreadyExists = errors.New("resource already exists")
	errAlreadyExistsLiteral  = "Error from server (AlreadyExists)"
)

// Kubectl is a wrapper for kubectl command.
type Kubectl struct {
	execPath string
}

// NewKubectl creates a new Kubectl instance.
func NewKubectl(path string) *Kubectl {
	return &Kubectl{
		execPath: path,
	}
}

// Apply runs kubectl apply command with the given manifest.
func (c *Kubectl) Apply(ctx context.Context, kubeconfig, namespace string, manifest Manifest) (err error) {
	// TODO: record the metrics for the kubectl apply command.

	data, err := manifest.YamlBytes()
	if err != nil {
		return err
	}

	args := make([]string, 0, 8)
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}

	args = append(args, "apply")
	if annotation := manifest.body.GetAnnotations()[LabelServerSideApply]; annotation == UseServerSideApply {
		args = append(args, "--server-side")
	}
	args = append(args, "-f", "-")

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	r := bytes.NewReader(data)
	cmd.Stdin = r

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply: %s (%w)", string(out), err)
	}
	return nil
}

// Create runs kubectl create command with the given manifest.
func (c *Kubectl) Create(ctx context.Context, kubeconfig, namespace string, manifest Manifest) (err error) {
	// TODO: record the metrics for the kubectl create command.

	data, err := manifest.YamlBytes()
	if err != nil {
		return err
	}

	args := make([]string, 0, 7)
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	args = append(args, "create", "-f", "-")

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	r := bytes.NewReader(data)
	cmd.Stdin = r

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create: %s (%w)", string(out), err)
	}
	return nil
}

// Replace runs kubectl replace command with the given manifest.
func (c *Kubectl) Replace(ctx context.Context, kubeconfig, namespace string, manifest Manifest) (err error) {
	// TODO: record the metrics for the kubectl replace command.

	data, err := manifest.YamlBytes()
	if err != nil {
		return err
	}

	args := make([]string, 0, 7)
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	args = append(args, "replace", "-f", "-")

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	r := bytes.NewReader(data)
	cmd.Stdin = r

	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	if strings.Contains(string(out), errorNotFoundLiteral) {
		return errorReplaceNotFound
	}

	return fmt.Errorf("failed to replace: %s (%w)", string(out), err)
}

// ForceReplace runs kubectl replace --force command with the given manifest.
func (c *Kubectl) ForceReplace(ctx context.Context, kubeconfig, namespace string, manifest Manifest) (err error) {
	// TODO: record the metrics for the kubectl replace --force command.

	data, err := manifest.YamlBytes()
	if err != nil {
		return err
	}

	args := make([]string, 0, 7)
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	args = append(args, "replace", "--force", "-f", "-")

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	r := bytes.NewReader(data)
	cmd.Stdin = r

	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	if strings.Contains(string(out), errorNotFoundLiteral) {
		return errorReplaceNotFound
	}

	return fmt.Errorf("failed to replace: %s (%w)", string(out), err)
}

// Delete runs kubectl delete command with the given resource key.
func (c *Kubectl) Delete(ctx context.Context, kubeconfig, namespace string, r ResourceKey) (err error) {
	// TODO: record the metrics for the kubectl delete command.

	args := make([]string, 0, 7)
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	args = append(args, "delete", r.Kind(), r.Name())

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	out, err := cmd.CombinedOutput()

	if strings.Contains(string(out), "(NotFound)") {
		return fmt.Errorf("failed to delete: %s, (%w), %v", string(out), ErrNotFound, err)
	}
	if err != nil {
		return fmt.Errorf("failed to delete: %s, %v", string(out), err)
	}
	return nil
}

// Get runs kubectl get with the gibven resource key.
func (c *Kubectl) Get(ctx context.Context, kubeconfig, namespace string, r ResourceKey) (m Manifest, err error) {
	// TODO: record the metrics for the kubectl get command.

	args := make([]string, 0, 7)
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	args = append(args, "get", r.Kind(), r.Name(), "-o", "yaml")

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	out, err := cmd.CombinedOutput()

	if strings.Contains(string(out), "(NotFound)") {
		return Manifest{}, fmt.Errorf("not found manifest %v, (%w), %v", r, ErrNotFound, err)
	}
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to get: %s, %v", string(out), err)
	}
	ms, err := ParseManifests(string(out))
	if err != nil {
		return Manifest{}, fmt.Errorf("failed to parse manifests %v: %v", r, err)
	}
	if len(ms) == 0 {
		return Manifest{}, fmt.Errorf("not found manifest %v, (%w)", r, ErrNotFound)
	}
	return ms[0], nil
}

// getClusterScopedAPIResources retrieves the list of available API resources from the Kubernetes cluster.
// It runs the `kubectl api-resources` command with the specified kubeconfig and returns the
// names of the resources that support the "list", "get", and "delete" verbs, and are cluster-scoped.
func (c *Kubectl) getClusterScopedAPIResources(ctx context.Context, kubeconfig string) ([]string, error) {
	args := []string{"api-resources", "--namespaced=false", "--verbs=list,get,delete", "--output=name"}
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	cmd := exec.CommandContext(ctx, c.execPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get API resources: %s, %v", string(out), err)
	}
	lines := strings.Split(string(out), "\n")
	resources := make([]string, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			resources = append(resources, line)
		}
	}

	return resources, nil
}

// GetAllClusterScoped retrieves all cluster-scoped resources from the Kubernetes cluster
// using the provided kubeconfig and optional selectors. It returns a slice of Manifests
// representing the resources or an error if the operation fails.
func (c *Kubectl) GetAllClusterScoped(ctx context.Context, kubeconfig string, selector ...string) ([]Manifest, error) {
	resources, err := c.getClusterScopedAPIResources(ctx, kubeconfig)
	if err != nil {
		return nil, err
	}

	args := make([]string, 0, 7)
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	args = append(args, "get", strings.Join(resources, ","), "-o", "yaml", "--selector", strings.Join(selector, ","))
	cmd := exec.CommandContext(ctx, c.execPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster-scoped resources: %s, %v", string(out), err)
	}

	// Unmarshal the output to the list of manifests.
	var list v1.List
	if err := yaml.Unmarshal(out, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the output: %w", err)
	}

	ms := make([]Manifest, 0, len(list.Items))
	for _, item := range list.Items {
		m, err := ParseManifests(string(item.Raw))
		if err != nil {
			return nil, fmt.Errorf("failed to parse the manifest: %w", err)
		}
		ms = append(ms, m...)
	}

	return ms, nil
}

// getNamespaceScopedAPIResources retrieves the list of available API resources from the Kubernetes cluster.
// It runs the `kubectl api-resources` command with the specified kubeconfig and returns the
// names of the resources that support the "list", "get", and "delete" verbs, and are namespace-scoped.
func (c *Kubectl) getNamespaceScopedAPIResources(ctx context.Context, kubeconfig string) ([]string, error) {
	args := []string{"api-resources", "--namespaced=true", "--verbs=list,get,delete", "--output=name"}
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	cmd := exec.CommandContext(ctx, c.execPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get API resources: %s, %v", string(out), err)
	}
	lines := strings.Split(string(out), "\n")
	resources := make([]string, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			resources = append(resources, line)
		}
	}

	return resources, nil
}

// GetAll retrieves all Kubernetes resources in the specified namespace and matching the given selector.
// It returns a list of manifests or an error if the retrieval or unmarshalling fails.
// If no resources are found, it returns nil without an error.
func (c *Kubectl) GetAll(ctx context.Context, kubeconfig, namespace string, selector ...string) (ms []Manifest, err error) {
	resources, err := c.getNamespaceScopedAPIResources(ctx, kubeconfig)
	if err != nil {
		return nil, err
	}

	args := make([]string, 0, 8)
	args = append(args, "get", strings.Join(resources, ","), "-o", "yaml", "--selector", strings.Join(selector, ","))

	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	if namespace == "" {
		args = append(args, "--all-namespaces")
	} else {
		args = append(args, "--namespace", namespace)
	}

	cmd := exec.CommandContext(ctx, c.execPath, args...)

	// use only stdout for output to prevent manifest parsing failures caused by messages on stderr, like warning messages.
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace-scoped resources: stdout: %s, stderr: %s, %v", string(out), stderr.String(), err)
	}

	if strings.Contains(stderr.String(), "(NotFound)") {
		// No resources found. Return nil. This is not an error.
		return nil, nil
	}

	// Unmarshal the output to the list of manifests.
	var list v1.List
	if err := yaml.Unmarshal(out, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the output: %w", err)
	}

	ms = make([]Manifest, 0, len(list.Items))
	for _, item := range list.Items {
		m, err := ParseManifests(string(item.Raw))
		if err != nil {
			return nil, fmt.Errorf("failed to parse the manifest: %w", err)
		}
		ms = append(ms, m...)
	}

	return ms, nil
}

// CreateNamespace runs kubectl create namespace with the given namespace.
func (c *Kubectl) CreateNamespace(ctx context.Context, kubeconfig, namespace string) (err error) {
	// TODO: record the metrics for the kubectl create namespace command.

	args := make([]string, 0, 7)
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	args = append(args, "create", "namespace", namespace)

	cmd := exec.CommandContext(ctx, c.execPath, args...)
	out, err := cmd.CombinedOutput()

	if strings.Contains(string(out), errAlreadyExistsLiteral) {
		return errResourceAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("failed to create namespace: %s, %v", string(out), err)
	}
	return nil
}
