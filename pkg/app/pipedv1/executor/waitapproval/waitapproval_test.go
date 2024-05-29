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

package waitapproval

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/executor"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/metadatastore"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeLogPersister struct{}

func (l *fakeLogPersister) Write(_ []byte) (int, error)         { return 0, nil }
func (l *fakeLogPersister) Info(_ string)                       {}
func (l *fakeLogPersister) Infof(_ string, _ ...interface{})    {}
func (l *fakeLogPersister) Success(_ string)                    {}
func (l *fakeLogPersister) Successf(_ string, _ ...interface{}) {}
func (l *fakeLogPersister) Error(_ string)                      {}
func (l *fakeLogPersister) Errorf(_ string, _ ...interface{})   {}

type metadata map[string]string

type fakeAPIClient struct {
	shared metadata
	stages map[string]metadata
}

func (c *fakeAPIClient) SaveDeploymentMetadata(_ context.Context, req *pipedservice.SaveDeploymentMetadataRequest, _ ...grpc.CallOption) (*pipedservice.SaveDeploymentMetadataResponse, error) {
	md := make(map[string]string, len(c.shared)+len(req.Metadata))
	for k, v := range c.shared {
		md[k] = v
	}
	for k, v := range req.Metadata {
		md[k] = v
	}
	c.shared = md
	return &pipedservice.SaveDeploymentMetadataResponse{}, nil
}

func (c *fakeAPIClient) SaveStageMetadata(_ context.Context, req *pipedservice.SaveStageMetadataRequest, _ ...grpc.CallOption) (*pipedservice.SaveStageMetadataResponse, error) {
	ori := c.stages[req.StageId]
	md := make(map[string]string, len(ori)+len(req.Metadata))
	for k, v := range ori {
		md[k] = v
	}
	for k, v := range req.Metadata {
		md[k] = v
	}
	c.stages[req.StageId] = md
	return &pipedservice.SaveStageMetadataResponse{}, nil
}

type fakeNotifier struct{}

func (n *fakeNotifier) Notify(_ model.NotificationEvent) {}

func TestValidateApproverNum(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ac := &fakeAPIClient{
		shared: make(map[string]string, 0),
		stages: make(map[string]metadata, 0),
	}
	testcases := []struct {
		name           string
		approver       string
		minApproverNum int
		executor       *Executor
		want           bool
	}{
		{
			name:           "return the person who just approved",
			approver:       "user-1",
			minApproverNum: 0,
			executor: &Executor{
				Input: executor.Input{
					Stage: &model.PipelineStage{
						Id: "stage-1",
					},
					LogPersister: &fakeLogPersister{},
					MetadataStore: metadatastore.NewMetadataStore(ac, &model.Deployment{
						Stages: []*model.PipelineStage{
							{
								Id:       "stage-1",
								Metadata: map[string]string{},
							},
						},
					}),
					Notifier: &fakeNotifier{},
				},
			},
			want: true,
		},
		{
			name:           "return an empty string because number of current approver is not enough",
			approver:       "user-1",
			minApproverNum: 2,
			executor: &Executor{
				Input: executor.Input{
					Stage: &model.PipelineStage{
						Id: "stage-1",
					},
					LogPersister: &fakeLogPersister{},
					MetadataStore: metadatastore.NewMetadataStore(ac, &model.Deployment{
						Stages: []*model.PipelineStage{
							{
								Id:       "stage-1",
								Metadata: map[string]string{},
							},
						},
					}),
				},
			},
			want: false,
		},
		{
			name:           "return an empty string because current approver is same as an approver in metadata",
			approver:       "user-1",
			minApproverNum: 2,
			executor: &Executor{
				Input: executor.Input{
					Stage: &model.PipelineStage{
						Id: "stage-1",
					},
					LogPersister: &fakeLogPersister{},
					MetadataStore: metadatastore.NewMetadataStore(ac, &model.Deployment{
						Stages: []*model.PipelineStage{
							{
								Id: "stage-1",
								Metadata: map[string]string{
									approvedByKey: "user-1",
								},
							},
						},
					}),
					Notifier: &fakeNotifier{},
				},
			},
			want: false,
		},
		{
			name:           "return an empty string because number of current approver and approvers in metadata is not enough",
			approver:       "user-2",
			minApproverNum: 3,
			executor: &Executor{
				Input: executor.Input{
					Stage: &model.PipelineStage{
						Id: "stage-1",
					},
					LogPersister: &fakeLogPersister{},
					MetadataStore: metadatastore.NewMetadataStore(ac, &model.Deployment{
						Stages: []*model.PipelineStage{
							{
								Id: "stage-1",
								Metadata: map[string]string{
									approvedByKey: "user-1",
								},
							},
						},
					}),
					Notifier: &fakeNotifier{},
				},
			},
			want: false,
		},
		{
			name:           "return all approvers",
			approver:       "user-2",
			minApproverNum: 2,
			executor: &Executor{
				Input: executor.Input{
					Stage: &model.PipelineStage{
						Id: "stage-1",
					},
					LogPersister: &fakeLogPersister{},
					MetadataStore: metadatastore.NewMetadataStore(ac, &model.Deployment{
						Stages: []*model.PipelineStage{
							{
								Id: "stage-1",
								Metadata: map[string]string{
									approvedByKey: "user-1",
								},
							},
						},
					}),
					Notifier: &fakeNotifier{},
				},
			},
			want: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.executor.validateApproverNum(ctx, tc.approver, tc.minApproverNum)
			assert.Equal(t, tc.want, got)
		})
	}
}
