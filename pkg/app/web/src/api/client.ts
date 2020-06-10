import * as grpcWeb from "grpc-web";
import { WebServiceClient } from "pipe/pkg/app/web/api_client/service_grpc_web_pb";

let apiEndpoint = `${location.protocol}//${location.host}`;

if (process.env.NODE_ENV === "development") {
  apiEndpoint = "/api";
}

export const apiClient = new WebServiceClient(apiEndpoint, null, {
  withCredentials: "true",
});

interface ApiCallback<Res> {
  (err: grpcWeb.Error, response: { toObject: () => Res }): void;
}

export async function apiRequest<Req, Res>(
  request: Req,
  api: { (request: Req, meta: {}, callback: ApiCallback<Res>): void }
): Promise<Res> {
  return new Promise((resolve, reject) => {
    api.bind(apiClient)(request, {}, (err, response) => {
      if (err) {
        reject(err);
      } else {
        resolve(response.toObject());
      }
    });
  });
}
