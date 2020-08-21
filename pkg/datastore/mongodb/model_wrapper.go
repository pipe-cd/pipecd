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

// modelWrapper wraps a model representing a BSON document so that the model comes with "_id".
type modelWrapper interface {
	// setID populates the given id to its own "_id".
	setID(id string)
	// storeModel stores the unwrapped model in the value pointed to by v.
	storeModel(v interface{}) error
}

func newWrapper(entity interface{}) (modelWrapper, error) {
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
		return nil, fmt.Errorf("the given entity is unknown type")
	}
}

type application struct {
	model.Application `bson:",inline"`
	ID                string `bson:"_id"`
}

func (a *application) setID(id string) {
	a.ID = id
}

func (a *application) storeModel(v interface{}) error {
	if app, ok := v.(*model.Application); ok {
		*app = a.Application
		return nil
	}
	return fmt.Errorf("the given v is not a pointer to model.Application")
}

type command struct {
	model.Command `bson:",inline"`
	ID            string `bson:"_id"`
}

func (c *command) setID(id string) {
	c.ID = id
}

func (c *command) storeModel(v interface{}) error {
	if command, ok := v.(*model.Command); ok {
		*command = c.Command
		return nil
	}
	return fmt.Errorf("the given v is not a pointer to model.Command")
}

type deployment struct {
	model.Deployment `bson:",inline"`
	ID               string `bson:"_id"`
}

func (d *deployment) setID(id string) {
	d.ID = id
}

func (d *deployment) storeModel(v interface{}) error {
	if deployment, ok := v.(*model.Deployment); ok {
		*deployment = d.Deployment
		return nil
	}
	return fmt.Errorf("the given v is not a pointer to model.Deployment")
}

type environment struct {
	model.Environment `bson:",inline"`
	ID                string `bson:"_id"`
}

func (e *environment) setID(id string) {
	e.ID = id
}

func (e *environment) storeModel(v interface{}) error {
	if environment, ok := v.(*model.Environment); ok {
		*environment = e.Environment
		return nil
	}
	return fmt.Errorf("the given v is not a pointer to model.Environment")
}

type piped struct {
	model.Piped `bson:",inline"`
	ID          string `bson:"_id"`
}

func (p *piped) setID(id string) {
	p.ID = id
}

func (p *piped) storeModel(v interface{}) error {
	if piped, ok := v.(*model.Piped); ok {
		*piped = p.Piped
		return nil
	}
	return fmt.Errorf("the given v is not a pointer to model.Piped")
}

type project struct {
	model.Project `bson:",inline"`
	ID            string `bson:"_id"`
}

func (p *project) setID(id string) {
	p.ID = id
}

func (p *project) storeModel(v interface{}) error {
	if project, ok := v.(*model.Project); ok {
		*project = p.Project
		return nil
	}
	return fmt.Errorf("the given v is not a pointer to model.Project")
}
