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

package samplecli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/api/service/pipedservice"
	"github.com/pipe-cd/pipe/pkg/app/api/service/webservice"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcclient"
)

type samplecli struct {
	pipedAPIAddress string
	webAPIAddress   string
	name            string
	function        string
	requestPayload  string
}

func NewCommand() *cobra.Command {
	s := &samplecli{
		pipedAPIAddress: "localhost:9080",
		webAPIAddress:   "localhost:9081",
		name:            "samplecli",
	}
	cmd := &cobra.Command{
		Use:   "samplecli",
		Short: "Start running sample client to api service",
		RunE:  cli.WithContext(s.run),
	}
	cmd.Flags().StringVar(&s.pipedAPIAddress, "piped-api-address", s.pipedAPIAddress, "The address to piped service.")
	cmd.Flags().StringVar(&s.webAPIAddress, "web-api-address", s.webAPIAddress, "The address to web service.")
	cmd.Flags().StringVar(&s.name, "name", s.name, "The name to be sent.")
	cmd.Flags().StringVar(&s.function, "function", s.function, "The function name.")
	cmd.Flags().StringVar(&s.requestPayload, "request-payload", s.requestPayload, "The json file that binds to request proto message.")
	return cmd
}

func (s *samplecli) run(ctx context.Context, t cli.Telemetry) error {
	webCli, err := s.createWebServiceClient(ctx, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create web service client", zap.Error(err))
		return err
	}
	defer webCli.Close()

	pipedCli, err := s.createPipedServiceClient(ctx, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create piped service client", zap.Error(err))
		return err
	}
	defer pipedCli.Close()

	data, err := ioutil.ReadFile(s.requestPayload)
	if err != nil {
		return err
	}

	switch s.function {
	// PipedService
	case "CreateDeployment":
		return s.createDeployment(ctx, pipedCli, data, t.Logger)
	case "ListApplications":
		return s.listApplications(ctx, pipedCli, data, t.Logger)
	case "ListNotCompletedDeployments":
		return s.listNotCompletedDeployments(ctx, pipedCli, data, t.Logger)
	case "ReportDeploymentPlanned":
		return s.reportDeploymentPlanned(ctx, pipedCli, data, t.Logger)
	case "ReportDeploymentStatusChanged":
		return s.reportDeploymentStatusChanged(ctx, pipedCli, data, t.Logger)
	case "ReportDeploymentCompleted":
		return s.reportDeploymentCompleted(ctx, pipedCli, data, t.Logger)
	case "SaveDeploymentMetadata":
		return s.saveDeploymentMetadata(ctx, pipedCli, data, t.Logger)
	case "SaveStageMetadata":
		return s.saveStageMetadata(ctx, pipedCli, data, t.Logger)
	case "ReportStageStatusChanged":
		return s.reportStageStatusChanged(ctx, pipedCli, data, t.Logger)
	case "ReportStageLogs":
		return s.reportStageLogs(ctx, pipedCli, data, t.Logger)
	case "ReportStageLogsFromLastCheckpoint":
		return s.reportStageLogsFromLastCheckpoint(ctx, pipedCli, data, t.Logger)
	// WebService
	case "ListDeployments":
		return s.listDeployments(ctx, webCli, data, t.Logger)
	case "GetDeployment":
		return s.getDeployment(ctx, webCli, data, t.Logger)
	case "GetStageLog":
		return s.getStageLog(ctx, webCli, data, t.Logger)
	case "GetApplicationLiveState":
		return s.getApplicationLiveState(ctx, webCli, data, t.Logger)
	default:
		return fmt.Errorf("invalid function name: %s", s.function)
	}

	return nil
}

func (s *samplecli) createWebServiceClient(ctx context.Context, logger *zap.Logger) (webservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	options := []rpcclient.DialOption{
		rpcclient.WithBlock(),
		rpcclient.WithStatsHandler(),
		rpcclient.WithInsecure(),
	}
	client, err := webservice.NewClient(ctx, s.webAPIAddress, options...)
	if err != nil {
		logger.Error("failed to create WebService client", zap.Error(err))
		return nil, err
	}
	return client, nil
}

func (s *samplecli) createPipedServiceClient(ctx context.Context, logger *zap.Logger) (pipedservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	options := []rpcclient.DialOption{
		rpcclient.WithBlock(),
		rpcclient.WithStatsHandler(),
		rpcclient.WithInsecure(),
	}
	client, err := pipedservice.NewClient(ctx, s.pipedAPIAddress, options...)
	if err != nil {
		logger.Error("failed to create PipedService client", zap.Error(err))
		return nil, err
	}
	return client, nil
}

func (s *samplecli) listDeployments(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.ListDeploymentsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.ListDeployments(ctx, &req)
	if err != nil {
		logger.Error("failed to run ListDeployments", zap.Error(err))
		return err
	}
	logger.Info("successfully run ListDeployments")
	fmt.Printf("deployments: %+v\n", resp)
	return nil
}

func (s *samplecli) getDeployment(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.GetDeploymentRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.GetDeployment(ctx, &req)
	if err != nil {
		logger.Error("failed to run GetDeployment", zap.Error(err))
		return err
	}
	logger.Info("successfully run GetDeployment")
	fmt.Printf("deployment: %+v\n", resp)
	return nil
}

func (s *samplecli) getStageLog(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.GetStageLogRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.GetStageLog(ctx, &req)
	if err != nil {
		logger.Error("failed to run GetStageLog", zap.Error(err))
		return err
	}
	logger.Info("successfully run GetStageLog")
	fmt.Printf("deployment: %+v\n", resp)
	return nil
}

func (s *samplecli) getApplicationLiveState(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.GetApplicationLiveStateRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.GetApplicationLiveState(ctx, &req)
	if err != nil {
		logger.Error("failed to run GetApplicationLiveState", zap.Error(err))
		return err
	}
	logger.Info("successfully run GetApplicationLiveState")
	fmt.Printf("state: %+v\n", resp)
	return nil
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
