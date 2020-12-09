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

package ecr

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/stretchr/testify/assert"
)

func TestLatestBySemver(t *testing.T) {
	testcases := []struct {
		name    string
		ids     []*ecr.ImageIdentifier
		want    string
		wantErr bool
	}{
		{
			name: "include tags tha don't accordance semver",
			ids: []*ecr.ImageIdentifier{
				{
					ImageTag: aws.String("latest"),
				},
				{
					ImageTag: aws.String("0.1"),
				},
			},
			wantErr: true,
		},
		{
			name: "only major versions",
			ids: []*ecr.ImageIdentifier{
				{
					ImageTag: aws.String("5"),
				},
				{
					ImageTag: aws.String("3"),
				},
				{
					ImageTag: aws.String("4"),
				},
				{
					ImageTag: aws.String("8"),
				},
			},
			want:    "8.0.0",
			wantErr: false,
		},
		{
			name: "various versions",
			ids: []*ecr.ImageIdentifier{
				{
					ImageTag: aws.String("5.0.1"),
				},
				{
					ImageTag: aws.String("v3.0"),
				},
				{
					ImageTag: aws.String("4"),
				},
				{
					ImageTag: aws.String("8.10"),
				},
			},
			want:    "8.10.0",
			wantErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := latestBySemver(tc.ids)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
