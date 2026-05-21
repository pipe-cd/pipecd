// Copyright 2025 The PipeCD Authors.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pipe-cd/pipecd/pkg/app/piped/executor"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes"
	"github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes/kubernetestest"
)

type fakeAppLiveResourceLister struct {
	resources []provider.Manifest
	ok        bool
}

func (l *fakeAppLiveResourceLister) ListKubernetesResources() ([]provider.Manifest, bool) {
	return l.resources, l.ok
}

func TestRollbackExecutor_pruneResourcesNotInRunningCommit(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	runningManifests := []provider.Manifest{
		provider.MakeManifest(provider.ResourceKey{
			APIVersion: "apps/v1",
			Kind:       provider.KindDeployment,
			Name:       "simple",
		}, &unstructured.Unstructured{}),
	}
	orphanServiceKey := provider.ResourceKey{
		APIVersion: "v1",
		Kind:       provider.KindService,
		Name:       "simple",
	}

	testcases := []struct {
		name    string
		lister  executor.AppLiveResourceLister
		applier provider.Applier
		wantErr bool
	}{
		{
			name: "skip when live resource lister has no data",
			lister: &fakeAppLiveResourceLister{
				ok: false,
			},
			applier: kubernetestest.NewMockApplier(ctrl),
		},
		{
			name: "no resource to remove",
			lister: &fakeAppLiveResourceLister{
				ok: true,
				resources: []provider.Manifest{
					provider.MakeManifest(runningManifests[0].Key, &unstructured.Unstructured{}),
				},
			},
			applier: kubernetestest.NewMockApplier(ctrl),
		},
		{
			name: "remove resources that are not defined in running commit",
			lister: &fakeAppLiveResourceLister{
				ok: true,
				resources: []provider.Manifest{
					provider.MakeManifest(runningManifests[0].Key, &unstructured.Unstructured{}),
					provider.MakeManifest(orphanServiceKey, &unstructured.Unstructured{}),
				},
			},
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().Delete(gomock.Any(), orphanServiceKey).Return(nil)
				return p
			}(),
		},
		{
			name: "return error when deletion fails",
			lister: &fakeAppLiveResourceLister{
				ok: true,
				resources: []provider.Manifest{
					provider.MakeManifest(orphanServiceKey, &unstructured.Unstructured{}),
				},
			},
			applier: func() provider.Applier {
				p := kubernetestest.NewMockApplier(ctrl)
				p.EXPECT().Delete(gomock.Any(), orphanServiceKey).Return(fmt.Errorf("unexpected error"))
				return p
			}(),
			wantErr: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			e := &rollbackExecutor{
				Input: executor.Input{
					LogPersister:          &fakeLogPersister{},
					AppLiveResourceLister: tc.lister,
				},
			}
			err := e.pruneResourcesNotInRunningCommit(
				context.Background(),
				&applierGroup{defaultApplier: tc.applier},
				runningManifests,
			)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
