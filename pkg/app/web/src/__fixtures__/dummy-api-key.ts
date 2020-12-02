import { APIKeyModel } from "../modules/api-keys";

export const dummyAPIKey = {
  id: "api-key-1",
  name: "API_KEY_1",
  keyHash: "KEY_HASH",
  projectId: "pipecd",
  role: APIKeyModel.Role.READ_WRITE,
  creator: "user",
  disabled: false,
  createdAt: 0,
  updatedAt: 0,
};
