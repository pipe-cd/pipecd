import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import { Deployment as DeploymentModel } from "pipe/pkg/app/web/model/deployment_pb";
import { getDeployment } from "../api/deployments";

export type Deployment = Required<DeploymentModel.AsObject>;

export const deploymentsAdapter = createEntityAdapter<Deployment>({});

export const { selectById } = deploymentsAdapter.getSelectors();

export const fetchDeploymentById = createAsyncThunk<Deployment, string>(
  "deployments/fetchById",
  async (deploymentId) => {
    const { deployment } = await getDeployment({ deploymentId });
    return deployment as Deployment;
  }
);

export const deploymentsSlice = createSlice({
  name: "deployments",
  initialState: deploymentsAdapter.getInitialState({
    loading: {} as Record<string, boolean>,
  }),
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchDeploymentById.pending, (state, action) => {
        state.loading[action.meta.arg] = true;
      })
      .addCase(fetchDeploymentById.fulfilled, (state, action) => {
        state.loading[action.meta.arg] = false;
        if (action.payload) {
          deploymentsAdapter.addOne(state, action.payload);
        }
      })
      .addCase(fetchDeploymentById.rejected, (state, action) => {
        state.loading[action.meta.arg] = false;
      });
  },
});
