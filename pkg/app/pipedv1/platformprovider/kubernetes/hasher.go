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

/*
Copyright 2017 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubernetes

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// HashManifests computes the hash of a list of manifests.
func HashManifests(manifests []Manifest) (string, error) {
	if len(manifests) == 0 {
		return "", errors.New("no manifest to hash")
	}

	hasher := sha256.New()
	for _, m := range manifests {
		var encoded string
		var err error

		switch {
		case m.Key.IsConfigMap():
			obj := &v1.ConfigMap{}
			if err := m.ConvertToStructuredObject(obj); err != nil {
				return "", err
			}
			encoded, err = encodeConfigMap(obj)
		case m.Key.IsSecret():
			obj := &v1.Secret{}
			if err := m.ConvertToStructuredObject(obj); err != nil {
				return "", err
			}
			encoded, err = encodeSecret(obj)
		default:
			var encodedBytes []byte
			encodedBytes, err = m.MarshalJSON()
			encoded = string(encodedBytes)
		}

		if err != nil {
			return "", err
		}
		if _, err := hasher.Write([]byte(encoded)); err != nil {
			return "", err
		}
	}

	hex := fmt.Sprintf("%x", hasher.Sum(nil))
	return encodeHash(hex)
}

// Borrowed from https://github.com/kubernetes/kubernetes/blob/
// ea0764452222146c47ec826977f49d7001b0ea8c/staging/src/k8s.io/kubectl/pkg/util/hash/hash.go
// encodeHash extracts the first 40 bits of the hash from the hex string
// (1 hex char represents 4 bits), and then maps vowels and vowel-like hex
// characters to consonants to prevent bad words from being formed (the theory
// is that no vowels makes it really hard to make bad words). Since the string
// is hex, the only vowels it can contain are 'a' and 'e'.
// We picked some arbitrary consonants to map to from the same character set as GenerateName.
// See: https://github.com/kubernetes/apimachinery/blob/dc1f89aff9a7509782bde3b68824c8043a3e58cc/pkg/util/rand/rand.go#L75
// If the hex string contains fewer than ten characters, returns an error.
func encodeHash(hex string) (string, error) {
	if len(hex) < 10 {
		return "", errors.New("the hex string must contain at least 10 characters")
	}
	enc := []rune(hex[:10])
	for i := range enc {
		switch enc[i] {
		case '0':
			enc[i] = 'g'
		case '1':
			enc[i] = 'h'
		case '3':
			enc[i] = 'k'
		case 'a':
			enc[i] = 'm'
		case 'e':
			enc[i] = 't'
		}
	}
	return string(enc), nil
}

// Borrowed from https://github.com/kubernetes/kubernetes/blob/
// ea0764452222146c47ec826977f49d7001b0ea8c/staging/src/k8s.io/kubectl/pkg/util/hash/hash.go
// encodeConfigMap encodes a ConfigMap.
// Data, Kind, and Name are taken into account.
func encodeConfigMap(cm *v1.ConfigMap) (string, error) {
	// json.Marshal sorts the keys in a stable order in the encoding
	m := map[string]interface{}{
		"kind": "ConfigMap",
		"name": cm.Name,
		"data": cm.Data,
	}
	if cm.Immutable != nil {
		m["immutable"] = *cm.Immutable
	}
	if len(cm.BinaryData) > 0 {
		m["binaryData"] = cm.BinaryData
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Borrowed from https://github.com/kubernetes/kubernetes/blob/
// ea0764452222146c47ec826977f49d7001b0ea8c/staging/src/k8s.io/kubectl/pkg/util/hash/hash.go
// encodeSecret encodes a Secret.
// Data, Kind, Name, and Type are taken into account.
func encodeSecret(sec *v1.Secret) (string, error) {
	m := map[string]interface{}{
		"kind": "Secret",
		"type": sec.Type,
		"name": sec.Name,
		"data": sec.Data,
	}
	if sec.Immutable != nil {
		m["immutable"] = *sec.Immutable
	}
	// json.Marshal sorts the keys in a stable order in the encoding
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
