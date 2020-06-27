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

package kubernetes

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/piped/toolregistry"
)

func TestMain(m *testing.M) {
	binDir := "/tmp/piped-bin"
	if err := toolregistry.InitDefaultRegistry(binDir, zap.NewNop()); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestTemplateLocalChart(t *testing.T) {
	var (
		ctx       = context.Background()
		appName   = "testapp"
		appDir    = "testdata"
		chartPath = "testchart"
	)

	// TODO: Preinstall a helm version inside CI runner to avoid installing.
	helmPath, _, err := toolregistry.DefaultRegistry().Helm(ctx, "")
	require.NoError(t, err)

	helm := NewHelm("", helmPath, zap.NewNop())
	out, err := helm.TemplateLocalChart(ctx, appName, appDir, chartPath, nil)
	require.NoError(t, err)

	out = strings.TrimPrefix(out, "---")
	manifests := strings.Split(out, "---")
	assert.Equal(t, 3, len(manifests))
}
