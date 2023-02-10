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

package firestoreindexensurer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed indexes.json
var indexesJSON []byte

// index represents a Firestore composite index.
type index struct {
	CollectionGroup string  `json:"collectionGroup"`
	QueryScope      string  `json:"queryScope"`
	Fields          []field `json:"fields"`
}

type field struct {
	FieldPath   string `json:"fieldPath"`
	Order       string `json:"order"`
	ArrayConfig string `json:"arrayConfig"`
}

func (idx *index) validate() error {
	if idx.CollectionGroup == "" {
		return fmt.Errorf("index must have a collection group")
	}
	if len(idx.Fields) == 0 {
		return fmt.Errorf("index must have a field")
	}
	for _, f := range idx.Fields {
		if f.FieldPath == "" {
			return fmt.Errorf("field must have a fieldPath")
		}
		if f.Order == "" && f.ArrayConfig == "" {
			return fmt.Errorf("field must have one of order or arrayConfig")
		}
	}
	return nil
}

// id builds a unique string based on its fields.
func (idx *index) id() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s/%s", idx.CollectionGroup, idx.QueryScope))

	fields := idx.Fields
	for _, f := range fields {
		b.WriteString(fmt.Sprintf("/field-path:%s", f.FieldPath))
		if f.Order != "" {
			b.WriteString(fmt.Sprintf("/order:%s", f.Order))
		}
		if f.ArrayConfig != "" {
			b.WriteString(fmt.Sprintf("/array-config:%s", f.ArrayConfig))
		}
	}
	return b.String()
}

// parseIndexes gives back indexes made of well-defined JSON file.
func parseIndexes() ([]index, error) {
	var out []index
	if err := json.Unmarshal(indexesJSON, &out); err != nil {
		return nil, fmt.Errorf("failed to parse indexes list: %w", err)
	}
	for _, i := range out {
		if err := i.validate(); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// filterIndexes gives back a new slice excluding the given one.
// TODO: Enable to extract indexes to be deleted
func filterIndexes(indexes []index, excludes []index) []index {
	if len(indexes) == 0 || len(excludes) == 0 {
		return indexes
	}

	blackList := make(map[string]struct{}, len(excludes))
	for _, e := range excludes {
		blackList[e.id()] = struct{}{}
	}

	filtered := make([]index, 0, len(indexes))
	for _, idx := range indexes {
		if _, ok := blackList[idx.id()]; !ok {
			filtered = append(filtered, idx)
		}
	}
	return filtered
}

// prefixIndexes adds the given prefix to all indexes.
// Currently, it just prefixes the CollectionGroup.
func prefixIndexes(indexes []index, prefix string) {
	for i := range indexes {
		indexes[i].CollectionGroup = prefix + indexes[i].CollectionGroup
	}
}
