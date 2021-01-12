import { rest } from "msw";
import {
  AddApplicationResponse,
  DeleteApplicationResponse,
  DisableApplicationResponse,
  EnableApplicationResponse,
  ListApplicationsResponse,
  SyncApplicationResponse,
  UpdateApplicationResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { dummyApplication } from "../../__fixtures__/dummy-application";
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
  rest.post<Uint8Array>(createMask("/EnableApplication"), (req, res, ctx) => {
    const response = new EnableApplicationResponse();

    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
  rest.post<Uint8Array>(createMask("/DisableApplication"), (req, res, ctx) => {
    const response = new DisableApplicationResponse();

    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
  rest.post<Uint8Array>(createMask("/DeleteApplication"), (req, res, ctx) => {
    const response = new DeleteApplicationResponse();

    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
  rest.post<Uint8Array>(createMask("/AddApplication"), (req, res, ctx) => {
    const response = new AddApplicationResponse();
    response.setApplicationId(dummyApplication.id);

    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
  rest.post<Uint8Array>(createMask("/UpdateApplication"), (req, res, ctx) => {
    const response = new UpdateApplicationResponse();
    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
  rest.post<Uint8Array>(createMask("/ListApplications"), (req, res, ctx) => {
    const response = new ListApplicationsResponse();
    response.setApplicationsList([]);
    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
];
