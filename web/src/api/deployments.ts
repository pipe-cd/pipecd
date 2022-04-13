import { apiClient, apiRequest } from "./client";
import {
  GetDeploymentRequest,
  GetDeploymentResponse,
  ListDeploymentsRequest,
  ListDeploymentsResponse,
  CancelDeploymentRequest,
  CancelDeploymentResponse,
  ApproveStageRequest,
  ApproveStageResponse,
  SkipStageRequest,
  SkipStageResponse,
} from "pipecd/web/api_client/service_pb";

export const getDeployment = ({
  deploymentId,
}: GetDeploymentRequest.AsObject): Promise<GetDeploymentResponse.AsObject> => {
  const req = new GetDeploymentRequest();
  req.setDeploymentId(deploymentId);
  return apiRequest(req, apiClient.getDeployment);
};

export const getDeployments = ({
  options,
  pageSize,
  cursor,
  pageMinUpdatedAt,
}: ListDeploymentsRequest.AsObject): Promise<
  ListDeploymentsResponse.AsObject
> => {
  const req = new ListDeploymentsRequest();
  if (options) {
    const opts = new ListDeploymentsRequest.Options();
    opts.setEnvIdsList(options.envIdsList);
    opts.setApplicationIdsList(options.applicationIdsList);
    opts.setKindsList(options.kindsList);
    opts.setStatusesList(options.statusesList);
    opts.setApplicationName(options.applicationName);
    req.setOptions(opts);
    req.setPageSize(pageSize);
    req.setCursor(cursor);
    req.setPageMinUpdatedAt(pageMinUpdatedAt);
    for (const label of options.labelsMap) {
      opts.getLabelsMap().set(label[0], label[1]);
    }
  }
  return apiRequest(req, apiClient.listDeployments);
};

export const cancelDeployment = ({
  deploymentId,
  forceRollback,
  forceNoRollback,
}: CancelDeploymentRequest.AsObject): Promise<
  CancelDeploymentResponse.AsObject
> => {
  const req = new CancelDeploymentRequest();
  req.setDeploymentId(deploymentId);
  req.setForceRollback(forceRollback);
  req.setForceNoRollback(forceNoRollback);
  return apiRequest(req, apiClient.cancelDeployment);
};

export const approveStage = ({
  deploymentId,
  stageId,
}: ApproveStageRequest.AsObject): Promise<ApproveStageResponse.AsObject> => {
  const req = new ApproveStageRequest();
  req.setDeploymentId(deploymentId);
  req.setStageId(stageId);
  return apiRequest(req, apiClient.approveStage);
};

export const skipStage = ({
  deploymentId,
  stageId,
}: SkipStageRequest.AsObject): Promise<SkipStageResponse.AsObject> => {
  const req = new SkipStageRequest();
  req.setDeploymentId(deploymentId);
  req.setStageId(stageId);
  return apiRequest(req, apiClient.skipStage);
};
