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

package lambda

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/lambda"
	"github.com/pipe-cd/pipecd/pkg/diff"
	"github.com/stretchr/testify/assert"
)

func TestIgnoreAndSortParameters(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name       string
		liveSpec   provider.FunctionManifestSpec
		headSpec   provider.FunctionManifestSpec
		expectDiff bool
	}{
		{
			name: "Ignore packaging types S3 and SourceCode",
			liveSpec: provider.FunctionManifestSpec{
				ImageURI: "test-image-uri",
			},
			headSpec: provider.FunctionManifestSpec{
				ImageURI:        "test-image-uri",
				S3Bucket:        "test-bucket",
				S3Key:           "test-key",
				S3ObjectVersion: "test-version",
				SourceCode: provider.SourceCode{
					Git:  "https://test-repo.git",
					Ref:  "test-ref",
					Path: "test-path",
				},
			},
			expectDiff: false,
		},
		{
			name: "Ignore not sorted fields and pipecd managed tags",
			liveSpec: provider.FunctionManifestSpec{
				Architectures: []provider.Architecture{
					{Name: string(types.ArchitectureArm64)},
					{Name: string(types.ArchitectureX8664)},
				},
				Environments: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				Tags: map[string]string{
					"key1":                    "value1",
					"key2":                    "value2",
					provider.LabelApplication: "test-app",
					provider.LabelCommitHash:  "test-hash",
					provider.LabelManagedBy:   "piped",
					provider.LabelPiped:       "test-piped",
				},
				VPCConfig: &provider.VPCConfig{
					SubnetIDs: []string{"subnet-1", "subnet-2"},
				},
			},
			headSpec: provider.FunctionManifestSpec{
				Architectures: []provider.Architecture{
					{Name: string(types.ArchitectureX8664)},
					{Name: string(types.ArchitectureArm64)},
				},
				Environments: map[string]string{
					"key2": "value2",
					"key1": "value1",
				},
				Tags: map[string]string{
					"key2": "value2",
					"key1": "value1",
				},
				VPCConfig: &provider.VPCConfig{
					SubnetIDs: []string{"subnet-2", "subnet-1"},
				},
			},
			expectDiff: false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ignoreAndSortParameters(&tc.liveSpec, &tc.headSpec)
			result, err := provider.Diff(
				provider.FunctionManifest{Spec: tc.liveSpec},
				provider.FunctionManifest{Spec: tc.headSpec},
				diff.WithEquateEmpty(),
				diff.WithIgnoreAddingMapKeys(),
				diff.WithCompareNumberAndNumericString(),
			)
			assert.NoError(t, err)
			fmt.Println(result.Render(provider.DiffRenderOptions{}))
			assert.Equal(t, tc.expectDiff, result.Diff.HasDiff())
		})
	}
}
