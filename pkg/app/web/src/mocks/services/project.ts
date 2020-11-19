import { rest } from "msw";
import { serialize } from "../serializer";
import { createMask } from "../utils";
import {
  UpdateProjectStaticAdminResponse,
  GetProjectResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { Project } from "pipe/pkg/app/web/model/project_pb";

export const projectHandlers = [
  rest.post(createMask("/UpdateProjectStaticAdmin"), (req, res, ctx) => {
    const r = new UpdateProjectStaticAdminResponse();
    const data = serialize(r.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
  rest.post(createMask("/GetProject"), (req, res, ctx) => {
    const r = new GetProjectResponse();
    const p = new Project();
    r.setProject(p);
    const data = serialize(r.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
];
