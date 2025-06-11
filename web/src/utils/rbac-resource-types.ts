import { RBAC_RESOURCE_TYPE_TEXT } from "~/constants/project";

export const rbacResourceTypes = (): string[] => {
  const resp: string[] = [];
  Object.values(RBAC_RESOURCE_TYPE_TEXT).map((v) => {
    resp.push(v);
  });
  return resp;
};
