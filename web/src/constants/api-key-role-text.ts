import { APIKey } from "pipecd/web/model/apikey_pb";

export const API_KEY_ROLE_TEXT: Record<APIKey.Role, string> = {
  [APIKey.Role.READ_ONLY]: "Read Only",
  [APIKey.Role.READ_WRITE]: "Read/Write",
};
