import { Environment } from "~/modules/environments";
import { createRandTimes, randomUUID, randomWords } from "./utils";

export const createEnvFromObject = (o: Environment.AsObject): Environment => {
  const env = new Environment();
  env.setCreatedAt(o.createdAt);
  env.setDesc(o.desc);
  env.setName(o.name);
  env.setProjectId(o.projectId);
  env.setUpdatedAt(o.updatedAt);
  env.setDeletedAt(o.deletedAt);
  env.setId(o.id);
  env.setDeleted(o.deleted);
  env.setDisabled(o.disabled);
  return env;
};

const [createdAt, updatedAt, deletedAt] = createRandTimes(3);

export const dummyEnv: Environment.AsObject = {
  id: randomUUID(),
  desc: randomWords(8),
  name: "staging",
  projectId: "project-1",
  disabled: false,
  deleted: false,
  deletedAt: deletedAt.unix(),
  updatedAt: updatedAt.unix(),
  createdAt: createdAt.unix(),
};

export const deletedDummyEnv: Environment.AsObject = {
  ...dummyEnv,
  id: randomUUID(),
  deleted: true,
  disabled: true,
};
