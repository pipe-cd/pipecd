import { Piped } from "../modules/pipeds";
import { dummyEnv } from "./dummy-environment";
import { createApplicationGitRepository, dummyRepo } from "./dummy-repo";
import faker from "faker";
import { createdRandTime, subtractRandTimeFrom } from "./utils";

faker.seed(1);

const updatedAt = createdRandTime();
const startedAt = subtractRandTimeFrom(updatedAt);
const createdAt = subtractRandTimeFrom(startedAt);

export const dummyPiped: Piped.AsObject = {
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
  desc: faker.lorem.text(1),
  disabled: false,
  id: faker.random.uuid(),
  name: "dummy-piped",
  projectId: "project-1",
  repositoriesList: [dummyRepo],
  createdAt: createdAt.unix(),
  startedAt: startedAt.unix(),
  updatedAt: updatedAt.unix(),
  version: "v0.1",
  status: Piped.ConnectionStatus.ONLINE,
  keyHash: "12345",
  keysList: [],
  envIdsList: [dummyEnv.id],
  sealedSecretEncryption: {
    encryptServiceAccount: "",
    publicKey: "",
    type: "",
  },
};

function createCloudProviderFromObject(
  o: Piped.CloudProvider.AsObject
): Piped.CloudProvider {
  const cp = new Piped.CloudProvider();
  cp.setName(o.name);
  cp.setType(o.type);
  return cp;
}

export function createPipedFromObject(o: Piped.AsObject): Piped {
  const piped = new Piped();
  piped.setId(o.id);
  piped.setDesc(o.desc);
  piped.setName(o.name);
  piped.setVersion(o.version);
  piped.setProjectId(o.projectId);
  piped.setCreatedAt(o.createdAt);
  piped.setStartedAt(o.startedAt);
  piped.setUpdatedAt(o.updatedAt);
  piped.setDisabled(o.disabled);
  piped.setEnvIdsList(o.envIdsList);
  piped.setRepositoriesList(
    o.repositoriesList.map(createApplicationGitRepository)
  );
  piped.setCloudProvidersList(
    o.cloudProvidersList.map(createCloudProviderFromObject)
  );
  return piped;
}
