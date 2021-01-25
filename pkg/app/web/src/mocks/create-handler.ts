import { rest } from "msw";
import { serialize } from "./serializer";
import { createMask } from "./utils";

const DummyHandler = rest.post("", (_0, res, ctx) => {
  return res(ctx.status(200));
});

export function createHandler<T extends { serializeBinary(): Uint8Array }>(
  serviceName: string,
  getResponseMessage: (requestData: Uint8Array) => T
): typeof DummyHandler {
  return rest.post(createMask(serviceName), (req, res, ctx) => {
    const message = getResponseMessage(req.body?.slice(5));
    const data = serialize(message.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  });
}
