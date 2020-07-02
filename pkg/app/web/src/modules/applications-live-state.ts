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
>("applicationLiveState/fetchById", async (applicationId, thunkApi) => {
  try {
    const { snapshot } = await getApplicationLiveState({
      applicationId,
    });
    return snapshot as ApplicationLiveState;
  } catch (error) {
    return thunkApi.rejectWithValue(error);
  }
});

export const applicationLiveStateSlice = createSlice({
  name: "applicationLiveState",
  initialState: applicationLiveStateAdapter.getInitialState({}),
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchApplicationStateById.fulfilled, (state, action) => {
      if (action.payload) {
        applicationLiveStateAdapter.upsertOne(state, action.payload);
      }
    });
  },
});

export { ApplicationLiveStateSnapshot as ApplicationLiveStateSnapshotModel } from "pipe/pkg/app/web/model/application_live_state_pb";
