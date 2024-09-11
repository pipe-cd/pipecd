// Copyright 2024 The PipeCD Authors.
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
	"fmt"
	"os"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
)

const (
	VersionV1Beta1       = "pipecd.dev/v1beta1"
	FunctionManifestKind = "LambdaFunction"
	// Memory and Timeout lower and upper limit as noted via
	// https://docs.aws.amazon.com/sdk-for-go/api/service/lambda/#UpdateFunctionConfigurationInput
	memoryLowerLimit           = 1
	timeoutLowerLimit          = 1
	timeoutUpperLimit          = 900
	ephemeralStorageLowerLimit = 512
	ephemeralStorageUpperLimit = 10240
)

type FunctionManifest struct {
	Kind       string               `json:"kind"`
	APIVersion string               `json:"apiVersion,omitempty"`
	Spec       FunctionManifestSpec `json:"spec"`
}

func (fm *FunctionManifest) validate() error {
	if fm.APIVersion != VersionV1Beta1 {
		return fmt.Errorf("unsupported version: %s", fm.APIVersion)
	}
	if fm.Kind != FunctionManifestKind {
		return fmt.Errorf("invalid manifest kind given: %s", fm.Kind)
	}
	if err := fm.Spec.validate(); err != nil {
		return err
	}
	return nil
}

// FunctionManifestSpec contains configuration for LambdaFunction.
type FunctionManifestSpec struct {
	Name             string            `json:"name"`
	Role             string            `json:"role"`
	ImageURI         string            `json:"image"`
	S3Bucket         string            `json:"s3Bucket"`
	S3Key            string            `json:"s3Key"`
	S3ObjectVersion  string            `json:"s3ObjectVersion"`
	SourceCode       SourceCode        `json:"source"`
	Handler          string            `json:"handler"`
	Architectures    []Architecture    `json:"architectures,omitempty"`
	EphemeralStorage *EphemeralStorage `json:"ephemeralStorage,omitempty"`
	Runtime          string            `json:"runtime"`
	Memory           int32             `json:"memory"`
	Timeout          int32             `json:"timeout"`
	Tags             map[string]string `json:"tags,omitempty"`
	Environments     map[string]string `json:"environments,omitempty"`
	VPCConfig        *VPCConfig        `json:"vpcConfig,omitempty"`
	// A list of function layers used in the function.
	// Specify each layer by its ARN including the version.
	// You can use layers only with Lambda functions deployed as a .zip file archive. Layers are ignored for a container image.
	// See https://docs.aws.amazon.com/lambda/latest/dg/chapter-layers.html.
	Layers []string `json:"layers,omitempty"`
}

type VPCConfig struct {
	SecurityGroupIDs []string `json:"securityGroupIds,omitempty"`
	SubnetIDs        []string `json:"subnetIds,omitempty"`
}

func (fmp FunctionManifestSpec) validate() error {
	if fmp.Name == "" {
		return fmt.Errorf("lambda function is missing")
	}
	if fmp.ImageURI == "" && fmp.S3Bucket == "" {
		if err := fmp.SourceCode.validate(); err != nil {
			return err
		}
	}
	if fmp.ImageURI == "" {
		if fmp.Handler == "" {
			return fmt.Errorf("handler is missing")
		}
		if fmp.Runtime == "" {
			return fmt.Errorf("runtime is missing")
		}
	}
	for _, arch := range fmp.Architectures {
		if err := arch.validate(); err != nil {
			return fmt.Errorf("architecture is invalid: %w", err)
		}
	}
	if fmp.EphemeralStorage != nil {
		if err := fmp.EphemeralStorage.validate(); err != nil {
			return fmt.Errorf("ephemeral storage is invalid: %w", err)
		}
	}
	if fmp.Role == "" {
		return fmt.Errorf("role is missing")
	}
	if fmp.Memory < memoryLowerLimit {
		return fmt.Errorf("memory is missing")
	}
	if fmp.Timeout < timeoutLowerLimit || fmp.Timeout > timeoutUpperLimit {
		return fmt.Errorf("timeout is missing or out of range")
	}
	return nil
}

type SourceCode struct {
	Git  string `json:"git"`
	Ref  string `json:"ref"`
	Path string `json:"path"`
}

func (sc SourceCode) validate() error {
	if sc.Git == "" {
		return fmt.Errorf("remote git source is missing")
	}
	if sc.Ref == "" {
		return fmt.Errorf("source ref is missing")
	}
	return nil
}

