package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyCustomtemplatingArgs(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		appDir        string
		valueFilePath string
		wantErr       bool
	}{
		{
			name:          "Input arg is a path inside the app dir",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "test.yaml",
			wantErr:       false,
		},
		{
			name:          "Input arg is a path inside the app dir and does not have extention",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "test",
			wantErr:       false,
		},
		{
			name:          "Input arg is a path inside the app dir (with ..)",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "../../../testdata/testcustomtemplating/appconfdir/test.yaml",
			wantErr:       false,
		},
		{
			name:          "Input arg is a path under the app dir",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "dir/test.yaml",
			wantErr:       false,
		},
		{
			name:          "Input arg is a path under the app dir (with ..)",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "../../../testdata/testcustomtemplating/appconfdir/dir/test.yaml",
			wantErr:       false,
		},
		{
			name:          "Input arg is a path outside the app dir",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "/etc/hosts",
			wantErr:       true,
		},
		{
			name:          "Input arg is a path outside the app dir (with ..)",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "../../../../../../../../../../../../etc/hosts",
			wantErr:       true,
		},
		{
			name:          "Input arg is remote URL (http)",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "http://exmaple.com/test.yaml",
			wantErr:       false,
		},
		{
			name:          "Input arg is remote URL (https)",
			appDir:        "testdata/testcustomtemplatingtemplating/appconfdir",
			valueFilePath: "https://exmaple.com/test.yaml",
			wantErr:       false,
		},
		{
			name:          "Input arg is  disallowed remote URL (ftp)",
			appDir:        "testdata/testcustomtemplatingtemplating/appconfdir",
			valueFilePath: "ftp://exmaple.com/test.yaml",
			wantErr:       true,
		},
		{
			name:          "Input arg is symlink targeting validtests file",
			appDir:        "testdata/testcustomtemplatingtemplating/appconfdir",
			valueFilePath: "valid-symlink",
			wantErr:       false,
		},
		{
			name:          "Input arg is symlink targeting invalidtests file",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "invalid-symlink",
			wantErr:       true,
		},
		{
			name:          "Input arg is ./...",
			appDir:        "testdata/testcustomtemplating/appconfdir",
			valueFilePath: "./...",
			wantErr:       false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyCustomtemplatingArgs(tc.appDir, tc.valueFilePath)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
