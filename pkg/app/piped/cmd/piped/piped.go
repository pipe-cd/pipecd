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

package piped

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/credentials"
	"sigs.k8s.io/yaml"

	"github.com/pipe-cd/pipecd/pkg/admin"
	"github.com/pipe-cd/pipecd/pkg/app/piped/apistore/analysisresultstore"
	"github.com/pipe-cd/pipecd/pkg/app/piped/apistore/applicationstore"
	"github.com/pipe-cd/pipecd/pkg/app/piped/apistore/commandstore"
	"github.com/pipe-cd/pipecd/pkg/app/piped/apistore/deploymentstore"
	"github.com/pipe-cd/pipecd/pkg/app/piped/apistore/eventstore"
	"github.com/pipe-cd/pipecd/pkg/app/piped/appconfigreporter"
	"github.com/pipe-cd/pipecd/pkg/app/piped/chartrepo"
	"github.com/pipe-cd/pipecd/pkg/app/piped/controller"
	"github.com/pipe-cd/pipecd/pkg/app/piped/controller/controllermetrics"
	"github.com/pipe-cd/pipecd/pkg/app/piped/driftdetector"
	"github.com/pipe-cd/pipecd/pkg/app/piped/eventwatcher"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatereporter"
	"github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore"
	k8slivestatestoremetrics "github.com/pipe-cd/pipecd/pkg/app/piped/livestatestore/kubernetes/kubernetesmetrics"
	"github.com/pipe-cd/pipecd/pkg/app/piped/notifier"
	"github.com/pipe-cd/pipecd/pkg/app/piped/planpreview"
	"github.com/pipe-cd/pipecd/pkg/app/piped/planpreview/planpreviewmetrics"
	k8scloudprovidermetrics "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/kubernetes/kubernetesmetrics"
	"github.com/pipe-cd/pipecd/pkg/app/piped/statsreporter"
	"github.com/pipe-cd/pipecd/pkg/app/piped/toolregistry"
	"github.com/pipe-cd/pipecd/pkg/app/piped/trigger"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/pipedservice"
	"github.com/pipe-cd/pipecd/pkg/cache/memorycache"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/crypto"
	"github.com/pipe-cd/pipecd/pkg/git"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
	"github.com/pipe-cd/pipecd/pkg/version"

	// Import to preload all built-in executors to the default registry.
	_ "github.com/pipe-cd/pipecd/pkg/app/piped/executor/registry"
	// Import to preload all planners to the default registry.
	_ "github.com/pipe-cd/pipecd/pkg/app/piped/planner/registry"
)

const (
	commandCheckPeriod time.Duration = 30 * time.Second
)

type piped struct {
	configFile      string
	configData      string
	configGCPSecret string

	insecure                             bool
	certFile                             string
	adminPort                            int
	toolsDir                             string
	enableDefaultKubernetesCloudProvider bool
	gracePeriod                          time.Duration
	addLoginUserToPasswd                 bool
	launcherVersion                      string
	maxRecvMsgSize                       int
}

func NewCommand() *cobra.Command {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to detect the current user's home directory: %v", err))
	}
	p := &piped{
		adminPort:      9085,
		toolsDir:       path.Join(home, ".piped", "tools"),
		gracePeriod:    30 * time.Second,
		maxRecvMsgSize: 1024 * 1024 * 10, // 10MB
	}
	cmd := &cobra.Command{
		Use:   "piped",
		Short: "Start running piped.",
		RunE:  cli.WithContext(p.run),
	}

	cmd.Flags().StringVar(&p.configFile, "config-file", p.configFile, "The path to the configuration file.")
	cmd.Flags().StringVar(&p.configData, "config-data", p.configData, "The base64 encoded string of the configuration data.")
	cmd.Flags().StringVar(&p.configGCPSecret, "config-gcp-secret", p.configGCPSecret, "The resource ID of secret that contains Piped config and be stored in GCP SecretManager.")

	cmd.Flags().BoolVar(&p.insecure, "insecure", p.insecure, "Whether disabling transport security while connecting to control-plane.")
	cmd.Flags().StringVar(&p.certFile, "cert-file", p.certFile, "The path to the TLS certificate file.")
	cmd.Flags().IntVar(&p.adminPort, "admin-port", p.adminPort, "The port number used to run a HTTP server for admin tasks such as metrics, healthz.")

	cmd.Flags().StringVar(&p.toolsDir, "tools-dir", p.toolsDir, "The path to directory where to install needed tools such as kubectl, helm, kustomize.")
	cmd.Flags().BoolVar(&p.enableDefaultKubernetesCloudProvider, "enable-default-kubernetes-cloud-provider", p.enableDefaultKubernetesCloudProvider, "Whether the default kubernetes provider is enabled or not. This feature is deprecated.")
	cmd.Flags().BoolVar(&p.addLoginUserToPasswd, "add-login-user-to-passwd", p.addLoginUserToPasswd, "Whether to add login user to $HOME/passwd. This is typically for applications running as a random user ID.")
	cmd.Flags().DurationVar(&p.gracePeriod, "grace-period", p.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().StringVar(&p.launcherVersion, "launcher-version", p.launcherVersion, "The version of launcher which initialized this Piped.")

	return cmd
}

