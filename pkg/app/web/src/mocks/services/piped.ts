import {
  GenerateApplicationSealedSecretResponse,
  ListPipedsResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import {
  createPipedFromObject,
  dummyPiped,
} from "../../__fixtures__/dummy-piped";
import { createHandler } from "../create-handler";

export const pipedHandlers = [
  createHandler<GenerateApplicationSealedSecretResponse>(
    "/GenerateApplicationSealedSecret",
    () => {
      const response = new GenerateApplicationSealedSecretResponse();
      response.setData("xxxxx");
      return response;
    }
  ),
  createHandler<ListPipedsResponse>("/ListPipeds", () => {
    const response = new ListPipedsResponse();
    response.setPipedsList([createPipedFromObject(dummyPiped)]);
    return response;
  }),
];
