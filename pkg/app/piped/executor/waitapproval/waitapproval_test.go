// Copyright 2021 The PipeCD Authors.
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

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/app/piped/executor"
	"github.com/pipe-cd/pipe/pkg/app/piped/metadatastore"
	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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


func TestValidateApproverNum(t *testing.T) {
	ac := &fakeAPIClient{
		shared: make(map[string]string, 0),
		stages: make(map[string]metadata, 0),
	}
	testcases := []struct {
		name          string
		approver      string
		executor      *Executor
		wantBool      bool
		wantApprovers string
	}{
		{
			name:     "return the person who just approved it and true",
			approver: "user-1",
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
									minApproverNum: "0",
								},
							},
						},
					}),
				},
			},
			wantBool:      true,
			wantApprovers: "user-1",
		},
		{
			name:     "return the person who just approved it and false",
			approver: "user-1",
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
									minApproverNum: "2",
								},
							},
						},
					}),
				},
			},
			wantBool:      false,
			wantApprovers: "user-1",
		},
		{
			name:     "return nothing and false",
			approver: "user-1",
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
									minApproverNum: "2",
									approversKey:   "user-1",
								},
							},
						},
					}),
				},
			},
			wantBool:      false,
			wantApprovers: "",
		},
		{
			name:     "return all approvers and false",
			approver: "user-2",
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
									minApproverNum: "3",
									approversKey:   "user-1",
								},
							},
						},
					}),
				},
			},
			wantBool:      false,
			wantApprovers: "user-1,user-2",
		},
		{
			name:     "return all approvers and true",
			approver: "user-2",
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
									minApproverNum: "2",
									approversKey:   "user-1",
								},
							},
						},
					}),
				},
			},
			wantBool:      true,
			wantApprovers: "user-1,user-2",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			as, ok := tc.executor.validateApproverNum(tc.approver)
			assert.Equal(t, tc.wantBool, ok)
			assert.Equal(t, tc.wantApprovers, as)
		})
	}
}
