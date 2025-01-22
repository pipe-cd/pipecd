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

package lifecycle

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGracefulStopCommand(t *testing.T) {
	testcases := []struct {
		name      string
		stopAfter time.Duration
	}{
		{
			name:      "graceful stop after very short time",
			stopAfter: time.Nanosecond,
		},
		{
			name:      "graceful stop after second",
			stopAfter: time.Second,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := RunBinary(context.TODO(), "sh", []string{"sleep", "1m"})
			require.NoError(t, err)
			require.NotNil(t, cmd)

			assert.True(t, cmd.IsRunning())
			cmd.GracefulStop(tc.stopAfter)
			assert.False(t, cmd.IsRunning())
		})
	}
}

func TestGracefulStopCommandResult(t *testing.T) {
	testcases := []struct {
		name      string
		exitCode  int
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "successfully exit",
			exitCode:  0,
			assertion: assert.NoError,
		},
		{
			name:      "exit with an error",
			exitCode:  1,
			assertion: assert.Error,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := RunBinary(context.TODO(), "sh", []string{"-c", "exit " + strconv.Itoa(tc.exitCode)})
			require.NoError(t, err)
			require.NotNil(t, cmd)

			time.Sleep(100 * time.Millisecond) // to avoid GracefulStop executed before the command exits
			tc.assertion(t, cmd.GracefulStop(time.Second))
			assert.False(t, cmd.IsRunning())
		})
	}
}

func TestDownloadBinary(t *testing.T) {
	server := httpTestServer()
	defer server.Close()

	logger := zaptest.NewLogger(t)

	t.Run("successful download", func(t *testing.T) {
		destDir := t.TempDir()
		destFile := "test-binary"
		url := server.URL + "/binary"
		path, err := DownloadBinary(url, destDir, destFile, logger)
		require.NoError(t, err)
		assert.FileExists(t, path)
	})

	t.Run("file already exists", func(t *testing.T) {
		destDir := t.TempDir()
		destFile := "test-binary"
		url := server.URL + "/binary"
		path, err := DownloadBinary(url, destDir, destFile, logger)
		require.NoError(t, err)
		assert.FileExists(t, path)

		// Try downloading again, should not error and file should still exist
		path, err = DownloadBinary(url, destDir, destFile, logger)
		require.NoError(t, err)
		assert.FileExists(t, path)
	})

	t.Run("file on local", func(t *testing.T) {
		sourceDir := t.TempDir()
		sourceFile := "test-binary"
		sourcePath := path.Join(sourceDir, sourceFile)
		err := os.WriteFile(sourcePath, []byte("test binary content"), 0755)
		require.NoError(t, err)

		destDir := t.TempDir()
		destFile := "test-binary"
		url := "file://" + sourcePath

		path, err := DownloadBinary(url, destDir, destFile, logger)
		require.NoError(t, err)
		assert.FileExists(t, path)
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "test binary content", string(content))
	})

	t.Run("not valid source url given", func(t *testing.T) {
		destDir := t.TempDir()
		destFile := "test-binary"
		url := "ftp://invalid-url"

		path, err := DownloadBinary(url, destDir, destFile, logger)
		require.Error(t, err)
		assert.Empty(t, path)
	})
}

func httpTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/binary" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test binary content"))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}
