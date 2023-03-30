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

package git

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/pipe-cd/pipecd/pkg/config"
)

const (
	defaultHost = "github.com"

	sshConfigTemplate = `
Host {{ .Host }}
    Hostname {{ .HostName }}
    User git
    IdentityFile {{ .IdentityFile }}
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
`
)

var (
	sshConfigTmpl = template.Must(template.New("ssh-config").Parse(sshConfigTemplate))
)

type sshConfig struct {
	Host         string
	HostName     string
	IdentityFile string
}

func AddSSHConfig(cfg config.PipedGit) error {
	cfgPath := cfg.SSHConfigFilePath
	if cfgPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to detect the current user's home directory: %w", err)
		}
		cfgPath = path.Join(home, ".ssh", "config")
	}
	sshDir := filepath.Dir(cfgPath)

	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create a directory %s: %v", sshDir, err)
	}

	sshKeys, err := cfg.LoadSSHKeys()
	if err != nil {
		return err
	}

	for _, sshKey := range sshKeys {
		// TODO: Ensure that names do not conflict
		sshKeyFile, err := os.CreateTemp(sshDir, "piped-ssh-key-*")
		if err != nil {
			return err
		}

		// TODO: Remove this key file when Piped terminating.
		if _, err := sshKeyFile.Write(sshKey); err != nil {
			return err
		}

		configData, err := generateSSHConfig(cfg, sshKeyFile.Name())
		if err != nil {
			return err
		}

		f, err := os.OpenFile(cfgPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("could not create/append to %s: %v", cfgPath, err)
		}
		defer f.Close()

		if _, err := f.Write([]byte(configData)); err != nil {
			return fmt.Errorf("failed to write sshConfig to %s: %v", cfgPath, err)
		}
	}
	return nil
}

func generateSSHConfig(cfg config.PipedGit, sshKeyFile string) (string, error) {
	var (
		buffer bytes.Buffer
		data   = sshConfig{
			Host:         defaultHost,
			IdentityFile: sshKeyFile,
		}
	)

	if cfg.Host != "" {
		data.Host = cfg.Host
	}
	if cfg.HostName != "" {
		data.HostName = cfg.HostName
	} else {
		data.HostName = data.Host
	}

	if err := sshConfigTmpl.Execute(&buffer, data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
