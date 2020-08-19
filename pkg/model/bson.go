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

package model

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

const mongodbPrimaryKey = "_id"

// MarshalBSON overrides bson.Marshal(). This overwritten aims to populate
// its own id to "_id" and instead remove id which duplicates "_id".
func (a *Application) MarshalBSON() ([]byte, error) {
	// TODO: Replace "fatih/structs" with "mitchellh/mapstructure"
	//   Currently using fatih/structs temporarily to easily convert all nested structures into a map.
	m := structs.Map(a)
	m[mongodbPrimaryKey] = a.Id
	delete(m, "Id")

	return bson.Marshal(m)
}

// UnmarshalBSON overrides bson.Unmarshal().
// This overwritten aims to populate "_id" to its own id.
func (a *Application) UnmarshalBSON(b []byte) error {
	m := make(map[string]interface{})
	if err := bson.Unmarshal(b, m); err != nil {
		return err
	}

	if err := mapstructure.Decode(m, a); err != nil {
		return fmt.Errorf("failed to decode map: %w", err)
	}
	if id, ok := m[mongodbPrimaryKey]; ok {
		a.Id = id.(string)
	}

	return nil
}

// MarshalBSON overrides bson.Marshal(). This overwritten aims to populate
// its own id to "_id" and instead remove id which duplicates "_id".
func (c *Command) MarshalBSON() ([]byte, error) {
	m := structs.Map(c)
	m[mongodbPrimaryKey] = c.Id
	delete(m, "Id")

	return bson.Marshal(m)
}

// UnmarshalBSON overrides bson.Unmarshal().
// This overwritten aims to populate "_id" to its own id.
func (c *Command) UnmarshalBSON(b []byte) error {
	m := make(map[string]interface{})
	if err := bson.Unmarshal(b, m); err != nil {
		return err
	}

	if err := mapstructure.Decode(m, c); err != nil {
		return fmt.Errorf("failed to decode map: %w", err)
	}
	if id, ok := m[mongodbPrimaryKey]; ok {
		c.Id = id.(string)
	}

	return nil
}

// MarshalBSON overrides bson.Marshal(). This overwritten aims to populate
// its own id to "_id" and instead remove id which duplicates "_id".
func (d *Deployment) MarshalBSON() ([]byte, error) {
	m := structs.Map(d)
	m[mongodbPrimaryKey] = d.Id
	delete(m, "Id")

	return bson.Marshal(m)
}

// UnmarshalBSON overrides bson.Unmarshal().
// This overwritten aims to populate "_id" to its own id.
func (d *Deployment) UnmarshalBSON(b []byte) error {
	m := make(map[string]interface{})
	if err := bson.Unmarshal(b, m); err != nil {
		return err
	}

	if err := mapstructure.Decode(m, d); err != nil {
		return fmt.Errorf("failed to decode map: %w", err)
	}
	if id, ok := m[mongodbPrimaryKey]; ok {
		d.Id = id.(string)
	}

	return nil
}

// MarshalBSON overrides bson.Marshal(). This overwritten aims to populate
// its own id to "_id" and instead remove id which duplicates "_id".
func (e *Environment) MarshalBSON() ([]byte, error) {
	m := structs.Map(e)
	m[mongodbPrimaryKey] = e.Id
	delete(m, "Id")

	return bson.Marshal(m)
}

// UnmarshalBSON overrides bson.Unmarshal().
// This overwritten aims to populate "_id" to its own id.
func (e *Environment) UnmarshalBSON(b []byte) error {
	m := make(map[string]interface{})
	if err := bson.Unmarshal(b, m); err != nil {
		return err
	}

	if err := mapstructure.Decode(m, e); err != nil {
		return fmt.Errorf("failed to decode map: %w", err)
	}
	if id, ok := m[mongodbPrimaryKey]; ok {
		e.Id = id.(string)
	}

	return nil
}

// MarshalBSON overrides bson.Marshal(). This overwritten aims to populate
// its own id to "_id" and instead remove id which duplicates "_id".
func (p *Piped) MarshalBSON() ([]byte, error) {
	m := structs.Map(p)
	m[mongodbPrimaryKey] = p.Id
	delete(m, "Id")

	return bson.Marshal(m)
}

// UnmarshalBSON overrides bson.Unmarshal().
// This overwritten aims to populate "_id" to its own id.
func (p *Piped) UnmarshalBSON(b []byte) error {
	m := make(map[string]interface{})
	if err := bson.Unmarshal(b, m); err != nil {
		return err
	}

	if err := mapstructure.Decode(m, p); err != nil {
		return fmt.Errorf("failed to decode map: %w", err)
	}
	if id, ok := m[mongodbPrimaryKey]; ok {
		p.Id = id.(string)
	}

	return nil
}

// MarshalBSON overrides bson.Marshal(). This overwritten aims to populate
// its own id to "_id" and instead remove id which duplicates "_id".
func (p *Project) MarshalBSON() ([]byte, error) {
	m := structs.Map(p)
	m[mongodbPrimaryKey] = p.Id
	delete(m, "Id")

	return bson.Marshal(m)
}

// UnmarshalBSON overrides bson.Unmarshal().
// This overwritten aims to populate "_id" to its own id.
func (p *Project) UnmarshalBSON(b []byte) error {
	m := make(map[string]interface{})
	if err := bson.Unmarshal(b, m); err != nil {
		return err
	}

	if err := mapstructure.Decode(m, p); err != nil {
		return fmt.Errorf("failed to decode map: %w", err)
	}
	if id, ok := m[mongodbPrimaryKey]; ok {
		p.Id = id.(string)
	}

	return nil
}
