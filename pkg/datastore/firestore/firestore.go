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

package firestore

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

const (
	defaultNamespace = "pipecd"
)

type FireStore struct {
	client      *firestore.Client
	namespace   string
	environment string

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

func (s *FireStore) Find(ctx context.Context, kind string, opts datastore.ListOptions) (datastore.Iterator, error) {
	if opts.Cursor != "" && len(opts.Orders) == 0 {
		return nil, errors.New("opts.Cursor also requires Orders to be set")
	}

	cursorSnapshot, err := s.fetchCursorDocumentSnapshot(ctx, kind, opts)
	if err != nil {
		return nil, err
	}

	q := s.client.Collection(s.namespace).Doc(s.environment).Collection(kind).Query
	for _, f := range opts.Filters {
		q = q.Where(f.Field, f.Operator, f.Value)
	}
	for _, o := range opts.Orders {
		q = q.OrderBy(o.Field, convertToDirection(o.Direction))
	}
	// Note: opts.Page parameter does not use in Cloud Firestore. Firestore cannot do paging like general NoSQL.
	// Instead of general paging, it will be a workload like infinite scroll.
	// The pseudo cursor points one behind of the target document.
	// See more: https://cloud.google.com/firestore/docs/query-data/query-cursors?hl=ja
	if cursorSnapshot != nil {
		q = q.StartAfter(cursorSnapshot.Data()[opts.Orders[0].Field])
	}

	if opts.PageSize > 0 {
		q = q.Limit(opts.PageSize)
	}
	return &Iterator{
		it: q.Documents(ctx),
	}, nil
}

func (s *FireStore) Get(ctx context.Context, kind, id string, v interface{}) error {
	ds, err := s.client.Collection(s.namespace).Doc(s.environment).Collection(kind).Doc(id).Get(ctx)
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

func (s *FireStore) Create(ctx context.Context, kind, id string, entity interface{}) error {
	ref := s.client.Collection(s.namespace).Doc(s.environment).Collection(kind).Doc(id)
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

func (s *FireStore) Put(ctx context.Context, kind, id string, entity interface{}) error {
	col := s.client.Collection(s.namespace).Doc(s.environment).Collection(kind)
	if _, err := col.Doc(id).Set(ctx, entity); err != nil {
		s.logger.Info("failed to put entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *FireStore) Update(ctx context.Context, kind, id string, factory datastore.Factory, updater datastore.Updater) error {
	ref := s.client.Collection(s.namespace).Doc(s.environment).Collection(kind).Doc(id)
	err := s.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		entity := factory()
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

func (s *FireStore) fetchCursorDocumentSnapshot(ctx context.Context, kind string, opts datastore.ListOptions) (*firestore.DocumentSnapshot, error) {
	if opts.Cursor == "" {
		return nil, nil
	}
	return s.client.Collection(s.namespace).Doc(s.environment).Collection(kind).Doc(opts.Cursor).Get(ctx)
}

func convertToDirection(od datastore.OrderDirection) firestore.Direction {
	if od == datastore.Asc {
		return firestore.Asc
	}
	return firestore.Desc
}
