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
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/lambda"
	"github.com/pipe-cd/pipecd/pkg/git"
)

func TestConfigureTrafficRouting(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name      string
		version   string
		percent   int
		primary   *provider.VersionTraffic
		secondary *provider.VersionTraffic
		out       bool
	}{
		{
			name:      "failed on invalid routing config: primary is missing",
			version:   "2",
			percent:   100,
			primary:   nil,
			secondary: nil,
			out:       false,
		},
		{
			name:    "configure successfully in case only primary provided",
			version: "2",
			percent: 100,
			primary: &provider.VersionTraffic{
				Version: "1",
				Percent: 100,
			},
			secondary: nil,
			out:       true,
		},
		{
			name:    "configure successfully in case set new primary lower than 100 percent",
			version: "2",
			percent: 70,
			primary: &provider.VersionTraffic{
				Version: "1",
				Percent: 100,
			},
			secondary: nil,
			out:       true,
		},
		{
			name:    "configure successfully in case set new primary lower than 100 percent and currently 2 versions is set",
			version: "3",
			percent: 70,
			primary: &provider.VersionTraffic{
				Version: "2",
				Percent: 50,
			},
			secondary: &provider.VersionTraffic{
				Version: "1",
				Percent: 50,
			},
			out: true,
		},
		{
			name:    "configure successfully in case set new primary to 100 percent and currently 2 versions is set",
			version: "3",
			percent: 100,
			primary: &provider.VersionTraffic{
				Version: "2",
				Percent: 50,
			},
			secondary: &provider.VersionTraffic{
				Version: "1",
				Percent: 50,
			},
			out: true,
		},
		{
			name:    "configure successfully in case new primary is the same as current primary",
			version: "2",
			percent: 100,
			primary: &provider.VersionTraffic{
				Version: "2",
				Percent: 50,
			},
			secondary: &provider.VersionTraffic{
				Version: "1",
				Percent: 50,
			},
			out: true,
		},
		{
			name:    "configure successfully in case new primary is the same as current secondary",
			version: "2",
			percent: 100,
			primary: &provider.VersionTraffic{
				Version: "1",
				Percent: 50,
			},
			secondary: &provider.VersionTraffic{
				Version: "2",
				Percent: 50,
			},
			out: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			trafficCfg := make(map[provider.TrafficConfigKeyName]provider.VersionTraffic)
			if tc.primary != nil {
				trafficCfg[provider.TrafficPrimaryVersionKeyName] = *tc.primary
			}
			if tc.secondary != nil {
				trafficCfg[provider.TrafficSecondaryVersionKeyName] = *tc.secondary
			}
			ok := configureTrafficRouting(trafficCfg, tc.version, tc.percent)
			assert.Equal(t, tc.out, ok)
			if primary, ok := trafficCfg[provider.TrafficPrimaryVersionKeyName]; ok {
				assert.Equal(t, tc.version, primary.Version)
				assert.Equal(t, float64(tc.percent), primary.Percent)
				if secondary, ok := trafficCfg[provider.TrafficSecondaryVersionKeyName]; ok {
					assert.Equal(t, float64(100-tc.percent), secondary.Percent)
				}
			}
		})
	}
}

type fakeRepo struct {
	git.Repo
	source string
}

func (m *fakeRepo) GetPath() string {
	return m.source
}

func (m *fakeRepo) Checkout(_ context.Context, _ string) error {
	return nil
}

func (m *fakeRepo) Clean() error {
	return nil
}

type fakeGitClient struct {
	repo git.Repo
}

func (g *fakeGitClient) Clone(_ context.Context, _, _, _, _ string) (git.Repo, error) {
	return g.repo, nil
}

func TestPrepareZipFromSource(t *testing.T) {
	t.Parallel()

	gc := &fakeGitClient{
		repo: &fakeRepo{
			source: "testdata/raw",
		},
	}
	fm := provider.FunctionManifest{}
	r, err := prepareZipFromSource(context.Background(), gc, fm)
	require.Nil(t, err)

	data, err := io.ReadAll(r)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(data))
}