func (p *piped) run(ctx context.Context, input cli.Input) (runErr error) {
	// Make a cancellable context.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	if p.addLoginUserToPasswd {
		if err := p.insertLoginUserToPasswd(ctx); err != nil {
			return fmt.Errorf("failed to insert logged-in user to passwd: %w", err)
		}
	}

	// Load piped configuration from the specified source.
	cfg, err := p.loadConfig(ctx)
	if err != nil {
		input.Logger.Error("failed to load piped configuration", zap.Error(err))
		return err
	}

	// Register all metrics.
	registry := registerMetrics(cfg.PipedID, cfg.ProjectID, p.launcherVersion)

	// Configure SSH config if needed.
	if cfg.Git.ShouldConfigureSSHConfig() {
		if err := git.AddSSHConfig(cfg.Git); err != nil {
			input.Logger.Error("failed to configure ssh-config", zap.Error(err))
			return err
		}
		input.Logger.Info("successfully configured ssh-config")
	}

	// Initialize default tool registry.
	if err := toolregistry.InitDefaultRegistry(p.toolsDir, input.Logger); err != nil {
		input.Logger.Error("failed to initialize default tool registry", zap.Error(err))
		return err
	}

	// Add configured Helm chart repositories.
	if repos := cfg.HTTPHelmChartRepositories(); len(repos) > 0 {
		reg := toolregistry.DefaultRegistry()
		if err := chartrepo.Add(ctx, repos, reg, input.Logger); err != nil {
			input.Logger.Error("failed to add configured chart repositories", zap.Error(err))
			return err
		}
		if err := chartrepo.Update(ctx, reg, input.Logger); err != nil {
			input.Logger.Error("failed to update Helm chart repositories", zap.Error(err))
			return err
		}
	}

	// Login to chart registries.
	if regs := cfg.ChartRegistries; len(regs) > 0 {
		reg := toolregistry.DefaultRegistry()
		helm, _, err := reg.Helm(ctx, "")
		if err != nil {
			return fmt.Errorf("failed to find helm while login to chart registries (%w)", err)
		}

		for _, r := range regs {
			switch r.Type {
			case config.OCIHelmChartRegistry:
				if r.Username == "" || r.Password == "" {
					continue
				}

				if err := loginToOCIRegistry(ctx, helm, r.Address, r.Username, r.Password); err != nil {
					input.Logger.Error(fmt.Sprintf("failed to login to %s Helm chart registry", r.Address), zap.Error(err))
					return err
				}
				input.Logger.Info("successfully logged in to Helm chart registry", zap.String("address", r.Address))

			default:
				return fmt.Errorf("unsupported Helm chart registry type: %s", r.Type)
			}
		}
	}

	pipedKey, err := cfg.LoadPipedKey()
	if err != nil {
		input.Logger.Error("failed to load piped key", zap.Error(err))
		return err
	}

	// Make gRPC client and connect to the API.
	apiClient, err := p.createAPIClient(ctx, cfg.APIAddress, cfg.ProjectID, cfg.PipedID, pipedKey, input.Logger)
	if err != nil {
		input.Logger.Error("failed to create gRPC client to control plane", zap.Error(err))
		return err
	}

	// Send the newest piped meta to the control-plane.
	if err := p.sendPipedMeta(ctx, apiClient, cfg, input.Logger); err != nil {
		input.Logger.Error("failed to report piped meta to control-plane", zap.Error(err))
		return err
	}

	// Initialize notifier and add piped events.
	notifier, err := notifier.NewNotifier(cfg, input.Logger)
	if err != nil {
		input.Logger.Error("failed to initialize notifier", zap.Error(err))
		return err
	}
	group.Go(func() error {
		return notifier.Run(ctx)
	})

	// Start running admin server.
	{
		var (
			ver   = []byte(version.Get().Version)
			admin = admin.NewAdmin(p.adminPort, p.gracePeriod, input.Logger)
		)

		admin.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
			w.Write(ver)
		})
		admin.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		admin.Handle("/metrics", input.PrometheusMetricsHandlerFor(registry))
		admin.HandleFunc("/debug/pprof/", pprof.Index)
		admin.HandleFunc("/debug/pprof/profile", pprof.Profile)
		admin.HandleFunc("/debug/pprof/trace", pprof.Trace)

		group.Go(func() error {
			return admin.Run(ctx)
		})
	}

	// Start running stats reporter.
	{
		url := fmt.Sprintf("http://localhost:%d/metrics", p.adminPort)
		r := statsreporter.NewReporter(url, apiClient, input.Logger)
		group.Go(func() error {
			return r.Run(ctx)
		})
	}

	// Initialize git client.
	gitOptions := []git.Option{
		git.WithUserName(cfg.Git.Username),
		git.WithEmail(cfg.Git.Email),
		git.WithAutoDetach(true),
		git.WithLogger(input.Logger),
	}
	for _, repo := range cfg.GitHelmChartRepositories() {
		if f := repo.SSHKeyFile; f != "" {
			// Configure git client to use the specified SSH key while fetching private Helm charts.
			env := fmt.Sprintf("GIT_SSH_COMMAND=ssh -i %s -o StrictHostKeyChecking=no -F /dev/null", f)
			gitOptions = append(gitOptions, git.WithGitEnvForRepo(repo.GitRemote, env))
		}
	}
	gitClient, err := git.NewClient(gitOptions...)
	if err != nil {
		input.Logger.Error("failed to initialize git client", zap.Error(err))
		return err
	}
	defer func() {
		if err := gitClient.Clean(); err != nil {
			input.Logger.Error("had an error while cleaning gitClient", zap.Error(err))
			return
		}
		input.Logger.Info("successfully cleaned gitClient")
	}()

	// Start running application store.
	var applicationLister applicationstore.Lister
	{
		store := applicationstore.NewStore(apiClient, p.gracePeriod, input.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		applicationLister = store.Lister()
	}

	// Start running deployment store.
	var deploymentLister deploymentstore.Lister
	{
		store := deploymentstore.NewStore(apiClient, p.gracePeriod, input.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		deploymentLister = store.Lister()
	}

	// Start running command store.
	var commandLister commandstore.Lister
	{
		store := commandstore.NewStore(apiClient, p.gracePeriod, input.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		commandLister = store.Lister()
	}

	// Start running event store.
	var eventLister eventstore.Lister
	{
		store := eventstore.NewStore(apiClient, p.gracePeriod, input.Logger)
		group.Go(func() error {
			return store.Run(ctx)
		})
		eventLister = store.Lister()
	}

	analysisResultStore := analysisresultstore.NewStore(apiClient, input.Logger)

	// Create memory caches.
	appManifestsCache := memorycache.NewTTLCache(ctx, time.Hour, time.Minute)

	var liveStateGetter livestatestore.Getter
	// Start running application live state store.
	{
		s := livestatestore.NewStore(ctx, cfg, applicationLister, p.gracePeriod, input.Logger)
		group.Go(func() error {
			return s.Run(ctx)
		})
		liveStateGetter = s.Getter()
	}

	// Start running application live state reporter.
	{
		r := livestatereporter.NewReporter(applicationLister, liveStateGetter, apiClient, cfg, input.Logger)
		group.Go(func() error {
			return r.Run(ctx)
		})
	}

	decrypter, err := p.initializeSecretDecrypter(cfg)
	if err != nil {
		input.Logger.Error("failed to initialize secret decrypter", zap.Error(err))
		return err
	}

	// Start running application application drift detector.
	{
		d, err := driftdetector.NewDetector(
			applicationLister,
			gitClient,
			liveStateGetter,
			apiClient,
			appManifestsCache,
			cfg,
			decrypter,
			input.Logger,
		)
		if err != nil {
			input.Logger.Error("failed to initialize application drift detector", zap.Error(err))
			return err
		}

		group.Go(func() error {
			return d.Run(ctx)
		})
	}

	// Start running deployment controller.
	{
		c := controller.NewController(
			apiClient,
			gitClient,
			deploymentLister,
			commandLister,
			applicationLister,
			livestatestore.LiveResourceLister{Getter: liveStateGetter},
			analysisResultStore,
			notifier,
			decrypter,
			cfg,
			appManifestsCache,
			p.gracePeriod,
			input.Logger,
		)

		group.Go(func() error {
			return c.Run(ctx)
		})
	}

	// Start running deployment trigger.
	var lastTriggeredCommitGetter trigger.LastTriggeredCommitGetter
	{
		tr, err := trigger.NewTrigger(
			apiClient,
			gitClient,
			applicationLister,
			commandLister,
			notifier,
			cfg,
			p.gracePeriod,
			input.Logger,
		)
		if err != nil {
			input.Logger.Error("failed to initialize trigger", zap.Error(err))
			return err
		}
		lastTriggeredCommitGetter = tr.GetLastTriggeredCommitGetter()

		group.Go(func() error {
			return tr.Run(ctx)
		})
	}

	// Start running event watcher.
	{
		w := eventwatcher.NewWatcher(
			cfg,
			eventLister,
			gitClient,
			apiClient,
			input.Logger,
		)
		group.Go(func() error {
			return w.Run(ctx)
		})
	}

	// Start running planpreview handler.
	{
		// Initialize a dedicated git client for plan-preview feature.
		// Basically, this feature is an utility so it should not share any resource with the main components of piped.
		gc, err := git.NewClient(
			git.WithUserName(cfg.Git.Username),
			git.WithEmail(cfg.Git.Email),
			git.WithAutoDetach(true),
			git.WithLogger(input.Logger),
		)
		if err != nil {
			input.Logger.Error("failed to initialize git client for plan-preview", zap.Error(err))
			return err
		}
		defer func() {
			if err := gc.Clean(); err != nil {
				input.Logger.Error("had an error while cleaning gitClient for plan-preview", zap.Error(err))
				return
			}
			input.Logger.Info("successfully cleaned gitClient for plan-preview")
		}()

		h := planpreview.NewHandler(
			gc,
			apiClient,
			commandLister,
			applicationLister,
			lastTriggeredCommitGetter,
			decrypter,
			appManifestsCache,
			cfg,
			planpreview.WithLogger(input.Logger),
		)
		group.Go(func() error {
			return h.Run(ctx)
		})
	}

	// Start running app-config-reporter.
	{
		r := appconfigreporter.NewReporter(
			apiClient,
			gitClient,
			applicationLister,
			cfg,
			p.gracePeriod,
			input.Logger,
		)

		group.Go(func() error {
			return r.Run(ctx)
		})
	}

	// Check for stop command.
	{
		group.Go(func() error {
			input.Logger.Info("start running piped stop checker")
			ticker := time.NewTicker(commandCheckPeriod)
			for {
				select {
				case <-ticker.C:
					shouldStop, err := stopCommandHandler(ctx, commandLister, input.Logger)
					// Don't return an error to continue this goroutine execution.
					if err != nil {
						input.Logger.Error("failed to check/handle piped stop command", zap.Error(err))
					}
					if shouldStop {
						input.Logger.Info("stop piped due to restart piped requested")
						cancel()
					}
				case <-ctx.Done():
					input.Logger.Info("piped stop checker has been stopped")
					return nil
				}
			}
		})
	}

	// Wait until all piped components have finished.
	// A terminating signal or a finish of any components
	// could trigger the finish of piped.
	// This ensures that all components are good or no one.
	if err := group.Wait(); err != nil {
		input.Logger.Error("failed while running", zap.Error(err))
		return err
	}
	return nil
}

// createAPIClient makes a gRPC client to connect to the API.
func (p *piped) createAPIClient(ctx context.Context, address, projectID, pipedID string, pipedKey []byte, logger *zap.Logger) (pipedservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var (
		token   = rpcauth.MakePipedToken(projectID, pipedID, string(pipedKey))
		creds   = rpcclient.NewPerRPCCredentials(token, rpcauth.PipedTokenCredentials, !p.insecure)
		options = []rpcclient.DialOption{
			rpcclient.WithBlock(),
			rpcclient.WithPerRPCCredentials(creds),
			rpcclient.WithMaxRecvMsgSize(p.maxRecvMsgSize),
		}
	)

	if !p.insecure {
		if p.certFile != "" {
			options = append(options, rpcclient.WithTLS(p.certFile))
		} else {
			config := &tls.Config{}
			options = append(options, rpcclient.WithTransportCredentials(credentials.NewTLS(config)))
		}
	} else {
		options = append(options, rpcclient.WithInsecure())
	}

	client, err := pipedservice.NewClient(ctx, address, options...)
	if err != nil {
		logger.Error("failed to create api client", zap.Error(err))
		return nil, err
	}
	return client, nil
}

// loadConfig reads the Piped configuration data from the specified source.
func (p *piped) loadConfig(ctx context.Context) (*config.PipedSpec, error) {
	if p.configFile != "" && p.configGCPSecret != "" {
		return nil, fmt.Errorf("only config-file or config-gcp-secret could be set")
	}

	extract := func(cfg *config.Config) (*config.PipedSpec, error) {
		if cfg.Kind != config.KindPiped {
			return nil, fmt.Errorf("wrong configuration kind for piped: %v", cfg.Kind)
		}
		if p.enableDefaultKubernetesCloudProvider {
			cfg.PipedSpec.EnableDefaultKubernetesPlatformProvider()
		}
		return cfg.PipedSpec, nil
	}

	if p.configFile != "" {
		cfg, err := config.LoadFromYAML(p.configFile)
		if err != nil {
			return nil, err
		}
		return extract(cfg)
	}

	if p.configData != "" {
		data, err := base64.StdEncoding.DecodeString(p.configData)
		if err != nil {
			return nil, fmt.Errorf("the given config-data isn't base64 encoded: %w", err)
		}

		cfg, err := config.DecodeYAML(data)
		if err != nil {
			return nil, err
		}
		return extract(cfg)
	}

	if p.configGCPSecret != "" {
		data, err := p.getConfigDataFromSecretManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from SecretManager (%w)", err)
		}
		cfg, err := config.DecodeYAML(data)
		if err != nil {
			return nil, err
		}
		return extract(cfg)
	}

	return nil, fmt.Errorf("either config-file or config-gcp-secret must be set")
}

func (p *piped) initializeSecretDecrypter(cfg *config.PipedSpec) (crypto.Decrypter, error) {
	sm := cfg.SecretManagement
	if sm == nil {
		return nil, nil
	}

	switch sm.Type {
	case model.SecretManagementTypeNone:
		return nil, nil

	case model.SecretManagementTypeKeyPair:
		key, err := sm.KeyPair.LoadPrivateKey()
		if err != nil {
			return nil, err
		}
		decrypter, err := crypto.NewHybridDecrypter(key)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize decrypter (%w)", err)
		}
		return decrypter, nil

	case model.SecretManagementTypeGCPKMS:
		return nil, fmt.Errorf("type %q is not implemented yet", sm.Type.String())

	case model.SecretManagementTypeAWSKMS:
		return nil, fmt.Errorf("type %q is not implemented yet", sm.Type.String())

	default:
		return nil, fmt.Errorf("unsupported secret management type: %s", sm.Type.String())
	}
}

