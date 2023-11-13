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

package config

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"

	"github.com/pipe-cd/pipecd/pkg/model"
)

// ControlPlaneSpec defines all configuration for all control-plane components.
type ControlPlaneSpec struct {
	// The address to the control plane.
	// This is required if SSO is enabled.
	Address string `json:"address"`
	// A randomly generated string used to sign oauth state.
	StateKey string `json:"stateKey"`
	// The configuration of datastore for control plane.
	Datastore ControlPlaneDataStore `json:"datastore"`
	// The configuration of filestore for control plane.
	Filestore ControlPlaneFileStore `json:"filestore"`
	// The configuration of cache for control plane.
	Cache ControlPlaneCache `json:"cache"`
	// The configuration of insight collector.
	InsightCollector ControlPlaneInsightCollector `json:"insightCollector"`
	// List of debugging/quickstart projects defined in Control Plane configuration.
	// Please note that do not use this to configure the projects running in the production.
	Projects []ControlPlaneProject `json:"projects"`
	// List of shared SSO configurations that can be used by any projects.
	SharedSSOConfigs []SharedSSOConfig `json:"sharedSSOConfigs"`
}

func (s *ControlPlaneSpec) Validate() error {
	return nil
}

type ControlPlaneProject struct {
	// The unique identifier of the project.
	ID string `json:"id"`
	// The description about the project.
	Desc string `json:"desc"`
	// Static admin account of the project.
	StaticAdmin ProjectStaticUser `json:"staticAdmin"`
}

type ProjectStaticUser struct {
	// The username string.
	Username string `json:"username"`
	// The bcrypt hashsed value of the password string.
	PasswordHash string `json:"passwordHash"`
}

type SharedSSOConfig struct {
	model.ProjectSSOConfig `json:",inline"`
	Name                   string `json:"name"`
}

func (s *SharedSSOConfig) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	provider := m["provider"].(string)
	v, ok := model.ProjectSSOConfig_Provider_value[provider]
	if !ok {
		return fmt.Errorf("unsupported provider %s", provider)
	}
	m["provider"] = v

	name, ok := m["name"]
	if !ok {
		return fmt.Errorf("name field in SharedSSOConfig is required")
	}
	s.Name = name.(string)
	delete(m, "name")

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Using jsonpb instead of the standard json to unmarshal because
	// json is unmarshaling with the underscored tags.
	// https://github.com/golang/protobuf/issues/183
	if err := jsonpb.UnmarshalString(string(data), &s.ProjectSSOConfig); err != nil {
		return err
	}
	return nil
}

// FindProject finds and returns a specific project in the configured list.
func (s *ControlPlaneSpec) FindProject(id string) (ControlPlaneProject, bool) {
	for i := range s.Projects {
		if s.Projects[i].ID != id {
			continue
		}
		return s.Projects[i], true
	}
	return ControlPlaneProject{}, false
}

func (s *ControlPlaneSpec) ProjectMap() map[string]ControlPlaneProject {
	m := make(map[string]ControlPlaneProject, len(s.Projects))
	for i := range s.Projects {
		m[s.Projects[i].ID] = s.Projects[i]
	}
	return m
}

func (s *ControlPlaneSpec) SharedSSOConfigMap() map[string]*model.ProjectSSOConfig {
	m := make(map[string]*model.ProjectSSOConfig, len(s.SharedSSOConfigs))
	for i := range s.SharedSSOConfigs {
		m[s.SharedSSOConfigs[i].Name] = &s.SharedSSOConfigs[i].ProjectSSOConfig
	}
	return m
}

type ControlPlaneDataStore struct {
	// The datastore type.
	Type model.DataStoreType

	// The configuration in the case of Cloud Firestore.
	FirestoreConfig *DataStoreFireStoreConfig
	// The configuration in the case of general MySQL.
	MySQLConfig *DataStoreMySQLConfig
}

type genericControlPlaneDataStore struct {
	Type   model.DataStoreType `json:"type"`
	Config json.RawMessage     `json:"config"`
}

func (d *ControlPlaneDataStore) UnmarshalJSON(data []byte) error {
	var err error
	gc := genericControlPlaneDataStore{}
	if err = json.Unmarshal(data, &gc); err != nil {
		return err
	}
	d.Type = gc.Type

	switch d.Type {
	case model.DataStoreFirestore:
		d.FirestoreConfig = &DataStoreFireStoreConfig{}
		if len(gc.Config) > 0 {
			err = json.Unmarshal(gc.Config, d.FirestoreConfig)
		}
	case model.DataStoreMySQL:
		d.MySQLConfig = &DataStoreMySQLConfig{}
		if len(gc.Config) > 0 {
			err = json.Unmarshal(gc.Config, d.MySQLConfig)
		}
	case model.DataStoreFileDB:
		// The FILEDB datastore using the same configuration with filestore
		// so there will be no `datastore.config` required for now.
		err = nil
	default:
		// Left comment out for mock response.
		// err = fmt.Errorf("unsupported datastore type: %s", d.Type)
		err = nil
	}
	return err
}

