// Copyright 2025 The PipeCD Authors.
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

package oci

const (
	MediaTypePipedPlugin = "application/vnd.pipecd.piped.plugin"
)

type PushOptions struct {
	insecure bool
}

type PushOption interface {
	applyPushOption(*PushOptions)
}

type PullOptions struct {
	insecure   bool
	targetOS   string
	targetArch string
	mediaType  string
}

type PullOption interface {
	applyPullOption(*PullOptions)
}

type Option interface {
	PushOption
	PullOption
}

type insecureOption bool

func (o insecureOption) applyPushOption(opts *PushOptions) {
	opts.insecure = bool(o)
}

func (o insecureOption) applyPullOption(opts *PullOptions) {
	opts.insecure = bool(o)
}

func WithInsecure() Option {
	return insecureOption(true)
}

type targetOSOption string

func (o targetOSOption) applyPullOption(opts *PullOptions) {
	opts.targetOS = string(o)
}

func WithTargetOS(os string) PullOption {
	return targetOSOption(os)
}

type targetArchOption string

func (o targetArchOption) applyPullOption(opts *PullOptions) {
	opts.targetArch = string(o)
}

func WithTargetArch(arch string) PullOption {
	return targetArchOption(arch)
}

type mediaTypeOption string

func (o mediaTypeOption) applyPullOption(opts *PullOptions) {
	opts.mediaType = string(o)
}

func WithMediaType(mediaType string) PullOption {
	return mediaTypeOption(mediaType)
}
