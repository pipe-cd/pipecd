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

package firestore

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipecd/pkg/datastore"
)

const (
	defaultNamespace = "pipecd"
)

var operatorMap = map[datastore.Operator]string{
	datastore.OperatorEqual:              "==",
	datastore.OperatorNotEqual:           "!=",
	datastore.OperatorIn:                 "in",
	datastore.OperatorNotIn:              "not-in",
	datastore.OperatorGreaterThan:        ">",
	datastore.OperatorGreaterThanOrEqual: ">=",
	datastore.OperatorLessThan:           "<",
	datastore.OperatorLessThanOrEqual:    "<=",
	datastore.OperatorContains:           "array-contains",
}

type FireStore struct {
	client               *firestore.Client
	namespace            string
	environment          string
	collectionNamePrefix string

	credentialsFile string
	logger          *zap.Logger
}

func NewFireStore(ctx context.Context, projectID, namespace, environment string, opts ...Option) (*FireStore, error) {
	s := &FireStore{
		namespace:   namespace,
		environment: environment,
		logger:      zap.NewNop(),
	}
	if s.namespace == "" {
		s.namespace = defaultNamespace
	}
	for _, opt := range opts {
		opt(s)
	}
	s.logger = s.logger.Named("firestore")

	var options []option.ClientOption
	if s.credentialsFile != "" {
		options = append(options, option.WithCredentialsFile(s.credentialsFile))
	}

	client, err := firestore.NewClient(ctx, projectID, options...)
	if err != nil {
		return nil, err
	}
	s.client = client
	return s, nil
}

type Option func(*FireStore)

func WithCredentialsFile(path string) Option {
	return func(s *FireStore) {
		s.credentialsFile = path
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *FireStore) {
		s.logger = logger
	}
}

func WithCollectionNamePrefix(prefix string) Option {
	return func(s *FireStore) {
		s.collectionNamePrefix = prefix
	}
}

func (s *FireStore) Find(ctx context.Context, col datastore.Collection, opts datastore.ListOptions) (datastore.Iterator, error) {
	kind := col.Kind()
	if opts.Cursor != "" && len(opts.Orders) == 0 {
		return nil, errors.New("opts.Cursor also requires Orders to be set")
	}

	colName := makeCollectionName(s.collectionNamePrefix, kind)

	q := s.client.Collection(s.namespace).Doc(s.environment).Collection(colName).Query
	for _, f := range opts.Filters {
		op, ok := operatorMap[f.Operator]
		if !ok {
			return nil, fmt.Errorf("unsupported operator given: %v", f.Operator)
		}
		q = q.Where(f.Field, op, f.Value)
	}
	for _, o := range opts.Orders {
		q = q.OrderBy(o.Field, convertToDirection(o.Direction))
	}

	// The pseudo cursor points one behind of the target document.
	// See more: https://cloud.google.com/firestore/docs/query-data/query-cursors?hl=ja
	if opts.Cursor != "" {
		values, err := makeCursorValues(opts)
		if err != nil {
			return nil, err
		}
		q = q.StartAfter(values...)
	}

	if opts.Limit > 0 {
		q = q.Limit(opts.Limit)
	}
	return &Iterator{
		it:     q.Documents(ctx),
		orders: opts.Orders,
	}, nil
}

func (s *FireStore) Get(ctx context.Context, col datastore.Collection, id string, v interface{}) error {
	kind := col.Kind()
	colName := makeCollectionName(s.collectionNamePrefix, kind)
	ds, err := s.client.Collection(s.namespace).Doc(s.environment).Collection(colName).Doc(id).Get(ctx)
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			return datastore.ErrNotFound
		}
		s.logger.Error("failed to retrieve entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	if err := ds.DataTo(v); err != nil {
		s.logger.Error("failed to unmarshal entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *FireStore) Create(ctx context.Context, col datastore.Collection, id string, entity interface{}) error {
	kind := col.Kind()
	colName := makeCollectionName(s.collectionNamePrefix, kind)
	ref := s.client.Collection(s.namespace).Doc(s.environment).Collection(colName).Doc(id)
	err := s.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := tx.Get(ref)
		if err == nil {
			return datastore.ErrAlreadyExists
		}
		if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
			s.logger.Error("failed to retrieve entity",
				zap.String("id", id),
				zap.String("kind", kind),
				zap.Error(err),
			)
			return err
		}
		return tx.Set(ref, entity)
	})
	if err != nil {
		s.logger.Error("failed to create entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *FireStore) Update(ctx context.Context, col datastore.Collection, id string, updater datastore.Updater) error {
	kind := col.Kind()
	colName := makeCollectionName(s.collectionNamePrefix, kind)
	ref := s.client.Collection(s.namespace).Doc(s.environment).Collection(colName).Doc(id)
	err := s.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		entity := col.Factory()()
		ds, err := tx.Get(ref)
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				return datastore.ErrNotFound
			}
			s.logger.Error("failed to retrieve entity",
				zap.String("id", id),
				zap.String("kind", kind),
				zap.Error(err),
			)
			return err
		}
		if err := ds.DataTo(entity); err != nil {
			s.logger.Error("failed to unmarshal entity",
				zap.String("id", id),
				zap.String("kind", kind),
				zap.Error(err),
			)
			return err
		}
		if err := updater(entity); err != nil {
			s.logger.Error("failed to run updater to update entity",
				zap.String("id", id),
				zap.String("kind", kind),
				zap.Error(err),
			)
			return err
		}
		return tx.Set(ref, entity)
	})
	if err != nil {
		s.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *FireStore) Close() error {
	return s.client.Close()
}

func makeCursorValues(opts datastore.ListOptions) ([]interface{}, error) {
	// Decode last object of previous page stored as opts.Cursor to string.
	data, err := base64.StdEncoding.DecodeString(opts.Cursor)
	if err != nil {
		return nil, err
	}
	// Encode cursor data string to map[string]interface{} format for further process.
	obj := make(map[string]interface{})
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	cursorVals := make([]interface{}, 0, len(opts.Orders))
	hasIDFieldInOrdering := false
	for _, o := range opts.Orders {
		if o.Field == "Id" {
			hasIDFieldInOrdering = true
		}
		val, ok := obj[o.Field]
		if !ok {
			return nil, errors.New("cursor does not contain values that match to ordering field")
		}
		cursorVals = append(cursorVals, val)
	}
	// The Id field is a required field to keep the sorted query result stable.
	// We also do not use `id => snapshot doc => snapshot doc value` pattern since the snapshot here is not
	// real snapshot (stable, unchanged on query) but just a single query to get one document by id, which could
	// lead us to unstable/unpredictable results if the value of that "snapshot" doc is changed since previous
	// query.
	// Read more: https://cloud.google.com/firestore/docs/query-data/query-cursors#set_cursor_based_on_multiple_fields
	if !hasIDFieldInOrdering {
		return nil, errors.New("id field is required as ordering field")
	}

	return cursorVals, nil
}

func convertToDirection(od datastore.OrderDirection) firestore.Direction {
	if od == datastore.Asc {
		return firestore.Asc
	}
	return firestore.Desc
}

func makeCollectionName(prefix, kind string) string {
	if prefix == "" {
		return kind
	}
	return prefix + kind
}
