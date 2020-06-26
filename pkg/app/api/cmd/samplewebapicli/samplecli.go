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

package samplewebapicli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/pipe-cd/pipe/pkg/app/api/service/webservice"
	"github.com/pipe-cd/pipe/pkg/cli"
	"github.com/pipe-cd/pipe/pkg/rpc/rpcclient"
)

type samplecli struct {
	apiAddress     string
	function       string
	requestPayload string
}

func NewCommand() *cobra.Command {
	s := &samplecli{
		apiAddress: "localhost:9081",
	}
	cmd := &cobra.Command{
		Use:   "samplewebapicli",
		Short: "Start running sample client to web api service",
		RunE:  cli.WithContext(s.run),
	}
	cmd.Flags().StringVar(&s.apiAddress, "api-address", s.apiAddress, "The address to web api service.")
	cmd.Flags().StringVar(&s.function, "function", s.function, "The function name.")
	cmd.Flags().StringVar(&s.requestPayload, "request-payload", s.requestPayload, "The json file that binds to request proto message.")
	return cmd
}

func (s *samplecli) run(ctx context.Context, t cli.Telemetry) error {
	data, err := ioutil.ReadFile(s.requestPayload)
	if err != nil {
		return err
	}

	cli, err := s.createServiceClient(ctx, t.Logger)
	if err != nil {
		t.Logger.Error("failed to create web service client", zap.Error(err))
		return err
	}
	defer cli.Close()

	switch s.function {
	case "AddEnvironment":
		return s.addEnvironment(ctx, cli, data, t.Logger)
	case "ListEnvironments":
		return s.listEnvironments(ctx, cli, data, t.Logger)
	case "RegisterPiped":
		return s.registerPiped(ctx, cli, data, t.Logger)
	case "DisablePiped":
		return s.disablePiped(ctx, cli, data, t.Logger)
	case "ListPipeds":
		return s.listPipeds(ctx, cli, data, t.Logger)
	case "GetPiped":
		return s.getPiped(ctx, cli, data, t.Logger)
	case "AddApplication":
		return s.addApplication(ctx, cli, data, t.Logger)
	case "ListApplications":
		return s.listApplications(ctx, cli, data, t.Logger)
	case "GetApplication":
		return s.getApplication(ctx, cli, data, t.Logger)
	case "SyncApplication":
		return s.syncApplication(ctx, cli, data, t.Logger)
	case "ListDeployments":
		return s.listDeployments(ctx, cli, data, t.Logger)
	case "GetDeployment":
		return s.getDeployment(ctx, cli, data, t.Logger)
	case "GetStageLog":
		return s.getStageLog(ctx, cli, data, t.Logger)
	case "CancelDeployment":
		return s.cancelDeployment(ctx, cli, data, t.Logger)
	case "ApproveStage":
		return s.approveStage(ctx, cli, data, t.Logger)
	case "GetApplicationLiveState":
		return s.getApplicationLiveState(ctx, cli, data, t.Logger)
	}
	return nil
}

func (s *samplecli) createServiceClient(ctx context.Context, logger *zap.Logger) (webservice.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	options := []rpcclient.DialOption{
		rpcclient.WithBlock(),
		rpcclient.WithInsecure(),
	}
	client, err := webservice.NewClient(ctx, s.apiAddress, options...)
	if err != nil {
		logger.Error("failed to create WebService client", zap.Error(err))
		return nil, err
	}
	return client, nil
}

func (s *samplecli) addEnvironment(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.AddEnvironmentRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.AddEnvironment(ctx, &req); err != nil {
		logger.Error("failed to run AddEnvironment", zap.Error(err))
		return err
	}
	logger.Info("successfully run AddEnvironment")
	return nil
}

func (s *samplecli) listEnvironments(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.ListEnvironmentsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.ListEnvironments(ctx, &req)
	if err != nil {
		logger.Error("failed to run ListEnvironments", zap.Error(err))
		return err
	}
	logger.Info("successfully run ListEnvironments")
	for _, app := range resp.Environments {
		fmt.Printf("environment: %+v\n", app)
	}
	return nil
}

func (s *samplecli) registerPiped(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.RegisterPipedRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.RegisterPiped(ctx, &req)
	if err != nil {
		logger.Error("failed to run RegisterPiped", zap.Error(err))
		return err
	}
	logger.Info("successfully run RegisterPiped")
	fmt.Printf("Id: %+v\n", resp.Id)
	fmt.Printf("key: %+v\n", resp.Key)
	return nil
}

func (s *samplecli) disablePiped(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.DisablePipedRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.DisablePiped(ctx, &req); err != nil {
		logger.Error("failed to run DisablePiped", zap.Error(err))
		return err
	}
	logger.Info("successfully run DisablePiped")
	return nil
}

func (s *samplecli) listPipeds(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.ListPipedsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.ListPipeds(ctx, &req)
	if err != nil {
		logger.Error("failed to run ListPipeds", zap.Error(err))
		return err
	}
	logger.Info("successfully run ListPipeds")
	fmt.Printf("Pipeds: %+v\n", resp.Pipeds)
	return nil
}

func (s *samplecli) getPiped(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.GetPipedRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.GetPiped(ctx, &req)
	if err != nil {
		logger.Error("failed to run GetPiped", zap.Error(err))
		return err
	}
	logger.Info("successfully run GetPiped")
	fmt.Printf("Piped: %+v\n", resp.Piped)
	return nil
}

func (s *samplecli) addApplication(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.AddApplicationRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if _, err := cli.AddApplication(ctx, &req); err != nil {
		logger.Error("failed to run AddApplication", zap.Error(err))
		return err
	}
	logger.Info("successfully run AddApplication")
	return nil
}

func (s *samplecli) listApplications(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.ListApplicationsRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.ListApplications(ctx, &req)
	if err != nil {
		logger.Error("failed to run ListApplications", zap.Error(err))
		return err
	}
	logger.Info("successfully run ListApplications")
	for _, app := range resp.Applications {
		fmt.Printf("application: %+v\n", app)
	}
	return nil
}

func (s *samplecli) getApplication(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.GetApplicationRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.GetApplication(ctx, &req)
	if err != nil {
		logger.Error("failed to run GetApplication", zap.Error(err))
		return err
	}
	logger.Info("successfully run GetApplication")
	fmt.Printf("application: %+v\n", resp.Application)
	return nil
}

func (s *samplecli) syncApplication(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.SyncApplicationRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.SyncApplication(ctx, &req)
	if err != nil {
		logger.Error("failed to run SyncApplication", zap.Error(err))
		return err
	}
	logger.Info("successfully run SyncApplication")
	fmt.Printf("commandID: %+v\n", resp.CommandId)
	return nil
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
	for _, app := range resp.Deployments {
		fmt.Printf("deployment: %+v\n", app)
	}
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

func (s *samplecli) cancelDeployment(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.CancelDeploymentRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.CancelDeployment(ctx, &req)
	if err != nil {
		logger.Error("failed to run CancelDeployment", zap.Error(err))
		return err
	}
	logger.Info("successfully run CancelDeployment")
	fmt.Printf("commandID: %+v\n", resp.CommandId)
	return nil
}

func (s *samplecli) approveStage(ctx context.Context, cli webservice.Client, payload []byte, logger *zap.Logger) error {
	req := webservice.ApproveStageRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	resp, err := cli.ApproveStage(ctx, &req)
	if err != nil {
		logger.Error("failed to run ApproveStage", zap.Error(err))
		return err
	}
	logger.Info("successfully run ApproveStage")
	fmt.Printf("commandID: %+v\n", resp.CommandId)
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
