import { Environment } from "../modules/environments";

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

export const dummyEnv: Environment.AsObject = {
  createdAt: 0,
  desc: "",
  name: "staging",
  projectId: "project-1",
  updatedAt: 0,
  id: "env-1",
};
