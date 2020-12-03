import { APIKeyModel } from "../modules/api-keys";

export const API_KEY_ROLE_TEXT: Record<APIKeyModel.Role, string> = {
  [APIKeyModel.Role.READ_ONLY]: "Read Only",
  [APIKeyModel.Role.READ_WRITE]: "Read/Write",
};