func (p *piped) sendPipedMeta(ctx context.Context, client pipedservice.Client, cfg *config.PipedSpec, logger *zap.Logger) error {
	repos := make([]*model.ApplicationGitRepository, 0, len(cfg.Repositories))
	for _, r := range cfg.Repositories {
		repos = append(repos, &model.ApplicationGitRepository{
			Id:     r.RepoID,
			Remote: r.Remote,
			Branch: r.Branch,
		})
	}

	cloneCfg, err := cfg.Clone()
	if err != nil {
		return err
	}

	cloneCfg.Mask()
	maskedCfg, err := yaml.Marshal(cloneCfg)
	if err != nil {
		return err
	}

	var (
		req = &pipedservice.ReportPipedMetaRequest{
			Version:           version.Get().Version,
			Config:            string(maskedCfg),
			Repositories:      repos,
			PlatformProviders: make([]*model.Piped_PlatformProvider, 0, len(cfg.PlatformProviders)),
		}
		retry = pipedservice.NewRetry(5)
	)

	// Configure the list of specified platform providers.
	for _, cp := range cfg.PlatformProviders {
		req.PlatformProviders = append(req.PlatformProviders, &model.Piped_PlatformProvider{
			Name: cp.Name,
			Type: cp.Type.String(),
		})
	}

	// Configure secret management.
	if sm := cfg.SecretManagement; sm != nil && sm.Type == model.SecretManagementTypeKeyPair {
		publicKey, err := sm.KeyPair.LoadPublicKey()
		if err != nil {
			return fmt.Errorf("failed to read public key for secret management (%w)", err)
		}
		req.SecretEncryption = &model.Piped_SecretEncryption{
			Type:      sm.Type.String(),
			PublicKey: string(publicKey),
		}
	}
	if req.SecretEncryption == nil {
		req.SecretEncryption = &model.Piped_SecretEncryption{
			Type: model.SecretManagementTypeNone.String(),
		}
	}

	for retry.WaitNext(ctx) {
		if res, err := client.ReportPipedMeta(ctx, req); err == nil {
			cfg.Name = res.Name
			if cfg.WebAddress == "" {
				cfg.WebAddress = res.WebBaseUrl
			}
			return nil
		}
		logger.Warn("failed to report piped meta to control-plane, wait to the next retry",
			zap.Int("calls", retry.Calls()),
			zap.Error(err),
		)
	}

	return err
}

