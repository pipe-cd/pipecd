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

package ecs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/sts"
	"go.uber.org/zap"
)

type client struct {
	region string
	client *ecs.ECS
	logger *zap.Logger
}

func newClient(region, profile, credentialsFile, roleARN, tokenPath string, logger *zap.Logger) (*client, error) {
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

	// Piped attempts to retrieve credentials in the following order:
	// 1. from the environment variables. Available environment variables are:
	//   - AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
	//   - AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
	// 2. from the given credentials file.
	// 3. from the pod running in EKS cluster via STS (SecurityTokenService) as WebIdentityRole.
	// 4. from the EC2 Instance Role.
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{
				Filename: credentialsFile,
				Profile:  profile,
			},
			// roleSessionName specifies the IAM role session name to use when assuming a role.
			// it will be generated automatically in case of empty string passed.
			// ref: https://github.com/aws/aws-sdk-go/blob/0dd12669013412980b665d4f6e2947d57b1cd062/aws/credentials/stscreds/web_identity_provider.go#L116-L121
			stscreds.NewWebIdentityRoleProvider(sts.New(sess), roleARN, "", tokenPath),
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		},
	)
	cfg := aws.NewConfig().WithRegion(c.region).WithCredentials(creds)
	c.client = ecs.New(sess, cfg)

	return c, nil
}

func (c *client) ServiceExist(ctx context.Context, clusterName string, services []string) (bool, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(clusterName),
		Services: aws.StringSlice(services),
	}
	_, err := c.client.DescribeServicesWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				return false, fmt.Errorf("aws ecs service encountered an internal error: %w", err)
			case ecs.ErrCodeClientException:
				return false, fmt.Errorf("aws ecs service encountered an client error: %w", err)
			case ecs.ErrCodeInvalidParameterException:
				return false, fmt.Errorf("invalid parameter given: %w", err)
			case ecs.ErrCodeClusterNotFoundException:
				return false, fmt.Errorf("aws ecs cluster not found: %w", err)
			}
		}
		return false, fmt.Errorf("unknown error given: %w", err)
	}
	return true, nil
}
