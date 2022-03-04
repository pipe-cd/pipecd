import { StatusCode } from "grpc-web";
import { GetStageLogResponse } from "pipecd/pkg/app/web/api_client/service_pb";
import {
  createLogBlockFromObject,
  dummyLogBlocks,
} from "~/__fixtures__/dummy-stage-log";
import { createHandler, createHandlerWithError } from "../create-handler";

export const getStageLogHandler = createHandler<GetStageLogResponse>(
  "/GetStageLog",
  () => {
    const response = new GetStageLogResponse();
    response.setBlocksList(dummyLogBlocks.map(createLogBlockFromObject));
    response.setCompleted(true);

    return response;
  }
);

export const getStageLogNotFoundHandler = createHandlerWithError(
  "/GetStageLog",
  StatusCode.NOT_FOUND
);

export const getStageLogInternalErrorHandler = createHandlerWithError(
  "/GetStageLog",
  StatusCode.INTERNAL
);

export const stageLogHandlers = [getStageLogHandler];
