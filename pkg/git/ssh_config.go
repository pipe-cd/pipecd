// Copyright 2020 The PipeCD Authors.
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
	"path/filepath"
	"text/template"

	"github.com/pipe-cd/pipe/pkg/config"
)

const (
	defaultSSHConfigFilePath = "/etc/ssh/ssh_config"
	defaultHost              = "github.com"

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
	// Check the existence of the specified private SSH key file.
	if _, err := os.Stat(cfg.SSHKeyFile); os.IsNotExist(err) {
		return fmt.Errorf("the specified private SSH key at %s was not found", cfg.SSHKeyFile)
	}

	configData, err := generateSSHConfig(cfg)
	if err != nil {
		return err
	}

	path := cfg.SSHConfigFilePath
	if path == "" {
		path = defaultSSHConfigFilePath
	}
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create a directory %s: %v", dir, err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not create/append to %s: %v", path, err)
	}
	defer f.Close()

	if _, err := f.Write([]byte(configData)); err != nil {
		return fmt.Errorf("failed to write sshConfig to %s: %v", path, err)
	}

	return nil
}

func generateSSHConfig(cfg config.PipedGit) (string, error) {
	var (
		buffer bytes.Buffer
		data   = sshConfig{
			Host:         defaultHost,
			IdentityFile: cfg.SSHKeyFile,
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
