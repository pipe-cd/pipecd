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

package gcs

import (
	"context"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/kapetaniosci/pipe/pkg/filestore"
)

type Store struct {
	client             *storage.Client
	bucket             string
	useCredentialsFile bool
	credentialsFile    string
	logger             *zap.Logger
}

type Option func(*Store)

func WithCredentialsFile(path string) Option {
	return func(s *Store) {
		s.useCredentialsFile = true
		s.credentialsFile = path
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
	if s.useCredentialsFile {
		options = append(options, option.WithCredentialsFile(s.credentialsFile))
	}
	client, err := storage.NewClient(ctx, options...)
	if err != nil {
		return nil, err
	}
	s.client = client
	return s, nil
}

func (s *Store) NewReader(ctx context.Context, path string) (rc io.ReadCloser, err error) {
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

func (s *Store) NewWriter(ctx context.Context, path string) io.WriteCloser {
	return s.client.Bucket(s.bucket).Object(path).NewWriter(ctx)
}

func (s *Store) GetObject(ctx context.Context, path string) (object filestore.Object, err error) {
	object.Path = path
	rc, err := s.NewReader(ctx, path)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(rc)
	if err != nil {
		rc.Close()
		return
	}
	err = rc.Close()
	if err != nil {
		return
	}
	object.Content = content
	object.Size = int64(len(content))
	return
}

func (s *Store) PutObject(ctx context.Context, path string, content []byte) error {
	wc := s.NewWriter(ctx, path)
	if _, err := wc.Write(content); err != nil {
		wc.Close()
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Store) ListObjects(ctx context.Context, prefix string) ([]filestore.Object, error) {
	var objects []filestore.Object
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
		object := filestore.Object{
			Path:    attrs.Name,
			Size:    attrs.Size,
			Content: []byte{},
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func (s *Store) Close() error {
	return s.client.Close()
}
