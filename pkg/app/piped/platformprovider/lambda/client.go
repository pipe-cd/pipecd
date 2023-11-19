// Copyright 2023 The PipeCD Authors.
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
	"io"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipecd/pkg/backoff"
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
	client *lambda.Client
	logger *zap.Logger
}

func newClient(region, profile, credentialsFile, roleARN, tokenPath string, logger *zap.Logger) (*client, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required field")
	}

	c := &client{
		logger: logger.Named("lambda"),
	}

	optFns := []func(*config.LoadOptions) error{config.WithRegion(region)}
	if credentialsFile != "" {
		optFns = append(optFns, config.WithSharedCredentialsFiles([]string{credentialsFile}))
	}
	if profile != "" {
		optFns = append(optFns, config.WithSharedConfigProfile(profile))
	}
	if tokenPath != "" && roleARN != "" {
		optFns = append(optFns, config.WithWebIdentityRoleCredentialOptions(func(v *stscreds.WebIdentityRoleOptions) {
			v.RoleARN = roleARN
			v.TokenRetriever = stscreds.IdentityTokenFile(tokenPath)
		}))
	}

	// When you initialize an aws.Config instance using config.LoadDefaultConfig, the SDK uses its default credential chain to find AWS credentials.
	// This default credential chain looks for credentials in the following order:
	//
	// 1. Environment variables.
	//   1. Static Credentials (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN)
	//   2. Web Identity Token (AWS_WEB_IDENTITY_TOKEN_FILE)
	// 2. Shared configuration files.
	//   1. SDK defaults to credentials file under .aws folder that is placed in the home folder on your computer.
	//   2. SDK defaults to config file under .aws folder that is placed in the home folder on your computer.
	// 3. If your application uses an ECS task definition or RunTask API operation, IAM role for tasks.
	// 4. If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.
	// ref: https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials
	cfg, err := config.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load config to create lambda client: %w", err)
	}
	c.client = lambda.NewFromConfig(cfg)

	return c, nil
}

func (c *client) IsFunctionExist(ctx context.Context, name string) (bool, error) {
	input := &lambda.GetFunctionInput{
		FunctionName: aws.String(name),
	}
	_, err := c.client.GetFunction(ctx, input)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			// Only in case ResourceNotFound error occurred, the FunctionName is available for create so do not raise error.
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *client) CreateFunction(ctx context.Context, fm FunctionManifest) error {
	input := &lambda.CreateFunctionInput{
		FunctionName: aws.String(fm.Spec.Name),
		Role:         aws.String(fm.Spec.Role),
		MemorySize:   aws.Int32(fm.Spec.Memory),
		Timeout:      aws.Int32(fm.Spec.Timeout),
		Tags:         fm.Spec.Tags,
		Environment: &types.Environment{
			Variables: fm.Spec.Environments,
		},
	}
	if len(fm.Spec.Architectures) != 0 {
		var architectures []types.Architecture
		for _, arch := range fm.Spec.Architectures {
			architectures = append(architectures, types.Architecture(arch.Name))
		}
		input.Architectures = architectures
	}
	if fm.Spec.EphemeralStorage.Size != 0 {
		input.EphemeralStorage.Size = aws.Int32(fm.Spec.EphemeralStorage.Size)
	}
	if fm.Spec.VPCConfig != nil {
		input.VpcConfig = &types.VpcConfig{
			SecurityGroupIds: fm.Spec.VPCConfig.SecurityGroupIDs,
			SubnetIds:        fm.Spec.VPCConfig.SubnetIDs,
		}
	}
	// Container image packing.
	if fm.Spec.ImageURI != "" {
		input.PackageType = types.PackageTypeImage
		input.Code = &types.FunctionCode{
			ImageUri: aws.String(fm.Spec.ImageURI),
		}
	}
	// Zip packing which stored in s3.
	if fm.Spec.S3Bucket != "" {
		input.PackageType = types.PackageTypeZip
		input.Code = &types.FunctionCode{
			S3Bucket:        aws.String(fm.Spec.S3Bucket),
			S3Key:           aws.String(fm.Spec.S3Key),
			S3ObjectVersion: aws.String(fm.Spec.S3ObjectVersion),
		}
		input.Handler = aws.String(fm.Spec.Handler)
		input.Runtime = types.Runtime(fm.Spec.Runtime)
	}
	_, err := c.client.CreateFunction(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create Lambda function %s: %w", fm.Spec.Name, err)
	}
	return nil
}

