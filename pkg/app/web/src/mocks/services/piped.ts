import {
  GenerateApplicationSealedSecretResponse,
  ListPipedsResponse,
  RecreatePipedKeyResponse,
  RegisterPipedResponse,
  DeleteOldPipedKeysResponse,
} from "pipecd/pkg/app/web/api_client/service_pb";
import { createPipedFromObject, dummyPiped } from "~/__fixtures__/dummy-piped";
import { randomKeyHash, randomUUID } from "~/__fixtures__/utils";
import { createHandler } from "../create-handler";

export const generateApplicationSealedSecretHandler = createHandler<
  GenerateApplicationSealedSecretResponse
>("/GenerateApplicationSealedSecret", () => {
  const response = new GenerateApplicationSealedSecretResponse();
  response.setData(randomKeyHash());
  return response;
});

export const pipedHandlers = [
  createHandler<RegisterPipedResponse>("/RegisterPiped", () => {
    const response = new RegisterPipedResponse();
    response.setId(randomUUID());
    response.setKey(randomKeyHash());
    return response;
  }),
  createHandler<RecreatePipedKeyResponse>("/RecreatePipedKey", () => {
    const response = new RecreatePipedKeyResponse();
    response.setKey(randomKeyHash());
    return response;
  }),
  createHandler<DeleteOldPipedKeysResponse>("/DeleteOldPipedKeys", () => {
    const response = new DeleteOldPipedKeysResponse();
    return response;
  }),
  generateApplicationSealedSecretHandler,
  createHandler<ListPipedsResponse>("/ListPipeds", () => {
    const response = new ListPipedsResponse();
    response.setPipedsList([createPipedFromObject(dummyPiped)]);
    return response;
  }),
];
