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

package deployment

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/kubernetes/provider"
)

type mockStageLogPersister struct {
	logs      []string
	completed bool
}

func (m *mockStageLogPersister) Write(log []byte) (int, error) {
	m.logs = append(m.logs, string(log))
	return len(log), nil
}

func (m *mockStageLogPersister) Info(log string) {
	m.logs = append(m.logs, log)
}

func (m *mockStageLogPersister) Infof(format string, a ...interface{}) {
	m.logs = append(m.logs, fmt.Sprintf(format, a...))
}

func (m *mockStageLogPersister) Success(log string) {
	m.logs = append(m.logs, log)
}

func (m *mockStageLogPersister) Successf(format string, a ...interface{}) {
	m.logs = append(m.logs, fmt.Sprintf(format, a...))
}

func (m *mockStageLogPersister) Error(log string) {
	m.logs = append(m.logs, log)
}

func (m *mockStageLogPersister) Errorf(format string, a ...interface{}) {
	m.logs = append(m.logs, fmt.Sprintf(format, a...))
}

func (m *mockStageLogPersister) Complete(timeout time.Duration) error {
	m.completed = true
	return nil
}

type mockApplier struct {
	applyErr        error
	forceReplaceErr error
	replaceErr      error
	createErr       error
}

func (m *mockApplier) ApplyManifest(ctx context.Context, manifest provider.Manifest) error {
	return m.applyErr
}

func (m *mockApplier) ForceReplaceManifest(ctx context.Context, manifest provider.Manifest) error {
	return m.forceReplaceErr
}

func (m *mockApplier) ReplaceManifest(ctx context.Context, manifest provider.Manifest) error {
	return m.replaceErr
}

func (m *mockApplier) CreateManifest(ctx context.Context, manifest provider.Manifest) error {
	return m.createErr
}

func TestApplyManifests(t *testing.T) {
	tests := []struct {
		name      string
		manifests []provider.Manifest
		namespace string
		applier   *mockApplier
		wantErr   bool
	}{
		{
			name: "apply manifests successfully",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				applyErr: nil,
			},
			wantErr: false,
		},
		{
			name: "force replace manifests successfully",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/force-sync-by-replace: "enabled"
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				forceReplaceErr: nil,
			},
			wantErr: false,
		},
		{
			name: "replace manifests successfully",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/sync-by-replace: "enabled"
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				replaceErr: nil,
			},
			wantErr: false,
		},
		{
			name: "apply manifests with error",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				applyErr: errors.New("apply error"),
			},
			wantErr: true,
		},
		{
			name: "force replace manifests with error",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/force-sync-by-replace: "enabled"
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				forceReplaceErr: errors.New("force replace error"),
			},
			wantErr: true,
		},
		{
			name: "replace manifests with error",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/sync-by-replace: "enabled"
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				replaceErr: errors.New("replace error"),
			},
			wantErr: true,
		},
		{
			name: "create manifests successfully after force replace not found",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/force-sync-by-replace: "enabled"
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				forceReplaceErr: provider.ErrNotFound,
				createErr:       nil,
			},
			wantErr: false,
		},
		{
			name: "create manifests successfully after replace not found",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/sync-by-replace: "enabled"
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				replaceErr: provider.ErrNotFound,
				createErr:  nil,
			},
			wantErr: false,
		},
		{
			name: "create manifests with error after force replace not found",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/force-sync-by-replace: "enabled"
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				forceReplaceErr: provider.ErrNotFound,
				createErr:       errors.New("create error"),
			},
			wantErr: true,
		},
		{
			name: "create manifests with error after replace not found",
			manifests: mustParseManifests(t, `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config
  annotations:
    pipecd.dev/sync-by-replace: "enabled"
data:
  key: value
`),
			namespace: "",
			applier: &mockApplier{
				replaceErr: provider.ErrNotFound,
				createErr:  errors.New("create error"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := new(mockStageLogPersister)
			err := applyManifests(context.Background(), tt.applier, tt.manifests, tt.namespace, lp)
			if (err != nil) != tt.wantErr {
				t.Errorf("applyManifests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
