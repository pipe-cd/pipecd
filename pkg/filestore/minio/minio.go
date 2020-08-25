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

package minio

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/filestore"
)

type Store struct {
	client *minio.Client
	bucket string

	logger *zap.Logger
}

type Option func(*Store)

func WithLogger(logger *zap.Logger) Option {
	return func(s *Store) {
		s.logger = logger.Named("minio")
	}
}

func NewStore(endpoint, bucket, accessKeyFile, secretKeyFile string, opts ...Option) (*Store, error) {
	s := &Store{
		bucket: bucket,
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(s)
	}

	var useSSL bool
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the given endpoint: %w", err)
	}
	if u.Scheme == "https" {
		useSSL = true
	}

	accessKey, err := ioutil.ReadFile(accessKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read access key file: %w", err)
	}
	secretKey, err := ioutil.ReadFile(secretKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret key file: %w", err)
	}
	client, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(string(accessKey), string(secretKey), ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	s.client = client
	return s, nil
}

func (s *Store) NewReader(ctx context.Context, path string) (rc io.ReadCloser, err error) {
	return
}

func (s *Store) NewWriter(ctx context.Context, path string) io.WriteCloser {
	return nil
}

func (s *Store) GetObject(ctx context.Context, path string) (object filestore.Object, err error) {
	return
}

func (s *Store) PutObject(ctx context.Context, path string, content []byte) error {
	return nil
}

func (s *Store) ListObjects(ctx context.Context, prefix string) ([]filestore.Object, error) {
	return nil, nil
}

func (s *Store) Close() error {
	return nil
}
