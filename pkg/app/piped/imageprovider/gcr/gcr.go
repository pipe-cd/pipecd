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

/*
import (
	"context"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/model"
)

type GCR struct {
	name               string
	serviceAccountFile string
	// If nil, treated as an anonymous user.
	authenticator authn.Authenticator
	logger        *zap.Logger
}

type Option func(*GCR)

func WithServiceAccountFile(path string) Option {
	return func(e *GCR) {
		e.serviceAccountFile = path
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(e *GCR) {
		e.logger = logger
	}
}

// NewGCR generates a GCR client with an anonymous user if no authenticate method set.
func NewGCR(name string, opts ...Option) (*GCR, error) {
	g := &GCR{
		name:   name,
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(g)
	}
	g.logger = g.logger.Named("gcr-provider")

	if g.serviceAccountFile != "" {
		b, err := ioutil.ReadFile(g.serviceAccountFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open the service account file: %w", err)
		}
		g.authenticator = google.NewJSONKeyAuthenticator(string(b))
	}
	return g, nil
}

func (g *GCR) Name() string {
	return g.name
}

func (g *GCR) Type() model.ImageProviderType {
	return model.ImageProviderTypeGCR
}

func (g *GCR) ParseImage(image string) (*model.ImageName, error) {
	ss := strings.SplitN(image, "/", 2)
	if len(ss) < 2 {
		return nil, fmt.Errorf("invalid image format (e.g. gcr.io/pipecd/helloworld)")
	}
	return &model.ImageName{
		Domain: ss[0],
		Repo:   ss[1],
	}, nil
}

func (g *GCR) GetLatestImage(ctx context.Context, image *model.ImageName) (*model.ImageRef, error) {
	repo, err := name.NewRepository(image.String())
	if err != nil {
		return nil, fmt.Errorf("%s is invalid repository: %w", image, err)
	}
	options := []google.ListerOption{
		google.WithContext(ctx),
	}
	if g.authenticator != nil {
		options = append(options, google.WithAuth(g.authenticator))
	}
	// TODO: Use pagination to retrieve image tags from GCR
	// Currently, the result could be quite large size if there are a lot of tags.
	// "google/go-containerregistry" doesn't provide any option to paginate.
	// We can propose it to them, or just borrow and modify for us.
	// See more: https://docs.docker.com/registry/spec/api/#listing-image-tags
	res, err := google.List(repo, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	if len(res.Manifests) == 0 {
		return nil, fmt.Errorf("no manifests found in %s", repo.Name())
	}

	// Determine the latest by sorting by the uploaded time.
	manifests := make([]google.ManifestInfo, 0, len(res.Manifests))
	for _, m := range res.Manifests {
		manifests = append(manifests, m)
	}
	sort.Slice(manifests, func(i, j int) bool {
		return manifests[i].Uploaded.After(manifests[j].Uploaded)
	})
	latest := manifests[0]
	if len(latest.Tags) == 0 {
		return nil, fmt.Errorf("no tag is associated to the latest image")
	}
	return &model.ImageRef{
		ImageName: *image,
		// TODO: Enable to specify the tag if multiple tags are associated to an image
		Tag: latest.Tags[0],
	}, nil
}
*/
