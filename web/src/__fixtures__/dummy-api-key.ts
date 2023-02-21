import { APIKey } from "~/modules/api-keys";
import { createRandTime, randomKeyHash, randomUUID } from "./utils";

const createdAt = createRandTime();

export const dummyAPIKey: APIKey.AsObject = {
  id: randomUUID(),
  name: "API_KEY_1",
  keyHash: randomKeyHash(),
  projectId: "pipecd",
  role: APIKey.Role.READ_WRITE,
  creator: "user",
  disabled: false,
  createdAt: createdAt.unix(),
  updatedAt: createdAt.unix(),
  lastUsedAt: createdAt.unix(),
};

export function createAPIKeyFromObject(o: APIKey.AsObject): APIKey {
  const key = new APIKey();

  key.setId(o.id);
  key.setName(o.name);
  key.setKeyHash(o.keyHash);
  key.setProjectId(o.projectId);
  key.setRole(o.role);
  key.setCreator(o.creator);
  key.setDisabled(o.disabled);
  key.setCreatedAt(o.createdAt);
  key.setUpdatedAt(o.updatedAt);

  return key;
}
