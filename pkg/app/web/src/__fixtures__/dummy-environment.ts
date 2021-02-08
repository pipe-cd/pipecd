import faker from "faker";
import { Environment } from "../modules/environments";
import { createdRandTime, subtractRandTimeFrom } from "./utils";

faker.seed(1);

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

const updatedAt = createdRandTime();
const createdAt = subtractRandTimeFrom(updatedAt);

export const dummyEnv: Environment.AsObject = {
  id: faker.random.uuid(),
  desc: faker.lorem.words(8),
  name: "staging",
  projectId: "project-1",
  updatedAt: updatedAt.unix(),
  createdAt: createdAt.unix(),
};
