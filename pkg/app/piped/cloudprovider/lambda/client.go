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
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"go.uber.org/zap"
)

const (
	defaultAliasName = "Service"
	// RequestRetryTime represents the number of times calling to AWS resource control.
	RequestRetryTime = 3
	// RetryIntervalDuration represents duration time between retry.
	RetryIntervalDuration = 1 * time.Minute
)

// ErrNotFound lambda resource occurred.
var ErrNotFound = errors.New("lambda resource not found")

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

func (c *client) IsFunctionExist(ctx context.Context, name string) (bool, error) {
	input := &lambda.GetFunctionInput{
		FunctionName: aws.String(name),
	}
	_, err := c.client.GetFunctionWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeInvalidParameterValueException:
				return false, fmt.Errorf("invalid parameter given: %w", err)
			case lambda.ErrCodeServiceException:
				return false, fmt.Errorf("aws lambda service encountered an internal error: %w", err)
			case lambda.ErrCodeTooManyRequestsException:
				return false, fmt.Errorf("request throughput limit was exceeded: %w", err)
			// Only in case ResourceNotFound error occurred, the FunctionName is available for create so do not raise error.
			case lambda.ErrCodeResourceNotFoundException:
				return false, nil
			}
		}
		return false, fmt.Errorf("unknown error given: %w", err)
	}
	return true, nil
}

func (c *client) CreateFunction(ctx context.Context, fm FunctionManifest) error {
	input := &lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ImageUri: aws.String(fm.Spec.ImageURI),
		},
		PackageType:  aws.String("Image"),
		Role:         aws.String(fm.Spec.Role),
		FunctionName: aws.String(fm.Spec.Name),
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
			case lambda.ErrCodeResourceNotFoundException, lambda.ErrCodeResourceNotReadyException:
				return fmt.Errorf("resource error occurred: %w", err)
			case lambda.ErrCodeTooManyRequestsException:
				return fmt.Errorf("request throughput limit was exceeded: %w", err)
			}
		}
		return fmt.Errorf("unknown error given: %w", err)
	}
	return nil
}

func (c *client) UpdateFunction(ctx context.Context, fm FunctionManifest) error {
	codeInput := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(fm.Spec.Name),
		ImageUri:     aws.String(fm.Spec.ImageURI),
	}
	_, err := c.client.UpdateFunctionCodeWithContext(ctx, codeInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeInvalidParameterValueException:
				return fmt.Errorf("invalid parameter given: %w", err)
			case lambda.ErrCodeServiceException:
				return fmt.Errorf("aws lambda service encountered an internal error: %w", err)
			case lambda.ErrCodeCodeStorageExceededException:
				return fmt.Errorf("total code size per account exceeded: %w", err)
			case lambda.ErrCodeTooManyRequestsException:
				return fmt.Errorf("request throughput limit was exceeded: %w", err)
			case lambda.ErrCodeResourceConflictException:
				return fmt.Errorf("resource already existed or in progress: %w", err)
			}
		}
		return fmt.Errorf("unknown error given: %w", err)
	}

	// TODO: Support more configurable fields using Lambda.UpdateFunctionConfiguration.
	// https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#UpdateFunctionConfiguration

	return nil
}

func (c *client) PublishFunction(ctx context.Context, fm FunctionManifest) (version string, err error) {
	input := &lambda.PublishVersionInput{
		FunctionName: aws.String(fm.Spec.Name),
	}
	cfg, err := c.client.PublishVersionWithContext(ctx, input)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok {
			err = fmt.Errorf("unknown error given: %w", err)
			return
		}
		switch aerr.Code() {
		case lambda.ErrCodeInvalidParameterValueException:
			err = fmt.Errorf("invalid parameter given: %w", err)
		case lambda.ErrCodeServiceException:
			err = fmt.Errorf("aws lambda service encountered an internal error: %w", err)
		case lambda.ErrCodeTooManyRequestsException:
			err = fmt.Errorf("request throughput limit was exceeded: %w", err)
		case lambda.ErrCodeCodeStorageExceededException:
			err = fmt.Errorf("total code size per account exceeded: %w", err)
		case lambda.ErrCodeResourceNotFoundException:
			err = fmt.Errorf("resource not found: %w", err)
		case lambda.ErrCodeResourceConflictException:
			err = fmt.Errorf("resource already existed or in progress: %w", err)
		}
		return
	}
	version = *cfg.Version
	return
}

