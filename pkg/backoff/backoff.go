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

package backoff

import (
	"context"
	"time"
)

type Error struct {
	Err       error
	Retriable bool
}

func NewError(err error, retriable bool) *Error {
	return &Error{
		Err:       err,
		Retriable: retriable,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

type Backoff interface {
	Next() time.Duration
	Calls() int
	Reset()
	Clone() Backoff
}

type Retry interface {
	Do(ctx context.Context, operation func() (interface{}, error)) (interface{}, error)
	WaitNext(ctx context.Context) bool
	Calls() int
}

func NewRetry(max int, backoff Backoff) Retry {
	return &retry{
		max:     max,
		backoff: backoff,
	}
}

type retry struct {
	max     int
	calls   int
	ctx     context.Context
	backoff Backoff
}

// TODO: Find all using of WaitNext and replace by Do to avoid panic.
func (r *retry) WaitNext(ctx context.Context) bool {
	defer func() {
		r.calls++
	}()

	if r.calls >= r.max {
		return false
	}

	d := r.backoff.Next()
	if d == 0 {
		select {
		case <-ctx.Done():
			return false
		default:
			return true
		}
	}

	t := time.NewTimer(d)
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

// Do executes and returns the results of the given operation.
// It automatically retries if the operation returns an error.
// To control which error should be retriable or not,
// you can wrap the error from operation with NewError function.
func (r *retry) Do(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
	var err error

	for r.WaitNext(ctx) {
		var data interface{}
		data, err = operation()
		if err == nil {
			return data, nil
		}
		if e, ok := err.(*Error); ok {
			err = e.Err
			if !e.Retriable {
				return nil, err
			}
		}
	}
	if err != nil {
		return nil, err
	}

	// Operation was not executed due to context error.
	return nil, ctx.Err()
}

func (r *retry) Calls() int {
	return r.calls
}
