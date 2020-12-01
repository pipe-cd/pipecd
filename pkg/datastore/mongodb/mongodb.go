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
	"fmt"
	"io/ioutil"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

const (
	// Scram type ref https://github.com/mongodb/mongo-go-driver/blob/9e2aca8afd8821e6b068cc2f25192bc640d90a0d/mongo/client_examples_test.go#L119
	// Default
	scram string = "SCRAM"
)

type MongoDB struct {
	ctx           context.Context
	client        *mongo.Client
	database      string
	logger        *zap.Logger
	authMechanism string
	usernameFile  string
	passwordFile  string
}

func NewMongoDB(ctx context.Context, url, database string, opts ...Option) (*MongoDB, error) {
	m := &MongoDB{
		ctx:      ctx,
		database: database,
		logger:   zap.NewNop(),
	}
	for _, opt := range opts {
		opt(m)
	}
	m.logger = m.logger.Named("mongodb")

	clientOpts := options.Client().ApplyURI(url)

	if m.authMechanism != "" {
		credential, err := m.determineCredential()
		if err != nil {
			return nil, err
		}
		clientOpts.SetAuth(*credential)
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	m.client = client
	return m, nil
}

type Option func(*MongoDB)

func WithAuthenticationFile(usernameFile, passwordFile string) Option {
	return func(m *MongoDB) {
		m.usernameFile = usernameFile
		m.passwordFile = passwordFile
		m.authMechanism = scram
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *MongoDB) {
		s.logger = logger
	}
}

func (m *MongoDB) Find(ctx context.Context, kind string, opts datastore.ListOptions) (datastore.Iterator, error) {
	col := m.client.Database(m.database).Collection(kind)
	ops := make([]bson.M, len(opts.Filters))
	for i, f := range opts.Filters {
		op, err := convertToMongoDBOperator(f.Operator)
		if err != nil {
			return nil, err
		}
		ops[i] = bson.M{
			// Note: The field name of protobuf is saved in lower case by default in mongodb.
			// e.g. Name => name, ProjectId => projectid, CreatedAt => createdat
			strings.ToLower(f.Field): bson.M{
				op: f.Value,
			},
		}
	}
	query := bson.M{}
	if len(ops) > 0 {
		query["$and"] = ops
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
	wrapper, err := wrapModel(v)
	if err != nil {
		return err
	}

	col := m.client.Database(m.database).Collection(kind)
	err = col.FindOne(ctx, makePrimaryKeyFilter(id)).Decode(wrapper)
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

	return extractModel(wrapper, v)
}

func (m *MongoDB) Create(ctx context.Context, kind, id string, entity interface{}) error {
	wrapper, err := wrapModel(entity)
	if err != nil {
		return err
	}

	// TODO: Support updating process with using transaction on mongoDB cluster
	//   err := m.client.UseSession(ctx, func(sessCtx mongo.SessionContext) error { }
	//   See the example at: https://godoc.org/go.mongodb.org/mongo-driver/mongo#Client.UseSessionWithOptions
	//   NOTE:
	//   - Multi-document transactions are only available in version 4.0 or later.
	//   - Also available for replica set deployments only.
	//   - Available even on a standalone server but need to configure it as a replica set (with just one node)
	col := m.client.Database(m.database).Collection(kind)
	err = col.FindOne(ctx, makePrimaryKeyFilter(id), options.FindOne()).Err()
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

	if _, err := col.InsertOne(ctx, wrapper); err != nil {
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
	wrapper, err := wrapModel(entity)
	if err != nil {
		return err
	}
	col := m.client.Database(m.database).Collection(kind)
	if _, err := col.UpdateOne(ctx, makePrimaryKeyFilter(id), wrapper); err != nil {
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
	col := m.client.Database(m.database).Collection(kind)
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
	wrapper, err := wrapModel(entity)
	if err != nil {
		return err
	}
	update := bson.D{{"$set", wrapper}}
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

func (m *MongoDB) determineCredential() (*options.Credential, error) {
	switch m.authMechanism {
	case scram:
		username, err := ioutil.ReadFile(m.usernameFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read username file: %w", err)
		}
		password, err := ioutil.ReadFile(m.passwordFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read password file: %w", err)
		}

		return &options.Credential{
			Username:   strings.TrimRight(string(username), "\n"),
			Password:   strings.TrimRight(string(password), "\n"),
			AuthSource: m.database,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported %q authMechanism credential", m.authMechanism)
	}
}