// insertLoginUserToPasswd adds the logged-in user to /etc/passwd.
// It requires nss_wrapper (https://cwrap.org/nss_wrapper.html)
// to get the operation done.
//
// This is a workaround to deal with OpenShift less than 4.2
// See more: https://github.com/pipe-cd/pipecd/issues/1905
func (p *piped) insertLoginUserToPasswd(ctx context.Context) error {
	var stdout, stderr bytes.Buffer

	// Use the id command so that it gets proper ids even in pure Go.
	cmd := exec.CommandContext(ctx, "id", "-u")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get uid: %s", &stderr)
	}
	uid := strings.TrimSpace(stdout.String())

	stdout.Reset()
	stderr.Reset()

	cmd = exec.CommandContext(ctx, "id", "-g")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get gid: %s", &stderr)
	}
	gid := strings.TrimSpace(stdout.String())

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to detect the current user's home directory: %w", err)
	}

	// echo "default:x:${USER_ID}:${GROUP_ID}:Dynamically created user:${HOME}:/sbin/nologin" >> "$HOME/passwd"
	entry := fmt.Sprintf("\ndefault:x:%s:%s:Dynamically created user:%s:/sbin/nologin", uid, gid, home)
	nssPasswdPath := filepath.Join(home, "passwd")
	f, err := os.OpenFile(nssPasswdPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", nssPasswdPath, err)
	}
	defer f.Close()
	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("failed to append entry to %q: %w", nssPasswdPath, err)
	}

	return nil
}

