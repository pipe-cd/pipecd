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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sts"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/backoff"
)

type TrafficConfigKeyName string

const (
	defaultAliasName = "Service"
	// RequestRetryTime represents the number of times calling to AWS resource control.
	RequestRetryTime = 3
	// RetryIntervalDuration represents duration time between retry.
	RetryIntervalDuration = 1 * time.Minute

	// TrafficPrimaryVersionKeyName represents the key points to primary version config on traffic routing map.
	TrafficPrimaryVersionKeyName TrafficConfigKeyName = "primary"
	// TrafficSecondaryVersionKeyName represents the key points to primary version config on traffic routing map.
	TrafficSecondaryVersionKeyName TrafficConfigKeyName = "secondary"
)

// ErrNotFound lambda resource occurred.
var ErrNotFound = errors.New("lambda resource not found")

type client struct {
	region string
	client *lambda.Lambda
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

	retry := backoff.NewRetry(RequestRetryTime, backoff.NewConstant(RetryIntervalDuration))
	updateFunctionConfigurationSucceed := false
	for retry.WaitNext(ctx) {
		configInput := &lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(fm.Spec.Name),
			MemorySize:   aws.Int64(fm.Spec.Memory),
			Timeout:      aws.Int64(fm.Spec.Timeout),
			Environment: &lambda.Environment{
				Variables: aws.StringMap(fm.Spec.Environments),
			},
		}
		_, err := c.client.UpdateFunctionConfigurationWithContext(ctx, configInput)
		if err != nil {
			c.logger.Error("Failed to update function configuration")
		} else {
			updateFunctionConfigurationSucceed = true
			break
		}
	}
	if !updateFunctionConfigurationSucceed && configErr != nil {
		return fmt.Errorf("failed to update configuration for Lambda function %s: %w", fm.Spec.Name, err)
	}

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
	version = aws.StringValue(cfg.Version)
	return
}

// RoutingTrafficConfig presents a map of primary and secondary version traffic for lambda function alias.
type RoutingTrafficConfig map[TrafficConfigKeyName]VersionTraffic

func (c *RoutingTrafficConfig) Encode() (string, error) {
	out, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (c *RoutingTrafficConfig) Decode(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}

// VersionTraffic presents the version, and the percent of traffic that's routed to it.
type VersionTraffic struct {
	Version string  `json:"version"`
	Percent float64 `json:"percent"`
}

func (c *client) GetTrafficConfig(ctx context.Context, fm FunctionManifest) (routingTrafficCfg RoutingTrafficConfig, err error) {
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

	routingTrafficCfg = make(map[TrafficConfigKeyName]VersionTraffic)
	/* The current return value from GetAlias as below
	{
		"AliasArn": "arn:aws:lambda:ap-northeast-1:769161735124:function:SimpleCanaryFunction:Service",
		"Name": "Service",
		"FunctionVersion": "1",
		"Description": "",
		"RoutingConfig": {
			"AdditionalVersionWeights": {
				"3": 0.9
			}
		},
		"RevisionId": "fe08805f-9851-44fc-9a79-6e086aefc290"
	}
	Note:
	- In case RoutingConfig is nil, this mean 100% of traffic is handled by version represented by FunctionVersion value (PRIMARY version).
	- In case RoutingConfig is not nil, RoutingConfig.AdditionalVersionWeights is expected to have ONLY ONE key/value pair
	which presents the SECONDARY version handling traffic (represented by the value of the pair).
		in short
			_ version: 1 - FunctionVersion (the PRIMARY) handles (1 - 0.9) percentage of current traffic.
			_ version: 3 - AdditionalVersionWeights key (the SECONDARY) handles 0.9 percentage of current traffic.
	*/
	// In case RoutingConfig is nil, 100 percent of current traffic is handled by FunctionVersion version.
	if cfg.RoutingConfig == nil {
		routingTrafficCfg[TrafficPrimaryVersionKeyName] = VersionTraffic{
			Version: aws.StringValue(cfg.FunctionVersion),
			Percent: 100,
		}
		return
	}
	// In case RoutingConfig is provided, FunctionVersion value represents the primary version while
	// RoutingConfig.AdditionalVersionWeights key represents the secondary version.
	var secondaryVersionTraffic float64
	for version, weight := range cfg.RoutingConfig.AdditionalVersionWeights {
		secondaryVersionTraffic = percentageToPercent(aws.Float64Value(weight))
		routingTrafficCfg[TrafficSecondaryVersionKeyName] = VersionTraffic{
			Version: version,
			Percent: secondaryVersionTraffic,
		}
	}
	routingTrafficCfg[TrafficPrimaryVersionKeyName] = VersionTraffic{
		Version: aws.StringValue(cfg.FunctionVersion),
		Percent: 100 - secondaryVersionTraffic,
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

func (c *client) UpdateTrafficConfig(ctx context.Context, fm FunctionManifest, routingTraffic RoutingTrafficConfig) error {
	primary, ok := routingTraffic[TrafficPrimaryVersionKeyName]
	if !ok {
		return fmt.Errorf("invalid routing traffic configuration given: primary version not found")
	}

	input := &lambda.UpdateAliasInput{
		FunctionName:    aws.String(fm.Spec.Name),
		Name:            aws.String(defaultAliasName),
		FunctionVersion: aws.String(primary.Version),
	}

	if secondary, ok := routingTraffic[TrafficSecondaryVersionKeyName]; ok {
		routingTrafficMap := make(map[string]*float64)
		routingTrafficMap[secondary.Version] = aws.Float64(precentToPercentage(secondary.Percent))
		input.RoutingConfig = &lambda.AliasRoutingConfiguration{
			AdditionalVersionWeights: routingTrafficMap,
		}
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

func percentageToPercent(in float64) float64 {
	return in * 100
}
