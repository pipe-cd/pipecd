package git

import (
	"net/url"
	"reflect"
	"testing"
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
			want:    "",
			wantErr: true,
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
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeCommitURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MakeCommitURL() got = %v, want %v", got, tt.want)
			}
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
			name:    "ssh to unsupported git host",
			repoURL: "git@foo.com:org/repo.git",
			dir:     "path/to",
			branch:  "abc",
			want:    "",
			wantErr: true,
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
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeCommitURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MakeCommitURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_parseGitURL(t *testing.T) {
	tests := []struct {
		name    string
		rawurl  string
		wantU   *url.URL
		wantErr bool
	}{
		{
			name:   "SCP-like URL with user",
			rawurl: "user@host.xz:path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "ssh",
				User:   url.User("user"),
				Host:   "host.xz",
				Path:   "path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "SCP-like URL without user",
			rawurl: "host.xz:path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "ssh",
				User:   nil,
				Host:   "host.xz",
				Path:   "path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "SCP-like URL with prefix `/`",
			rawurl: "host.xz:/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "ssh",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "ssh with user",
			rawurl: "ssh://user@host.xz/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "ssh",
				User:   url.User("user"),
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "ssh with user with port",
			rawurl: "ssh://user@host.xz:1234/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "ssh",
				User:   url.User("user"),
				Host:   "host.xz:1234",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "git+ssh",
			rawurl: "git+ssh://host.xz/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "git+ssh",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "file scheme",
			rawurl: "file:///path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "file",
				User:   nil,
				Host:   "",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "rsync + ssh",
			rawurl: "rsync://host.xz/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "rsync",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "git scheme",
			rawurl: "git://host.xz/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "git",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "http scheme",
			rawurl: "http://host.xz/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "http",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "https scheme",
			rawurl: "https://host.xz/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "https",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "ftp scheme",
			rawurl: "ftp://host.xz/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "ftp",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
		{
			name:   "ftps scheme",
			rawurl: "ftps://host.xz/path/to/repo.git/",
			wantU: &url.URL{
				Scheme: "ftps",
				User:   nil,
				Host:   "host.xz",
				Path:   "/path/to/repo.git/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, err := parseGitURL(tt.rawurl)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGitURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotU, tt.wantU) {
				t.Errorf("parseGitURL() got = %#v, want %#v", gotU, tt.wantU)
			}
		})
	}
}
