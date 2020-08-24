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

package mongodb

import (
	"fmt"

	"github.com/pipe-cd/pipe/pkg/model"
)

// wrapModel returns a wrapper corresponding to the given entity.
// A wrapper wraps a model representing BSON a document so that the model comes with "_id".
func wrapModel(entity interface{}) (interface{}, error) {
	switch e := entity.(type) {
	case *model.Application:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &application{
			ID:          e.GetId(),
			Application: *e,
		}, nil
	case *model.Command:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &command{
			ID:      e.GetId(),
			Command: *e,
		}, nil
	case *model.Deployment:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &deployment{
			ID:         e.GetId(),
			Deployment: *e,
		}, nil
	case *model.Environment:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &environment{
			ID:          e.GetId(),
			Environment: *e,
		}, nil
	case *model.Piped:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &piped{
			ID:    e.GetId(),
			Piped: *e,
		}, nil
	case *model.Project:
		if e == nil {
			return nil, fmt.Errorf("nil entity given")
		}
		return &project{
			ID:      e.GetId(),
			Project: *e,
		}, nil
	default:
		return nil, fmt.Errorf("%T is not supported", e)
	}
}

// extractModel stores the unwrapped model in the value pointed to by e.
func extractModel(wrapper interface{}, e interface{}) error {
	msg := "entity type doesn't correspond to the wrapper type (%T)"

	switch w := wrapper.(type) {
	case *application:
		e, ok := e.(*model.Application)
		if !ok {
			return fmt.Errorf(msg, w)
		}
		*e = w.Application
	case *command:
		e, ok := e.(*model.Command)
		if !ok {
			return fmt.Errorf(msg, w)
		}
		*e = w.Command
	case *deployment:
		e, ok := e.(*model.Deployment)
		if !ok {
			return fmt.Errorf(msg, w)
		}
		*e = w.Deployment
	case *environment:
		e, ok := e.(*model.Environment)
		if !ok {
			return fmt.Errorf(msg, w)
		}
		*e = w.Environment
	case *piped:
		e, ok := e.(*model.Piped)
		if !ok {
			return fmt.Errorf(msg, w)
		}
		*e = w.Piped
	case *project:
		e, ok := e.(*model.Project)
		if !ok {
			return fmt.Errorf(msg, w)
		}
		*e = w.Project
	default:
		return fmt.Errorf("%T is not supported", w)
	}
	return nil
}

type application struct {
	model.Application `bson:",inline"`
	ID                string `bson:"_id"`
}

type command struct {
	model.Command `bson:",inline"`
	ID            string `bson:"_id"`
}

type deployment struct {
	model.Deployment `bson:",inline"`
	ID               string `bson:"_id"`
}

type environment struct {
	model.Environment `bson:",inline"`
	ID                string `bson:"_id"`
}

type piped struct {
	model.Piped `bson:",inline"`
	ID          string `bson:"_id"`
}

type project struct {
	model.Project `bson:",inline"`
	ID            string `bson:"_id"`
}
