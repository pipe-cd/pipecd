package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipe-cd/pipecd/pkg/app/piped/deploysource"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func TestLoadTerraformFiles(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		moduleDir   string
		expected    []File
		expectedErr bool
	}{
		{
			name:      "single module",
			moduleDir: "./testdata/single_module",
			expected: []File{
				{
					Modules: []*Module{
						{
							Name:    "helloworld",
							Source:  "./helloworld",
							Version: "v1.0.0",
							IsLocal: true,
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			name:      "multi modules",
			moduleDir: "./testdata/multi_modules",
			expected: []File{
				{
					Modules: []*Module{
						{
							Name:    "helloworld_01",
							Source:  "./helloworld",
							Version: "v1.0.0",
							IsLocal: true,
						},
						{
							Name:    "helloworld_02",
							Source:  "./helloworld",
							Version: "v0.9.0",
							IsLocal: true,
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			name:      "multi modules with multi files",
			moduleDir: "./testdata/multi_modules_with_multi_files",
			expected: []File{
				{
					Modules: []*Module{
						{
							Name:    "helloworld_01",
							Source:  "./helloworld",
							Version: "v1.0.0",
							IsLocal: true,
						},
					},
				},
				{
					Modules: []*Module{
						{
							Name:    "helloworld_02",
							Source:  "./helloworld",
							Version: "v0.9.0",
							IsLocal: true,
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
		gp          *model.ApplicationGitPath
		ds          *deploysource.DeploySource
		expected    []*model.ArtifactVersion
		expectedErr bool
	}{
		{
			name:      "single local module",
			moduleDir: "./testdata/single_module",
			gp: &model.ApplicationGitPath{
				Repo: &model.ApplicationGitRepository{
					Remote: "https://githuh.com/example",
					Branch: "main",
				},
			},
			ds: &deploysource.DeploySource{
				RepoDir: "/repo/example",
				AppDir:  "/repo/example",
			},
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
					Name:    "helloworld",
					Url:     "https://githuh.com/example/tree/main/helloworld",
					Version: "v1.0.0",
				},
			},
			expectedErr: false,
		},
		{
			name:      "multi local modules",
			moduleDir: "./testdata/multi_modules",
			gp: &model.ApplicationGitPath{
				Repo: &model.ApplicationGitRepository{
					Remote: "https://githuh.com/example",
					Branch: "main",
				},
			},
			ds: &deploysource.DeploySource{
				RepoDir: "/repo/example",
				AppDir:  "/repo/example",
			},
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
					Name:    "helloworld_01",
					Url:     "https://githuh.com/example/tree/main/helloworld",
					Version: "v1.0.0",
				},
				{
					Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
					Name:    "helloworld_02",
					Url:     "https://githuh.com/example/tree/main/helloworld",
					Version: "v0.9.0",
				},
			},
			expectedErr: false,
		},
		{
			name:      "multi local modules with multi files",
			moduleDir: "./testdata/multi_modules",
			gp: &model.ApplicationGitPath{
				Repo: &model.ApplicationGitRepository{
					Remote: "https://githuh.com/example",
					Branch: "main",
				},
			},
			ds: &deploysource.DeploySource{
				RepoDir: "/repo/example",
				AppDir:  "/repo/example",
			},
			expected: []*model.ArtifactVersion{
				{
					Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
					Name:    "helloworld_01",
					Url:     "https://githuh.com/example/tree/main/helloworld",
					Version: "v1.0.0",
				},
				{
					Kind:    model.ArtifactVersion_TERRAFORM_MODULE,
					Name:    "helloworld_02",
					Url:     "https://githuh.com/example/tree/main/helloworld",
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

			versions, err := FindArtifactVersions(tfs, tc.gp, tc.ds)
			assert.ElementsMatch(t, tc.expected, versions)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}

func TestMakeURL(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		lc          *LocalModuleSourceConverter
		moduleSrc   string
		expected    string
		expectedErr bool
	}{
		{
			name: "relative path ./",
			lc: &LocalModuleSourceConverter{
				GitURL:  "https://github.com/example/test",
				Branch:  "main",
				RepoDir: "/repo/test",
				AppDir:  "/repo/test/hoge",
			},
			moduleSrc:   "./",
			expected:    "https://github.com/example/test/tree/main/hoge",
			expectedErr: false,
		},
		{
			name: "relative path ../",
			lc: &LocalModuleSourceConverter{
				GitURL:  "https://github.com/example/test",
				Branch:  "main",
				RepoDir: "/repo/test",
				AppDir:  "/repo/test/hoge/fuga",
			},
			moduleSrc:   "../",
			expected:    "https://github.com/example/test/tree/main/hoge",
			expectedErr: false,
		},
		{
			name: "relative path ./hoge",
			lc: &LocalModuleSourceConverter{
				GitURL:  "https://github.com/example/test",
				Branch:  "main",
				RepoDir: "/repo/test",
				AppDir:  "/repo/test",
			},
			moduleSrc:   "./hoge",
			expected:    "https://github.com/example/test/tree/main/hoge",
			expectedErr: false,
		},
		{
			name: "relative path ../fuga",
			lc: &LocalModuleSourceConverter{
				GitURL:  "https://github.com/example/test",
				Branch:  "main",
				RepoDir: "/repo/test",
				AppDir:  "/repo/test/hoge",
			},
			moduleSrc:   "../fuga",
			expected:    "https://github.com/example/test/tree/main/fuga",
			expectedErr: false,
		},
		{
			name: "can't resolve path",
			lc: &LocalModuleSourceConverter{
				GitURL:  "https://github.com/example/test",
				Branch:  "main",
				RepoDir: "/repo/test",
				AppDir:  "/repo/test",
			},
			moduleSrc:   "../../",
			expected:    "",
			expectedErr: true,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			url, err := tc.lc.MakeURL(tc.moduleSrc)

			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, url, tc.expected)
		})
	}
}

func TestIsLocalModule(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name      string
		moduleSrc string
		expected  bool
	}{
		{
			name:      "start with ./",
			moduleSrc: "./test",
			expected:  true,
		},
		{
			name:      "start with ../",
			moduleSrc: "../test",
			expected:  true,
		},
		{
			name:      "not a relative path",
			moduleSrc: "test",
			expected:  false,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, isLocalModule(tc.moduleSrc), tc.expected)
		})
	}
}
