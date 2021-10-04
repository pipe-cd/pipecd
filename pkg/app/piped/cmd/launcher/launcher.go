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

package launcher

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/credentials"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/config"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcauth"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcclient"
	"github.com/pipe-cd/pipe/pkg/version"
)

const (
	pipedDownloadURL    = "https://github.com/pipe-cd/pipe/releases/download/%s/piped_%s_%s_amd64"
	pipedBinaryFileName = "piped"
	pipedConfigFileName = "piped-config.yaml"
)

type launcher struct {
	configFile              string
	configData              string
	configFromGCPSecret     bool
	gcpSecretID             string
	configFromGitRepo       bool
	gitRepoURL              string
	configFilePathInGitRepo string
	insecure                bool
	certFile                string
	homeDir                 string
	defaultVersion          string
	checkInterval           time.Duration
	gracePeriod             time.Duration

	runningVersion    string
	runningConfigData []byte

	clientKey string
	client    pipedservice.Client
}

func NewCommand() *cobra.Command {
	l := &launcher{
		checkInterval: time.Minute,
		gracePeriod:   30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "launcher",
		Short: "Start running Piped launcher.",
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		RunE: cli.WithContext(l.run),
	}

	cmd.Flags().StringVar(&l.configFile, "config-file", l.configFile, "The path to the configuration file.")
	cmd.Flags().StringVar(&l.configData, "config-data", l.configData, "The configuration data in YAML/JSON format.")

	cmd.Flags().BoolVar(&l.configFromGCPSecret, "config-from-gcp-secret", l.configFromGCPSecret, "Whether to load Piped config that is being stored in GCP SecretManager service.")
	cmd.Flags().StringVar(&l.gcpSecretID, "gcp-secret-id", l.gcpSecretID, "The resource ID of secret that contains Piped config in GCP SecretManager service.")

	cmd.Flags().BoolVar(&l.configFromGitRepo, "config-from-git-repo", l.configFromGitRepo, "Whether to load Piped config that is being stored in a Git repository.")
	cmd.Flags().StringVar(&l.gitRepoURL, "git-repo-url", l.gitRepoURL, "The remote URL of Git repository to fetch Piped configuration.")

	cmd.Flags().BoolVar(&l.insecure, "insecure", l.insecure, "Whether disabling transport security while connecting to control-plane.")
	cmd.Flags().StringVar(&l.certFile, "cert-file", l.certFile, "The path to the TLS certificate file.")

	cmd.Flags().StringVar(&l.homeDir, "home-dir", l.homeDir, "The working directory of Launcher.")
	cmd.Flags().StringVar(&l.defaultVersion, "default-version", l.defaultVersion, "The version should be run when no desired version was specified. Empty means using the same version with Launcher.")

	cmd.Flags().DurationVar(&l.checkInterval, "check-interval", l.checkInterval, "Interval to periodically check desired config/version to restart Piped. Default is 1m.")
	cmd.Flags().DurationVar(&l.gracePeriod, "grace-period", l.gracePeriod, "How long to wait for graceful shutdown.")

	return cmd
}

func (l *launcher) run(ctx context.Context, t cli.Telemetry) error {
	if l.homeDir == "" {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			t.Logger.Error("LAUNCHER: failed to get the user's cache directory", zap.Error(err))
			return err
		}
		l.homeDir = filepath.Join(userCacheDir, "piped-launcher")
	}

	var (
		runningPiped *command
		workingDir   = filepath.Join(l.homeDir, "piped")
		ticker       = time.NewTicker(l.checkInterval)
	)

	execute := func() error {
		version, config, relaunch, err := l.shouldRelaunch(ctx, t.Logger)
		if err != nil {
			t.Logger.Error("LAUNCHER: failed while checking desired version and config",
				zap.String("version", version),
				zap.Error(err),
			)
			return err
		}

		if !relaunch {
			if runningPiped != nil && runningPiped.IsRunning() {
				t.Logger.Info("LAUNCHER: everything up-to-date", zap.String("version", l.runningVersion))
				return nil
			}
			t.Logger.Warn("LAUNCHER: it seems the launched Piped has stopped unexpectedly")
		}
		t.Logger.Info("LAUNCHER: will relaunch a new Piped because some changes in version/config were detected")

		// Stop old piped process and clean its data.
		if err := l.cleanOldPiped(runningPiped, workingDir, t.Logger); err != nil {
			t.Logger.Error("LAUNCHER: failed while cleaning old Piped",
				zap.String("version", version),
				zap.Error(err),
			)
			return err
		}

		// Start new piped process.
		runningPiped, err = l.launchNewPiped(version, config, workingDir, t.Logger)
		if err != nil {
			t.Logger.Error("LAUNCHER: failed while launching new Piped", zap.Error(err))
			return err
		}

		l.runningVersion = version
		l.runningConfigData = config
		t.Logger.Info("LAUNCHER: successfully launched a new Piped", zap.String("version", version))
		return nil
	}

	// Execute the first time immediately.
	execute()

	for {
		select {
		case <-ticker.C:
			execute()

		case <-ctx.Done():
			// Stop old piped process and clean its data.
			if err := l.cleanOldPiped(runningPiped, workingDir, t.Logger); err != nil {
				t.Logger.Error("LAUNCHER: failed while cleaning old Piped",
					zap.String("version", l.runningVersion),
					zap.Error(err),
				)
				return err
			}
			return nil
		}
	}

	return nil
}

