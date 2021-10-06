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
	  "role": "arn:aws:iam::xxxxx:role/lambda-role",
	  "memory": 128,
	  "timeout": 5,
	  "image": "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1"
  }
}`,
			wantSpec: FunctionManifest{
				Kind:       "LambdaFunction",
				APIVersion: "pipecd.dev/v1beta1",
				Spec: FunctionManifestSpec{
					Name:     "SimpleFunction",
					Role:     "arn:aws:iam::xxxxx:role/lambda-role",
					Memory:   128,
					Timeout:  5,
					ImageURI: "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1",
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
		{
			name: "missing memory value",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "LambdaFunction",
  "spec": {
	  "name": "SimpleFunction",
	  "role": "arn:aws:iam::xxxxx:role/lambda-role",
	  "timeout": 5,
	  "image": "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1"
  }
}`,
			wantSpec: FunctionManifest{},
			wantErr:  true,
		},
		{
			name: "invalid timeout value",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "LambdaFunction",
  "spec": {
	  "name": "SimpleFunction",
	  "role": "arn:aws:iam::xxxxx:role/lambda-role",
	  "memory": 128,
	  "timeout": 1000,
	  "image": "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1"
  }
}`,
			wantSpec: FunctionManifest{},
			wantErr:  true,
		},
		{
			name: "no function code defined",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "LambdaFunction",
  "spec": {
	  "name": "SimpleFunction",
	  "role": "arn:aws:iam::xxxxx:role/lambda-role",
	  "memory": 128,
	  "timeout": 5
  }
}`,
			wantSpec: FunctionManifest{},
			wantErr:  true,
		},
		{
			name: "no error in case of multiple function code defined",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "LambdaFunction",
  "spec": {
	  "name": "SimpleFunction",
	  "role": "arn:aws:iam::xxxxx:role/lambda-role",
	  "memory": 128,
	  "timeout": 5,
	  "image": "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1",
	  "source": {
		"git": "git@remote-url",
		"branch": "master",
		"path": "./"
	  }
  }
}`,
			wantSpec: FunctionManifest{
				Kind:       "LambdaFunction",
				APIVersion: "pipecd.dev/v1beta1",
				Spec: FunctionManifestSpec{
					Name:     "SimpleFunction",
					Role:     "arn:aws:iam::xxxxx:role/lambda-role",
					Memory:   128,
					Timeout:  5,
					ImageURI: "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1",
					SourceCode: SourceCode{
						Git:    "git@remote-url",
						Branch: "master",
						Path:   "./",
					},
				},
			},
			wantErr: false,
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
