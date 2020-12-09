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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/model"
	"github.com/pipe-cd/pipe/pkg/semver"
)

// The maximum number of image results returned by the APIs about images.
// This value can be between 1 and 1000.
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

func NewECR(name string, region string, opts ...Option) (*ECR, error) {
	if region != "" {
		return nil, fmt.Errorf("region is required")
	}
	e := &ECR{
		name:   name,
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(e)
	}
	e.logger = e.logger.Named("ecr-provider")

	cfg := aws.NewConfig().WithRegion(e.region)
	if e.credentialsFile != "" {
		cfg = cfg.WithCredentials(credentials.NewSharedCredentials(e.credentialsFile, e.profile))
	} else {
		cfg = cfg.WithCredentials(credentials.NewEnvCredentials())
	}
	sess := session.Must(session.NewSession())
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
	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(image.Repo),
		Filter:         &ecr.ListImagesFilter{TagStatus: aws.String("TAGGED")},
		MaxResults:     aws.Int64(maxResults),
	}
	if e.registryID != "" {
		input.RegistryId = &e.registryID
	}

	// TODO: Consider the way to determine the latest tag other than fetching all tags
	// Iterate over the pages of a ListImages operation until the last page.
	// NOTE: A lot of requests may be issued if there are a lot of tags.
	// For instance, for 6k tags, it will issue 6 requests.
	imageIDs := make([]*ecr.ImageIdentifier, maxResults)
	err := e.client.ListImagesPagesWithContext(ctx, input, func(page *ecr.ListImagesOutput, lastPage bool) bool {
		imageIDs = append(imageIDs, page.ImageIds...)
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
			}
		}
		return nil, fmt.Errorf("unknow error given: %w", err)
	}
	if len(imageIDs) == 0 {
		return nil, fmt.Errorf("no ids found")
	}

	// To avoid reaching the API rate limit, determine by the semantic versioning as much as possible.
	latestTag, err := latestBySemver(imageIDs)
	if err != nil {
		e.logger.Info("it will try to determine the latest tag by the PushedAt due to the failure by semver")
		latestTag, err = e.latestByPushedAt(ctx, image.Repo, imageIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to determine the latest tag: %w", err)
		}
	}
	return &model.ImageRef{
		ImageName: *image,
		Tag:       latestTag,
	}, nil
}

// latestByPushedAt determines the latest tag by comparing the time pushed at.
// It first issues a request to the DescribeImages API to fetch the images' PushedAt.
func (e *ECR) latestByPushedAt(ctx context.Context, repo string, ids []*ecr.ImageIdentifier) (string, error) {
	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repo),
		Filter:         &ecr.DescribeImagesFilter{TagStatus: aws.String("TAGGED")},
		MaxResults:     aws.Int64(maxResults),
		ImageIds:       ids,
	}
	if e.registryID != "" {
		input.RegistryId = &e.registryID
	}
	// NOTE: A lot of requests may be issued if there are a lot of tags.
	// For instance, for 6k tags, it will issue 6 requests.
	imageDetails := make([]*ecr.ImageDetail, maxResults)
	err := e.client.DescribeImagesPagesWithContext(ctx, input, func(page *ecr.DescribeImagesOutput, lastPage bool) bool {
		imageDetails = append(imageDetails, page.ImageDetails...)
		return true
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				return "", fmt.Errorf("server-side issue occured: %w", err)
			case ecr.ErrCodeInvalidParameterException:
				return "", fmt.Errorf("invalid parameter given: %w", err)
			case ecr.ErrCodeRepositoryNotFoundException:
				return "", fmt.Errorf("repository not found: %w", err)
			case ecr.ErrCodeImageNotFoundException:
				return "", fmt.Errorf("image not found: %w", err)
			}
		}
		return "", fmt.Errorf("unknow error given: %w", err)
	}
	if len(imageDetails) == 0 {
		return "", fmt.Errorf("no images found")
	}

	sort.SliceStable(imageDetails, func(i, j int) bool {
		l, r := imageDetails[i], imageDetails[j]
		if l.ImagePushedAt == nil || r.ImagePushedAt == nil {
			return l.ImagePushedAt == nil && r.ImagePushedAt != nil
		}
		return l.ImagePushedAt.After(*r.ImagePushedAt)
	})
	if len(imageDetails[0].ImageTags) == 0 {
		return "", fmt.Errorf("no images tag is associated the image")
	}
	// NOTE: Even if the tags are different, they are managed as a single
	// image if the images' sha256 digests are identical, so there may
	// be multiple tags associated with a single image.
	latest := *imageDetails[0].ImageTags[0]
	return latest, nil
}

// latestBySemver gives back the latest tag after sorting tags by semver.
// Returns an error if one of any tag couldn't be parsed.
func latestBySemver(ids []*ecr.ImageIdentifier) (string, error) {
	length := len(ids)
	if length == 0 {
		return "", nil
	}
	tags := make([]*semver.Version, 0, length)
	for _, id := range ids {
		tag, err := semver.NewVersion(*id.ImageTag)
		if err != nil {
			return "", fmt.Errorf("failed to parse the tag: %w", err)
		}
		tags = append(tags, tag)
	}

	sort.Sort(semver.ByNewer(tags))
	return tags[0].String(), nil
}
