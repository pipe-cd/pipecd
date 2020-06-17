import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import { ApplicationLiveStateSnapshot } from "pipe/pkg/app/web/model/application_live_state_pb";
import { getApplicationLiveState } from "../api/applications";

export type ApplicationLiveState = Required<
  ApplicationLiveStateSnapshot.AsObject
>;

export const applicationLiveStateAdapter = createEntityAdapter<
  ApplicationLiveState
>({
  selectId: (liveState) => liveState.applicationId,
});

export const { selectById } = applicationLiveStateAdapter.getSelectors();

export const fetchApplicationById = createAsyncThunk<
  ApplicationLiveState,
  string
>("applicationLiveState/fetchById", async (applicationId) => {
  const { snapshot } = await getApplicationLiveState({ applicationId });
  return snapshot as ApplicationLiveState;
});

export const applicationLiveStateSlice = createSlice({
  name: "applicationLiveState",
  initialState: applicationLiveStateAdapter.getInitialState({}),
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchApplicationById.pending, (state, action) => {})
      .addCase(fetchApplicationById.fulfilled, (state, action) => {
        if (action.payload) {
          applicationLiveStateAdapter.addOne(state, action.payload);
        }
      })
      .addCase(fetchApplicationById.rejected, (state, action) => {});
  },
});
