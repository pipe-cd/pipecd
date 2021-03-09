// Copyright 2021 The PipeCD Authors.
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

package s3v2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type Store struct {
	client          *s3.Client
	bucket          string
	profile         string
	credentialsFile string
	roleARN         string
	tokenFile       string

	logger *zap.Logger
}

type Option func(*Store)

func WithLogger(logger *zap.Logger) Option {
	return func(s *Store) {
		s.logger = logger.Named("s3")
	}
}

func WithCredentialsFile(path, profile string) Option {
	return func(s *Store) {
		s.profile = profile
		s.credentialsFile = path
	}
}

func WithTokenFile(roleARN, path string) Option {
	return func(s *Store) {
		s.roleARN = roleARN
		s.tokenFile = path
	}
}

func NewStore(ctx context.Context, region, bucket string, opts ...Option) (*Store, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required field")
	}
	if bucket == "" {
		return nil, fmt.Errorf("bucket is required field")
	}

	s := &Store{
		bucket: bucket,
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(s)
	}

	// sess, err := session.NewSession()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create a session: %w", err)
	// }
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load config to create s3 client: %w", err)
	}

	// creds := credentials.NewChainCredentials(
	// 	[]credentials.Provider{
	// 		&credentials.EnvProvider{},
	// 		&credentials.SharedCredentialsProvider{
	// 			Filename: s.credentialsFile,
	// 			Profile:  s.profile,
	// 		},
	// 		// roleSessionName specifies the IAM role session name to use when assuming a role.
	// 		// it will be generated automatically in case of empty string passed.
	// 		// ref: https://github.com/aws/aws-sdk-go/blob/0dd12669013412980b665d4f6e2947d57b1cd062/aws/credentials/stscreds/web_identity_provider.go#L116-L121
	// 		stscreds.NewWebIdentityRoleProvider(sts.New(sess), s.roleARN, "", s.tokenFile),
	// 		&ec2rolecreds.EC2RoleProvider{
	// 			Client: ec2metadata.New(sess),
	// 		},
	// 	},
	// )
	// cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)
	s.client = s3.NewFromConfig(cfg)

	return s, nil
}

// func (s *Store) NewReader(ctx context.Context, path string) (rc io.ReadCloser, err error) {
// 	input := &s3.GetObjectInput{
// 		Bucket: aws.String(s.bucket),
// 		Key:    aws.String(path),
// 	}
// 	out, err := s.client.GetObjectWithContext(ctx, input)
// 	if err != nil {
// 		if aerr, ok := err.(awserr.Error); ok {
// 			switch aerr.Code() {
// 			case s3.ErrCodeNoSuchKey:
// 				err = filestore.ErrNotFound
// 				return
// 			case s3.ErrCodeInvalidObjectState:
// 				err = fmt.Errorf("invalid object state: %w", err)
// 				return
// 			default:
// 				err = fmt.Errorf("unexpected aws error given: %w", err)
// 				return
// 			}
// 		}
// 		err = fmt.Errorf("unknown error given: %w", err)
// 		return
// 	}
// 	return out.Body, nil
// }

// func (s *Store) GetObject(ctx context.Context, path string) (object filestore.Object, err error) {
// 	rc, err := s.NewReader(ctx, path)
// 	if err != nil {
// 		return
// 	}
// 	defer func() {
// 		if err := rc.Close(); err != nil {
// 			s.logger.Error("failed to close object reader")
// 		}
// 	}()

// 	content, err := ioutil.ReadAll(rc)
// 	if err != nil {
// 		return
// 	}

// 	object.Path = path
// 	object.Content = content
// 	object.Size = int64(len(content))
// 	return
// }

// func (s *Store) PutObject(ctx context.Context, path string, content []byte) error {
// 	input := &s3.PutObjectInput{
// 		Body:   aws.ReadSeekCloser(bytes.NewReader(content)),
// 		Bucket: aws.String(s.bucket),
// 		Key:    aws.String(path),
// 	}
// 	_, err := s.client.PutObjectWithContext(ctx, input)
// 	if err != nil {
// 		if aerr, ok := err.(awserr.Error); ok {
// 			switch aerr.Code() {
// 			default:
// 				return fmt.Errorf("error occured on aws side: %w", err)
// 			}
// 		}
// 		return fmt.Errorf("unknown error given: %w", err)
// 	}
// 	return nil
// }

// func (s *Store) ListObjects(ctx context.Context, prefix string) ([]filestore.Object, error) {
// 	var objects []filestore.Object
// 	input := &s3.ListObjectsV2Input{
// 		Bucket: aws.String(s.bucket),
// 		Prefix: aws.String(prefix),
// 	}

// 	err := s.client.ListObjectsV2PagesWithContext(ctx, input,
// 		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
// 			for _, obj := range page.Contents {
// 				objects = append(objects, filestore.Object{
// 					Path:    aws.StringValue(obj.Key),
// 					Size:    aws.Int64Value(obj.Size),
// 					Content: []byte{},
// 				})
// 			}
// 			return *page.IsTruncated
// 		},
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get list objects: %w", err)
// 	}
// 	return objects, nil
// }

// func (s *Store) Close() error {
// 	// aws client does not provide the way to close a connection via sdk
// 	return nil
// }
