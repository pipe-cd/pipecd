// Copyright 2023 The PipeCD Authors.
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

package git

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeCommitURL(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
		hash    string
		want    string
		wantErr bool
	}{
		{
			name:    "ssh to github.com",
			repoURL: "git@github.com:org/repo.git",
			hash:    "abc",
			want:    "https://github.com/org/repo/commit/abc",
			wantErr: false,
		},
		{
			name:    "ssh to gitlab.com",
			repoURL: "git@gitlab.com:org/repo.git",
			hash:    "abc",
			want:    "https://gitlab.com/org/repo/commit/abc",
			wantErr: false,
		},
		{
			name:    "ssh to bitbucket.org",
			repoURL: "git@bitbucket.org:org/repo.git",
			hash:    "abc",
			want:    "https://bitbucket.org/org/repo/commits/abc",
			wantErr: false,
		},
		{
			name:    "ssh to unsupported git host",
			repoURL: "git@foo.com:org/repo.git",
			hash:    "abc",
			want:    "https://foo.com/org/repo/commit/abc",
			wantErr: false,
		},
		{
			name:    "ssh to github.com without `.git` suffix",
			repoURL: "git@github.com:org/repo",
			hash:    "abc",
			want:    "https://github.com/org/repo/commit/abc",
			wantErr: false,
		},
		{
			name:    "ssh to github.com with `/` suffix",
			repoURL: "git@github.com:org/repo/",
			hash:    "abc",
			want:    "https://github.com/org/repo/commit/abc",
			wantErr: false,
		},
		{
			name:    "http to github.com",
			repoURL: "http://github.com/org/repo",
			hash:    "abc",
			want:    "http://github.com/org/repo/commit/abc",
			wantErr: false,
		},
		{
			name:    "unparseable url",
			repoURL: "1234abcd",
			hash:    "abc",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeCommitURL(tt.repoURL, tt.hash)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMakeDirURL(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
		dir     string
		branch  string
		want    string
		wantErr bool
	}{
		{
			name:    "ssh to github.com",
			repoURL: "git@github.com:org/repo.git",
			dir:     "path/to",
			branch:  "abc",
			want:    "https://github.com/org/repo/tree/abc/path/to",
			wantErr: false,
		},
		{
			name:    "ssh to gitlab.com",
			repoURL: "git@gitlab.com:org/repo.git",
			dir:     "path/to",
			branch:  "abc",
			want:    "https://gitlab.com/org/repo/tree/abc/path/to",
			wantErr: false,
		},
		{
			name:    "ssh to bitbucket.org",
			repoURL: "git@bitbucket.org:org/repo.git",
			dir:     "path/to",
			branch:  "abc",
			want:    "https://bitbucket.org/org/repo/src/abc/path/to",
			wantErr: false,
		},
		{
			name:    "ssh to unsupported git host",
			repoURL: "git@foo.com:org/repo.git",
			dir:     "path/to",
			branch:  "abc",
			want:    "https://foo.com/org/repo/tree/abc/path/to",
			wantErr: false,
		},
		{
			name:    "ssh to github.com without `.git` suffix",
			repoURL: "git@github.com:org/repo",
			dir:     "path/to",
			branch:  "abc",
			want:    "https://github.com/org/repo/tree/abc/path/to",
			wantErr: false,
		},
		{
			name:    "ssh to github.com with `/` suffix",
			repoURL: "git@github.com:org/repo/",
			dir:     "path/to",
			branch:  "abc",
			want:    "https://github.com/org/repo/tree/abc/path/to",
			wantErr: false,
		},
		{
			name:    "http to github.com",
			repoURL: "http://github.com/org/repo",
			dir:     "path/to",
			branch:  "abc",
			want:    "http://github.com/org/repo/tree/abc/path/to",
			wantErr: false,
		},
		{
			name:    "unparseable url",
			repoURL: "1234abcd",
			dir:     "path/to",
			branch:  "abc",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeDirURL(tt.repoURL, tt.dir, tt.branch)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMakeFileCreationURL(t *testing.T) {
	tests := []struct {
		name     string
		repoURL  string
		dir      string
		branch   string
		filename string
		value    string
		want     string
		wantErr  bool
	}{
		{
			name:     "given filename",
			repoURL:  "git@github.com:org/repo.git",
			dir:      "path/to",
			branch:   "abc",
			filename: "foo.txt",
			want:     "https://github.com/org/repo/new/abc/path/to/dummy?filename=foo.txt",
			wantErr:  false,
		},
		{
			name:    "given value",
			repoURL: "git@github.com:org/repo.git",
			dir:     "path/to",
			branch:  "abc",
			value: `# Comment
foo:
  bar:
    baz:
      - a
      - b
`,
			want:    "https://github.com/org/repo/new/abc/path/to?value=%23+Comment%0Afoo%3A%0A++bar%3A%0A++++baz%3A%0A++++++-+a%0A++++++-+b%0A",
			wantErr: false,
		},
		{
			name:     "given filename and value",
			repoURL:  "git@github.com:org/repo.git",
			dir:      "path/to",
			branch:   "abc",
			filename: "foo.txt",
			value: `# Comment
foo:
  bar:
    baz:
      - a
      - b
`,
			want:    "https://github.com/org/repo/new/abc/path/to/dummy?filename=foo.txt&value=%23+Comment%0Afoo%3A%0A++bar%3A%0A++++baz%3A%0A++++++-+a%0A++++++-+b%0A",
			wantErr: false,
		},
		{
			name:    "ssh to unsupported git host",
			repoURL: "git@foo.com:org/repo.git",
			dir:     "path/to",
			branch:  "abc",
			want:    "https://foo.com/org/repo/new/abc/path/to",
			wantErr: false,
		},
		{
			name:    "no branch given",
			repoURL: "1234abcd",
			dir:     "path/to",
			wantErr: true,
		},
		{
			name:    "unparseable url",
			repoURL: "1234abcd",
			dir:     "path/to",
			branch:  "abc",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeFileCreationURL(tt.repoURL, tt.dir, tt.branch, tt.filename, tt.value)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseGitURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantURL *url.URL
		wantErr bool
	}{
		{
			name:   "SCP-like URL with user",
			rawURL: "user@host.xz:path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "ssh",
				User:   url.User("user"),
				Host:   "host.xz",
				Path:   "path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "SCP-like URL without user",
			rawURL: "host.xz:path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "ssh",
				User:   nil,
				Host:   "host.xz",
				Path:   "path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "SCP-like URL with prefix `/`",
			rawURL: "host.xz:/path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "ssh",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "ssh with user",
			rawURL: "ssh://user@host.xz/path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "ssh",
				User:   url.User("user"),
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "ssh with user with port",
			rawURL: "ssh://user@host.xz:1234/path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "ssh",
				User:   url.User("user"),
				Host:   "host.xz:1234",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "git+ssh",
			rawURL: "git+ssh://host.xz/path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "git+ssh",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "file scheme",
			rawURL: "file:///path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "file",
				User:   nil,
				Host:   "",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "rsync + ssh",
			rawURL: "rsync://host.xz/path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "rsync",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "git scheme",
			rawURL: "git://host.xz/path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "git",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "http scheme",
			rawURL: "http://host.xz/path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "http",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "https scheme",
			rawURL: "https://host.xz/path/to/repo.git/",
			wantURL: &url.URL{
				Scheme: "https",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGitURL(tt.rawURL)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, got, tt.wantURL)
		})
	}
}
