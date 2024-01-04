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

package gcs

import (
	"context"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/pipe-cd/pipecd/pkg/filestore"
)

type Store struct {
	client          *storage.Client
	bucket          string
	credentialsFile string
	httpClient      *http.Client
	logger          *zap.Logger
}

type Option func(*Store)

func WithCredentialsFile(path string) Option {
	return func(s *Store) {
		s.credentialsFile = path
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(s *Store) {
		s.httpClient = client
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *Store) {
		s.logger = logger.Named("gcs")
	}
}

func NewStore(ctx context.Context, bucket string, opts ...Option) (*Store, error) {
	s := &Store{
		bucket: bucket,
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(s)
	}

	var options []option.ClientOption
	if s.credentialsFile != "" {
		options = append(options, option.WithCredentialsFile(s.credentialsFile))
	}
	if s.httpClient != nil {
		options = append(options, option.WithHTTPClient(s.httpClient))
	}
	client, err := storage.NewClient(ctx, options...)
	if err != nil {
		return nil, err
	}
	s.client = client
	return s, nil
}

func (s *Store) GetReader(ctx context.Context, path string) (rc io.ReadCloser, err error) {
	rc, err = s.client.Bucket(s.bucket).Object(path).NewReader(ctx)
	switch err {
	case nil:
	case storage.ErrObjectNotExist:
		err = filestore.ErrNotFound
		return
	default:
		s.logger.Error("failed to create GCS object reader", zap.String("path", path), zap.Error(err))
		return
	}
	return
}

func (s *Store) Get(ctx context.Context, path string) ([]byte, error) {
	rc, err := s.GetReader(ctx, path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			s.logger.Error("failed to close object reader")
		}
	}()

	return io.ReadAll(rc)
}

func (s *Store) Put(ctx context.Context, path string, content []byte) error {
	wc := s.client.Bucket(s.bucket).Object(path).NewWriter(ctx)
	if _, err := wc.Write(content); err != nil {
		wc.Close()
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, path string) error {
	return s.client.Bucket(s.bucket).Object(path).Delete(ctx)
}

func (s *Store) List(ctx context.Context, prefix string) ([]filestore.ObjectAttrs, error) {
	var objects []filestore.ObjectAttrs
	query := &storage.Query{
		Prefix: prefix,
	}
	it := s.client.Bucket(s.bucket).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			s.logger.Error("failed to iterate to the next object",
				zap.String("prefix", prefix),
				zap.Error(err),
			)
			return nil, err
		}
		object := filestore.ObjectAttrs{
			Path:      attrs.Name,
			Size:      attrs.Size,
			Etag:      attrs.Etag,
			UpdatedAt: attrs.Updated.Unix(),
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func (s *Store) Close() error {
	return s.client.Close()
}
