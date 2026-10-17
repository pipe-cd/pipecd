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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestWaitNext_ContextCancelDuringWait(t *testing.T) {
	t.Parallel()

	bo := NewConstant(time.Minute)
	r := NewRetry(3, bo)
	ctx, cancel := context.WithCancel(context.Background())

	// First call has 0 delay.
	ok := r.WaitNext(ctx)
	assert.True(t, ok)

	// Second call waits for 1 minute; cancel context during wait.
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	ok = r.WaitNext(ctx)
	assert.False(t, ok)
	assert.True(t, time.Since(start) < time.Second)
}

func TestDo(t *testing.T) {
	var calls int

	testcases := []struct {
		name          string
		canceled      bool
		operation     func() (any, error)
		expected      any
		expectedErr   error
		expectedCalls int
	}{
		{
			name:     "canceled context",
			canceled: true,
			operation: func() (any, error) {
				calls++
				return 1, nil
			},
			expected:      nil,
			expectedErr:   context.Canceled,
			expectedCalls: 0,
		},
		{
			name: "retriable error",
			operation: func() (any, error) {
				calls++
				return nil, NewError(errors.New("retriable-error"), true)
			},
			expected:      nil,
			expectedErr:   errors.New("retriable-error"),
			expectedCalls: 3,
		},
		{
			name: "non-retriable error",
			operation: func() (any, error) {
				calls++
				return nil, NewError(errors.New("non-retriable-error"), false)
			},
			expected:      nil,
			expectedErr:   errors.New("non-retriable-error"),
			expectedCalls: 1,
		},
		{
			name: "not using Error type",
			operation: func() (any, error) {
				calls++
				return nil, errors.New("test-error")
			},
			expected:      nil,
			expectedErr:   errors.New("test-error"),
			expectedCalls: 3,
		},
		{
			name: "ok",
			operation: func() (any, error) {
				calls++
				return 1, nil
			},
			expected:      1,
			expectedErr:   nil,
			expectedCalls: 1,
		},
		{
			name: "ok after a retry",
			operation: func() (any, error) {
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
		{
			name: "wrapped non-retriable error stops retry",
			operation: func() (any, error) {
				calls++
				return nil, fmt.Errorf("wrapped: %w", NewError(errors.New("non-retriable"), false))
			},
			expected:      nil,
			expectedErr:   fmt.Errorf("wrapped: %w", NewError(errors.New("non-retriable"), false)),
			expectedCalls: 1,
		},
		{
			name: "wrapped retriable error continues retry",
			operation: func() (any, error) {
				calls++
				return nil, fmt.Errorf("wrapped: %w", NewError(errors.New("retriable"), true))
			},
			expected:      nil,
			expectedErr:   fmt.Errorf("wrapped: %w", NewError(errors.New("retriable"), true)),
			expectedCalls: 3,
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
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedCalls, calls)
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()

	baseErr := errors.New("base error")
	bErr := NewError(baseErr, true)

	assert.Equal(t, baseErr, bErr.Unwrap())
	assert.True(t, errors.Is(bErr, baseErr))

	var target *Error
	assert.True(t, errors.As(bErr, &target))
	assert.Equal(t, true, target.Retriable)

	wrapped := fmt.Errorf("outer: %w", bErr)
	assert.True(t, errors.Is(wrapped, baseErr))
	var targetFromWrapped *Error
	assert.True(t, errors.As(wrapped, &targetFromWrapped))
	assert.Equal(t, true, targetFromWrapped.Retriable)
}

func TestError_NilErr(t *testing.T) {
	t.Parallel()

	bErr := NewError(nil, false)
	assert.Equal(t, "", bErr.Error())
	assert.Nil(t, bErr.Unwrap())
}

func TestRetry_Calls(t *testing.T) {
	t.Parallel()

	r := NewRetry(5, NewConstant(time.Millisecond))
	assert.Equal(t, 0, r.Calls())

	ctx := context.Background()
	r.WaitNext(ctx)
	assert.Equal(t, 1, r.Calls())
}

func TestConstant_Clone(t *testing.T) {
	t.Parallel()

	c := NewConstant(10 * time.Millisecond)
	_ = c.Next()
	assert.Equal(t, 1, c.Calls())

	clone := c.Clone()
	assert.Equal(t, 0, clone.Calls())
	assert.Equal(t, time.Duration(0), clone.Next())
	assert.Equal(t, 10*time.Millisecond, clone.Next())
}
