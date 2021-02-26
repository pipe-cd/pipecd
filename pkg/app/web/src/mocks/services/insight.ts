import { GetInsightDataResponse } from "pipe/pkg/app/web/api_client/service_pb";
import {
  InsightResultType,
  InsightSampleStream,
} from "pipe/pkg/app/web/model/insight_pb";
import {
  createDataPointsListFromObject,
  dummyDataPointsList,
} from "../../__fixtures__/dummy-insight";
import { createHandler } from "../create-handler";

export const insightHandlers = [
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
