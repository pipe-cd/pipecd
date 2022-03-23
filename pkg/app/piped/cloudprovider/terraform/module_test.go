package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestLoadTerraformFiles(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		moduleDir   string
		expected    []*File
		expectedErr bool
	}{
		{
			name:      "single module",
			moduleDir: "./testdata/single_module",
			expected: []*File{
				{
					Modules: []*Module{
						{
							Name:    "helloworld",
							Source:  "helloworld",
							Version: "v1.0.0",
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			name:      "multi modules",
			moduleDir: "./testdata/multi_modules",
			expected: []*File{
				{
					Modules: []*Module{
						{
							Name:    "helloworld_01",
							Source:  "helloworld",
							Version: "v1.0.0",
						},
						{
							Name:    "helloworld_02",
							Source:  "helloworld",
							Version: "v0.9.0",
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			name:      "multi modules with multi files",
			moduleDir: "./testdata/multi_modules_with_multi_files",
			expected: []*File{
				{
					Modules: []*Module{
						{
							Name:    "helloworld_01",
							Source:  "helloworld",
							Version: "v1.0.0",
						},
					},
				},
				{
					Modules: []*Module{
						{
							Name:    "helloworld_02",
							Source:  "helloworld",
							Version: "v0.9.0",
						},
					},
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tfs, err := LoadTerraformFiles(tc.moduleDir)
			if err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, tc.expected, tfs)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}

func TestFindArticatVersions(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		moduleDir   string
		expected    []*model.ArtifactVersion
		expectedErr bool
	}{
		{
			name:      "single module",
			moduleDir: "./testdata/single_module",
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
					Name:    "helloworld",
					Url:     "helloworld",
					Version: "v1.0.0",
				},
			},
			expectedErr: false,
		},
		{
			name:      "multi modules",
			moduleDir: "./testdata/multi_modules",
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
					Name:    "helloworld_01",
					Url:     "helloworld",
					Version: "v1.0.0",
				},
				{
					Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
					Name:    "helloworld_02",
					Url:     "helloworld",
					Version: "v0.9.0",
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tfs, err := LoadTerraformFiles(tc.moduleDir)
			require.NoError(t, err)

			versions, err := FindArtifactVersions(tfs)
			assert.ElementsMatch(t, tc.expected, versions)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}