// shouldRelaunch fetches the latest state of desired version and config
// to determine whether a new Piped should be launched or not.
// This also returns the desired version and config.
func (l *launcher) shouldRelaunch(ctx context.Context, logger *zap.Logger) (version string, config []byte, should bool, err error) {
	config, err = l.loadConfigData(ctx)
	if err != nil {
		return
	}

	cfg, err := parseConfig(config)
	if err != nil {
		return
	}

	pipedKey, err := cfg.LoadPipedKey()
	if err != nil {
		return
	}

	version, err = l.getDesiredVersion(ctx, cfg.APIAddress, cfg.ProjectID, cfg.PipedID, pipedKey, logger)
	if err != nil {
		return
	}

	should = version != l.runningVersion || !bytes.Equal(config, l.runningConfigData)
	return
}

func (l *launcher) cleanOldPiped(cmd *command, workingDir string, logger *zap.Logger) error {
	// Stop running Piped gracefully.
	if cmd != nil {
		if err := cmd.GracefulStop(l.gracePeriod); err != nil {
			// We just log the error and continue to the next part
			// because the process was killed after all even if an error occured.
			logger.Warn("LAUNCHER: received an error while shutting down old Piped", zap.Error(err))
		}
		logger.Info("LAUNCHER: piped has been stopped")
	}

	// Clean old data.
	if err := os.RemoveAll(workingDir); err != nil {
		return fmt.Errorf("could not clean working directory %s (%w)", workingDir, err)
	}

	return nil
}

func (l *launcher) launchNewPiped(version string, config []byte, workingDir string, logger *zap.Logger) (*command, error) {
	if err := os.MkdirAll(workingDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create working directory %s (%w)", workingDir, err)
	}

	// Download Piped binary into working directory.
	var (
		binaryDir   = filepath.Join(workingDir, "bin")
		downloadURL = makeDownloadURL(version)
	)
	pipedPath, err := downloadBinary(downloadURL, binaryDir, pipedBinaryFileName, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to download Piped from %s to %s (%w)", downloadURL, binaryDir, err)
	}
	logger.Info(fmt.Sprintf("LAUNCHER: downloaded Piped binary from %s to %s", downloadURL, pipedPath))

	// Write Piped configuration into working directory.
	var (
		configFileDir  = filepath.Join(workingDir, "config")
		configFilePath = filepath.Join(configFileDir, pipedConfigFileName)
	)
	if err := os.MkdirAll(configFileDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s (%w)", configFileDir, err)
	}
	if err := os.WriteFile(configFilePath, config, 0644); err != nil {
		return nil, fmt.Errorf("failed to write Piped config data to file at %s (%w)", configFilePath, err)
	}
	logger.Info(fmt.Sprintf("LAUNCHER: wrote Piped config to %s", configFilePath))

	args := makePipedArgs(os.Args[2:], configFilePath)
	logger.Info(fmt.Sprintf("LAUNCHER: start running Piped %s with args %v", version, args))

	return runBinary(pipedPath, args)
}

