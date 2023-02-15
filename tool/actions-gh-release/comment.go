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

package main

import (
	"fmt"
	"strings"
)

const (
	successBadgeURL = `<!-- RELEASE -->
[![RELEASE](https://img.shields.io/static/v1?label=GitHub&message=RELEASE&color=success&style=flat)](https://github.com/pipe-cd/actions-gh-release)

`
)

func makeCommentBody(proposals []ReleaseProposal, exists []ReleaseProposal) string {
	var b strings.Builder
	b.WriteString(successBadgeURL)

	if len(proposals) == 0 {
		if len(exists) == 0 {
			fmt.Fprintf(&b, "No GitHub releases will be created one this pull request got merged. Because this pull request did not modified any RELEASE files.\n")
			return b.String()
		}

		fmt.Fprintf(&b, "No GitHub releases will be created one this pull request got merged. Because the following tags were already created before.\n")
		for _, p := range exists {
			fmt.Fprintf(&b, "- %s\n", p.Tag)
		}
		return b.String()
	}

	b.WriteString(fmt.Sprintf("The following %d GitHub releases will be created once this pull request got merged.\n", len(proposals)))
	for _, p := range proposals {
		fmt.Fprintf(&b, "\n")
		fmt.Fprintf(&b, p.ReleaseNote)
		fmt.Fprintf(&b, "\n")
	}

	if len(exists) > 0 {
		fmt.Fprintf(&b, "The following %d releases will be skipped because they were already created before.\n", len(exists))
		for _, p := range exists {
			fmt.Fprintf(&b, "- %s\n", p.Tag)
		}
	}

	return b.String()
}
