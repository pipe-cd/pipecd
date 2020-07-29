import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { DeploymentStatus } from "./deployments";
import { ApplicationKind } from "./applications";

export type DeploymentFilterOptions = {
  statuses: DeploymentStatus[];
  kinds: ApplicationKind[];
  applicationIds: string[];
  envIds: string[];
};

const initialState: DeploymentFilterOptions = {
  statuses: [],
  kinds: [],
  applicationIds: [],
  envIds: [],
};

export const deploymentFilterOptionsSlice = createSlice({
  name: "deploymentFilterOptions",
  initialState,
  reducers: {
    updateDeploymentFilter(
      state,
      action: PayloadAction<Partial<DeploymentFilterOptions>>
    ) {
      return { ...state, ...action.payload };
    },
    clearDeploymentFilter() {
      return initialState;
    },
  },
});

export const {
  updateDeploymentFilter,
  clearDeploymentFilter,
} = deploymentFilterOptionsSlice.actions;
