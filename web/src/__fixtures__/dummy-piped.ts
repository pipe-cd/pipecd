import { Piped, PipedKey } from "~/modules/pipeds";
import { createApplicationGitRepository, dummyRepo } from "./dummy-repo";
import { createRandTimes, randomText, randomUUID } from "./utils";

const [createdAt, startedAt, updatedAt] = createRandTimes(3);

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
  platformProvidersList: [],
  pluginsList: [],
  desc: randomText(1),
  disabled: false,
  id: randomUUID(),
  name: "dummy-piped",
  projectId: "project-1",
  repositoriesList: [dummyRepo],
  createdAt: createdAt.unix(),
  startedAt: startedAt.unix(),
  updatedAt: updatedAt.unix(),
  version: "v0.1",
  desiredVersion: "v1.0.0",
  status: Piped.ConnectionStatus.ONLINE,
  config: "apiVersion: pipecd.dev/v1beta1",
  keysList: [
    { hash: "key-1", creator: "user", createdAt: createdAt.unix() },
    { hash: "key-2", creator: "user", createdAt: createdAt.unix() },
  ],
};

function createCloudProviderFromObject(
  o: Piped.CloudProvider.AsObject
): Piped.CloudProvider {
  const cp = new Piped.CloudProvider();
  cp.setName(o.name);
  cp.setType(o.type);
  return cp;
}

function createPipedKeyFromObject(o: PipedKey.AsObject): PipedKey {
  const key = new PipedKey();
  key.setHash(o.hash);
  key.setCreator(o.creator);
  key.setCreatedAt(o.createdAt);
  return key;
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
  piped.setRepositoriesList(
    o.repositoriesList.map(createApplicationGitRepository)
  );
  piped.setCloudProvidersList(
    o.cloudProvidersList.map(createCloudProviderFromObject)
  );
  piped.setKeysList(o.keysList.map(createPipedKeyFromObject));
  return piped;
}