func (c *client) CreateFunctionFromSource(ctx context.Context, fm FunctionManifest, zip io.Reader) error {
	data, err := io.ReadAll(zip)
	if err != nil {
		return err
	}

	input := &lambda.CreateFunctionInput{
		FunctionName: aws.String(fm.Spec.Name),
		PackageType:  types.PackageTypeZip,
		Code: &types.FunctionCode{
			ZipFile: data,
		},
		Handler:    aws.String(fm.Spec.Handler),
		Runtime:    types.Runtime(fm.Spec.Runtime),
		Role:       aws.String(fm.Spec.Role),
		MemorySize: aws.Int32(fm.Spec.Memory),
		Timeout:    aws.Int32(fm.Spec.Timeout),
		Tags:       fm.Spec.Tags,
		Environment: &types.Environment{
			Variables: fm.Spec.Environments,
		},
	}

	_, err = c.client.CreateFunction(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create Lambda function %s: %w", fm.Spec.Name, err)
	}
	return nil
}

func (c *client) UpdateFunction(ctx context.Context, fm FunctionManifest) error {
	// UpdateFunctionConfiguration must be called before UpdateFunctionCode.
	// Lambda has named by state.
	// If Lambda's state is pending, UpdateFunctionConfiguration is failed. This error is explained as a ResourceConflictException.
	// ref: https://docs.aws.amazon.com/lambda/latest/dg/troubleshooting-invocation.html
	// Update function configuration.
	if err := c.updateFunctionConfiguration(ctx, fm); err != nil {
		return err
	}
	// Update function code.
	codeInput := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(fm.Spec.Name),
	}
	// Container image packing.
	if fm.Spec.ImageURI != "" {
		codeInput.ImageUri = aws.String(fm.Spec.ImageURI)
	}
	// Zip packing which stored in s3.
	if fm.Spec.S3Bucket != "" {
		codeInput.S3Bucket = aws.String(fm.Spec.S3Bucket)
		codeInput.S3Key = aws.String(fm.Spec.S3Key)
		codeInput.S3ObjectVersion = aws.String(fm.Spec.S3ObjectVersion)
	}
	if len(fm.Spec.Architectures) != 0 {
		var architectures []types.Architecture
		for _, arch := range fm.Spec.Architectures {
			architectures = append(architectures, types.Architecture(arch.Name))
		}
		codeInput.Architectures = architectures
	}
	_, err := c.client.UpdateFunctionCode(ctx, codeInput)
	if err != nil {
		return fmt.Errorf("failed to update function code for Lambda function %s: %w", fm.Spec.Name, err)
	}

	// Tag/Untag function if necessary.
	return c.updateTagsConfig(ctx, fm)
}

func (c *client) UpdateFunctionFromSource(ctx context.Context, fm FunctionManifest, zip io.Reader) error {
	// Update function configuration.
	if err := c.updateFunctionConfiguration(ctx, fm); err != nil {
		return err
	}

	data, err := io.ReadAll(zip)
	if err != nil {
		return err
	}

	// Update function code.
	codeInput := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(fm.Spec.Name),
		ZipFile:      data,
	}
	if len(fm.Spec.Architectures) != 0 {
		var architectures []types.Architecture
		for _, arch := range fm.Spec.Architectures {
			architectures = append(architectures, types.Architecture(arch.Name))
		}
		codeInput.Architectures = architectures
	}
	_, err = c.client.UpdateFunctionCode(ctx, codeInput)
	if err != nil {
		return fmt.Errorf("failed to update function code for Lambda function %s: %w", fm.Spec.Name, err)
	}

	// Tag/Untag function if necessary.
	return c.updateTagsConfig(ctx, fm)
}

