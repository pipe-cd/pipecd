import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import {
  ApplicationLiveStateSnapshot,
  KubernetesResourceState,
} from "pipecd/web/model/application_live_state_pb";
import { getApplicationLiveState } from "~/api/applications";

export type ApplicationLiveState = Required<
  ApplicationLiveStateSnapshot.AsObject
>;

export const HealthStatus = KubernetesResourceState.HealthStatus;
export type HealthStatus = KubernetesResourceState.HealthStatus;

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

export {
  ApplicationLiveStateSnapshot,
  KubernetesResourceState,
  CloudRunResourceState,
  ECSResourceState,
} from "pipecd/web/model/application_live_state_pb";
