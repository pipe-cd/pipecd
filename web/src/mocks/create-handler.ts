import { rest } from "msw";
import { serialize } from "./serializer";
import { createMask } from "./utils";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const DummyHandler = rest.post<Uint8Array | string>("", (_0, res, ctx) => {
  return res(ctx.status(200));
});
type HandlerType = typeof DummyHandler;

const encoder = new TextEncoder();

export function createHandler<T extends { serializeBinary(): Uint8Array }>(
  serviceName: string,
  getResponseMessage: (requestData: Uint8Array) => T
): HandlerType {
  return rest.post<Uint8Array | string>(
    createMask(serviceName),
    (req, res, ctx) => {
      const arr: Uint8Array =
        typeof req.body === "string" ? encoder.encode(req.body) : req.body;
      const message = getResponseMessage(arr.slice(5));
      const data = serialize(message.serializeBinary());
      return res(
        ctx.status(200),
        ctx.set("Content-Type", "application/grpc-web+proto"),
        ctx.body(data)
      );
    }
  );
}

export function createHandlerWithError(
  serviceName: string,
  statusCode: number
): HandlerType {
  return rest.post(createMask(serviceName), (_, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.set({
        "content-length": "0",
        "content-type": "application/grpc-web+proto",
        "grpc-message": "Error Message",
        "grpc-status": `${statusCode}`,
      })
    );
  });
}
