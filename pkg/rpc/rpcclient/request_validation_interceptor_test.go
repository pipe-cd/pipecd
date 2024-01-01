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

package rpcclient

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type testValidator bool

func (v testValidator) Validate() error {
	if bool(v) {
		return nil
	}
	return errors.New("validation was not passed")
}

func TestRequestValidationUnaryClientInterceptor(t *testing.T) {
	in := RequestValidationUnaryClientInterceptor()

	testcases := []struct {
		name  string
		req   interface{}
		fails bool
	}{
		{
			name:  "passes because the request does not implement RequestValidator interface",
			req:   "string request",
			fails: false,
		},
		{
			name:  "fails because the validation was not passed",
			req:   testValidator(false),
			fails: true,
		},
		{
			name:  "passes",
			req:   testValidator(true),
			fails: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := in(
				context.TODO(),
				"method",
				tc.req,
				nil,
				nil,
				func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
					return nil
				},
			)
			assert.Equal(t, tc.fails, err != nil)
		})
	}
}
