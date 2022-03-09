import { StatusCode } from "grpc-web";
import {
  GetInsightDataResponse,
  GetInsightApplicationCountResponse,
} from "pipecd/web/api_client/service_pb";
import {
  InsightResultType,
  InsightSampleStream,
} from "pipecd/web/model/insight_pb";
import {
  dummyApplicationCounts,
  createInsightApplicationCountFromObject,
} from "~/__fixtures__/dummy-application-counts";
import {
  createDataPointsListFromObject,
  dummyDataPointsList,
} from "~/__fixtures__/dummy-insight";
import { createRandTime } from "~/__fixtures__/utils";
import { createHandler, createHandlerWithError } from "../create-handler";

export const getInsightApplicationCountHandler = createHandler<
  GetInsightApplicationCountResponse
>("/GetInsightApplicationCount", () => {
  const response = new GetInsightApplicationCountResponse();
  response.setUpdatedAt(createRandTime().unix());
  response.setCountsList(
    dummyApplicationCounts.map(createInsightApplicationCountFromObject)
  );
  return response;
});

export const getInsightApplicationCountNotFound = createHandlerWithError(
  "/GetInsightApplicationCount",
  StatusCode.NOT_FOUND
);

export const insightHandlers = [
  getInsightApplicationCountHandler,
  createHandler<GetInsightDataResponse>("/GetInsightData", () => {
    const response = new GetInsightDataResponse();
    response.setUpdatedAt(1);
    response.setMatrixList([]);
    const dataPointsList = createDataPointsListFromObject(dummyDataPointsList);
    const insightSampleStream = new InsightSampleStream();
    insightSampleStream.setDataPointsList(dataPointsList);
    response.setVectorList([]);
    response.setDataPointsList(dataPointsList);
    response.setType(InsightResultType.MATRIX);
    return response;
  }),
];