func (c *client) updateFunctionConfiguration(ctx context.Context, fm FunctionManifest) error {
	retry := backoff.NewRetry(RequestRetryTime, backoff.NewConstant(RetryIntervalDuration))
	updateFunctionConfigurationSucceed := false
	var err error
	for retry.WaitNext(ctx) {
		configInput := &lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(fm.Spec.Name),
			Role:         aws.String(fm.Spec.Role),
			MemorySize:   aws.Int32(fm.Spec.Memory),
			Timeout:      aws.Int32(fm.Spec.Timeout),
			Runtime:      types.Runtime(fm.Spec.Runtime),
			Environment: &types.Environment{
				Variables: fm.Spec.Environments,
			},
		}
		// For zip packing Lambda function code, allow update the function handler
		// on update the function's manifest.
		if fm.Spec.Handler != "" {
			configInput.Handler = aws.String(fm.Spec.Handler)
		}
		if fm.Spec.EphemeralStorage.Size != 0 {
			configInput.EphemeralStorage.Size = aws.Int32(fm.Spec.EphemeralStorage.Size)
		}
		if fm.Spec.VPCConfig != nil {
			configInput.VpcConfig = &types.VpcConfig{
				SecurityGroupIds: fm.Spec.VPCConfig.SecurityGroupIDs,
				SubnetIds:        fm.Spec.VPCConfig.SubnetIDs,
			}
		}
		_, err = c.client.UpdateFunctionConfiguration(ctx, configInput)
		if err != nil {
			c.logger.Error("Failed to update function configuration")
		} else {
			updateFunctionConfigurationSucceed = true
			break
		}
	}
	if !updateFunctionConfigurationSucceed {
		return fmt.Errorf("failed to update configuration for Lambda function %s: %w", fm.Spec.Name, err)
	}

	// Wait until function updated successfully.
	retry = backoff.NewRetry(RequestRetryTime, backoff.NewConstant(RetryIntervalDuration))
	input := &lambda.GetFunctionInput{
		FunctionName: aws.String(fm.Spec.Name),
	}
	_, err = retry.Do(ctx, func() (any, error) {
		output, err := c.client.GetFunction(ctx, input)
		if err != nil {
			return nil, err
		}
		if output.Configuration.LastUpdateStatus != types.LastUpdateStatusSuccessful {
			return nil, fmt.Errorf("failed to update Lambda function %s, status code %v, error reason %s",
				fm.Spec.Name, output.Configuration.LastUpdateStatus, *output.Configuration.LastUpdateStatusReason)
		}
		return nil, nil
	})
	return err
}

func (c *client) PublishFunction(ctx context.Context, fm FunctionManifest) (string, error) {
	input := &lambda.PublishVersionInput{
		FunctionName: aws.String(fm.Spec.Name),
	}
	cfg, err := c.client.PublishVersion(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to publish new version for Lambda function %s: %w", fm.Spec.Name, err)
	}
	return aws.ToString(cfg.Version), nil
}

// GetTrafficConfig returns lambda provider.ErrNotFound in case remote traffic config is not existed.
func (c *client) GetTrafficConfig(ctx context.Context, fm FunctionManifest) (routingTrafficCfg RoutingTrafficConfig, err error) {
	input := &lambda.GetAliasInput{
		FunctionName: aws.String(fm.Spec.Name),
		Name:         aws.String(defaultAliasName),
	}

	cfg, err := c.client.GetAlias(ctx, input)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
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
			Version: aws.ToString(cfg.FunctionVersion),
			Percent: 100,
		}
		return
	}
	// In case RoutingConfig is provided, FunctionVersion value represents the primary version while
	// RoutingConfig.AdditionalVersionWeights key represents the secondary version.
	var secondaryVersionTraffic float64
	for version, weight := range cfg.RoutingConfig.AdditionalVersionWeights {
		secondaryVersionTraffic = percentageToPercent(weight)
		routingTrafficCfg[TrafficSecondaryVersionKeyName] = VersionTraffic{
			Version: version,
			Percent: secondaryVersionTraffic,
		}
	}
	routingTrafficCfg[TrafficPrimaryVersionKeyName] = VersionTraffic{
		Version: aws.ToString(cfg.FunctionVersion),
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
	_, err := c.client.CreateAlias(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create traffic config for Lambda function %s: %w", fm.Spec.Name, err)
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
		routingTrafficMap := make(map[string]float64)
		routingTrafficMap[secondary.Version] = precentToPercentage(secondary.Percent)
		input.RoutingConfig = &types.AliasRoutingConfiguration{
			AdditionalVersionWeights: routingTrafficMap,
		}
	}

	_, err := c.client.UpdateAlias(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update traffic config for Lambda function %s: %w", fm.Spec.Name, err)
	}
	return nil
}

