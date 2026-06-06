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

package git

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	configv1 "github.com/pipe-cd/pipecd/pkg/configv1"
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

func AddSSHConfig(cfg configv1.PipedGit) (string, error) {
	cfgPath := cfg.SSHConfigFilePath
	if cfgPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to detect the current user's home directory: %w", err)
		}
		cfgPath = path.Join(home, ".ssh", "config")
	}
	sshDir := filepath.Dir(cfgPath)

	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create a directory %s: %v", sshDir, err)
	}

	tempDir, err := os.MkdirTemp(sshDir, "piped-ssh-keys-*")
	if err != nil {
		return "", err
	}
	needCleanUp := false
	defer func() {
		if needCleanUp {
			os.RemoveAll(tempDir)
		}
	}()

	var configsData []byte

	if cfg.SSHKeyData != "" || cfg.SSHKeyFile != "" {
		sshKey, err := cfg.LoadSSHKey()
		if err != nil {
			needCleanUp = true
			return "", err
		}
		keyPath, err := writeSSHKeyFile(tempDir, sshKey)
		if err != nil {
			needCleanUp = true
			return "", err
		}
		configData, err := generateSSHConfig(sshConfig{
			Host:         cfg.Host,
			HostName:     cfg.HostName,
			IdentityFile: keyPath,
		})
		if err != nil {
			needCleanUp = true
			return "", err
		}
		configsData = append(configsData, []byte(configData)...)
	}

	for _, key := range cfg.SSHKeys {
		sshKey, err := key.LoadSSHKey()
		if err != nil {
			needCleanUp = true
			return "", err
		}
		keyPath, err := writeSSHKeyFile(tempDir, sshKey)
		if err != nil {
			needCleanUp = true
			return "", err
		}
		configData, err := generateSSHConfig(sshConfig{
			Host:         key.Host,
			HostName:     key.HostName,
			IdentityFile: keyPath,
		})
		if err != nil {
			needCleanUp = true
			return "", err
		}
		configsData = append(configsData, []byte(configData)...)
	}

	f, err := os.OpenFile(cfgPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		needCleanUp = true
		return "", fmt.Errorf("could not create/append to %s: %v", cfgPath, err)
	}
	defer f.Close()

	if _, err := f.Write(configsData); err != nil {
		needCleanUp = true
		return "", fmt.Errorf("failed to write sshConfig to %s: %v", cfgPath, err)
	}

	return tempDir, nil
}

func generateSSHConfig(data sshConfig) (string, error) {
	var buffer bytes.Buffer

	if data.Host == "" {
		data.Host = defaultHost
	}
	if data.HostName == "" {
		data.HostName = data.Host
	}

	if err := sshConfigTmpl.Execute(&buffer, data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func writeSSHKeyFile(dir string, key []byte) (string, error) {
	f, err := os.CreateTemp(dir, "piped-ssh-key-*")
	if err != nil {
		return "", err
	}
	if err := os.Chmod(f.Name(), 0600); err != nil {
		f.Close()
		return "", err
	}
	_, err = f.Write(key)
	f.Close()
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}
