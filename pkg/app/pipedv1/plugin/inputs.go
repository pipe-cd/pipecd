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

package plugin

import (
	"bytes"
	"encoding/json"

	"github.com/pipe-cd/pipecd/pkg/plugin/api/v1alpha1/platform"

	"github.com/creasty/defaults"
)

// DecodeApplicationSpec decodes the spec field of the given ApplicationConfig
func DecodeApplicationSpec[T any](src *platform.ApplicationConfig) (*T, error) {
	dec := json.NewDecoder(bytes.NewReader(src.GetSpec()))
	dec.DisallowUnknownFields()

	dest := new(T)

	if err := dec.Decode(dest); err != nil {
		return nil, err
	}

	if err := defaults.Set(dest); err != nil {
		return nil, err
	}

	return dest, nil
}
