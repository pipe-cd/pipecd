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

package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"go.uber.org/zap"
)

type client struct {
	region string
	client *lambda.Lambda
	logger *zap.Logger
}

func newClient(region, profile, credentialsFile string, logger *zap.Logger) (*client, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required field")
	}

	c := &client{
		region: region,
		logger: logger.Named("lambda"),
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create a session: %w", err)
	}
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{
				Filename: credentialsFile,
				Profile:  profile,
			},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		},
	)
	cfg := aws.NewConfig().WithRegion(c.region).WithCredentials(creds)
	c.client = lambda.New(sess, cfg)

	return c, nil
}

func (c *client) Apply(ctx context.Context, fm FunctionManifest, role string) error {
	if role == "" {
		return fmt.Errorf("role arn is required")
	}
	input := &lambda.CreateFunctionInput{
		Code:         &lambda.FunctionCode{ImageUri: &fm.Spec.ImageURI},
		Role:         &role,
		FunctionName: &fm.Spec.Name,
		Runtime:      &fm.Spec.Runtime,
		Handler:      &fm.Spec.Handler,
	}
	_, err := c.client.CreateFunctionWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeInvalidParameterValueException:
				return fmt.Errorf("invalid parameter given: %w", err)
			case lambda.ErrCodeServiceException:
				return fmt.Errorf("aws lambda service encountered an internal error: %w", err)
			case lambda.ErrCodeCodeStorageExceededException:
				return fmt.Errorf("total code size per account exceeded: %w", err)
			case lambda.ErrCodeResourceNotFoundException:
				fallthrough
			case lambda.ErrCodeResourceNotReadyException:
				return fmt.Errorf("resource error occurred: %w", err)
			}
		}
		return fmt.Errorf("unknown error given: %w", err)
	}
	return nil
}
