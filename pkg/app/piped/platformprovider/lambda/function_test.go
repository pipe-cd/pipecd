// Copyright 2023 The PipeCD Authors.
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

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestParseFunctionManifest(t *testing.T) {
	t.Parallel()

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
			name: "correct config for LambdaFunction with specifying architecture",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "LambdaFunction",
  "spec": {
	  "name": "SimpleFunction",
	  "role": "arn:aws:iam::xxxxx:role/lambda-role",
	  "memory": 128,
	  "timeout": 5,
	  "image": "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1",
      "architectures": [
        {
          "name": "x86_64",
        },
        {
          "name": "arm64",
        }
      ]
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
					Architectures: []Architecture{
						{
							Name: "x86_64",
						},
						{
							Name: "arm64",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "correct config for LambdaFunction with specifying vpc config",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "LambdaFunction",
  "spec": {
	  "name": "SimpleFunction",
	  "role": "arn:aws:iam::xxxxx:role/lambda-role",
	  "memory": 128,
	  "timeout": 5,
	  "image": "ecr.region.amazonaws.com/lambda-simple-function:v0.0.1",
      "vpcConfig": {
          securityGroupIds: ["sg-1234567890", "sg-0987654321"],
          subnetIds: ["subnet-1234567890", "subnet-0987654321"]
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
					VPCConfig: &VPCConfig{
						SecurityGroupIDs: []string{"sg-1234567890", "sg-0987654321"},
						SubnetIDs:        []string{"subnet-1234567890", "subnet-0987654321"},
					},
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
		"ref": "master",
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
						Git:  "git@remote-url",
						Ref:  "master",
						Path: "./",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing required values in case of using other than container image as function",
			data: `{
  "apiVersion": "pipecd.dev/v1beta1",
  "kind": "LambdaFunction",
  "spec": {
	  "name": "SimpleFunction",
	  "role": "arn:aws:iam::xxxxx:role/lambda-role",
	  "memory": 128,
	  "timeout": 10,
	  "s3Bucket": "pipecd-sample",
	  "s3Key": "function-code",
	  "s3ObjectVersion": "xyz"
  }
}`,
			wantSpec: FunctionManifest{},
			wantErr:  true,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fm, err := parseFunctionManifest([]byte(tc.data))
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.wantSpec, fm)
		})
	}
}

func TestFindArtifactVersions(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		input       []byte
		expected    []*model.ArtifactVersion
		expectedErr bool
	}{
		{
			name: "[From container image] ok: using container image",
			input: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleFunction
  image: ecr.ap-northeast-1.amazonaws.com/lambda-test:v0.0.1
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  memory: 512
  timeout: 30
`),
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
					Version: "v0.0.1",
					Name:    "lambda-test",
					Url:     "ecr.ap-northeast-1.amazonaws.com/lambda-test:v0.0.1",
				},
			},
			expectedErr: false,
		},
		{
			name: "[From container image] error: no image name",
			input: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleFunction
  image: ecr.ap-northeast-1.amazonaws.com/:v0.0.1
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  memory: 512
  timeout: 30
`),
			expected:    nil,
			expectedErr: true,
		},
		{
			name: "[From S3] ok: using s3 object",
			input: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleZipPackingS3Function
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  s3Bucket: pipecd-sample-lambda
  s3Key: pipecd-sample-src
  s3ObjectVersion: 1pTK9_v0Kd7I8Sk4n6abzCL
  handler: app.lambdaHandler
  runtime: nodejs14.x
  memory: 512
  timeout: 30
`),
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_S3_OBJECT,
					Version: "1pTK9_v0Kd7I8Sk4n6abzCL",
					Name:    "pipecd-sample-src",
					Url:     "https://console.aws.amazon.com/s3/object/pipecd-sample-lambda?prefix=pipecd-sample-src",
				},
			},
			expectedErr: false,
		},
		{
			name: "[From Source Code] ok: using source code",
			input: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleSourceCodeFunction
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  source:
    git: git@github.com:username/lambda-function-code.git
    ref: dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603
    path: hello-world
  handler: app.lambdaHandler
  runtime: nodejs14.x
  memory: 512
  timeout: 30
`),
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_GIT_SOURCE,
					Version: "dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603",
					Name:    "username/lambda-function-code",
					Url:     "https://github.com/username/lambda-function-code/commit/dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603",
				},
			},
			expectedErr: false,
		},
		{
			name: "[From Source Code] ok: using source code from gitlab",
			input: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleSourceCodeFunction
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  source:
    git: git@gitlab.com:username/lambda-function-code.git
    ref: dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603
    path: hello-world
  handler: app.lambdaHandler
  runtime: nodejs14.x
  memory: 512
  timeout: 30
`),
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_GIT_SOURCE,
					Version: "dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603",
					Name:    "username/lambda-function-code",
					Url:     "https://gitlab.com/username/lambda-function-code/commit/dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603",
				},
			},
			expectedErr: false,
		},
		{
			name: "[From Source Code] ok: using source code from bitbucket",
			input: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleSourceCodeFunction
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  source:
    git: git@bitbucket.org:username/lambda-function-code.git
    ref: dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603
    path: hello-world
  handler: app.lambdaHandler
  runtime: nodejs14.x
  memory: 512
  timeout: 30
`),
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_GIT_SOURCE,
					Version: "dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603",
					Name:    "username/lambda-function-code",
					Url:     "https://bitbucket.org/username/lambda-function-code/commits/dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603",
				},
			},
			expectedErr: false,
		},
		{
			name: "[From Source Code] ok: using source code from other git provider",
			input: []byte(`
apiVersion: pipecd.dev/v1beta1
kind: LambdaFunction
spec:
  name: SimpleSourceCodeFunction
  role: arn:aws:iam::76xxxxxxx:role/lambda-role
  source:
    git: git@ghe.github.com:username/lambda-function-code.git
    ref: dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603
    path: hello-world
  handler: app.lambdaHandler
  runtime: nodejs14.x
  memory: 512
  timeout: 30
`),
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_GIT_SOURCE,
					Version: "dede7cdea5bbd3fdbcc4674bfcd2b2f9e0579603",
					Name:    "",
					Url:     "",
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fm, _ := parseFunctionManifest(tc.input)
			versions, err := FindArtifactVersions(fm)

			assert.Equal(t, tc.expectedErr, err != nil)
			assert.ElementsMatch(t, tc.expected, versions)
		})
	}
}
