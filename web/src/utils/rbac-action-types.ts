import { RBAC_ACTION_TYPE_TEXT } from "~/constants/project";

export const rbacActionTypes = (): string[] => {
  const resp: string[] = [];
  Object.values(RBAC_ACTION_TYPE_TEXT).map((v) => {
    resp.push(v);
  });
  return resp;
};
