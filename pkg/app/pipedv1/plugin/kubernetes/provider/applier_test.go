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
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/config"
)

type mockKubectl struct {
	ApplyFunc           func(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error
	CreateFunc          func(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error
	ReplaceFunc         func(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error
	ForceReplaceFunc    func(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error
	DeleteFunc          func(ctx context.Context, kubeconfig, namespace string, key ResourceKey) error
	GetFunc             func(ctx context.Context, kubeconfig, namespace string, key ResourceKey) (Manifest, error)
	CreateNamespaceFunc func(ctx context.Context, kubeconfig, namespace string) error
}

var (
	errUnexpectedCall = errors.New("unexpected call")
)

func (m *mockKubectl) Apply(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error {
	if m.ApplyFunc != nil {
		return m.ApplyFunc(ctx, kubeconfig, namespace, manifest)
	}
	return errUnexpectedCall
}

func (m *mockKubectl) Create(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, kubeconfig, namespace, manifest)
	}
	return errUnexpectedCall
}

func (m *mockKubectl) Replace(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error {
	if m.ReplaceFunc != nil {
		return m.ReplaceFunc(ctx, kubeconfig, namespace, manifest)
	}
	return errUnexpectedCall
}

func (m *mockKubectl) ForceReplace(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error {
	if m.ForceReplaceFunc != nil {
		return m.ForceReplaceFunc(ctx, kubeconfig, namespace, manifest)
	}
	return errUnexpectedCall
}

func (m *mockKubectl) Delete(ctx context.Context, kubeconfig, namespace string, key ResourceKey) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, kubeconfig, namespace, key)
	}
	return errUnexpectedCall
}

func (m *mockKubectl) Get(ctx context.Context, kubeconfig, namespace string, key ResourceKey) (Manifest, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, kubeconfig, namespace, key)
	}
	return Manifest{}, errUnexpectedCall
}

func (m *mockKubectl) CreateNamespace(ctx context.Context, kubeconfig, namespace string) error {
	if m.CreateNamespaceFunc != nil {
		return m.CreateNamespaceFunc(ctx, kubeconfig, namespace)
	}
	return errUnexpectedCall
}

