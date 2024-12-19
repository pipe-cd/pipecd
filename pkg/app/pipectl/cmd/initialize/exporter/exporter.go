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

package exporter

import (
	"fmt"
	"os"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/cmd/initialize/prompt"
)

// Export the bytes to the path.
// If the path is empty or a directory, return an error.
// If the file already exists, ask if overwrite it.
func Export(bytes []byte, path string) error {
	if len(path) == 0 {
		return fmt.Errorf("path is not specified. Please specify a file path")
	}

	// Check if the file/directory already exists
	if fInfo, err := os.Stat(path); err == nil {
		if fInfo.IsDir() {
			// When the target is a directory.
			return fmt.Errorf("the path %s is a directory. Please specify a file path", path)
		}

		// When the file exists, ask if overwrite it.
		overwrite, err := askOverwrite()
		if err != nil {
			return fmt.Errorf("invalid input for overwrite(y/n): %v", err)
		}

		if !overwrite {
			return fmt.Errorf("cancelled exporting")
		}
	}

	fmt.Printf("Start exporting to %s\n", path)
	err := os.WriteFile(path, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to export to %s: %v", path, err)
	} else {
		fmt.Printf("Successfully exported to %s\n", path)
	}
	return nil
}

func askOverwrite() (overwrite bool, err error) {
	overwriteInput := prompt.Input{
		Message:       "The file already exists. Overwrite it? [y/n]",
		TargetPointer: &overwrite,
		Required:      false,
	}
	p := prompt.NewPrompt(os.Stdin)
	err = p.Run(overwriteInput)
	if err != nil {
		return false, err
	}
	return overwrite, nil
}
