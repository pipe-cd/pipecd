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

package ecr

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/model"
)

// The maximum number of image results returned by the APIs about images.
// The API allows this value to be between 1 and 1000.
// See more: https://pkg.go.dev/github.com/aws/aws-sdk-go/service/ecr#ListImagesInput
const maxResults = 1000

type ECR struct {
	name            string
	client          *ecr.ECR
	region          string
	credentialsFile string
	profile         string
	registryID      string

	logger *zap.Logger
}

type Option func(*ECR)

func WithRegistryID(id string) Option {
	return func(e *ECR) {
		e.registryID = id
	}
}

func WithCredentialsFile(path string) Option {
	return func(e *ECR) {
		e.credentialsFile = path
	}
}

func WithProfile(profile string) Option {
	return func(e *ECR) {
		e.profile = profile
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(e *ECR) {
		e.logger = logger
	}
}

// NewECR attempts to retrieve credentials in the following order:
// 1. from the environment variables. Available environment variables are:
//   - AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
//   - AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
// 2. from the given credentials file.
// 3. from the EC2 Instance Role
func NewECR(name string, region string, opts ...Option) (*ECR, error) {
	if region == "" {
		return nil, fmt.Errorf("region is required")
	}
	e := &ECR{
		name:   name,
		region: region,
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(e)
	}
	e.logger = e.logger.Named("ecr-provider")

	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create a session: %w", err)
	}
	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{
				Filename: e.credentialsFile,
				Profile:  e.profile,
			},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		},
	)
	cfg := aws.NewConfig().WithRegion(e.region).WithCredentials(creds)
	// TODO: Use ecrpublic package when public image given
	//   See more: https://docs.aws.amazon.com/sdk-for-go/api/service/ecrpublic
	e.client = ecr.New(sess, cfg)
	return e, nil
}

func (e *ECR) Name() string {
	return e.name
}

func (e *ECR) Type() model.ImageProviderType {
	return model.ImageProviderTypeECR
}

func (e *ECR) ParseImage(image string) (*model.ImageName, error) {
	ss := strings.SplitN(image, "/", 2)
	if len(ss) < 2 {
		return nil, fmt.Errorf("invalid image format (e.g. account-id.dkr.ecr.region.amazon.aws.com/pipecd/helloworld)")
	}
	return &model.ImageName{
		Domain: ss[0],
		Repo:   ss[1],
	}, nil
}

func (e *ECR) GetLatestImage(ctx context.Context, image *model.ImageName) (*model.ImageRef, error) {
	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(image.Repo),
		Filter:         &ecr.DescribeImagesFilter{TagStatus: aws.String("TAGGED")},
		MaxResults:     aws.Int64(maxResults),
	}
	if e.registryID != "" {
		input.RegistryId = &e.registryID
	}

	// TODO: Consider the way to determine the latest tag other than fetching all tags
	//
	// Iterate over the pages of a DescribeImages operation until the last page.
	// NOTE: A lot of requests may be issued if there are a lot of tags.
	// For instance, for 6k tags, it will issue 6 requests.
	imageDetails := make([]*ecr.ImageDetail, 0, maxResults)
	err := e.client.DescribeImagesPagesWithContext(ctx, input, func(page *ecr.DescribeImagesOutput, lastPage bool) bool {
		imageDetails = append(imageDetails, page.ImageDetails...)
		return true
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				return nil, fmt.Errorf("server-side issue occured: %w", err)
			case ecr.ErrCodeInvalidParameterException:
				return nil, fmt.Errorf("invalid parameter given: %w", err)
			case ecr.ErrCodeRepositoryNotFoundException:
				return nil, fmt.Errorf("repository not found: %w", err)
			case ecr.ErrCodeImageNotFoundException:
				return nil, fmt.Errorf("image not found: %w", err)
			}
		}
		return nil, fmt.Errorf("unknow error given: %w", err)
	}
	if len(imageDetails) == 0 {
		return nil, fmt.Errorf("no images found")
	}
	sort.Slice(imageDetails, func(i, j int) bool {
		l, r := imageDetails[i], imageDetails[j]
		if l.ImagePushedAt == nil || r.ImagePushedAt == nil {
			return l.ImagePushedAt == nil && r.ImagePushedAt != nil
		}
		return l.ImagePushedAt.After(*r.ImagePushedAt)
	})
	if len(imageDetails[0].ImageTags) == 0 {
		return nil, fmt.Errorf("no images tag is associated the image")
	}
	// NOTE: Even if the tags are different, they are managed as a single
	// image if the images' sha256 digests are identical, so there may
	// be multiple tags associated with a single image. That's why
	// an ImageDetail has multiple tags.
	latest := *imageDetails[0].ImageTags[0]
	return &model.ImageRef{
		ImageName: *image,
		Tag:       latest,
	}, nil
}
