import { rest } from "msw";
import { GetMeResponse } from "pipe/pkg/app/web/api_client/service_pb";
import { serialize } from "../serializer";
import { Role } from "../../modules/me";
import { createMask } from "../utils";

export const meHandlers = [
  rest.post(createMask("/GetMe"), (req, res, ctx) => {
    const me = new GetMeResponse();
    me.setAvatarUrl("https://test.pipecd.dev/avatar.jpg");
    me.setSubject("hello-pipecd");
    me.setProjectId("pipecd");
    me.setProjectRole(Role.ProjectRole.ADMIN);
    const data = serialize(me.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
];
