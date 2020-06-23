import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import {
  Deployment as DeploymentModel,
  PipelineStage,
} from "pipe/pkg/app/web/model/deployment_pb";
import { getDeployment, getDeployments } from "../api/deployments";
export { DeploymentStatus } from "pipe/pkg/app/web/model/deployment_pb";

export type Deployment = Required<DeploymentModel.AsObject>;
export type Stage = Required<PipelineStage.AsObject>;

export const deploymentsAdapter = createEntityAdapter<Deployment>({});

export const { selectById, selectAll } = deploymentsAdapter.getSelectors();

export const fetchDeploymentById = createAsyncThunk<Deployment, string>(
  "deployments/fetchById",
  async (deploymentId) => {
    const { deployment } = await getDeployment({ deploymentId });
    return deployment as Deployment;
  }
);

export const fetchDeployments = createAsyncThunk<Deployment[]>(
  "deployments/fetchList",
  async () => {
    const { deploymentsList } = await getDeployments({});
    return (deploymentsList as Deployment[]) || [];
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
      })
      .addCase(fetchDeployments.pending, (state, action) => {})
      .addCase(fetchDeployments.fulfilled, (state, action) => {
        if (action.payload.length > 0) {
          deploymentsAdapter.addMany(state, action.payload);
        }
      })
      .addCase(fetchDeployments.rejected, (state, action) => {});
  },
});
