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

package mysql

import (
	"encoding/json"
	"fmt"

	"github.com/pipe-cd/pipecd/pkg/model"
)

// wrapModel attaches an extra field named `_extra` to JSON data.
// Currently, on MySQL v8, functional indexes are not hit on `LIKE` queries,
// this added `_extra` field will be shadowed by table column `extra` so that
// we could create indexes (FULLTEXT or normal) on that column for search features.
// ref: https://github.com/pipe-cd/pipecd/blob/master/docs/rfcs/0003-sql-datastore.md#text-search-operations-on-specific-json-field
func wrapModel(entity interface{}) (interface{}, error) {
	switch e := entity.(type) {
	case *model.Project:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &project{
			Project: *e,
			Extra:   e.Id,
		}, nil
	case *model.Application:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &application{
			Application: *e,
			Extra:       e.Name,
		}, nil
	case *model.Command:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &command{
			Command: *e,
			Extra:   e.Id,
		}, nil
	case *model.Deployment:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &deployment{
			Deployment: *e,
			Extra:      e.ApplicationName,
		}, nil
	case *model.Piped:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &piped{
			Piped: *e,
			Extra: e.Name,
		}, nil
	case *model.APIKey:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &apiKey{
			APIKey: *e,
			Extra:  e.Name,
		}, nil
	case *model.Event:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &event{
			Event: *e,
			Extra: e.Name,
		}, nil
	case *model.DeploymentChain:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &deploymentChain{
			DeploymentChain: *e,
			Extra:           e.Id,
		}, nil
	default:
		return nil, fmt.Errorf("%T is not supported", e)
	}
}

func encodeJSONValue(entity interface{}) (string, error) {
	wrapper, err := wrapModel(entity)
	if err != nil {
		return "", err
	}
	encodedEntity, err := json.Marshal(wrapper)
	if err != nil {
		return "", err
	}
	return string(encodedEntity), nil
}

func decodeJSONValue(val string, target interface{}) error {
	return json.Unmarshal([]byte(val), target)
}

type project struct {
	model.Project `json:",inline"`
	Extra         string `json:"_extra"`
}

type application struct {
	model.Application `json:",inline"`
	Extra             string `json:"_extra"`
}

type command struct {
	model.Command `json:",inline"`
	Extra         string `json:"_extra"`
}

type deployment struct {
	model.Deployment `json:",inline"`
	Extra            string `json:"_extra"`
}

type piped struct {
	model.Piped `json:",inline"`
	Extra       string `json:"_extra"`
}

type apiKey struct {
	model.APIKey `json:",inline"`
	Extra        string `json:"_extra"`
}

type event struct {
	model.Event `json:",inline"`
	Extra       string `json:"_extra"`
}

type deploymentChain struct {
	model.DeploymentChain `json:",inline"`
	Extra                 string `json:"_extra"`
}
