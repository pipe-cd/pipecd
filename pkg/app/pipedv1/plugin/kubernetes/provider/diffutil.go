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
	"bytes"
	"encoding/json"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

// All functions in this file is borrowed from argocd/gitops-engine and modified
// All function except `remarshal` is borrowed from
// https://github.com/argoproj/gitops-engine/blob/0bc2f8c395f67123156d4ce6b667bf730618307f/pkg/utils/json/json.go
// and `remarshal` function is borrowed from
// https://github.com/argoproj/gitops-engine/blob/b0c5e00ccfa5d1e73087a18dc59e2e4c72f5f175/pkg/diff/diff.go#L685-L723

// https://github.com/ksonnet/ksonnet/blob/master/pkg/kubecfg/diff.go
func removeFields(config, live interface{}) interface{} {
	switch c := config.(type) {
	case map[string]interface{}:
		l, ok := live.(map[string]interface{})
		if ok {
			return removeMapFields(c, l)
		}
		return live
	case []interface{}:
		l, ok := live.([]interface{})
		if ok {
			return removeListFields(c, l)
		}
		return live
	default:
		return live
	}

}

// removeMapFields remove all non-existent fields in the live that don't exist in the config
func removeMapFields(config, live map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v1 := range config {
		v2, ok := live[k]
		if !ok {
			continue
		}
		if v2 != nil {
			v2 = removeFields(v1, v2)
		}
		result[k] = v2
	}
	return result
}

func removeListFields(config, live []interface{}) []interface{} {
	// If live is longer than config, then the extra elements at the end of the
	// list will be returned as-is so they appear in the diff.
	result := make([]interface{}, 0, len(live))
	for i, v2 := range live {
		if len(config) > i {
			if v2 != nil {
				v2 = removeFields(config[i], v2)
			}
			result = append(result, v2)
		} else {
			result = append(result, v2)
		}
	}
	return result
}

// remarshal checks resource kind and version and re-marshal using corresponding struct custom marshaller.
// This ensures that expected resource state is formatter same as actual resource state in kubernetes
// and allows to find differences between actual and target states more accurately.
// Remarshalling also strips any type information (e.g. float64 vs. int) from the unstructured
// object. This is important for diffing since it will cause godiff to report a false difference.
func remarshal(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	item, err := scheme.Scheme.New(obj.GroupVersionKind())
	if err != nil {
		// This is common. the scheme is not registered
		return nil, err
	}
	// This will drop any omitempty fields, perform resource conversion etc...
	unmarshalledObj := reflect.New(reflect.TypeOf(item).Elem()).Interface()
	// Unmarshal data into unmarshalledObj, but detect if there are any unknown fields that are not
	// found in the target GVK object.
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&unmarshalledObj); err != nil {
		// Likely a field present in obj that is not present in the GVK type, or user
		// may have specified an invalid spec in git, so return original object
		return nil, err
	}
	unstrBody, err := runtime.DefaultUnstructuredConverter.ToUnstructured(unmarshalledObj)
	if err != nil {
		return nil, err
	}
	// Remove all default values specified by custom formatter (e.g. creationTimestamp)
	unstrBody = removeMapFields(obj.Object, unstrBody)
	return &unstructured.Unstructured{Object: unstrBody}, nil
}
