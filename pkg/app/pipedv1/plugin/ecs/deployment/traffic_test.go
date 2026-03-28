// Copyright 2026 The PipeCD Authors.
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
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/ecs/provider"
)

// mockMetadataStore implements metadataStore for testing.
type mockMetadataStore struct {
	PutMultiFunc func(ctx context.Context, metadata map[string]string) error
	GetFunc      func(ctx context.Context, key string) (string, bool, error)
	PutFunc      func(ctx context.Context, key string, value string) error
}

func (m *mockMetadataStore) PutDeploymentPluginMetadataMulti(ctx context.Context, metadata map[string]string) error {
	return m.PutMultiFunc(ctx, metadata)
}
func (m *mockMetadataStore) GetDeploymentPluginMetadata(ctx context.Context, key string) (string, bool, error) {
	return m.GetFunc(ctx, key)
}
func (m *mockMetadataStore) PutDeploymentPluginMetadata(ctx context.Context, key string, value string) error {
	return m.PutFunc(ctx, key, value)
}

var _ metadataStore = (*mockMetadataStore)(nil)

func happyMetadataStore() *mockMetadataStore {
	return &mockMetadataStore{
		PutMultiFunc: func(_ context.Context, _ map[string]string) error { return nil },
		GetFunc:      func(_ context.Context, _ string) (string, bool, error) { return "", false, nil },
		PutFunc:      func(_ context.Context, _, _ string) error { return nil },
	}
}

