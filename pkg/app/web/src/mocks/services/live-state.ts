import { GetApplicationLiveStateResponse } from "pipe/pkg/app/web/api_client/service_pb";
import {
  createLiveStateSnapshotFromObject,
  dummyApplicationLiveState,
} from "../../__fixtures__/dummy-application-live-state";
import { createHandler } from "../create-handler";

export const liveStateHandlers = [
  createHandler<GetApplicationLiveStateResponse>(
    "/GetApplicationLiveState",
    () => {
      const response = new GetApplicationLiveStateResponse();
      response.setSnapshot(
        createLiveStateSnapshotFromObject(dummyApplicationLiveState)
      );
      return response;
    }
  ),
];
