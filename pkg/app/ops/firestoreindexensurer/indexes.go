// Copyright 2021 The PipeCD Authors.
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
	"encoding/json"
	"fmt"
)

// indexes is a mapper for the indexes JSON file emitted by "firebase firestore:indexes".
type indexes struct {
	Indexes []index `json:"indexes"`
}

type index struct {
	CollectionGroup string  `json:"collectionGroup"`
	Fields          []field `json:"fields"`
}

type field struct {
	FieldPath   string `json:"fieldPath"`
	Order       string `json:"order"`
	ArrayConfig string `json:"arrayConfig"`
}

func (i *indexes) validate() error {
	for _, idx := range i.Indexes {
		if idx.CollectionGroup == "" {
			return fmt.Errorf("index must have a collectionGroup")
		}
		for _, f := range idx.Fields {
			if f.FieldPath == "" {
				return fmt.Errorf("field must have a fieldPath")
			}
			if f.Order == "" && f.ArrayConfig == "" {
				return fmt.Errorf("field must have one of order or arrayConfig")
			}
		}
	}
	return nil
}

// parseIndexes gives back indexes made of well-defined JSON file.
func parseIndexes() (*indexes, error) {
	var out indexes
	if err := json.Unmarshal(indexesJSON, &out); err != nil {
		return nil, fmt.Errorf("failed to parse indexes list: %w", err)
	}
	if err := out.validate(); err != nil {
		return nil, err
	}
	return &out, nil
}