func TestRouting(t *testing.T) {
	t.Parallel()

	var (
		primaryARN   = "arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/primary/aaa"
		canaryARN    = "arn:aws:elasticloadbalancing:us-east-1:123:targetgroup/canary/bbb"
		listenerARN1 = "arn:aws:elasticloadbalancing:us-east-1:123:listener/app/my-alb/aaa/bbb"
		listenerARN2 = "arn:aws:elasticloadbalancing:us-east-1:123:listener/app/my-alb/aaa/ccc"
		primaryTG    = types.LoadBalancer{TargetGroupArn: aws.String(primaryARN)}
		canaryTG     = types.LoadBalancer{TargetGroupArn: aws.String(canaryARN)}
	)

	testcases := []struct {
		name       string
		options    config.ECSTrafficRoutingStageOptions
		metadata   *mockMetadataStore
		client     *mockECSClient
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:     "success: listener ARNs not cached, fetched from AWS",
			options:  config.ECSTrafficRoutingStageOptions{Canary: 20},
			metadata: happyMetadataStore(),
			client: func() *mockECSClient {
				m := &mockECSClient{}
				m.GetListenerArnsFunc = func(_ context.Context, _ types.LoadBalancer) ([]string, error) {
					return []string{listenerARN1}, nil
				}
				m.ModifyListenersFunc = func(_ context.Context, listenerArns []string, cfg provider.RoutingTrafficConfig) ([]string, error) {
					assert.Equal(t, []string{listenerARN1}, listenerArns)
					assert.Equal(t, provider.RoutingTrafficConfig{
						{TargetGroupArn: primaryARN, Weight: 80},
						{TargetGroupArn: canaryARN, Weight: 20},
					}, cfg)
					return []string{"rule-1"}, nil
				}
				return m
			}(),
		},
		{
			name:    "success: listener ARNs cached in metadata, skip GetListenerArns",
			options: config.ECSTrafficRoutingStageOptions{Canary: 20},
			metadata: func() *mockMetadataStore {
				m := happyMetadataStore()
				m.GetFunc = func(_ context.Context, key string) (string, bool, error) {
					if key == currentListenersKey {
						return listenerARN1 + "," + listenerARN2, true, nil
					}
					return "", false, nil
				}
				return m
			}(),
			client: func() *mockECSClient {
				m := &mockECSClient{}
				m.GetListenerArnsFunc = func(_ context.Context, _ types.LoadBalancer) ([]string, error) {
					t.Error("GetListenerArns should not be called when cached")
					return nil, nil
				}
				m.ModifyListenersFunc = func(_ context.Context, listenerArns []string, _ provider.RoutingTrafficConfig) ([]string, error) {
					assert.Equal(t, []string{listenerARN1, listenerARN2}, listenerArns)
					return []string{"rule-1", "rule-2"}, nil
				}
				return m
			}(),
		},
		{
			name:     "success: primary=100 routes all traffic to primary",
			options:  config.ECSTrafficRoutingStageOptions{Primary: 100},
			metadata: happyMetadataStore(),
			client: func() *mockECSClient {
				m := &mockECSClient{}
				m.GetListenerArnsFunc = func(_ context.Context, _ types.LoadBalancer) ([]string, error) {
					return []string{listenerARN1}, nil
				}
				m.ModifyListenersFunc = func(_ context.Context, _ []string, cfg provider.RoutingTrafficConfig) ([]string, error) {
					assert.Equal(t, 100, cfg[0].Weight)
					assert.Equal(t, 0, cfg[1].Weight)
					return []string{"rule-1"}, nil
				}
				return m
			}(),
		},
		{
			name:     "success: no options set defaults to primary=100",
			options:  config.ECSTrafficRoutingStageOptions{},
			metadata: happyMetadataStore(),
			client: func() *mockECSClient {
				m := &mockECSClient{}
				m.GetListenerArnsFunc = func(_ context.Context, _ types.LoadBalancer) ([]string, error) {
					return []string{listenerARN1}, nil
				}
				m.ModifyListenersFunc = func(_ context.Context, _ []string, cfg provider.RoutingTrafficConfig) ([]string, error) {
					assert.Equal(t, 100, cfg[0].Weight)
					assert.Equal(t, 0, cfg[1].Weight)
					return []string{"rule-1"}, nil
				}
				return m
			}(),
		},
		{
			name:    "fail: PutDeploymentPluginMetadataMulti error",
			options: config.ECSTrafficRoutingStageOptions{Canary: 20},
			metadata: func() *mockMetadataStore {
				m := happyMetadataStore()
				m.PutMultiFunc = func(_ context.Context, _ map[string]string) error {
					return errors.New("put multi error")
				}
				return m
			}(),
			client:     &mockECSClient{},
			wantErr:    true,
			wantErrMsg: "Failed to store percentage metadata",
		},
		{
			name:    "fail: GetDeploymentPluginMetadata error",
			options: config.ECSTrafficRoutingStageOptions{Canary: 20},
			metadata: func() *mockMetadataStore {
				m := happyMetadataStore()
				m.GetFunc = func(_ context.Context, _ string) (string, bool, error) {
					return "", false, errors.New("get error")
				}
				return m
			}(),
			client:     &mockECSClient{},
			wantErr:    true,
			wantErrMsg: "Failed to get current listener arns",
		},
		{
			name:     "fail: GetListenerArns error when not cached",
			options:  config.ECSTrafficRoutingStageOptions{Canary: 20},
			metadata: happyMetadataStore(),
			client: func() *mockECSClient {
				m := &mockECSClient{}
				m.GetListenerArnsFunc = func(_ context.Context, _ types.LoadBalancer) ([]string, error) {
					return nil, errors.New("describe listeners error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "Failed to get current active listeners",
		},
		{
			name:    "fail: PutDeploymentPluginMetadata error when saving listeners",
			options: config.ECSTrafficRoutingStageOptions{Canary: 20},
			metadata: func() *mockMetadataStore {
				m := happyMetadataStore()
				m.PutFunc = func(_ context.Context, key, _ string) error {
					if key == currentListenersKey {
						return errors.New("put listeners error")
					}
					return nil
				}
				return m
			}(),
			client: func() *mockECSClient {
				m := &mockECSClient{}
				m.GetListenerArnsFunc = func(_ context.Context, _ types.LoadBalancer) ([]string, error) {
					return []string{listenerARN1}, nil
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "Failed to store listeners to metadata store",
		},
		{
			name:     "fail: ModifyListeners error",
			options:  config.ECSTrafficRoutingStageOptions{Canary: 20},
			metadata: happyMetadataStore(),
			client: func() *mockECSClient {
				m := &mockECSClient{}
				m.GetListenerArnsFunc = func(_ context.Context, _ types.LoadBalancer) ([]string, error) {
					return []string{listenerARN1}, nil
				}
				m.ModifyListenersFunc = func(_ context.Context, _ []string, _ provider.RoutingTrafficConfig) ([]string, error) {
					return nil, errors.New("modify listeners error")
				}
				return m
			}(),
			wantErr:    true,
			wantErrMsg: "Failed to routing traffic",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := routing(context.Background(), &fakeLogPersister{}, tc.metadata, tc.client, primaryTG, canaryTG, tc.options)

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}
