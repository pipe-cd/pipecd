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
	// MediaTypePipedPlugin is the media type for PipeCD Agent plugins.
	MediaTypePipedPlugin = "application/vnd.pipecd.piped.plugin"
)

// PushOptions holds options for pushing to an OCI registry.
type PushOptions struct {
	insecure bool
}

// PushOption is an interface for applying push options.
type PushOption interface {
	applyPushOption(*PushOptions)
}

// PullOptions holds options for pulling from an OCI registry.
type PullOptions struct {
	insecure     bool
	targetOS     string
	targetArch   string
	mediaType    string
	artifactType string
}

// PullOption is an interface for applying pull options.
type PullOption interface {
	applyPullOption(*PullOptions)
}

// Option is an interface that combines PushOption and PullOption.
type Option interface {
	PushOption
	PullOption
}

// insecureOption is an option to enable insecure connections.
type insecureOption bool

// applyPushOption applies the insecure option to PushOptions.
func (o insecureOption) applyPushOption(opts *PushOptions) {
	opts.insecure = bool(o)
}

// applyPullOption applies the insecure option to PullOptions.
func (o insecureOption) applyPullOption(opts *PullOptions) {
	opts.insecure = bool(o)
}

// WithInsecure returns an Option that enables insecure connections.
func WithInsecure() Option {
	return insecureOption(true)
}

// targetOSOption is an option to specify the target OS for pulling.
type targetOSOption string

// applyPullOption applies the target OS option to PullOptions.
func (o targetOSOption) applyPullOption(opts *PullOptions) {
	opts.targetOS = string(o)
}

// WithTargetOS returns a PullOption that sets the target OS.
func WithTargetOS(os string) PullOption {
	return targetOSOption(os)
}

// targetArchOption is an option to specify the target architecture for pulling.
type targetArchOption string

// applyPullOption applies the target architecture option to PullOptions.
func (o targetArchOption) applyPullOption(opts *PullOptions) {
	opts.targetArch = string(o)
}

// WithTargetArch returns a PullOption that sets the target architecture.
func WithTargetArch(arch string) PullOption {
	return targetArchOption(arch)
}

// mediaTypeOption is an option to specify the media type for pulling.
type mediaTypeOption string

// applyPullOption applies the media type option to PullOptions.
func (o mediaTypeOption) applyPullOption(opts *PullOptions) {
	opts.mediaType = string(o)
}

// WithMediaType returns a PullOption that sets the media type.
func WithMediaType(mediaType string) PullOption {
	return mediaTypeOption(mediaType)
}

// artifactTypeOption is an option to specify the artifact type for pulling.
type artifactTypeOption string

// applyPullOption applies the artifact type option to PullOptions.
func (o artifactTypeOption) applyPullOption(opts *PullOptions) {
	opts.artifactType = string(o)
}

// WithArtifactType returns a PullOption that sets the artifact type.
func WithArtifactType(artifactType string) PullOption {
	return artifactTypeOption(artifactType)
}
