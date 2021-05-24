import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import * as deploymentConfigAPI from "../api/deployment-config";
import {
  DeploymentConfigTemplateLabel,
  DeploymentConfigTemplate,
} from "pipe/pkg/app/web/api_client/service_pb";
import { addApplication } from "./applications";
import type { AppState } from "../store";

export interface DeploymentConfigsState {
  templates: Record<string, DeploymentConfigTemplate.AsObject[]>;
  targetApplicationId: string | null;
}
export type DeploymentConfigTemplateLabelKey = keyof typeof DeploymentConfigTemplateLabel;
const initialState: DeploymentConfigsState = {
  templates: {},
  targetApplicationId: null,
};

export const fetchTemplateList = createAsyncThunk<
  DeploymentConfigTemplate.AsObject[],
  { labels: DeploymentConfigTemplateLabel[] },
  { state: AppState }
>("deploymentConfigs/fetchTemplates", async ({ labels }, thunkAPI) => {
  const { targetApplicationId } = thunkAPI.getState().deploymentConfigs;
  if (targetApplicationId === null) {
    throw new Error("target application is null.");
  }

  const {
    templatesList,
  } = await deploymentConfigAPI.getDeploymentConfigTemplates({
    applicationId: targetApplicationId,
    labelsList: labels,
  });
  return templatesList;
});

export const deploymentConfigsSlice = createSlice({
  name: "deploymentConfigs",
  initialState,
  reducers: {
    clearTemplateTarget(state) {
      state.targetApplicationId = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchTemplateList.fulfilled, (state, action) => {
        if (state.targetApplicationId) {
          state.templates[state.targetApplicationId] = action.payload;
        }
      })
      .addCase(addApplication.fulfilled, (state, action) => {
        state.targetApplicationId = action.payload;
      });
  },
});

export const selectTemplatesByAppId = (
  state: DeploymentConfigsState
): DeploymentConfigTemplate.AsObject[] | null => {
  if (!state.targetApplicationId) {
    return null;
  }

  const templates = state.templates[state.targetApplicationId];

  if (templates === undefined) {
    return null;
  }

  return templates;
};

export const { clearTemplateTarget } = deploymentConfigsSlice.actions;

export {
  DeploymentConfigTemplateLabel,
  DeploymentConfigTemplate,
} from "pipe/pkg/app/web/api_client/service_pb";
