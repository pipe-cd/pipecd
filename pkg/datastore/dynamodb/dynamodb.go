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

package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sts"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/datastore"
)

// DynamoDB is wrapper for AWS dynamoDB client
type DynamoDB struct {
	client          *dynamodb.DynamoDB
	profile         string
	credentialsFile string
	roleARN         string
	tokenFile       string
	endpoint        string

	logger *zap.Logger
}

// Option for create DynamoDB typed instance
type Option func(*DynamoDB)

// WithLogger returns logger setup function
func WithLogger(logger *zap.Logger) Option {
	return func(s *DynamoDB) {
		s.logger = logger
	}
}

// WithCredentialsFile returns credentials infor setup function
func WithCredentialsFile(profile, path string) Option {
	return func(s *DynamoDB) {
		s.profile = profile
		s.credentialsFile = path
	}
}

// WithTokenFile returns authenticate with token setup function
func WithTokenFile(roleARN, path string) Option {
	return func(s *DynamoDB) {
		s.roleARN = roleARN
		s.tokenFile = path
	}
}

// NewDynamoDB returns new DynamoDB instance
func NewDynamoDB(region, endpoint string, opts ...Option) (*DynamoDB, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required field")
	}
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required field")
	}

	s := &DynamoDB{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(s)
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create a session: %w", err)
	}
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{
				Filename: s.credentialsFile,
				Profile:  s.profile,
			},
			// roleSessionName specifies the IAM role session name to use when assuming a role.
			// it will be generated automatically in case of empty string passed.
			// ref: https://github.com/aws/aws-sdk-go/blob/0dd12669013412980b665d4f6e2947d57b1cd062/aws/credentials/stscreds/web_identity_provider.go#L116-L121
			stscreds.NewWebIdentityRoleProvider(sts.New(sess), s.roleARN, "", s.tokenFile),
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		},
	)
	cfg := aws.NewConfig().WithRegion(region).WithEndpoint(endpoint).WithCredentials(creds)
	s.client = dynamodb.New(sess, cfg)

	return s, nil
}

// Find implementation for DynamoDB
func (s *DynamoDB) Find(ctx context.Context, kind string, opts datastore.ListOptions) (datastore.Iterator, error) {
	expr, err := buildDynamoDBExpression(opts)
	if err != nil {
		s.logger.Error("failed to build query",
			zap.String("kind", kind),
			zap.Error(err),
		)
		return nil, err
	}
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(kind),
	}
	// TODO: Support pagination with ListOptions.Page and ListOptions.PageSize for DynamoDB.
	// TODO: Support ordering with ListOptions.Orders for DynamoDB.
	var items []map[string]*dynamodb.AttributeValue
	err = s.client.ScanPagesWithContext(ctx, input,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			items = append(items, page.Items...)
			// The only way to know when you have reached
			// the end of the result set is when LastEvaluatedKey is empty.
			return len(page.LastEvaluatedKey) != 0
		},
	)
	if err != nil {
		s.logger.Error("failed to get cursor",
			zap.String("kind", kind),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get cursor: %w", err)
	}
	return &Iterator{
		datapool: items,
	}, nil
}

// Get implementation for DynamoDB
func (s *DynamoDB) Get(ctx context.Context, kind, id string, v interface{}) error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(kind),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}
	result, err := s.client.GetItemWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return datastore.ErrNotFound
			default:
				s.logger.Error("failed to retrieve entity: aws error",
					zap.String("id", id),
					zap.String("kind", kind),
					zap.Error(err),
				)
				return err
			}
		}
		s.logger.Error("failed to retrieve entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	return dynamodbattribute.UnmarshalMap(result.Item, v)
}

// Create implementation for DynamoDB
func (s *DynamoDB) Create(ctx context.Context, kind, id string, entity interface{}) error {
	av, err := dynamodbattribute.MarshalMap(entity)
	if err != nil {
		s.logger.Error("failed to marshal entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	// Only create in case the data item which contains key with specificed value not existed.
	expr, err := buildDynamoDBKeyNotExistedExpression("Id", id)
	if err != nil {
		s.logger.Error("failed to build condition expresion for create",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName:                 aws.String(kind),
		Item:                      av,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}
	_, err = s.client.PutItemWithContext(ctx, input)
	if err != nil {
		s.logger.Error("failed to insert entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// Put implementation for DynamoDB
func (s *DynamoDB) Put(ctx context.Context, kind, id string, entity interface{}) error {
	av, err := dynamodbattribute.MarshalMap(entity)
	if err != nil {
		s.logger.Error("failed to marshal entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(kind),
		Item:      av,
	}
	_, err = s.client.PutItemWithContext(ctx, input)
	if err != nil {
		s.logger.Error("failed to put entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	return nil
}

// Update implementation for DynamoDB
func (s *DynamoDB) Update(ctx context.Context, kind, id string, factory datastore.Factory, updater datastore.Updater) error {
	// Get existing data item.
	entity := factory()
	err := s.Get(ctx, kind, id, entity)
	if errors.Is(err, datastore.ErrNotFound) {
		return datastore.ErrNotFound
	}
	if err != nil {
		s.logger.Error("failed to retrieve entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	// Update entity by updater.
	if err := updater(entity); err != nil {
		s.logger.Error("failed to update entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	// Put back the updated data to database.
	av, err := dynamodbattribute.MarshalMap(entity)
	if err != nil {
		s.logger.Error("failed to marshal entity",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}
	// Only update in case the data item which contains key with specificed value existed.
	expr, err := buildDynamoDBKeyExistedExpression("Id", id)
	if err != nil {
		s.logger.Error("failed to build condition expresion for update",
			zap.String("id", id),
			zap.String("kind", kind),
			zap.Error(err),
		)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName:                 aws.String(kind),
		Item:                      av,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}
	_, err = s.client.PutItemWithContext(ctx, input)
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

// Close implementation for DynamoDB
func (s *DynamoDB) Close() error {
	// Connection is initialized on use, so we could not close
	// Besides, AWS session will be handled by AWS itself, it could be reused by others or cleaned by AWS
	return nil
}
