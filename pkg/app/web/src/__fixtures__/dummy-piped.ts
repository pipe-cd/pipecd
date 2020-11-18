import { Piped, PipedModel } from "../modules/pipeds";
import { dummyEnv } from "./dummy-environment";

export const dummyPiped: Piped = {
  cloudProvidersList: [
    {
      name: "kubernetes-default",
      type: "KUBERNETES",
    },

    {
      name: "terraform-default",
      type: "TERRAFORM",
    },
  ],
  createdAt: 0,
  desc: "",
  disabled: false,
  id: "piped-1",
  name: "dummy piped",
  projectId: "project-1",
  repositoriesList: [
    {
      id: "debug-repo",
      remote: "git@github.com:pipe-cd/debug.git",
      branch: "master",
    },
  ],
  startedAt: 0,
  updatedAt: 0,
  version: "v0.1",
  status: PipedModel.ConnectionStatus.ONLINE,
  keyHash: "12345",
  keysList: [],
  envIdsList: [dummyEnv.id],
  sealedSecretEncryption: {
    encryptServiceAccount: "",
    publicKey: "",
    type: "",
  },
};
