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

package lambda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestparseFunctionManifest(t *testing.T) {
	testcases := []struct {
		name     string
		data     string
		wantSpec interface{}
		wantErr  bool
	}{
		{
			name: "correct config for LambdaFunction",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "LambdaFunction",
  "spec": {
	  "name": "SimpleFunction",
	  "runtime": "nodejs12.x",
	  "handler": "SampleFunction",
	  "image": "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1"
  }
}`,
			wantSpec: FunctionManifest{
				Kind:       "LambdaFunction",
				APIVersion: "pipecd.dev/v1beta1",
				Spec: FunctionManifestSpec{
					Name:     "SimpleFunction",
					ImageURI: "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1",
					Runtime:  "nodejs12.x",
					Handler:  "SampleFunction",
				},
			},
			wantErr: false,
		},
		{
			name: "missing required fields",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "LambdaFunction",
  "spec": {}
}`,
			wantSpec: FunctionManifest{},
			wantErr:  true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			fm, err := parseFunctionManifest([]byte(tc.data))
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.wantSpec, fm)
		})
	}
}
