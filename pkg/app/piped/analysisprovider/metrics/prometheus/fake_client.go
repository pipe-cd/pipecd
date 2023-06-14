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

package prometheus

import (
	"context"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type fakeClient struct {
	value    model.Value
	err      error
	warnings v1.Warnings
}

func (f fakeClient) QueryRange(_ context.Context, _ string, _ v1.Range) (model.Value, v1.Warnings, error) {
	if f.err != nil {
		return nil, f.warnings, f.err
	}
	return f.value, f.warnings, nil
}
