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
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// MakeCommitURL builds a link to the HTML page of the commit, using the given repoURL and hash.
func MakeCommitURL(repoURL, hash string) (string, error) {
	u, err := parseGitURL(repoURL)
	if err != nil {
		return "", err
	}

	scheme := "https"
	if u.Scheme != "ssh" {
		scheme = u.Scheme
	}

	repoPath := strings.Trim(u.Path, "/")
	repoPath = strings.TrimSuffix(repoPath, ".git")

	subPath := ""
	switch u.Host {
	case "github.com", "gitlab.com":
		subPath = "commit"
	case "bitbucket.org":
		subPath = "commits"
	default:
		// TODO: Allow users to specify git host
		//   Currently, the same subPath as Github is applied for all of unsupported hosts,
		//   to support GHE where its host could be customized
		subPath = "commit"
	}

	return fmt.Sprintf("%s://%s/%s/%s/%s", scheme, u.Host, repoPath, subPath, hash), nil
}

// MakeDirURL builds a link to the HTML page of the directory.
func MakeDirURL(repoURL, dir, branch string) (string, error) {
	if branch == "" {
		return "", fmt.Errorf("no branch given")
	}
	u, err := parseGitURL(repoURL)
	if err != nil {
		return "", err
	}

	scheme := "https"
	if u.Scheme != "ssh" {
		scheme = u.Scheme
	}

	repoPath := strings.Trim(u.Path, "/")
	repoPath = strings.TrimSuffix(repoPath, ".git")

	subPath := ""
	switch u.Host {
	case "github.com", "gitlab.com":
		subPath = "tree"
	case "bitbucket.org":
		subPath = "src"
	default:
		// TODO: Allow users to specify git host
		subPath = "tree"
	}

	dir = strings.Trim(dir, "/")

	return fmt.Sprintf("%s://%s/%s/%s/%s/%s", scheme, u.Host, repoPath, subPath, branch, dir), nil
}

// MakeFileCreationURL builds a link to create a file under the given directory.
func MakeFileCreationURL(repoURL, dir, branch, filename, value string) (string, error) {
	if branch == "" {
		return "", fmt.Errorf("no branch given")
	}
	u, err := parseGitURL(repoURL)
	if err != nil {
		return "", err
	}

	if u.Scheme == "ssh" {
		u.Scheme = "https"
		u.User = nil
	}
	repoPath := strings.TrimSuffix(strings.Trim(u.Path, "/"), ".git")
	dir = strings.Trim(dir, "/")

	switch u.Host {
	case "github.com":
		u.Path = fmt.Sprintf("%s/%s/%s/%s", repoPath, "new", branch, dir)
		params := &url.Values{}
		if filename != "" {
			// NOTE: We're getting an issue with specifying a filename: https://github.com/isaacs/github/issues/1527
			//   In short, the last path part of the URL is ignored.
			//   For now, appending dummy path as a workaround.
			u.Path += "/dummy"
			params.Add("filename", filename)
		}
		if value != "" {
			params.Add("value", value)
		}
		u.RawQuery = params.Encode()
	default:
		// TODO: Allow users to specify git host
		u.Path = fmt.Sprintf("%s/%s/%s/%s", repoPath, "new", branch, dir)
	}

	return u.String(), nil
}

var (
	knownSchemes = map[string]interface{}{
		"ssh":     struct{}{},
		"git":     struct{}{},
		"git+ssh": struct{}{},
		"http":    struct{}{},
		"https":   struct{}{},
		"rsync":   struct{}{},
		"file":    struct{}{},
	}
	scpRegex = regexp.MustCompile(`^([a-zA-Z0-9_]+@)?([a-zA-Z0-9._-]+):(.*)$`)
)

// parseGitURL parses git url into a URL structure.
func parseGitURL(rawURL string) (u *url.URL, err error) {
	u, err = parseTransport(rawURL)
	if err == nil {
		return
	}
	return parseScp(rawURL)
}

// ParseGitURL parses git url into a URL structure.
func ParseGitURL(rawURL string) (u *url.URL, err error) {
	return parseGitURL(rawURL)
}

// Return a structured URL only when scheme is a known Git transport.
func parseTransport(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse git url: %w", err)
	}
	if _, ok := knownSchemes[u.Scheme]; !ok {
		return nil, fmt.Errorf("unknown scheme %q", u.Scheme)
	}
	return u, nil
}

// Return a structured URL only when the rawURL is an SCP-like URL.
func parseScp(rawURL string) (*url.URL, error) {
	match := scpRegex.FindAllStringSubmatch(rawURL, -1)
	if len(match) == 0 {
		return nil, fmt.Errorf("no scp URL found in %q", rawURL)
	}
	m := match[0]
	user := strings.TrimRight(m[1], "@")
	var userinfo *url.Userinfo
	if user != "" {
		userinfo = url.User(user)
	}
	return &url.URL{
		Scheme: "ssh",
		User:   userinfo,
		Host:   m[2],
		Path:   m[3],
	}, nil
}