func TestApplier_ApplyManifest(t *testing.T) {
	t.Parallel()

	var (
		errNamespaceCreation = errors.New("namespace creation error")
		errApply             = errors.New("apply error")
	)

	testCases := []struct {
		name                string
		autoCreateNamespace bool
		createNamespaceErr  error
		applyErr            error
		expectedErr         error
	}{
		{
			name:                "successful apply without namespace creation",
			autoCreateNamespace: false,
			expectedErr:         nil,
		},
		{
			name:                "successful apply with namespace creation",
			autoCreateNamespace: true,
			expectedErr:         nil,
		},
		{
			name:                "namespace creation error",
			autoCreateNamespace: true,
			createNamespaceErr:  errNamespaceCreation,
			expectedErr:         errNamespaceCreation,
		},
		{
			name:                "apply error",
			autoCreateNamespace: false,
			applyErr:            errApply,
			expectedErr:         errApply,
		},
		{
			name:                "successful apply with existing namespace",
			autoCreateNamespace: true,
			createNamespaceErr:  errResourceAlreadyExists,
			expectedErr:         nil,
		},
		{
			name:                "apply error after successful namespace creation",
			autoCreateNamespace: true,
			applyErr:            errApply,
			expectedErr:         errApply,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockKubectl := &mockKubectl{
				CreateNamespaceFunc: func(ctx context.Context, kubeconfig, namespace string) error {
					return tc.createNamespaceErr
				},
				ApplyFunc: func(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error {
					return tc.applyErr
				},
			}

			applier := NewApplier(
				mockKubectl,
				config.KubernetesDeploymentInput{
					AutoCreateNamespace: tc.autoCreateNamespace,
				},
				config.KubernetesDeployTargetConfig{
					KubeConfigPath: "test-kubeconfig",
				},
				zap.NewNop(),
			)

			manifest := Manifest{
				Key: ResourceKey{
					Namespace: "test-namespace",
				},
			}

			err := applier.ApplyManifest(context.Background(), manifest)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestApplier_CreateManifest(t *testing.T) {
	t.Parallel()

	var (
		errNamespaceCreation = errors.New("namespace creation error")
		errCreate            = errors.New("create error")
	)

	testCases := []struct {
		name                string
		autoCreateNamespace bool
		createNamespaceErr  error
		createErr           error
		expectedErr         error
	}{
		{
			name:                "successful create without namespace creation",
			autoCreateNamespace: false,
			expectedErr:         nil,
		},
		{
			name:                "successful create with namespace creation",
			autoCreateNamespace: true,
			expectedErr:         nil,
		},
		{
			name:                "namespace creation error",
			autoCreateNamespace: true,
			createNamespaceErr:  errNamespaceCreation,
			expectedErr:         errNamespaceCreation,
		},
		{
			name:                "create error",
			autoCreateNamespace: false,
			createErr:           errCreate,
			expectedErr:         errCreate,
		},
		{
			name:                "successful create with existing namespace",
			autoCreateNamespace: true,
			createNamespaceErr:  errResourceAlreadyExists,
			expectedErr:         nil,
		},
		{
			name:                "create error after successful namespace creation",
			autoCreateNamespace: true,
			createErr:           errCreate,
			expectedErr:         errCreate,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockKubectl := &mockKubectl{
				CreateNamespaceFunc: func(ctx context.Context, kubeconfig, namespace string) error {
					return tc.createNamespaceErr
				},
				CreateFunc: func(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error {
					return tc.createErr
				},
			}

			applier := NewApplier(
				mockKubectl,
				config.KubernetesDeploymentInput{
					AutoCreateNamespace: tc.autoCreateNamespace,
				},
				config.KubernetesDeployTargetConfig{
					KubeConfigPath: "test-kubeconfig",
				},
				zap.NewNop(),
			)

			manifest := Manifest{
				Key: ResourceKey{
					Namespace: "test-namespace",
				},
			}

			err := applier.CreateManifest(context.Background(), manifest)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestApplier_ReplaceManifest(t *testing.T) {
	t.Parallel()

	var (
		errReplace = errors.New("replace error")
	)

	testCases := []struct {
		name        string
		replaceErr  error
		expectedErr error
	}{
		{
			name:        "successful replace",
			expectedErr: nil,
		},
		{
			name:        "replace error",
			replaceErr:  errReplace,
			expectedErr: errReplace,
		},
		{
			name:        "replace not found error",
			replaceErr:  errorReplaceNotFound,
			expectedErr: ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockKubectl := &mockKubectl{
				ReplaceFunc: func(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error {
					return tc.replaceErr
				},
			}

			applier := NewApplier(
				mockKubectl,
				config.KubernetesDeploymentInput{},
				config.KubernetesDeployTargetConfig{
					KubeConfigPath: "test-kubeconfig",
				},
				zap.NewNop(),
			)

			manifest := Manifest{
				Key: ResourceKey{
					Namespace: "test-namespace",
				},
			}

			err := applier.ReplaceManifest(context.Background(), manifest)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestApplier_ForceReplaceManifest(t *testing.T) {
	t.Parallel()

	var (
		errForceReplace = errors.New("force replace error")
	)

	testCases := []struct {
		name            string
		forceReplaceErr error
		expectedErr     error
	}{
		{
			name:        "successful force replace",
			expectedErr: nil,
		},
		{
			name:            "force replace error",
			forceReplaceErr: errForceReplace,
			expectedErr:     errForceReplace,
		},
		{
			name:            "force replace not found error",
			forceReplaceErr: errorReplaceNotFound,
			expectedErr:     ErrNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockKubectl := &mockKubectl{
				ForceReplaceFunc: func(ctx context.Context, kubeconfig, namespace string, manifest Manifest) error {
					return tc.forceReplaceErr
				},
			}

			applier := NewApplier(
				mockKubectl,
				config.KubernetesDeploymentInput{},
				config.KubernetesDeployTargetConfig{
					KubeConfigPath: "test-kubeconfig",
				},
				zap.NewNop(),
			)

			manifest := Manifest{
				Key: ResourceKey{
					Namespace: "test-namespace",
				},
			}

			err := applier.ForceReplaceManifest(context.Background(), manifest)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestApplier_Delete(t *testing.T) {
	t.Parallel()

	var (
		errGet    = errors.New("get error")
		errDelete = errors.New("delete error")
	)

	testCases := []struct {
		name        string
		getErr      error
		deleteErr   error
		manifest    string
		resourceKey ResourceKey
		expectedErr error
	}{
		{
			name: "successful delete",
			manifest: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/resource-key: "v1:ConfigMap::test-config"
`,
			resourceKey: ResourceKey{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Namespace:  "",
				Name:       "test-config",
			},
			expectedErr: nil,
		},
		{
			name:   "get error",
			getErr: errGet,
			manifest: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/resource-key: "v1:ConfigMap::test-config"
`,
			resourceKey: ResourceKey{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Namespace:  "",
				Name:       "test-config",
			},
			expectedErr: errGet,
		},
		{
			name:      "delete error",
			deleteErr: errDelete,
			manifest: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/resource-key: "v1:ConfigMap::test-config"
`,
			resourceKey: ResourceKey{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Namespace:  "",
				Name:       "test-config",
			},
			expectedErr: errDelete,
		},
		{
			name: "resource key mismatch",
			manifest: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/resource-key: "v1:ConfigMap::test-config"
`,
			resourceKey: ResourceKey{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Namespace:  "",
				Name:       "another-config",
			},
			expectedErr: ErrNotFound,
		},
		{
			name: "successful delete with namespace",
			manifest: `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  namespace: test-namespace
  annotations:
    pipecd.dev/resource-key: "v1:ConfigMap:test-namespace:test-config"
`,
			resourceKey: ResourceKey{
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Namespace:  "test-namespace",
				Name:       "test-config",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			manifests := mustParseManifests(t, tc.manifest)
			manifest := manifests[0]

			mockKubectl := &mockKubectl{
				GetFunc: func(ctx context.Context, kubeconfig, namespace string, key ResourceKey) (Manifest, error) {
					if tc.getErr != nil {
						return Manifest{}, tc.getErr
					}
					return manifest, nil
				},
				DeleteFunc: func(ctx context.Context, kubeconfig, namespace string, key ResourceKey) error {
					return tc.deleteErr
				},
			}

			applier := NewApplier(
				mockKubectl,
				config.KubernetesDeploymentInput{},
				config.KubernetesDeployTargetConfig{
					KubeConfigPath: "test-kubeconfig",
				},
				zap.NewNop(),
			)

			err := applier.Delete(context.Background(), tc.resourceKey)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestApplier_getNamespaceToRun(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		inputNamespace string
		resourceKey    ResourceKey
		expected       string
	}{
		{
			name:           "input namespace is used",
			inputNamespace: "input-namespace",
			resourceKey: ResourceKey{
				Namespace: "resource-namespace",
			},
			expected: "input-namespace",
		},
		{
			name:           "resource key namespace is used",
			inputNamespace: "",
			resourceKey: ResourceKey{
				Namespace: "resource-namespace",
			},
			expected: "resource-namespace",
		},
		{
			name:           "both namespaces are empty",
			inputNamespace: "",
			resourceKey: ResourceKey{
				Namespace: "",
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			applier := &Applier{
				input: config.KubernetesDeploymentInput{
					Namespace: tc.inputNamespace,
				},
			}

			result := applier.getNamespaceToRun(tc.resourceKey)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
