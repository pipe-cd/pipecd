import { Piped, PipedModel } from "../modules/pipeds";

export const dummyPiped: Piped = {
  cloudProvidersList: [],
  createdAt: 0,
  desc: "",
  disabled: false,
  id: "piped-1",
  name: "demo piped",
  projectId: "project-1",
  repositoriesList: [],
  startedAt: 0,
  updatedAt: 0,
  version: "v0.1",
  status: PipedModel.ConnectionStatus.ONLINE,
  keyHash: "12345",
};