type Architecture struct {
	Name string `json:"name"`
}

func (a Architecture) validate() error {
	if a.Name != "x86_64" && a.Name != "arm64" {
		return fmt.Errorf("architecture is invalid")
	}
	return nil
}

type EphemeralStorage struct {
	Size int32 `json:"size,omitempty"`
}

func (es EphemeralStorage) validate() error {
	if es.Size < ephemeralStorageLowerLimit || es.Size > ephemeralStorageUpperLimit {
		return fmt.Errorf("ephemeral storage is out of range")
	}
	return nil
}

func loadFunctionManifest(path string) (FunctionManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FunctionManifest{}, err
	}
	return parseFunctionManifest(data)
}

func parseFunctionManifest(data []byte) (FunctionManifest, error) {
	var obj FunctionManifest
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return FunctionManifest{}, err
	}
	if err := obj.validate(); err != nil {
		return FunctionManifest{}, err
	}
	return obj, nil
}

// DecideRevisionName returns revision name to apply.
func DecideRevisionName(fm FunctionManifest, commit string) (string, error) {
	tag, err := FindImageTag(fm)
	if err != nil {
		return "", err
	}
	tag = strings.ReplaceAll(tag, ".", "")

	if len(commit) > 7 {
		commit = commit[:7]
	}
	return fmt.Sprintf("%s-%s-%s", fm.Spec.Name, tag, commit), nil
}

// FindImageTag parses image tag from given LambdaFunction manifest.
func FindImageTag(fm FunctionManifest) (string, error) {
	name, tag := parseContainerImage(fm.Spec.ImageURI)
	if name == "" {
		return "", fmt.Errorf("image name could not be empty")
	}
	return tag, nil
}

func parseContainerImage(image string) (name, tag string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		tag = parts[1]
	}
	paths := strings.Split(parts[0], "/")
	name = paths[len(paths)-1]
	return
}

// FindArtifactVersions parses artifact versions from function.yaml.
func FindArtifactVersions(fm FunctionManifest) ([]*model.ArtifactVersion, error) {
	// Extract container image tag as application version.
	if fm.Spec.ImageURI != "" {
		name, tag := parseContainerImage(fm.Spec.ImageURI)
		if name == "" {
			return nil, fmt.Errorf("image name could not be empty")
		}

		return []*model.ArtifactVersion{
			{
				Kind:    model.ArtifactVersion_CONTAINER_IMAGE,
				Version: tag,
				Name:    name,
				Url:     fm.Spec.ImageURI,
			},
		}, nil
	}

	// Extract s3 object version as application version.
	if fm.Spec.S3ObjectVersion != "" {
		return []*model.ArtifactVersion{
			{
				Kind:    model.ArtifactVersion_S3_OBJECT,
				Version: fm.Spec.S3ObjectVersion,
				Name:    fm.Spec.S3Key,
				Url:     fmt.Sprintf("https://console.aws.amazon.com/s3/object/%s?prefix=%s", fm.Spec.S3Bucket, fm.Spec.S3Key),
			},
		}, nil
	}

	// Extract source code commihash as application version.
	if fm.Spec.SourceCode.Ref != "" {
		u, err := git.ParseGitURL(fm.Spec.SourceCode.Git)
		if err != nil {
			return nil, err
		}

		scheme := "https"
		if u.Scheme != "ssh" {
			scheme = u.Scheme
		}

		repoPath := strings.Trim(u.Path, "/")
		repoPath = strings.TrimSuffix(repoPath, ".git")

		var gitURL string
		switch u.Host {
		case "github.com", "gitlab.com":
			gitURL = fmt.Sprintf("%s://%s/%s/commit/%s", scheme, u.Host, repoPath, fm.Spec.SourceCode.Ref)
		case "bitbucket.org":
			gitURL = fmt.Sprintf("%s://%s/%s/commits/%s", scheme, u.Host, repoPath, fm.Spec.SourceCode.Ref)
		default:
			// TODO: Show repo name with commit link for other git provider
			gitURL = ""
			repoPath = ""
		}

		return []*model.ArtifactVersion{
			{
				Kind:    model.ArtifactVersion_GIT_SOURCE,
				Version: fm.Spec.SourceCode.Ref,
				Name:    repoPath,
				Url:     gitURL,
			},
		}, nil
	}

	return nil, fmt.Errorf("couldn't determine artifact versions")
}
