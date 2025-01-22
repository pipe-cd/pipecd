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

package pipedservice

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/backoff"
)

// Retriable checks whether the caller should retry the api call for the given error.
func Retriable(err error) bool {
	switch status.Code(err) {
	case codes.OK:
		return false
	case codes.InvalidArgument:
		return false
	case codes.NotFound:
		return false
	case codes.AlreadyExists:
		return false
	case codes.PermissionDenied:
		return false
	case codes.FailedPrecondition:
		return false
	case codes.Unimplemented:
		return false
	case codes.Unauthenticated:
		return false
	default:
		return true
	}
}

// NewRetriableErr returns a new backoff.Error for the given error.
// Based on the error code, it determines whether the error is retriable or not.
func NewRetriableErr(err error) *backoff.Error {
	return backoff.NewError(err, Retriable(err))
}

// NewRetry returns a new backoff.Retry for piped API caller.
// 0s 997.867435ms 2.015381172s 3.485134345s 4.389600179s 18.118099328s 48.73058264s
func NewRetry(maxRetries int) backoff.Retry {
	bo := backoff.NewExponential(2*time.Second, time.Minute)
	return backoff.NewRetry(maxRetries, bo)
}