func (l *launcher) loadConfigData(ctx context.Context) ([]byte, error) {
	// Load config data from the specified file.
	if l.configFile != "" {
		return os.ReadFile(l.configFile)
	}

	// Return config data passed directly.
	if l.configData != "" {
		return []byte(l.configData), nil
	}

	// Load config data from a secret which is stored in Google Cloud Secret Manager service.
	if l.configFromGCPSecret {
		if l.gcpSecretID == "" {
			return nil, fmt.Errorf("gcp-secret-id is required to load Piped config from GCP SecretManager service")
		}

		client, err := secretmanager.NewClient(ctx)
		if err != nil {
			return nil, err
		}
		defer client.Close()

		req := &secretmanagerpb.AccessSecretVersionRequest{
			Name: l.gcpSecretID,
		}
		resp, err := client.AccessSecretVersion(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp.Payload.Data, nil
	}

	if l.configFromGitRepo {
		// TODO: Support loading config data from a Git repository.
	}

	return nil, fmt.Errorf("either [%s] must be set", strings.Join([]string{
		"config-file",
		"config-data",
		"config-from-gcp-secret",
		"config-from-git-repo",
	}, ", "))
}

func (l *launcher) getDesiredVersion(ctx context.Context, address, projectID, pipedID string, pipedKey []byte, logger *zap.Logger) (string, error) {
	clientKey := fmt.Sprintf("%s,%s,%s,%s", address, projectID, pipedID, string(pipedKey))

	// In order to reduce the time of initializing gRPC client
	// we reuse the client when no configuration changes occurred.
	if clientKey != l.clientKey {
		client, err := l.createAPIClient(ctx, address, projectID, pipedID, pipedKey)
		if err != nil {
			logger.Error("LAUNCHER: failed to create api client", zap.Error(err))
			return "", err
		}
		l.clientKey = clientKey
		l.client = client
	}

	resp, err := l.client.GetDesiredVersion(ctx, &pipedservice.GetDesiredVersionRequest{})
	if err != nil {
		return "", err
	}
	if resp.Version != "" {
		return resp.Version, nil
	}

	if l.defaultVersion != "" {
		return l.defaultVersion, nil
	}
// Using launcher version if there is no runner version is set.
	return version.Get().Version, nil
}

func (l *launcher) createAPIClient(ctx context.Context, address, projectID, pipedID string, pipedKey []byte) (pipedservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var (
		token   = rpcauth.MakePipedToken(projectID, pipedID, string(pipedKey))
		creds   = rpcclient.NewPerRPCCredentials(token, rpcauth.PipedTokenCredentials, !l.insecure)
		options = []rpcclient.DialOption{
			rpcclient.WithBlock(),
			rpcclient.WithPerRPCCredentials(creds),
		}
	)

	if !l.insecure {
		if l.certFile != "" {
			options = append(options, rpcclient.WithTLS(l.certFile))
		} else {
			config := &tls.Config{}
			options = append(options, rpcclient.WithTransportCredentials(credentials.NewTLS(config)))
		}
	} else {
		options = append(options, rpcclient.WithInsecure())
	}

	return pipedservice.NewClient(ctx, address, options...)
}

// makePipedArgs generates arguments for Piped from the ones passed to Launcher.
func makePipedArgs(launcherArgs []string, configFile string) []string {
	ignoreArgs := map[string]struct{}{
		"config-file":            struct{}{},
		"config-data":            struct{}{},
		"config-from-gcp-secret": struct{}{},
		"gcp-secret-id":          struct{}{},
		"config-from-git-repo":   struct{}{},
		"git-repo-url":           struct{}{},
		"home-dir":               struct{}{},
		"default-version":        struct{}{},
		"check-interval":         struct{}{},
	}

	pipedArgs := make([]string, 0, len(launcherArgs)+2)
	pipedArgs = append(pipedArgs, "piped", "--config-file="+configFile)

	for _, a := range launcherArgs {
		normalizedArg := strings.TrimLeft(a, "-")
		parts := strings.SplitN(normalizedArg, "=", 2)
		name := parts[0]

		if _, ok := ignoreArgs[name]; !ok {
			pipedArgs = append(pipedArgs, a)
		}
	}

	return pipedArgs
}

func parseConfig(data []byte) (*config.PipedSpec, error) {
	cfg, err := config.DecodeYAML(data)
	if err != nil {
		return nil, err
	}
	if cfg.Kind != config.KindPiped {
		return nil, fmt.Errorf("wrong configuration kind for piped: %v", cfg.Kind)
	}
	return cfg.PipedSpec, nil
}

func makeDownloadURL(version string) string {
	return fmt.Sprintf(pipedDownloadURL, version, version, runtime.GOOS)
}
