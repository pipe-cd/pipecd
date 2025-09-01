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

package datastore

import (
	"context"
	"errors"
)

type OrderDirection int

const (
	// Asc sorts results from smallest to largest.
	Asc OrderDirection = iota + 1
	// Desc sorts results from largest to smallest.
	Desc
)

type Operator int

const (
	// Operation to find the field is equal to the specified value.
	OperatorEqual = iota + 1
	// Operation to find the field isn't equal to the specified value.
	OperatorNotEqual
	// Operation to find ones that contain any one of the multiple values.
	OperatorIn
	// Operation to find ones that do not contain any of the specified multiple values.
	OperatorNotIn
	// Operation to find ones the field is greater than the specified value.
	OperatorGreaterThan
	// Operation to find ones the field is greater or equal than the specified value.
	OperatorGreaterThanOrEqual
	// Operation to find ones the field is less than the specified value.
	OperatorLessThan
	// Operation to find ones the field is less or equal than the specified value.
	OperatorLessThanOrEqual
	// Operation to find ones that have a specified value in its array.
	OperatorContains
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidCursor   = errors.New("invalid cursor")
	ErrIteratorDone    = errors.New("iterator is done")
	ErrInternal        = errors.New("internal")
	ErrUnimplemented   = errors.New("unimplemented")
	ErrUnsupported     = errors.New("unsupported")
	ErrUserDefined     = errors.New("user defined error")
)

type Factory func() interface{}
type Updater func(interface{}) error

type Collection interface {
	Kind() string
	Factory() Factory
}

type DataStore interface {
	// Find finds the documents matched given criteria.
	Find(ctx context.Context, col Collection, opts ListOptions) (Iterator, error)
	// Get gets one document specified with ID, and unmarshal it to typed struct.
	// If the document can not be found in datastore, ErrNotFound will be returned.
	Get(ctx context.Context, col Collection, id string, entity interface{}) error
	// Create saves a new entity to the datastore.
	// If an entity with the same ID is already existing, ErrAlreadyExists will be returned.
	Create(ctx context.Context, col Collection, id string, entity interface{}) error
	// Update updates an existing entity in the datastore.
	// If updating entity was not found in the datastore, ErrNotFound will be returned.
	Update(ctx context.Context, col Collection, id string, updater Updater) error
	// Close closes datastore resources held by the client.
	Close() error
}

type Iterator interface {
	Next(dst interface{}) error
	Cursor() (string, error)
}

type ListFilter struct {
	Field    string
	Operator Operator
	Value    interface{}
}

type Order struct {
	Field     string
	Direction OrderDirection
}

type ListOptions struct {
	Limit   int
	Filters []ListFilter
	Orders  []Order
	Cursor  string
}

type backend struct {
	ds  DataStore
	col Collection
}
