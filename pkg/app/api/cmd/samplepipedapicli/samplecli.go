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

package samplepipedapicli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcauth"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcclient"
)

type samplecli struct {
	apiAddress     string
	name           string
	function       string
	requestPayload string
	tls            bool
	certFile       string
	projectID      string
	pipedID        string
	pipedKey       string
}

func NewCommand() *cobra.Command {
	s := &samplecli{
		apiAddress: "localhost:9080",
	}
	cmd := &cobra.Command{
		Use:   "samplepipedapicli",
		Short: "Start running sample client to piped api service",
		RunE:  cli.WithContext(s.run),
	}
	cmd.Flags().StringVar(&s.apiAddress, "api-address", s.apiAddress, "The address to piped api service.")
	cmd.Flags().StringVar(&s.function, "function", s.function, "The function name.")
	cmd.Flags().StringVar(&s.requestPayload, "request-payload", s.requestPayload, "The json file that binds to request proto message.")
	cmd.Flags().BoolVar(&s.tls, "tls", s.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&s.certFile, "cert-file", s.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&s.projectID, "project-id", s.projectID, "The project ID.")
	cmd.Flags().StringVar(&s.pipedID, "piped-id", s.pipedID, "The piped ID for using API.")
	cmd.Flags().StringVar(&s.pipedKey, "piped-key", s.pipedKey, "The piped key for using API.")

	cmd.MarkFlagRequired("function")
	cmd.MarkFlagRequired("request-payload")
	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("piped-id")
	cmd.MarkFlagRequired("piped-key")
	return cmd
}

func (s *samplecli) run(ctx context.Context, t cli.Telemetry) error {
	data, err := ioutil.ReadFile(s.requestPayload)
	if err != nil {
		return err
	}

	cli, err := s.createServiceClient(ctx, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create piped service client", zap.Error(err))
		return err
	}
	defer cli.Close()

	switch s.function {
	// PipedService
	case "CreateDeployment":
		return s.createDeployment(ctx, cli, data, t.Logger)
	case "ReportPipedMeta":
		return s.reportPipedMeta(ctx, cli, data, t.Logger)
	case "ListApplications":
		return s.listApplications(ctx, cli, data, t.Logger)
	case "ReportApplicationSyncState":
		return s.reportApplicationSyncState(ctx, cli, data, t.Logger)
	case "ReportApplicationMostRecentDeployment":
		return s.reportMostRecentlySuccessfulDeployment(ctx, cli, data, t.Logger)
	case "ListNotCompletedDeployments":
		return s.listNotCompletedDeployments(ctx, cli, data, t.Logger)
	case "GetApplicationMostRecentDeployment":
		return s.getMostRecentDeployment(ctx, cli, data, t.Logger)
	case "ReportDeploymentPlanned":
		return s.reportDeploymentPlanned(ctx, cli, data, t.Logger)
	case "ReportDeploymentStatusChanged":
		return s.reportDeploymentStatusChanged(ctx, cli, data, t.Logger)
	case "ReportDeploymentCompleted":
		return s.reportDeploymentCompleted(ctx, cli, data, t.Logger)
	case "SaveDeploymentMetadata":
		return s.saveDeploymentMetadata(ctx, cli, data, t.Logger)
	case "SaveStageMetadata":
		return s.saveStageMetadata(ctx, cli, data, t.Logger)
	case "ReportStageStatusChanged":
		return s.reportStageStatusChanged(ctx, cli, data, t.Logger)
	case "ReportStageLogs":
		return s.reportStageLogs(ctx, cli, data, t.Logger)
	case "ReportStageLogsFromLastCheckpoint":
		return s.reportStageLogsFromLastCheckpoint(ctx, cli, data, t.Logger)
	case "ListUnhandledCommands":
		return s.listUnhandledCommands(ctx, cli, data, t.Logger)
	case "ReportCommandHandled":
		return s.reportCommandHandled(ctx, cli, data, t.Logger)
	case "ReportApplicationLiveState":
		return s.reportApplicationLiveState(ctx, cli, data, t.Logger)
	case "ReportApplicationLiveStateEvents":
		return s.reportApplicationLiveStateEvents(ctx, cli, data, t.Logger)
	default:
		return fmt.Errorf("invalid function name: %s", s.function)
	}

	return nil
}

func (s *samplecli) createServiceClient(ctx context.Context, logger *zap.Logger) (pipedservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var (
		token   = rpcauth.MakePipedToken(s.projectID, s.pipedID, string(s.pipedKey))
		creds   = rpcclient.NewPerRPCCredentials(token, rpcauth.PipedTokenCredentials, s.tls)
		options = []rpcclient.DialOption{
			rpcclient.WithBlock(),
			rpcclient.WithPerRPCCredentials(creds),
		}
	)

	if s.tls {
		options = append(options, rpcclient.WithTLS(s.certFile))
	} else {
		options = append(options, rpcclient.WithInsecure())
	}
	client, err := pipedservice.NewClient(ctx, s.apiAddress, options...)
	if err != nil {
		logger.Error("failed to create PipedService client", zap.Error(err))
		return nil, err
	}
	return client, nil
}

func (s *samplecli) createDeployment(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.CreateDeploymentRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.CreateDeployment(ctx, &req); err != nil {
		logger.Error("failed to run CreateDeployment", zap.Error(err))
		return err
	}
	logger.Info("successfully run CreateDeployment")
	return nil
}

