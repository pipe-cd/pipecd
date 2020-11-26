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
	"net/http"
	"net/url"

	"github.com/docker/distribution/registry/client/auth/challenge"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/imageprovider/gcr"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

// Provider acs as a container registry client.
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
		return gcr.NewProvider(cfg.Name, cfg.GCRConfig, doChallenge, logger)
	case model.ImageProviderTypeDockerhub:
		return nil, fmt.Errorf("not implemented yet")
	case model.ImageProviderTypeECR:
		return nil, fmt.Errorf("not implemented yet")
	default:
		return nil, fmt.Errorf("unknown image provider type: %s", cfg.Type)
	}
}

func doChallenge(manager challenge.Manager, tx http.RoundTripper, domain string) (*url.URL, error) {
	registryURL := url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   "/v2/",
	}
	cs, err := manager.GetChallenges(registryURL)
	if err != nil {
		return nil, err
	}
	if len(cs) == 0 {
		// TODO: Handle referring to https://github.com/fluxcd/flux/blob/72743f209207453a4326757ba89fb03cb514b34d/pkg/registry/client_factory.go#L64-L91
	}

	return &registryURL, nil
}
