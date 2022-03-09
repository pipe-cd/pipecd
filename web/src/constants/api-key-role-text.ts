import { APIKey } from "~/modules/api-keys";

export const API_KEY_ROLE_TEXT: Record<APIKey.Role, string> = {
  [APIKey.Role.READ_ONLY]: "Read Only",
  [APIKey.Role.READ_WRITE]: "Read/Write",
};