// VersionTraffic presents the version, and the percent of traffic that's routed to it.
type VersionTraffic struct {
	Version string
	Percent float64
}

func (c *client) GetTrafficConfig(ctx context.Context, fm FunctionManifest) (routingTrafficCfg []VersionTraffic, err error) {
	input := &lambda.GetAliasInput{
		FunctionName: aws.String(fm.Spec.Name),
		Name:         aws.String(defaultAliasName),
	}

	cfg, err := c.client.GetAliasWithContext(ctx, input)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok {
			err = fmt.Errorf("unknown error given: %w", err)
			return
		}
		switch aerr.Code() {
		case lambda.ErrCodeInvalidParameterValueException:
			err = fmt.Errorf("invalid parameter given: %w", err)
		case lambda.ErrCodeServiceException:
			err = fmt.Errorf("aws lambda service encountered an internal error: %w", err)
		case lambda.ErrCodeTooManyRequestsException:
			err = fmt.Errorf("request throughput limit was exceeded: %w", err)
		case lambda.ErrCodeResourceNotFoundException:
			err = ErrNotFound
		}
		return
	}

	// TODO: Fix Lambda.AliasConfiguration.RoutingConfig nil value.
	if cfg.RoutingConfig == nil {
		return
	}
	routingTrafficCfg = make([]VersionTraffic, 0, len(cfg.RoutingConfig.AdditionalVersionWeights))
	for version := range cfg.RoutingConfig.AdditionalVersionWeights {
		routingTrafficCfg = append(routingTrafficCfg, VersionTraffic{
			Version: version,
			Percent: *cfg.RoutingConfig.AdditionalVersionWeights[version],
		})
	}
	return
}

func (c *client) CreateTrafficConfig(ctx context.Context, fm FunctionManifest, version string) error {
	input := &lambda.CreateAliasInput{
		FunctionName:    aws.String(fm.Spec.Name),
		FunctionVersion: aws.String(version),
		Name:            aws.String(defaultAliasName),
	}
	_, err := c.client.CreateAliasWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeInvalidParameterValueException:
				return fmt.Errorf("invalid parameter given: %w", err)
			case lambda.ErrCodeServiceException:
				return fmt.Errorf("aws lambda service encountered an internal error: %w", err)
			case lambda.ErrCodeTooManyRequestsException:
				return fmt.Errorf("request throughput limit was exceeded: %w", err)
			case lambda.ErrCodeResourceNotFoundException:
				return fmt.Errorf("resource not found: %w", err)
			case lambda.ErrCodeResourceConflictException:
				return fmt.Errorf("resource already existed or in progress: %w", err)
			}
		}
		return fmt.Errorf("unknown error given: %w", err)
	}
	return nil
}

func (c *client) UpdateTrafficConfig(ctx context.Context, fm FunctionManifest, routingTraffic []VersionTraffic) error {
	routingTrafficMap := make(map[string]*float64)
	switch len(routingTraffic) {
	case 2:
		routingTrafficMap[routingTraffic[0].Version] = aws.Float64(precentToPercentage(routingTraffic[0].Percent))
		routingTrafficMap[routingTraffic[1].Version] = aws.Float64(precentToPercentage(routingTraffic[1].Percent))
	case 1:
		routingTrafficMap[routingTraffic[0].Version] = aws.Float64(precentToPercentage(routingTraffic[0].Percent))
	default:
		return fmt.Errorf("invalid routing traffic configuration given")
	}

	input := &lambda.UpdateAliasInput{
		FunctionName: aws.String(fm.Spec.Name),
		Name:         aws.String(defaultAliasName),
		RoutingConfig: &lambda.AliasRoutingConfiguration{
			AdditionalVersionWeights: routingTrafficMap,
		},
	}

	_, err := c.client.UpdateAliasWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeInvalidParameterValueException:
				return fmt.Errorf("invalid parameter given: %w", err)
			case lambda.ErrCodeServiceException:
				return fmt.Errorf("aws lambda service encountered an internal error: %w", err)
			case lambda.ErrCodeTooManyRequestsException:
				return fmt.Errorf("request throughput limit was exceeded: %w", err)
			case lambda.ErrCodeResourceNotFoundException:
				return fmt.Errorf("resource not found: %w", err)
			case lambda.ErrCodeResourceConflictException:
				return fmt.Errorf("resource already existed or in progress: %w", err)
			}
		}
		return fmt.Errorf("unknown error given: %w", err)
	}
	return nil
}

func precentToPercentage(in float64) float64 {
	return in / 100.0
}
