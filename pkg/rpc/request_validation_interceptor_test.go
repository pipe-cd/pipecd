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

package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testValidator bool

func (v testValidator) Validate() error {
	if bool(v) {
		return nil
	}
	return errors.New("validation was not passed")
}

func TestRequestValidationUnaryServerInterceptor(t *testing.T) {
	in := RequestValidationUnaryServerInterceptor()

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
			_, err := in(context.TODO(), tc.req, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			})
			assert.Equal(t, tc.fails, err != nil)
		})
	}
}
