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

package lambda

import (
	"fmt"
	"io/ioutil"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

// FunctionManifest contains configuration for LambdaFunction.
type FunctionManifest struct {
	Name     string
	ImageURI string
}

func loadFunctionManifest(path string) (FunctionManifest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return FunctionManifest{}, err
	}
	return parseFunctionManifest(data)
}

func parseFunctionManifest(data []byte) (FunctionManifest, error) {
	var obj unstructured.Unstructured
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return FunctionManifest{}, err
	}

	imageURI, ok, err := unstructured.NestedString(obj.Object, "spec", "template", "spec", "image")
	if err != nil {
		return FunctionManifest{}, err
	}
	if !ok || imageURI == "" {
		return FunctionManifest{}, fmt.Errorf("spec.template.spec.image is missing")
	}

	return FunctionManifest{
		Name:     obj.GetName(),
		ImageURI: imageURI,
	}, nil
}

// DecideRevisionName returns revision name to apply.
func DecideRevisionName(fm FunctionManifest, commit string) (string, error) {
	tag, err := FindImageTag(fm)
	if err != nil {
		return "", err
	}
	tag = strings.ReplaceAll(tag, ".", "")

	if len(commit) > 7 {
		commit = commit[:7]
	}
	return fmt.Sprintf("%s-%s-%s", fm.Name, tag, commit), nil
}

// FindImageTag parses image tag from given LambdaFunction manifest.
func FindImageTag(fm FunctionManifest) (string, error) {
	name, tag := parseContainerImage(fm.ImageURI)
	if name == "" {
		return "", fmt.Errorf("image name could not be empty")
	}
	return tag, nil
}

func parseContainerImage(image string) (name, tag string) {
	parts := strings.Split(image, ":")
	if len(parts) == 2 {
		tag = parts[1]
	}
	paths := strings.Split(parts[0], "/")
	name = paths[len(paths)-1]
	return
}