func (s *samplecli) reportPipedMeta(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportPipedMetaRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportPipedMeta(ctx, &req); err != nil {
		logger.Error("failed to run ReportPipedMeta", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportPipedMeta")
	return nil
}

func (s *samplecli) listApplications(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ListApplicationsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.ListApplications(ctx, &req)
	if err != nil {
		logger.Error("failed to run ListApplications", zap.Error(err))
		return err
	}
	logger.Info("successfully run ListApplications", zap.Int("count", len(resp.Applications)))
	for _, app := range resp.Applications {
		fmt.Printf("application: %+v\n", app)
	}
	return nil
}

func (s *samplecli) reportApplicationSyncState(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportApplicationSyncStateRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportApplicationSyncState(ctx, &req); err != nil {
		logger.Error("failed to run ReportApplicationSyncState", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportApplicationSyncState")
	return nil
}

func (s *samplecli) reportMostRecentlySuccessfulDeployment(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportApplicationMostRecentDeploymentRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportApplicationMostRecentDeployment(ctx, &req); err != nil {
		logger.Error("failed to run ReportApplicationMostRecentDeployment", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportApplicationMostRecentDeployment")
	return nil
}

func (s *samplecli) listNotCompletedDeployments(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ListNotCompletedDeploymentsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.ListNotCompletedDeployments(ctx, &req)
	if err != nil {
		logger.Error("failed to run ListNotCompletedDeployments", zap.Error(err))
		return err
	}
	logger.Info("successfully run ListNotCompletedDeployments", zap.Int("count", len(resp.Deployments)))
	for _, app := range resp.Deployments {
		fmt.Printf("application: %+v\n", app)
	}
	return nil
}

func (s *samplecli) getMostRecentDeployment(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.GetApplicationMostRecentDeploymentRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.GetApplicationMostRecentDeployment(ctx, &req)
	if err != nil {
		logger.Error("failed to run GetApplicationMostRecentDeployment", zap.Error(err))
		return err
	}
	logger.Info("successfully run GetApplicationMostRecentDeployment")
	fmt.Printf("deployment: %+v\n", resp.Deployment)
	return nil
}

func (s *samplecli) reportDeploymentPlanned(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportDeploymentPlannedRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportDeploymentPlanned(ctx, &req); err != nil {
		logger.Error("failed to run ReportDeploymentPlanned", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportDeploymentPlanned")
	return nil
}

func (s *samplecli) reportDeploymentStatusChanged(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportDeploymentStatusChangedRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportDeploymentStatusChanged(ctx, &req); err != nil {
		logger.Error("failed to run ReportDeploymentStatusChanged", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportDeploymentStatusChanged")
	return nil
}

func (s *samplecli) reportDeploymentCompleted(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportDeploymentCompletedRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportDeploymentCompleted(ctx, &req); err != nil {
		logger.Error("failed to run ReportDeploymentCompleted", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportDeploymentCompleted")
	return nil
}

func (s *samplecli) saveStageMetadata(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.SaveStageMetadataRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.SaveStageMetadata(ctx, &req); err != nil {
		logger.Error("failed to run SaveStageMetadata", zap.Error(err))
		return err
	}
	logger.Info("successfully run SaveStageMetadata")
	return nil
}

func (s *samplecli) saveDeploymentMetadata(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.SaveDeploymentMetadataRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.SaveDeploymentMetadata(ctx, &req); err != nil {
		logger.Error("failed to run SaveDeploymentMetadata", zap.Error(err))
		return err
	}
	logger.Info("successfully run SaveDeploymentMetadata")
	return nil
}

func (s *samplecli) reportStageStatusChanged(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportStageStatusChangedRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportStageStatusChanged(ctx, &req); err != nil {
		logger.Error("failed to run ReportStageStatusChanged", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportStageStatusChanged")
	return nil
}

func (s *samplecli) reportStageLogs(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportStageLogsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportStageLogs(ctx, &req); err != nil {
		logger.Error("failed to run ReportStageLogs", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportStageLogs")
	return nil
}

func (s *samplecli) reportStageLogsFromLastCheckpoint(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportStageLogsFromLastCheckpointRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportStageLogsFromLastCheckpoint(ctx, &req); err != nil {
		logger.Error("failed to run ReportStageLogsFromLastCheckpoint", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportStageLogsFromLastCheckpoint")
	return nil
}

func (s *samplecli) listUnhandledCommands(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ListUnhandledCommandsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.ListUnhandledCommands(ctx, &req)
	if err != nil {
		logger.Error("failed to run ListUnhandledCommands", zap.Error(err))
		return err
	}
	logger.Info("successfully run ListUnhandledCommands")
	for _, cmd := range resp.Commands {
		fmt.Printf("command: %+v\n", cmd)
	}
	return nil
}

func (s *samplecli) reportApplicationLiveState(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportApplicationLiveStateRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportApplicationLiveState(ctx, &req); err != nil {
		logger.Error("failed to run ReportApplicationLiveState", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportApplicationLiveState")
	return nil
}

func (s *samplecli) reportCommandHandled(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportCommandHandledRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportCommandHandled(ctx, &req); err != nil {
		logger.Error("failed to run ReportCommandHandled", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportCommandHandled")
	return nil
}

func (s *samplecli) reportApplicationLiveStateEvents(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportApplicationLiveStateEventsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportApplicationLiveStateEvents(ctx, &req); err != nil {
		logger.Error("failed to run ReportApplicationLiveStateEvents", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportApplicationLiveStateEvents")
	return nil
}
