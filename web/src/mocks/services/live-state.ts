import {
  GetApplicationLiveStateRequest,
  GetApplicationLiveStateResponse,
} from "pipecd/web/api_client/service_pb";
import {
  createLiveStateSnapshotFromObject,
  dummyApplicationLiveState,
  dummyLiveStates,
} from "~/__fixtures__/dummy-application-live-state";
import { createHandler } from "../create-handler";

export const liveStateHandlers = [
  createHandler<GetApplicationLiveStateResponse>(
    "/GetApplicationLiveState",
    (requestBody) => {
      const response = new GetApplicationLiveStateResponse();
      const params = GetApplicationLiveStateRequest.deserializeBinary(
        requestBody
      );
      const appId = params.getApplicationId();
      const findState = Object.values(dummyLiveStates).find(
        (state) => state.applicationId === appId
      );
      response.setSnapshot(
        createLiveStateSnapshotFromObject(
          findState ?? dummyApplicationLiveState
        )
      );
      return response;
    }
  ),
];
