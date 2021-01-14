import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import {
  ApplicationLiveStateSnapshot as ApplicationLiveStateSnapshotModel,
  KubernetesResourceState as KubernetesResourceStateModel,
} from "pipe/pkg/app/web/model/application_live_state_pb";
import { getApplicationLiveState } from "../api/applications";

export type ApplicationLiveState = Required<
  ApplicationLiveStateSnapshotModel.AsObject
>;

export type KubernetesResourceState = Required<
  KubernetesResourceStateModel.AsObject
>;

export const HealthStatus = KubernetesResourceStateModel.HealthStatus;
export type HealthStatus = KubernetesResourceStateModel.HealthStatus;

export const applicationLiveStateAdapter = createEntityAdapter<
  ApplicationLiveState
>({
  selectId: (liveState) => liveState.applicationId,
});

export const { selectById } = applicationLiveStateAdapter.getSelectors();

export const fetchApplicationStateById = createAsyncThunk<
  ApplicationLiveState,
  string
>("applicationLiveState/fetchById", async (applicationId) => {
  const { snapshot } = await getApplicationLiveState({
    applicationId,
  });
  return snapshot as ApplicationLiveState;
});

const initialState = applicationLiveStateAdapter.getInitialState<{
  loading: Record<string, boolean>;
  hasError: Record<string, boolean>;
}>({
  loading: {},
  hasError: {},
});

export type ApplicationLiveStateState = typeof initialState;

export const selectHasError = (
  state: ApplicationLiveStateState,
  applicationId: string
): boolean => {
  return state.hasError[applicationId] || false;
};

export const selectLoadingById = (
  state: ApplicationLiveStateState,
  applicationId: string
): boolean => {
  return state.loading[applicationId] || false;
};

export const applicationLiveStateSlice = createSlice({
  name: "applicationLiveState",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchApplicationStateById.pending, (state, action) => {
        state.loading[action.meta.arg] = true;
        state.hasError[action.meta.arg] = false;
      })
      .addCase(fetchApplicationStateById.fulfilled, (state, action) => {
        state.loading[action.meta.arg] = false;
        state.hasError[action.meta.arg] = false;
        if (action.payload) {
          applicationLiveStateAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(fetchApplicationStateById.rejected, (state, action) => {
        state.loading[action.meta.arg] = false;
        state.hasError[action.meta.arg] = true;
      });
  },
});

export { ApplicationLiveStateSnapshot as ApplicationLiveStateSnapshotModel } from "pipe/pkg/app/web/model/application_live_state_pb";
