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

package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
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

func newClient(ctx context.Context, region, profile, credentialsFile string, logger *zap.Logger) (*client, error) {
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

func (c *client) Apply(ctx context.Context) error {
	// TODO implement
	return nil
}
