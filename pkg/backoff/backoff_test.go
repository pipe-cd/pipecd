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

package backoff

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitNext(t *testing.T) {
	var (
		bo          = NewConstant(time.Millisecond)
		r           = NewRetry(10, bo)
		ctx, cancel = context.WithCancel(context.TODO())
	)
	ok := r.WaitNext(ctx)
	assert.Equal(t, true, ok)

	cancel()
	ok = r.WaitNext(ctx)
	assert.Equal(t, false, ok)
}

func TestWaitNextCancel(t *testing.T) {
	var (
		bo          = NewConstant(time.Minute)
		r           = NewRetry(3, bo)
		ctx, cancel = context.WithCancel(context.TODO())
	)
	cancel()
	ok := r.WaitNext(ctx)
	assert.Equal(t, false, ok)
}

func TestDo(t *testing.T) {
	var calls int

	testcases := []struct {
		name          string
		canceled      bool
		operation     func() (interface{}, error)
		expected      interface{}
		expectedErr   error
		expectedCalls int
	}{
		{
			name:     "canceled context",
			canceled: true,
			operation: func() (interface{}, error) {
				calls++
				return 1, nil
			},
			expected:      nil,
			expectedErr:   context.Canceled,
			expectedCalls: 0,
		},
		{
			name: "retriable error",
			operation: func() (interface{}, error) {
				calls++
				return nil, NewError(errors.New("retriable-error"), true)
			},
			expected:      nil,
			expectedErr:   errors.New("retriable-error"),
			expectedCalls: 3,
		},
		{
			name: "non-retriable error",
			operation: func() (interface{}, error) {
				calls++
				return nil, NewError(errors.New("non-retriable-error"), false)
			},
			expected:      nil,
			expectedErr:   errors.New("non-retriable-error"),
			expectedCalls: 1,
		},
		{
			name: "not using Error type",
			operation: func() (interface{}, error) {
				calls++
				return nil, errors.New("test-error")
			},
			expected:      nil,
			expectedErr:   errors.New("test-error"),
			expectedCalls: 3,
		},
		{
			name: "ok",
			operation: func() (interface{}, error) {
				calls++
				return 1, nil
			},
			expected:      1,
			expectedErr:   nil,
			expectedCalls: 1,
		},
		{
			name: "ok after a retry",
			operation: func() (interface{}, error) {
				calls++
				if calls == 1 {
					return nil, NewError(errors.New("retriable-error"), true)
				}
				return "data", nil
			},
			expected:      "data",
			expectedErr:   nil,
			expectedCalls: 2,
		},
	}

	for _, tc := range testcases {
		calls = 0
		t.Run(tc.name, func(t *testing.T) {
			var (
				bo          = NewConstant(time.Millisecond)
				r           = NewRetry(3, bo)
				ctx, cancel = context.WithCancel(context.TODO())
			)

			if tc.canceled {
				cancel()
			}

			defer cancel()

			data, err := r.Do(ctx, tc.operation)
			assert.Equal(t, tc.expected, data)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedCalls, calls)
		})
	}
}
