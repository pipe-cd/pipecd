import { GetStageLogResponse } from "pipe/pkg/app/web/api_client/service_pb";
import {
  createLogBlockFromObject,
  dummyLogBlocks,
} from "../../__fixtures__/dummy-stage-log";
import { createHandler } from "../create-handler";

export const stageLogHandlers = [
  createHandler<GetStageLogResponse>("/GetStageLog", () => {
    const response = new GetStageLogResponse();
    response.setBlocksList(dummyLogBlocks.map(createLogBlockFromObject));
    response.setCompleted(true);

    return response;
  }),
];
