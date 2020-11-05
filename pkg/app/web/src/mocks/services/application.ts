import { rest } from "msw";
import { SyncApplicationResponse } from "pipe/pkg/app/web/api_client/service_pb";
import { serialize } from "../serializer";
import { createMask } from "../utils";

export const applicationHandlers = [
  rest.post<Uint8Array>(createMask("/SyncApplication"), (req, res, ctx) => {
    const response = new SyncApplicationResponse();

    response.setCommandId("sync-command");

    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
];
