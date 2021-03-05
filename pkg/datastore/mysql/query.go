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

package mysql

import "fmt"

func buildGetQuery(table string) string {
	// TODO: Make kinds from datastore package public
	if table == "Project" {
		return fmt.Sprintf("SELECT data FROM %s WHERE data->>\"$.id\" = ?", table)
	}
	return fmt.Sprintf("SELECT data FROM %s WHERE id = UUID_TO_BIN(?,true)", table)
}

func buildUpdateQuery(table string) string {
	if table == "Project" {
		return fmt.Sprintf("UPDATE %s SET data = ? WHERE data->>\"$.id\" = ?", table)
	}
	return fmt.Sprintf("UPDATE %s SET data = ? WHERE id = UUID_TO_BIN(?,true)", table)
}

func buildPutQuery(table string) string {
	return fmt.Sprintf("INSERT INTO %s (id, data) VALUE (UUID_TO_BIN(?,true), ?) ON DUPLICATE KEY UPDATE data = ?", table)
}

func buildCreateQuery(table string) string {
	return fmt.Sprintf("INSERT INTO %s (id, data) VALUE (UUID_TO_BIN(?,true), ?)", table)
}
