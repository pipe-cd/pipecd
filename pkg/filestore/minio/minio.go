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
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/url"
	"path/filepath"
	"strings"

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

// NewStore generates a minio client with the given params. Specify makeBucket, then it makes the given bucket if not exists.
func NewStore(ctx context.Context, endpoint, bucket, accessKeyFile, secretKeyFile string, makeBucket bool, opts ...Option) (*Store, error) {
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
		Creds:  credentials.NewStaticV4(strings.TrimRight(string(accessKey), "\n"), strings.TrimRight(string(secretKey), "\n"), ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	s.client = client

	if makeBucket {
		exists, err := s.client.BucketExists(ctx, bucket)
		if err != nil {
			return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
		}
		if !exists {
			if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return nil, fmt.Errorf("failed to make bucket: %w", err)
			}
		}
	}

	return s, nil
}

func (s *Store) NewReader(ctx context.Context, path string) (rc io.ReadCloser, err error) {
	// No error returned even if the object does not exist.
	return s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
}

func (s *Store) GetObject(ctx context.Context, path string) (object filestore.Object, err error) {
	rc, err := s.NewReader(ctx, path)
	if err != nil {
		return
	}
	defer func() {
		if err := rc.Close(); err != nil {
			s.logger.Error("failed to close object reader")
		}
	}()

	content, err := ioutil.ReadAll(rc)
	if err != nil {
		return
	}
	object.Path = path
	object.Content = content
	object.Size = int64(len(content))
	return
}

func (s *Store) PutObject(ctx context.Context, path string, content []byte) error {
	opts := minio.PutObjectOptions{}
	if opts.ContentType = mime.TypeByExtension(filepath.Ext(path)); opts.ContentType == "" {
		opts.ContentType = "application/octet-stream"
	}
	b := bytes.NewBuffer(content)

	_, err := s.client.PutObject(ctx, s.bucket, path, b, int64(b.Len()), opts)
	return err
}

func (s *Store) ListObjects(ctx context.Context, prefix string) ([]filestore.Object, error) {
	objectCh := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true})
	objects := make([]filestore.Object, 0, len(objectCh))
	for o := range objectCh {
		if o.Err != nil {
			return nil, fmt.Errorf("invalid object %q found: %w", o.Key, o.Err)
		}
		objects = append(objects, filestore.Object{
			Path: o.Key,
			Size: o.Size,
		})
	}
	return objects, nil
}

func (s *Store) Close() error {
	// No need to close the connection. Minio server automatically cleans
	// idle connections and properly gives back resources to kernel.
	return nil
}
