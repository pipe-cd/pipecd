import {
  ListAPIKeysResponse,
  GenerateAPIKeyResponse,
} from "pipecd/web/api_client/service_pb";
import {
  createAPIKeyFromObject,
  dummyAPIKey,
} from "~/__fixtures__/dummy-api-key";
import { createHandler } from "../create-handler";

export const apiKeyHandlers = [
  createHandler<ListAPIKeysResponse>("/ListAPIKeys", () => {
    const response = new ListAPIKeysResponse();
    response.setKeysList([createAPIKeyFromObject(dummyAPIKey)]);
    return response;
  }),
  createHandler<GenerateAPIKeyResponse>("/GenerateAPIKey", () => {
    const response = new GenerateAPIKeyResponse();
    response.setKey(dummyAPIKey.keyHash);
    return response;
  }),
];

export const getListAPIKeysEmpty = createHandler<ListAPIKeysResponse>(
  "/ListAPIKeys",
  () => {
    const response = new ListAPIKeysResponse();
    response.setKeysList([]);
    return response;
  }
);
