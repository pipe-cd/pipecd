import {
  ProjectRBACResource,
  ProjectRBACPolicy,
} from "pipecd/web/model/project_pb";
import { rbacActionTypes } from "~/utils/rbac-action-types";
import { rbacResourceTypes } from "~/utils/rbac-resource-types";

export const RBAC_RESOURCE_TYPE_TEXT: Record<
  ProjectRBACResource.ResourceType,
  string
> = {
  [ProjectRBACResource.ResourceType.ALL]: "*",
  [ProjectRBACResource.ResourceType.APPLICATION]: "application",
  [ProjectRBACResource.ResourceType.DEPLOYMENT]: "deployment",
  [ProjectRBACResource.ResourceType.EVENT]: "event",
  [ProjectRBACResource.ResourceType.PIPED]: "piped",
  [ProjectRBACResource.ResourceType.DEPLOYMENT_CHAIN]: "deploymentChain",
  [ProjectRBACResource.ResourceType.PROJECT]: "project",
  [ProjectRBACResource.ResourceType.API_KEY]: "apiKey",
  [ProjectRBACResource.ResourceType.INSIGHT]: "insight",
};

export const RBAC_ACTION_TYPE_TEXT: Record<ProjectRBACPolicy.Action, string> = {
  [ProjectRBACPolicy.Action.ALL]: "*",
  [ProjectRBACPolicy.Action.GET]: "get",
  [ProjectRBACPolicy.Action.LIST]: "list",
  [ProjectRBACPolicy.Action.CREATE]: "create",
  [ProjectRBACPolicy.Action.UPDATE]: "update",
  [ProjectRBACPolicy.Action.DELETE]: "delete",
};

export const RESOURCE_ACTION_SEPARATOR = ";";
export const KEY_VALUE_SEPARATOR = "=";
export const VALUES_SEPARATOR = ",";
export const RESOURCES_KEY = "resources";
export const ACTIONS_KEY = "actions";

export const RESOURCES_NAME_REGEX = /([^,{}]+(?:{[^}]*})?)/g;
export const RESOURCES_LABELS_REGEX = /\{([^}]*)\}/;

export const TEXT_TO_RBAC_RESOURCE_TYPE: Record<
  string,
  ProjectRBACResource.ResourceType
> = {
  "*": ProjectRBACResource.ResourceType.ALL,
  application: ProjectRBACResource.ResourceType.APPLICATION,
  deployment: ProjectRBACResource.ResourceType.DEPLOYMENT,
  event: ProjectRBACResource.ResourceType.EVENT,
  piped: ProjectRBACResource.ResourceType.PIPED,
  deploymentChain: ProjectRBACResource.ResourceType.DEPLOYMENT_CHAIN,
  project: ProjectRBACResource.ResourceType.PROJECT,
  apiKey: ProjectRBACResource.ResourceType.API_KEY,
  insight: ProjectRBACResource.ResourceType.INSIGHT,
};

export const TEXT_TO_RBAC_ACTION_TYPE: Record<
  string,
  ProjectRBACPolicy.Action
> = {
  "*": ProjectRBACPolicy.Action.ALL,
  get: ProjectRBACPolicy.Action.GET,
  list: ProjectRBACPolicy.Action.LIST,
  create: ProjectRBACPolicy.Action.CREATE,
  update: ProjectRBACPolicy.Action.UPDATE,
  delete: ProjectRBACPolicy.Action.DELETE,
};

// example: resources=(\*|application|deployment|event|piped|deploymentChain|project|apiKey|insight|,)+;\s*actions=(\*|get|list|create|update|delete|,)+
export const POLICIES_STRING_REGEX = new RegExp(
  "resources=(" +
    rbacResourceTypes()
      .map((v) => v.replace(/\*/, "\\*"))
      .join("|") +
    "|,)+;\\s*actions=(" +
    rbacActionTypes()
      .map((v) => v.replace(/\*/, "\\*"))
      .join("|") +
    "|,)+"
);
