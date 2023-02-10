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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

const defaultGcloudPath = "gcloud"

type gcloud struct {
	// The path to executable of gcloud.
	gcloudPath string
	// The Google Cloud Platform project ID
	projectID string
	// The path to the service account key file.
	serviceAccountFile string

	logger *zap.Logger
}

func newGcloud(gcloudPath, projectID, serviceAcccountFile string, logger *zap.Logger) *gcloud {
	if gcloudPath == "" {
		gcloudPath = defaultGcloudPath
	}
	return &gcloud{
		gcloudPath:         gcloudPath,
		projectID:          projectID,
		serviceAccountFile: serviceAcccountFile,
		logger:             logger.Named("gcloud-client"),
	}
}

func (c *gcloud) authorize(ctx context.Context) error {
	// Extracts client email from the service account file at first.
	serviceAccount := struct {
		ClientEmail string `json:"client_email"`
	}{}
	b, err := os.ReadFile(c.serviceAccountFile)
	if err != nil {
		return fmt.Errorf("failed to open service account file: %w", err)
	}
	if err := json.Unmarshal(b, &serviceAccount); err != nil {
		return fmt.Errorf("failed to parse service account: %w", err)
	}
	if serviceAccount.ClientEmail == "" {
		return fmt.Errorf("missing \"client_email\" field in service account file")
	}

	_, err = c.runGcloudCommand(ctx, "auth", "activate-service-account", serviceAccount.ClientEmail, "--key-file", c.serviceAccountFile)
	return err
}

func (c *gcloud) createIndex(ctx context.Context, idx *index) error {
	// TODO: Track the progress of Firebase indexes creation and ensure the operation in progress to complete
	// For that, seems like additional permission is required. We have to look out for.

	// Run gcloud command in async mode, which returns immediately without waiting for the operation in progress to complete.
	args := []string{
		"firestore", "indexes", "composite", "create",
		"--async",
		"--project", c.projectID,
		"--collection-group", idx.CollectionGroup,
	}
	for _, f := range idx.Fields {
		fieldCfg := fmt.Sprintf("field-path=%s", f.FieldPath)
		if f.Order != "" {
			fieldCfg += fmt.Sprintf(",order=%s", f.Order)
		}
		if f.ArrayConfig != "" {
			fieldCfg += fmt.Sprintf(",array-config=%s", f.ArrayConfig)
		}
		args = append(args, "--field-config", fieldCfg)
	}

	c.logger.Info("start creating a Firestore index", zap.Strings("command", args))
	if _, err := c.runGcloudCommand(ctx, args...); err != nil {
		return err
	}
	return nil
}

func (c *gcloud) listIndexes(ctx context.Context) ([]index, error) {
	type rawIndex struct {
		Name       string  `json:"name"`
		QueryScope string  `json:"queryScope"`
		Fields     []field `json:"fields"`
	}
	args := []string{
		"firestore", "indexes", "composite", "list",
		"--sort-by", "name",
		"--format", "json",
		"--project", c.projectID,
	}
	out, err := c.runGcloudCommand(ctx, args...)
	if err != nil {
		return nil, err
	}
	var rawIndexes []rawIndex
	if err := json.Unmarshal(out, &rawIndexes); err != nil {
		return nil, fmt.Errorf("failed to parse indexes list: %w", err)
	}

	// Start converting the raw indexes into our own index type.
	indexes := make([]index, 0, len(rawIndexes))
	for _, idx := range rawIndexes {
		// Supposed to be like "projects/project-a/databases/(default)/collectionGroups/CollectionA/indexes/CICAgLjRnZMK"
		name := strings.Split(idx.Name, "/")
		if len(name) < 6 {
			c.logger.Warn("index has unexpected name", zap.String("name", idx.Name))
			continue
		}
		// Ignore the "__name__" field which is automatically created by Firestore.
		fields := make([]field, 0, len(idx.Fields)-1)
		for _, f := range idx.Fields {
			if f.FieldPath == "__name__" {
				continue
			}
			fields = append(fields, f)
		}
		indexes = append(indexes, index{
			CollectionGroup: name[5],
			QueryScope:      idx.QueryScope,
			Fields:          fields,
		})
	}
	return indexes, nil
}

func (c *gcloud) runGcloudCommand(ctx context.Context, args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, c.gcloudPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run gcloud: stderr: %s: err: %w", stderr.String(), err)
	}
	return stdout.Bytes(), nil
}
