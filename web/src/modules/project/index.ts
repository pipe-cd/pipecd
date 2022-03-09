import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import {
  ProjectRBACConfig,
  ProjectSSOConfig,
} from "pipecd/web/model/project_pb";
import * as projectAPI from "~/api/project";
import type { AppState } from "~/store";

export type GitHubSSO = ProjectSSOConfig.GitHub.AsObject;
export type Teams = ProjectRBACConfig.AsObject;

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
}

const initialState: ProjectState = {
  id: null,
  desc: null,
  username: null,
  sharedSSO: null,
  staticAdminDisabled: false,
  isUpdatingStaticAdmin: false,
  isUpdatingGitHubSSO: false,
};

export const fetchProject = createAsyncThunk<{
  id: string | null;
  desc: string | null;
  username: string | null;
  teams: Teams | null;
  sharedSSO: string | null;
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
      sharedSSO: null,
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

export { ProjectSSOConfig } from "pipecd/web/model/project_pb";
