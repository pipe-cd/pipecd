import { apiClient, apiRequest } from "./client";
import {
  GetCommandRequest,
  GetCommandResponse,
} from "pipecd/pkg/app/web/api_client/service_pb";

export const getCommand = ({
  commandId,
}: GetCommandRequest.AsObject): Promise<GetCommandResponse.AsObject> => {
  const req = new GetCommandRequest();
  req.setCommandId(commandId);
  return apiRequest(req, apiClient.getCommand);
};