func (c *client) updateTagsConfig(ctx context.Context, fm FunctionManifest) error {
	getFuncInput := &lambda.GetFunctionInput{
		FunctionName: aws.String(fm.Spec.Name),
	}
	output, err := c.client.GetFunction(ctx, getFuncInput)
	if err != nil {
		return fmt.Errorf("error occurred on list tags of Lambda function %s: %w", fm.Spec.Name, err)
	}

	functionArn := aws.ToString(output.Configuration.FunctionArn)
	currentTags := output.Tags
	// Skip if there are no changes on tags.
	if reflect.DeepEqual(currentTags, fm.Spec.Tags) {
		return nil
	}

	newDefinedTags, updatedTags, removedTags := makeFlowControlTagsMaps(currentTags, fm.Spec.Tags)

	if len(newDefinedTags) > 0 {
		if err := c.tagFunction(ctx, functionArn, newDefinedTags); err != nil {
			return fmt.Errorf("failed on add new defined tags to Lambda function %s: %w", fm.Spec.Name, err)
		}
	}

	if len(updatedTags) > 0 {
		if err := c.untagFunction(ctx, functionArn, updatedTags); err != nil {
			return fmt.Errorf("failed on update changed tags to Lambda function %s: %w", fm.Spec.Name, err)
		}
		if err := c.tagFunction(ctx, functionArn, updatedTags); err != nil {
			return fmt.Errorf("failed on update changed tags to Lambda function %s: %w", fm.Spec.Name, err)
		}
	}

	if len(removedTags) > 0 {
		if err := c.untagFunction(ctx, functionArn, removedTags); err != nil {
			return fmt.Errorf("failed on remove tags for Lambda function %s: %w", fm.Spec.Name, err)
		}
	}

	return nil
}

func (c *client) tagFunction(ctx context.Context, functionArn string, tags map[string]string) error {
	tagInput := &lambda.TagResourceInput{
		Resource: aws.String(functionArn),
		Tags:     tags,
	}
	_, err := c.client.TagResource(ctx, tagInput)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) untagFunction(ctx context.Context, functionArn string, tags map[string]string) error {
	tagsKeys := make([]string, 0, len(tags))
	for k := range tags {
		tagsKeys = append(tagsKeys, k)
	}
	untagInput := &lambda.UntagResourceInput{
		Resource: aws.String(functionArn),
		TagKeys:  tagsKeys,
	}
	_, err := c.client.UntagResource(ctx, untagInput)
	if err != nil {
		return err
	}

	return nil
}

func makeFlowControlTagsMaps(remoteTags, definedTags map[string]string) (newDefinedTags, updatedTags, removedTags map[string]string) {
	newDefinedTags = make(map[string]string)
	updatedTags = make(map[string]string)
	removedTags = make(map[string]string)
	for k, v := range definedTags {
		val, ok := remoteTags[k]
		if !ok {
			newDefinedTags[k] = v
			continue
		}
		if val != v {
			updatedTags[k] = v
		}
	}
	for k, v := range remoteTags {
		_, ok := definedTags[k]
		if !ok {
			removedTags[k] = v
		}
	}
	return
}

func precentToPercentage(in float64) float64 {
	return in / 100.0
}

func percentageToPercent(in float64) float64 {
	return in * 100
}
