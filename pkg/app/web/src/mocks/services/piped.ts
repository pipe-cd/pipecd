import {
  GenerateApplicationSealedSecretResponse,
  ListPipedsResponse,
  RecreatePipedKeyResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import {
  createPipedFromObject,
  dummyPiped,
} from "../../__fixtures__/dummy-piped";
import { randomKeyHash } from "../../__fixtures__/utils";
import { createHandler } from "../create-handler";

export const generateApplicationSealedSecretHandler = createHandler<
  GenerateApplicationSealedSecretResponse
>("/GenerateApplicationSealedSecret", () => {
  const response = new GenerateApplicationSealedSecretResponse();
  response.setData(randomKeyHash());
  return response;
});

export const pipedHandlers = [
  createHandler<RecreatePipedKeyResponse>("/RecreatePipedKey", () => {
    const response = new RecreatePipedKeyResponse();
    response.setKey(randomKeyHash());
    return response;
  }),
  generateApplicationSealedSecretHandler,
  createHandler<ListPipedsResponse>("/ListPipeds", () => {
    const response = new ListPipedsResponse();
    response.setPipedsList([createPipedFromObject(dummyPiped)]);
    return response;
  }),
];
