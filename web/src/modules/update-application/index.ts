import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { ApplicationGitRepository } from "pipecd/web/model/common_pb";
import * as applicationAPI from "~/api/applications";
import { ApplicationKind } from "../applications";

const MODULE_NAME = "updateApplication";

export interface UpdateApplicationState {
  updating: boolean;
  targetId: string | null;
}

const initialState: UpdateApplicationState = {
  updating: false,
  targetId: null,
};

export const updateApplication = createAsyncThunk<
  void,
  {
    applicationId: string;
    name: string;
    pipedId: string;
    repo: ApplicationGitRepository.AsObject;
    repoPath: string;
    configFilename?: string;
    kind?: ApplicationKind;
    platformProvider?: string;
    deployTargets?: Array<{ pluginName: string; deployTarget: string }>;
  }
>(`${MODULE_NAME}/update`, async (values) => {
  const deployTargetsMap =
    values.deployTargets?.reduce((all, { pluginName, deployTarget }) => {
      if (!all[pluginName]) all[pluginName] = [];
      all[pluginName].push(deployTarget);
      return all;
    }, {} as Record<string, string[]>) || {};

  const deployTargetsByPluginMap = Object.entries(deployTargetsMap).map(
    ([pluginName, deployTargetsList]) => {
      return [pluginName, { deployTargetsList }] as [
        string,
        { deployTargetsList: string[] }
      ];
    }
  );

  await applicationAPI.updateApplication({
    applicationId: values.applicationId,
    name: values.name,
    pipedId: values.pipedId,
    platformProvider: values.platformProvider,
    kind: values.kind,
    deployTargetsByPluginMap,
    configFilename: values.configFilename || "",
  });
});

export const updateApplicationSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {
    setUpdateTargetId(state, action: PayloadAction<string>) {
      state.targetId = action.payload;
    },
    clearUpdateTarget(state) {
      state.targetId = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(updateApplication.pending, (state) => {
        state.updating = true;
      })
      .addCase(updateApplication.rejected, (state) => {
        state.updating = false;
      })
      .addCase(updateApplication.fulfilled, (state) => {
        state.updating = false;
        state.targetId = null;
      });
  },
});

export const {
  clearUpdateTarget,
  setUpdateTargetId,
} = updateApplicationSlice.actions;
