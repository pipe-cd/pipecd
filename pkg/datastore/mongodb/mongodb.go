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
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

type MongoDB struct {
	ctx       context.Context
	client    *mongo.Client
	namespace string
	direct    bool

	logger *zap.Logger
}

func NewMongoDB(ctx context.Context, url, namespace string, opts ...Option) (*MongoDB, error) {
	m := &MongoDB{
		ctx:       ctx,
		namespace: namespace,
		logger:    zap.NewNop(),
	}
	for _, opt := range opts {
		opt(m)
	}
	m.logger = m.logger.Named("mongodb")

	clientOpts := options.Client().SetDirect(m.direct).ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	m.client = client
	return m, nil
}

type Option func(*MongoDB)

func WithLogger(logger *zap.Logger) Option {
	return func(s *MongoDB) {
		s.logger = logger
	}
}

func WithDirect(direct bool) Option {
	return func(s *MongoDB) {
		s.direct = direct
	}
}

func (m *MongoDB) Find(ctx context.Context, kind string, opts datastore.ListOptions) (datastore.Iterator, error) {
	col := m.client.Database(m.namespace).Collection(kind)
	opes := make([]bson.M, len(opts.Filters))
	for i, f := range opts.Filters {
		ope, err := convertToMongoDBOperator(f.Operator)
		if err != nil {
			return nil, err
		}
		opes[i] = bson.M{
			// Note: The field name of protobuf is saved in lower case by default in mongodb.
			// e.g. Name => name, ProjectId => projectid, CreatedAt, createdat
			// The field name in mongodb can be set by bson tag in proto file, but setting this for all fields is very tedious.
			// For that reason, does change strings to lowercase before making a query. However, the exception is the "_id" field.
			strings.ToLower(f.Field): bson.M{
				ope: f.Value,
			},
		}
	}
	query := bson.M{}
	if len(opes) > 0 {
		query["$and"] = opes
	}

	findOpts := options.Find()
	if opts.PageSize > 0 {
		findOpts.SetLimit(int64(opts.PageSize))
		if opts.Page > 0 {
			findOpts.SetSkip(int64(opts.PageSize * opts.Page))
		}
	}

	cur, err := col.Find(ctx, query, findOpts)
	if err != nil {
		m.logger.Error("failed to get cursor",
			zap.String("kind", kind),
			zap.Error(err),
		)
		return nil, err
	}
	return &Iterator{
		ctx: ctx,
		cur: cur,
	}, nil
}

func (m *MongoDB) Get(ctx context.Context, kind, id string, v interface{}) error {
	col := m.client.Database(m.namespace).Collection(kind)
	err := col.FindOne(ctx, makePrimaryKeyFilter(id)).Decode(v)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return datastore.ErrNotFound
	}
	if err != nil {
		m.logger.Error("failed to retrieve entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (m *MongoDB) Create(ctx context.Context, kind, id string, entity interface{}) error {
	// TODO: Support updating process with using transaction on mongoDB cluster
	col := m.client.Database(m.namespace).Collection(kind)
	err := col.FindOne(ctx, makePrimaryKeyFilter(id), options.FindOne()).Err()
	if err == nil {
		return datastore.ErrAlreadyExists
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		m.logger.Error("failed to retrieve entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	if _, err := col.InsertOne(ctx, entity); err != nil {
		m.logger.Error("failed to insert entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (m *MongoDB) Put(ctx context.Context, kind, id string, entity interface{}) error {
	col := m.client.Database(m.namespace).Collection(kind)
	err := m.client.UseSession(ctx, func(sc mongo.SessionContext) error {
		if _, err := col.UpdateOne(sc, makePrimaryKeyFilter(id), entity); err != nil {
			return err
		}
		return sc.CommitTransaction(sc)
	})
	if err != nil {
		m.logger.Error("failed to put entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (m *MongoDB) Update(ctx context.Context, kind, id string, factory datastore.Factory, updater datastore.Updater) error {
	col := m.client.Database(m.namespace).Collection(kind)
	entity := factory()
	err := col.FindOne(ctx, makePrimaryKeyFilter(id)).Decode(entity)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return datastore.ErrNotFound
	}
	if err != nil {
		return err
	}
	if err := updater(entity); err != nil {
		m.logger.Error("failed to run updater to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	update := bson.D{{"$set", entity}}
	if _, err := col.UpdateOne(ctx, makePrimaryKeyFilter(id), update); err != nil {
		m.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (m *MongoDB) Close() error {
	return m.client.Disconnect(m.ctx)
}

func makePrimaryKeyFilter(id string) bson.D {
	return bson.D{{"_id", id}}
}
