import { rest } from "msw";
import { GenerateApplicationSealedSecretResponse } from "pipe/pkg/app/web/api_client/service_pb";
import { serialize } from "../serializer";
import { createMask } from "../utils";

export const pipedHandlers = [
  rest.post(createMask("/GenerateApplicationSealedSecret"), (req, res, ctx) => {
    const response = new GenerateApplicationSealedSecretResponse();
    response.setData("xxxxx");
    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
];
