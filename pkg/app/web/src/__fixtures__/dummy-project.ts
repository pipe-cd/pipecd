import { Project } from "pipe/pkg/app/web/model/project_pb";
import faker from "faker";
import { createdRandTime, subtractRandTimeFrom } from "./utils";

const updatedAt = createdRandTime();
const createdAt = subtractRandTimeFrom(updatedAt);

export const dummyProject: Project.AsObject = {
  id: faker.random.uuid(),
  desc: "",
  sharedSsoName: "",
  createdAt: createdAt.unix(),
  updatedAt: updatedAt.unix(),
  staticAdminDisabled: false,
};
