import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import {
  ProjectRBACConfig,
  ProjectSSOConfig,
  ProjectUserGroup,
  ProjectRBACRole,
  ProjectRBACPolicy,
  ProjectRBACResource,
} from "pipecd/web/model/project_pb";
import * as projectAPI from "~/api/project";
import type { AppState } from "~/store";

export type GitHubSSO = ProjectSSOConfig.GitHub.AsObject;
export type Teams = ProjectRBACConfig.AsObject;

const RBAC_RESOURCE_TYPE_TEXT: Record<
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

export const rbacResourceTypes = (): string[] => {
  const resp: string[] = [];
  Object.values(RBAC_RESOURCE_TYPE_TEXT).map((v) => {
    resp.push(v);
  });
  return resp;
};

const TEXT_TO_RBAC_RESOURCE_TYPE: Record<
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

const RBAC_ACTION_TYPE_TEXT: Record<ProjectRBACPolicy.Action, string> = {
  [ProjectRBACPolicy.Action.ALL]: "*",
  [ProjectRBACPolicy.Action.GET]: "get",
  [ProjectRBACPolicy.Action.LIST]: "list",
  [ProjectRBACPolicy.Action.CREATE]: "create",
  [ProjectRBACPolicy.Action.UPDATE]: "update",
  [ProjectRBACPolicy.Action.DELETE]: "delete",
};

export const rbacActionTypes = (): string[] => {
  const resp: string[] = [];
  Object.values(RBAC_ACTION_TYPE_TEXT).map((v) => {
    resp.push(v);
  });
  return resp;
};

const TEXT_TO_RBAC_ACTION_TYPE: Record<string, ProjectRBACPolicy.Action> = {
  "*": ProjectRBACPolicy.Action.ALL,
  get: ProjectRBACPolicy.Action.GET,
  list: ProjectRBACPolicy.Action.LIST,
  create: ProjectRBACPolicy.Action.CREATE,
  update: ProjectRBACPolicy.Action.UPDATE,
  delete: ProjectRBACPolicy.Action.DELETE,
};

const RESOURCE_ACTION_SEPARATOR = ";";
const KEY_VALUE_SEPARATOR = "=";
const VALUES_SEPARATOR = ",";
const RESOURCES_KEY = "resources";
const ACTIONS_KEY = "actions";

export const parseRBACPolicies = ({
  policies,
}: {
  policies: string;
}): ProjectRBACPolicy[] => {
  const ps = policies.split("\n\n").filter((p) => p);
  const ret: ProjectRBACPolicy[] = [];
  ps.map((p) => {
    p = p.replace(/\s/g, "");
    const policyResource: ProjectRBACPolicy = new ProjectRBACPolicy();

    // policy pattern
    // resources=RESOURCE_NAME{key1:value1,key2:value2};actions=ACTION
    const policy = p.split(RESOURCE_ACTION_SEPARATOR);

    if (
      policy.length !== 2 ||
      policy[0].startsWith(RESOURCES_KEY) === false ||
      policy[1].startsWith(ACTIONS_KEY) === false
    ) {
      return;
    }

    const resources = policy[0].split(KEY_VALUE_SEPARATOR);
    if (resources[0] == RESOURCES_KEY) {
      resources[1].split(VALUES_SEPARATOR).map((v) => {
        const res: ProjectRBACResource = new ProjectRBACResource();
        res.setType(TEXT_TO_RBAC_RESOURCE_TYPE[v]);
        res.clearLabelsMap(); // ensure no labels
        policyResource.addResources(res);
      });
    }

    const actions = policy[1].split(KEY_VALUE_SEPARATOR);
    if (actions[0] == ACTIONS_KEY) {
      actions[1].split(VALUES_SEPARATOR).map((v) => {
        policyResource.addActions(TEXT_TO_RBAC_ACTION_TYPE[v]);
      });
    }

    ret.push(policyResource);
  });
  return ret;
};

export const formalizePoliciesList = ({
  policiesList,
}: {
  policiesList: ProjectRBACPolicy.AsObject[];
}): string => {
  const policies: string[] = [];
  policiesList.map((policy) => {
    const resources: string[] = [];
    policy.resourcesList.map((resource) => {
      let rsc = RBAC_RESOURCE_TYPE_TEXT[resource.type];
      if (resource.labelsMap.length > 0) {
        rsc += "{";
        resource.labelsMap.map((label) => {
          rsc += label[0] + ":" + label[1] + ",";
        });
        rsc = rsc.slice(0, -1); // remove last comma
        rsc += "}";
      }
      resources.push(rsc);
    });

    const actions: string[] = [];
    policy.actionsList.map((action) => {
      actions.push(RBAC_ACTION_TYPE_TEXT[action]);
    });

    const resource =
      RESOURCES_KEY + KEY_VALUE_SEPARATOR + resources.join(VALUES_SEPARATOR);
    const action =
      ACTIONS_KEY + KEY_VALUE_SEPARATOR + actions.join(VALUES_SEPARATOR);
    policies.push(resource + RESOURCE_ACTION_SEPARATOR + action);
  });

  return policies.join("\n\n");
};

