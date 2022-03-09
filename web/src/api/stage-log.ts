import { apiClient, apiRequest } from "./client";
import {
  GetStageLogRequest,
  GetStageLogResponse,
} from "pipecd/web/api_client/service_pb";

export const getStageLog = ({
  deploymentId,
  offsetIndex,
  retriedCount,
  stageId,
}: GetStageLogRequest.AsObject): Promise<GetStageLogResponse.AsObject> => {
  const req = new GetStageLogRequest();
  req.setDeploymentId(deploymentId);
  req.setStageId(stageId);
  req.setOffsetIndex(offsetIndex);
  req.setRetriedCount(retriedCount);
  return apiRequest(req, apiClient.getStageLog);
};
