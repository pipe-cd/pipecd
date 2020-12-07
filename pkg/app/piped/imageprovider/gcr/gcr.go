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

package gcr

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/docker/distribution/registry/client"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/auth/challenge"
	"github.com/docker/distribution/registry/client/transport"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/model"
)

type Provider struct {
	name      string
	baseURL   url.URL
	transport http.RoundTripper

	logger *zap.Logger
}

type determineURL func(manager challenge.Manager, tx http.RoundTripper, domain string) (*url.URL, error)

func NewProvider(name string, cfg *config.ImageProviderGCRConfig, fn determineURL, logger *zap.Logger) (*Provider, error) {
	var tx http.RoundTripper = &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
		Proxy:           http.ProxyFromEnvironment,
	}
	manager := challenge.NewSimpleManager()

	u, err := fn(manager, tx, cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to determine registry URL: %w", err)
	}
	a := newAuthorizer(tx, manager)
	return &Provider{
		name:      name,
		baseURL:   *u,
		transport: transport.NewTransport(tx, a),
		logger:    logger.Named("gcr-provider"),
	}, nil
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Type() model.ImageProviderType {
	return model.ImageProviderTypeGCR
}

func (p *Provider) ParseImage(image string) (*model.ImageName, error) {
	ss := strings.SplitN(image, "/", 2)
	if len(ss) < 2 {
		return nil, fmt.Errorf("invalid image format (e.g. gcr.io/pipecd/helloworld)")
	}
	return &model.ImageName{
		Domain: ss[0],
		Repo:   ss[1],
	}, nil
}

func (p *Provider) GetLatestImage(ctx context.Context, image *model.ImageName) (*model.ImageRef, error) {
	repository, err := client.NewRepository(image, p.baseURL.String(), p.transport)
	if err != nil {
		return nil, err
	}
	// TODO: Stop listing all tags
	_, err = repository.Tags(ctx).All(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: Give back latest image from GCR
	return nil, nil
}

func newAuthorizer(tx http.RoundTripper, manager challenge.Manager) transport.RequestModifier {
	// TODO: Use credentials for GCR configured by user
	authHandlers := []auth.AuthenticationHandler{
		auth.NewTokenHandler(tx, nil, "", "pull"),
		auth.NewBasicHandler(nil),
	}
	return auth.NewAuthorizer(manager, authHandlers...)
}
