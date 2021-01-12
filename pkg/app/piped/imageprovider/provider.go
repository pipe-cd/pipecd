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

package imageprovider

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/imageprovider/ecr"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

// Provider acts as a container registry client.
type Provider interface {
	// Name gives back the provider name that is unique in the Piped.
	Name() string
	// Type indicates which container registry client to act as.
	Type() model.ImageProviderType
	// ParseImage converts the given string into structured image.
	ParseImage(image string) (*model.ImageName, error)
	// GetLatestImages gives back an image with the latest tag.
	GetLatestImage(ctx context.Context, image *model.ImageName) (*model.ImageRef, error)
}

// NewProvider yields an appropriate provider according to the given config.
func NewProvider(cfg *config.PipedImageProvider, logger *zap.Logger) (Provider, error) {
	switch cfg.Type {
	case model.ImageProviderTypeGCR:
		/*
			options := []gcr.Option{
				gcr.WithServiceAccountFile(cfg.GCRConfig.ServiceAccountFile),
				gcr.WithLogger(logger),
			}
			return gcr.NewGCR(cfg.Name, options...)
		*/
		return nil, nil
	case model.ImageProviderTypeECR:
		options := []ecr.Option{
			ecr.WithRegistryID(cfg.ECRConfig.RegistryID),
			ecr.WithCredentialsFile(cfg.ECRConfig.CredentialsFile),
			ecr.WithProfile(cfg.ECRConfig.Profile),
			ecr.WithLogger(logger),
		}
		return ecr.NewECR(cfg.Name, cfg.ECRConfig.Region, options...)
	case model.ImageProviderTypeDockerHub:
		return nil, fmt.Errorf("not implemented yet")
	default:
		return nil, fmt.Errorf("unknown image provider type: %s", cfg.Type)
	}
}