type ControlPlaneCache struct {
	TTL Duration `json:"ttl"`
}

type ControlPlaneInsightCollector struct {
	Application InsightCollectorApplication `json:"application"`
	Deployment  InsightCollectorDeployment  `json:"deployment"`
}

type InsightCollectorApplication struct {
	Enabled *bool `json:"enabled" default:"true"`
	// Default is running every hour.
	Schedule string `json:"schedule" default:"0 * * * *"`
}

type InsightCollectorDeployment struct {
	Enabled *bool `json:"enabled" default:"true"`
	// Default is running every hour.
	Schedule      string `json:"schedule" default:"30 * * * *"`
	ChunkMaxCount int    `json:"chunkMaxCount" default:"1000"`
}

func (c ControlPlaneCache) TTLDuration() time.Duration {
	const defaultTTL = 5 * time.Minute

	if c.TTL == 0 {
		return defaultTTL
	}
	return c.TTL.Duration()
}

type DataStoreFireStoreConfig struct {
	// The root path element considered as a logical namespace, e.g. `pipecd`.
	Namespace string `json:"namespace"`
	// The second path element considered as a logical environment, e.g. `dev`.
	// All pipecd collections will have path formatted according to `{namespace}/{environment}/{collection-name}`.
	Environment string `json:"environment"`
	// The prefix for collection name.
	// This can be used to avoid conflicts with existing collections in your Firestore database.
	CollectionNamePrefix string `json:"collectionNamePrefix"`
	// The name of GCP project hosting the Firestore.
	Project string `json:"project"`
	// The path to the service account file for accessing Firestores.
	CredentialsFile string `json:"credentialsFile"`
}

type DataStoreMySQLConfig struct {
	// The url of MySQL. All of credentials can be specified via this field.
	URL string `json:"url"`
	// The name of the database.
	// For those who don't want to include the database in the URL.
	Database string `json:"database"`
	// The path to the username file.
	// For those who don't want to include the username in the URL.
	UsernameFile string `json:"usernameFile"`
	// The path to the password file.
	// For those who don't want to include the password in the URL.
	PasswordFile string `json:"passwordFile"`
}

type ControlPlaneFileStore struct {
	// The filestore type.
	Type model.FileStoreType

	// The configuration in the case of Google Cloud Storage.
	GCSConfig *FileStoreGCSConfig `json:"gcs"`
	// The configuration in the case of Amazon S3.
	S3Config *FileStoreS3Config `json:"s3"`
	// The configuration in the case of Minio.
	MinioConfig *FileStoreMinioConfig `json:"minio"`
}

type genericControlPlaneFileStore struct {
	Type   model.FileStoreType `json:"type"`
	Config json.RawMessage     `json:"config"`
}

func (f *ControlPlaneFileStore) UnmarshalJSON(data []byte) error {
	var err error
	gf := genericControlPlaneFileStore{}
	if err = json.Unmarshal(data, &gf); err != nil {
		return err
	}
	f.Type = gf.Type

	switch f.Type {
	case model.FileStoreGCS:
		f.GCSConfig = &FileStoreGCSConfig{}
		if len(gf.Config) > 0 {
			err = json.Unmarshal(gf.Config, f.GCSConfig)
		}
	case model.FileStoreS3:
		f.S3Config = &FileStoreS3Config{}
		if len(gf.Config) > 0 {
			err = json.Unmarshal(gf.Config, f.S3Config)
		}
	case model.FileStoreMINIO:
		f.MinioConfig = &FileStoreMinioConfig{}
		if len(gf.Config) > 0 {
			err = json.Unmarshal(gf.Config, f.MinioConfig)
		}
	default:
		// Left comment out for mock response.
		// err = fmt.Errorf("unsupported filestore type: %s", f.Type)
		err = nil
	}
	return err
}

type FileStoreGCSConfig struct {
	// The bucket name to store artifacts and logs in the piped.
	Bucket string `json:"bucket"`
	// The path to the credentials file for accessing GCS.
	CredentialsFile string `json:"credentialsFile"`
}

type FileStoreS3Config struct {
	// The bucket name to store artifacts and logs in the piped.
	Bucket string `json:"bucket"`
	// The aws region of S3 bucket.
	Region string `json:"region"`
	// The aws profile name.
	Profile string `json:"profile"`
	// The path to the credentials file for accessing AWS.
	CredentialsFile string `json:"credentialsFile"`
	// The IAM role arn to use when assuming an role.
	RoleARN string `json:"roleARN"`
	// Path to the WebIdentity token the SDK should use to assume a role with.
	TokenFile string `json:"tokenFile"`
}

type FileStoreMinioConfig struct {
	// The address of Minio.
	Endpoint string `json:"endpoint"`
	// The bucket name to store.
	Bucket string `json:"bucket"`
	// The path to the access key file.
	AccessKeyFile string `json:"accessKeyFile"`
	// The path to the secret key file.
	SecretKeyFile string `json:"secretKeyFile"`
	// Whether the given bucket should be made automatically if not exists.
	AutoCreateBucket bool `json:"autoCreateBucket"`
}
