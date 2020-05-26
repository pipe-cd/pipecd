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

	"github.com/kapetaniosci/pipe/pkg/app/api/service/pipedservice"
	"github.com/kapetaniosci/pipe/pkg/cli"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcclient"
)

type samplecli struct {
	address string
	name    string

	function       string
	requestPayload string
}

func NewCommand() *cobra.Command {
	s := &samplecli{
		address: "localhost:9080",
		name:    "samplecli",
	}
	cmd := &cobra.Command{
		Use:   "samplecli",
		Short: "Start running sample client to api service",
		RunE:  cli.WithContext(s.run),
	}
	cmd.Flags().StringVar(&s.address, "address", s.address, "The address to HelloWorld service.")
	cmd.Flags().StringVar(&s.name, "name", s.name, "The name to be sent.")
	cmd.Flags().StringVar(&s.function, "function", s.function, "The function name.")
	cmd.Flags().StringVar(&s.requestPayload, "request-payload", s.requestPayload, "The json file that binds to request proto message.")
	return cmd
}

func (s *samplecli) run(ctx context.Context, t cli.Telemetry) error {
	cli, err := s.createPipedServiceClient(ctx, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create client", zap.Error(err))
		return err
	}
	defer cli.Close()

	data, err := ioutil.ReadFile(s.requestPayload)
	if err != nil {
		return err
	}

	switch s.function {
	case "CreateDeployment":
		return s.createDeployment(ctx, cli, data, t.Logger)
	case "ListApplications":
		return s.listApplications(ctx, cli, data, t.Logger)
	case "ReportDeploymentPlanned":
		return s.reportDeploymentPlanned(ctx, cli, data, t.Logger)
	case "ReportDeploymentRunning":
		return s.reportDeploymentRunning(ctx, cli, data, t.Logger)
	case "ReportDeploymentCompleted":
		return s.reportDeploymentCompleted(ctx, cli, data, t.Logger)
	case "SaveDeploymentMetadata":
		return s.saveDeploymentMetadata(ctx, cli, data, t.Logger)
	case "SaveStageMetadata":
		return s.saveStageMetadata(ctx, cli, data, t.Logger)
	case "ReportStageStatusChanged":
		return s.reportStageStatusChanged(ctx, cli, data, t.Logger)
	default:
		return fmt.Errorf("invalid function name: %s", s.function)
	}

	return nil
}

func (s *samplecli) createPipedServiceClient(ctx context.Context, logger *zap.Logger) (pipedservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	options := []rpcclient.DialOption{
		rpcclient.WithBlock(),
		rpcclient.WithStatsHandler(),
		rpcclient.WithInsecure(),
	}
	client, err := pipedservice.NewClient(ctx, s.address, options...)
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
		logger.Error("failure run CreateDeployment", zap.Error(err))
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
		logger.Error("failure run ListApplications", zap.Error(err))
		return err
	}
	logger.Info("successfully run ListApplications", zap.Int("count", len(resp.Applications)))
	for _, app := range resp.Applications {
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
		logger.Error("failure run ReportDeploymentPlanned", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportDeploymentPlanned")
	return nil
}

func (s *samplecli) reportDeploymentRunning(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportDeploymentRunningRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportDeploymentRunning(ctx, &req); err != nil {
		logger.Error("failure run ReportDeploymentRunning", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportDeploymentRunning")
	return nil
}

func (s *samplecli) reportDeploymentCompleted(ctx context.Context, cli pipedservice.Client, payload []byte, logger *zap.Logger) error {
	req := pipedservice.ReportDeploymentCompletedRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.ReportDeploymentCompleted(ctx, &req); err != nil {
		logger.Error("failure run ReportDeploymentCompleted", zap.Error(err))
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
		logger.Error("failure run SaveStageMetadata", zap.Error(err))
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
		logger.Error("failure run SaveDeploymentMetadata", zap.Error(err))
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
		logger.Error("failure run ReportStageStatusChanged", zap.Error(err))
		return err
	}
	logger.Info("successfully run ReportStageStatusChanged")
	return nil
}
