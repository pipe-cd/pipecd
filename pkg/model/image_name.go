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

package model

import (
	"fmt"
	"path"
)

// ImageName represents an untagged image. Note that images may have
// the domain omitted (e.g. Docker Hub). If they only have single path element,
// the prefix `library` is implied.
//
// Examples:
//   - alpine
//   - library/alpine
//   - gcr.io/pipecd/helloworld
type ImageName struct {
	Domain string
	Repo   string
}

func (i ImageName) String() string {
	return path.Join(i.Domain, i.Repo)
}

// Name gives back just repository name without domain.
func (i ImageName) Name() string {
	return i.Repo
}

// ImageRef represents a tagged image. The tag is allowed to be
// empty, though it is in general undefined what that means
//
// Examples:
//   - alpine:3.0
//   - library/alpine:3.0
//   - gcr.io/pipecd/helloworld:0.1.0
type ImageRef struct {
	ImageName
	Tag string
}

func (i ImageRef) String() string {
	if i.Tag == "" {
		return i.ImageName.String()
	}

	return fmt.Sprintf("%s:%s", i.ImageName.String(), i.Tag)
}
