// Copyright 2021 The PipeCD Authors.
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

package mysql

import (
	"fmt"

	"github.com/pipe-cd/pipe/pkg/model"
)

func wrapModel(entity interface{}) (interface{}, error) {
	switch e := entity.(type) {
	case *model.Application:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &application{
			ID:   e.GetId(),
			Data: e,
		}, nil
	case *model.Project:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &project{
			ID:   e.GetId(),
			Data: e,
		}, nil
	default:
		return nil, fmt.Errorf("%T is not supported", e)
	}
}

type application struct {
	ID   string
	Data *model.Application
}

type project struct {
	ID   string
	Data *model.Project
}
