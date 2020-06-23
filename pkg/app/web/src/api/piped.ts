import { apiClient, apiRequest } from "./client";
import {
  RegisterPipedRequest,
  RegisterPipedResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { Piped } from "../modules/pipeds";

// TODO: Replace mocks
// NOTE: This is mock function
export const getPipeds = ({}): Promise<Piped[]> => {
  return Promise.resolve([
    {
      id: "piped-1",
      desc: "mock piped",
      keyHash: "piped-key-hash",
      projectId: "debug-project",
      version: "v0.0.0",
      startedAt: 0,
      cloudProvidersList: [{ name: "kubernetes-default", type: "KUBERNETES" }],
      repositoryIdsList: ["repo-1", "repo-2"],
      disabled: false,
      createdAt: 0,
      updatedAt: 0,
    },
  ]);
};

export const registerPiped = ({
  desc,
}: RegisterPipedRequest.AsObject): Promise<RegisterPipedResponse.AsObject> => {
  const req = new RegisterPipedRequest();
  req.setDesc(desc);
  return apiRequest(req, apiClient.registerPiped);
};
