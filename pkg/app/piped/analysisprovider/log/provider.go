// Copyright 2020 The PipeCD Authors.
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

package log

import (
	"context"

	"github.com/kapetaniosci/pipe/pkg/app/piped/analysisprovider"
)

// Provider represents a client for log provider which provides logs for analysis.
type Provider interface {
	analysisprovider.Provider
	// RunQuery runs the given query against the log provider,
	// and then checks if there is at least one error log..
	RunQuery(ctx context.Context, query string) (result bool, err error)
}
