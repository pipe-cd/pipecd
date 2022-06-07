// Copyright 2022 The PipeCD Authors.
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

package githubstore

import (
	"context"

	"github.com/google/go-github/v29/github"
)

const (
	pipecdGithubOwner = "pipe-cd"
	pipecdGithubRepo  = "pipecd"
)

type Store interface {
	ListReleasedVersions(ctx context.Context) ([]string, error)
}

type store struct {
	*github.Client
}

func NewStore() Store {
	return &store{
		Client: github.NewClient(nil),
	}
}

func (s *store) ListReleasedVersions(ctx context.Context) ([]string, error) {
	releases, _, err := s.Repositories.ListReleases(ctx, pipecdGithubOwner, pipecdGithubRepo, nil)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return []string{}, nil
	}

	releasedVersions := make([]string, 0, len(releases))
	for _, release := range releases {
		// Ignore pre-release tagged release.
		if *release.Prerelease {
			continue
		}
		releasedVersions = append(releasedVersions, *release.TagName)
	}
	return releasedVersions, nil
}