func (p *piped) getConfigDataFromSecretManager(ctx context.Context) ([]byte, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: p.configGCPSecret,
	}

	resp, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Payload.Data, nil
}

func registerMetrics(pipedID, projectID, launcherVersion string) *prometheus.Registry {
	r := prometheus.NewRegistry()
	wrapped := prometheus.WrapRegistererWith(
		map[string]string{
			"pipecd_component": "piped",
			"piped":            pipedID,
			"piped_version":    version.Get().Version,
			"launcher_version": launcherVersion,
			"project":          projectID,
		},
		r,
	)
	wrapped.Register(collectors.NewGoCollector())
	wrapped.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	k8scloudprovidermetrics.Register(wrapped)
	k8slivestatestoremetrics.Register(wrapped)
	planpreviewmetrics.Register(wrapped)
	controllermetrics.Register(wrapped)

	return r
}

func loginToOCIRegistry(ctx context.Context, execPath, address, username, password string) error {
	args := []string{
		"registry",
		"login",
		"-u",
		username,
		"-p",
		password,
		address,
	}

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, execPath, args...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}
	return nil
}

func stopCommandHandler(ctx context.Context, cmdLister commandstore.Lister, logger *zap.Logger) (bool, error) {
	logger.Debug("fetch unhandled piped commands")

	commands := cmdLister.ListPipedCommands()
	if len(commands) == 0 {
		return false, nil
	}

	stopCmds := make([]model.ReportableCommand, 0, len(commands))
	for _, command := range commands {
		if command.IsRestartPipedCmd() {
			stopCmds = append(stopCmds, command)
		}
	}

	if len(stopCmds) == 0 {
		return false, nil
	}

	for _, command := range stopCmds {
		if err := command.Report(ctx, model.CommandStatus_COMMAND_SUCCEEDED, nil, []byte(command.Id)); err != nil {
			return false, fmt.Errorf("failed to report command %s: %w", command.Id, err)
		}
	}

	return true, nil
}
