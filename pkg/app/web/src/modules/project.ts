import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import * as projectAPI from "../api/project";
import { AppState } from "./";
import {
  ProjectSSOConfig,
  ProjectRBACConfig,
} from "pipe/pkg/app/web/model/project_pb";

export type GitHubSSO = ProjectSSOConfig.GitHub.AsObject;
export type Teams = ProjectRBACConfig.AsObject;

export interface ProjectState {
  id: string | null;
  desc: string | null;
  username: string | null;
  teams: Teams | null;
  github: GitHubSSO | null;
  staticAdminDisabled: boolean;
  isUpdatingPassword: boolean;
  isUpdatingUsername: boolean;
  isUpdatingGitHubSSO: boolean;
}

const initialState: ProjectState = {
  id: null,
  desc: null,
  username: null,
  teams: null,
  github: null,
  staticAdminDisabled: false,
  isUpdatingPassword: false,
  isUpdatingUsername: false,
  isUpdatingGitHubSSO: false,
};

export const fetchProject = createAsyncThunk<{
  id: string | null;
  desc: string | null;
  username: string | null;
  teams: Teams | null;
  staticAdminDisabled: boolean;
  github: GitHubSSO | null;
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
    };
  }

  return {
    id: project.id,
    desc: project.desc,
    staticAdminDisabled: project.staticAdminDisabled,
    username: project.staticAdmin?.username || "",
    teams: project.rbac ?? null,
    github: project.sso?.github ?? null,
  };
});

export const updateUsername = createAsyncThunk<void, { username: string }>(
  "project/updateUsername",
  async ({ username }) => {
    await projectAPI.updateStaticAdminUsername({ username });
  }
);
export const updatePassword = createAsyncThunk<void, { password: string }>(
  "project/updatePassword",
  async ({ password }) => {
    await projectAPI.updateStaticAdminPassword({ password });
  }
);

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
  if (project.teams === null) {
    throw new Error();
  }
  await projectAPI.updateRBAC(Object.assign({}, project.teams, params));
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
      })
      .addCase(fetchProject.rejected, (_, action) => {
        console.log(action);
      })
      .addCase(updateUsername.pending, (state) => {
        state.isUpdatingUsername = true;
      })
      .addCase(updateUsername.fulfilled, (state) => {
        state.isUpdatingUsername = false;
      })
      .addCase(updateUsername.rejected, (state, action) => {
        console.error(action);
        state.isUpdatingUsername = false;
      })
      .addCase(updatePassword.pending, (state) => {
        state.isUpdatingPassword = true;
      })
      .addCase(updatePassword.fulfilled, (state) => {
        state.isUpdatingPassword = false;
      })
      .addCase(updatePassword.rejected, (state, action) => {
        console.error(action);
        state.isUpdatingPassword = false;
      })
      .addCase(updateGitHubSSO.pending, (state) => {
        state.isUpdatingGitHubSSO = true;
      })
      .addCase(updateGitHubSSO.fulfilled, (state) => {
        state.isUpdatingGitHubSSO = false;
      })
      .addCase(updateGitHubSSO.rejected, (state, action) => {
        console.error(action);
        state.isUpdatingGitHubSSO = false;
      })
      .addCase(updateRBAC.rejected, (state, action) => {
        console.error(action);
      });
  },
});

export { ProjectSSOConfig } from "pipe/pkg/app/web/model/project_pb";
