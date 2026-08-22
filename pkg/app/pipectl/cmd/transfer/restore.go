// Copyright 2026 The PipeCD Authors.
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

package transfer

import (
	"github.com/spf13/cobra"
)

func newRestoreCommand(root *command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore piped and application data to the target control plane.",
		Long: `Restore re-creates pipeds and applications on the target control plane from a backup file.

The restore process requires two steps because the control plane validates that each
application's Git repository is registered in the target piped before the application
can be created. Repository registration only happens after the piped agent connects.

Two-step workflow:

  Step 1 - Register pipeds:
    pipectl transfer restore piped --input-file=backup.json --output-file=mapping.json
    Update each piped config with the new ID and key from mapping.json, then restart the piped agents.

  Step 2 - Restore applications (after pipeds have connected and registered their repos):
    pipectl transfer restore application --input-file=backup.json --piped-id-mapping-file=mapping.json`,
	}

	cmd.AddCommand(newRestorePipedCommand(root))
	cmd.AddCommand(newRestoreApplicationCommand(root))

	return cmd
}
