import * as grpcWeb from "grpc-web";
import { WebServiceClient } from "pipecd/web/api_client/service_grpc_web_pb";
import { apiEndpoint } from "~/constants/api-endpoint";

export const apiClient = new WebServiceClient(apiEndpoint, null, {
  withCredentials: "true",
});

interface ApiCallback<Res> {
  (err: grpcWeb.RpcError, response: { toObject: () => Res }): void;
}

export async function apiRequest<Req, Res>(
  request: Req,
  api: {
    (request: Req, meta: grpcWeb.Metadata, callback: ApiCallback<Res>): void;
  }
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