export interface ProjectState {
  id: string | null;
  desc: string | null;
  username: string | null;
  staticAdminDisabled: boolean;
  isUpdatingStaticAdmin: boolean;
  isUpdatingGitHubSSO: boolean;
  sharedSSO: string | null;
  teams?: Teams | null;
  github?: GitHubSSO | null;
  userGroups: ProjectUserGroup.AsObject[] | [];
  rbacRoles: ProjectRBACRole.AsObject[] | [];
}

const initialState: ProjectState = {
  id: null,
  desc: null,
  username: null,
  sharedSSO: null,
  staticAdminDisabled: false,
  isUpdatingStaticAdmin: false,
  isUpdatingGitHubSSO: false,
  userGroups: [],
  rbacRoles: [],
};

export const fetchProject = createAsyncThunk<{
  id: string | null;
  desc: string | null;
  username: string | null;
  teams: Teams | null;
  sharedSSO: string | null;
  staticAdminDisabled: boolean;
  github: GitHubSSO | null;
  userGroups: ProjectUserGroup.AsObject[] | [];
  rbacRoles: ProjectRBACRole.AsObject[] | [];
}>("project/fetchProject", async () => {
  const { project } = await projectAPI.getProject();

  if (!project) {
    return {
      id: null,
      desc: null,
      staticAdminDisabled: false,
      username: null,
      teams: null,
      github: null,
      sharedSSO: null,
      userGroups: [],
      rbacRoles: [],
    };
  }

  return {
    id: project.id,
    desc: project.desc,
    staticAdminDisabled: project.staticAdminDisabled,
    username: project.staticAdmin?.username || "",
    teams: project.rbac ?? null,
    github: project.sso?.github ?? null,
    sharedSSO: project.sharedSsoName,
    userGroups: project.userGroupsList,
    rbacRoles: project.rbacRolesList,
  };
});

export const updateStaticAdmin = createAsyncThunk<
  void,
  { username?: string; password?: string }
>("project/updateStaticAdmin", async (params) => {
  await projectAPI.updateStaticAdmin(params);
});

export const toggleAvailability = createAsyncThunk<
  void,
  void,
  { state: AppState }
>("project/toggleAvailability", async (_, thunkAPI) => {
  const s = thunkAPI.getState();

  if (s.project.staticAdminDisabled) {
    await projectAPI.enableStaticAdmin();
  } else {
    await projectAPI.disableStaticAdmin();
  }
});

export const updateGitHubSSO = createAsyncThunk<
  void,
  Partial<GitHubSSO> & { clientId: string; clientSecret: string }
>("project/updateGitHubSSO", async (params) => {
  await projectAPI.updateGitHubSSO(params);
});

export const updateRBAC = createAsyncThunk<
  void,
  Partial<Teams>,
  { state: AppState }
>("project/updateRBAC", async (params, thunkAPI) => {
  const project = thunkAPI.getState().project;
  await projectAPI.updateRBAC(Object.assign({}, project.teams, params));
});

export const addUserGroup = createAsyncThunk<
  void,
  { ssoGroup: string; role: string }
>("project/addUserGroup", async (params) => {
  await projectAPI.addUserGroup(params);
});

export const deleteUserGroup = createAsyncThunk<void, { ssoGroup: string }>(
  "project/deleteUserGroup",
  async (params) => {
    await projectAPI.deleteUserGroup(params);
  }
);

export const addRBACRole = createAsyncThunk<
  void,
  { name: string; policies: ProjectRBACPolicy[] }
>("project/addRBACRole", async (params) => {
  await projectAPI.addRBACRole(params);
});

export const deleteRBACRole = createAsyncThunk<void, { name: string }>(
  "project/deleteRBACRole",
  async (params) => {
    await projectAPI.deleteRBACRole(params);
  }
);

export const updateRBACRole = createAsyncThunk<
  void,
  { name: string; policies: ProjectRBACPolicy[] }
>("project/updateRBACRole", async (params) => {
  await projectAPI.updateRBACRole(params);
});

export const projectSlice = createSlice({
  name: "project",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      // .addCase(fetchProject.pending, () => {})
      .addCase(fetchProject.fulfilled, (state, action) => {
        state.id = action.payload.id;
        state.desc = action.payload.desc;
        state.username = action.payload.username;
        state.staticAdminDisabled = action.payload.staticAdminDisabled;
        state.teams = action.payload.teams;
        state.github = action.payload.github;
        state.sharedSSO = action.payload.sharedSSO;
        state.userGroups = action.payload.userGroups;
        state.rbacRoles = action.payload.rbacRoles;
      })
      // .addCase(fetchProject.rejected, (_, action) => {})
      .addCase(updateStaticAdmin.pending, (state) => {
        state.isUpdatingStaticAdmin = true;
      })
      .addCase(updateStaticAdmin.fulfilled, (state) => {
        state.isUpdatingStaticAdmin = false;
      })
      .addCase(updateStaticAdmin.rejected, (state) => {
        state.isUpdatingStaticAdmin = false;
      })
      .addCase(updateGitHubSSO.pending, (state) => {
        state.isUpdatingGitHubSSO = true;
      })
      .addCase(updateGitHubSSO.fulfilled, (state) => {
        state.isUpdatingGitHubSSO = false;
      })
      .addCase(updateGitHubSSO.rejected, (state) => {
        state.isUpdatingGitHubSSO = false;
      });
    // .addCase(updateRBAC.rejected, (state, action) => {});
  },
});
