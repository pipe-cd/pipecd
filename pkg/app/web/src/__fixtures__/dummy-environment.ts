import { Environment } from "../modules/environments";
import { createRandTimes, randomUUID, randomWords } from "./utils";

export const createEnvFromObject = (o: Environment.AsObject): Environment => {
  const env = new Environment();
  env.setCreatedAt(o.createdAt);
  env.setDesc(o.desc);
  env.setName(o.name);
  env.setProjectId(o.projectId);
  env.setUpdatedAt(o.updatedAt);
  env.setId(o.id);
  return env;
};

const [createdAt, updatedAt] = createRandTimes(2);

export const dummyEnv: Environment.AsObject = {
  id: randomUUID(),
  desc: randomWords(8),
  name: "staging",
  projectId: "project-1",
  updatedAt: updatedAt.unix(),
  createdAt: createdAt.unix(),
};
