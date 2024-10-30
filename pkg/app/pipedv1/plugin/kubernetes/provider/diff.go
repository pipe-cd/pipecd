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

package provider

import (
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pipe-cd/pipecd/pkg/plugin/diff"
)


func Diff(old, new Manifest, logger *zap.Logger, opts ...diff.Option) (*diff.Result, error) {
	if old.Key.IsSecret() && new.Key.IsSecret() {
		var err error
		old.Body, err = normalizeNewSecret(old.Body, new.Body)
		if err != nil {
			return nil, err
		}
	}

	key := old.Key.String()

	normalizedOld, err := remarshal(old.Body)
	if err != nil {
		logger.Info("compare manifests directly since it was unable to remarshal old Kubernetes manifest to normalize special fields", zap.Error(err))
		return diff.DiffUnstructureds(*old.Body, *new.Body, key, opts...)
	}

	normalizedNew, err := remarshal(new.Body)
	if err != nil {
		logger.Info("compare manifests directly since it was unable to remarshal new Kubernetes manifest to normalize special fields", zap.Error(err))
		return diff.DiffUnstructureds(*old.Body, *new.Body, key, opts...)
	}

	return diff.DiffUnstructureds(*normalizedOld, *normalizedNew, key, opts...)
}

func normalizeNewSecret(old, new *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	var o, n v1.Secret
	runtime.DefaultUnstructuredConverter.FromUnstructured(old.Object, &o)
	runtime.DefaultUnstructuredConverter.FromUnstructured(new.Object, &n)

	// Move as much as possible fields from `o.Data` to `o.StringData` to make `o` close to `n` to minimize the diff.
	for k, v := range o.Data {
		// Skip if the field also exists in StringData.
		if _, ok := o.StringData[k]; ok {
			continue
		}

		if _, ok := n.StringData[k]; !ok {
			continue
		}

		if o.StringData == nil {
			o.StringData = make(map[string]string)
		}

		// If the field is existing in `n.StringData`, we should move that field from `o.Data` to `o.StringData`
		o.StringData[k] = string(v)
		delete(o.Data, k)
	}

	newO, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&o)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: newO}, nil
}
