import { GenerateApplicationSealedSecretResponse } from "pipe/pkg/app/web/api_client/service_pb";
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
];
