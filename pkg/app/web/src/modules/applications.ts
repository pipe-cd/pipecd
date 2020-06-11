import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import { ApplicationLiveStateSnapshot } from "pipe/pkg/app/web/model/application_live_state_pb";
import { getApplicationLiveState } from "../api/applications";

export type ApplicationLiveState = Required<
  ApplicationLiveStateSnapshot.AsObject
>;

export const applicationsAdapter = createEntityAdapter<ApplicationLiveState>({
  selectId: (application) => application.applicationId,
});

export const { selectById } = applicationsAdapter.getSelectors();

export const fetchApplicationById = createAsyncThunk<
  ApplicationLiveState,
  string
>("applications/fetchById", async (applicationId) => {
  const { snapshot } = await getApplicationLiveState({ applicationId });
  return snapshot as ApplicationLiveState;
});

export const applicationsSlice = createSlice({
  name: "applications",
  initialState: applicationsAdapter.getInitialState({}),
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchApplicationById.pending, (state, action) => {})
      .addCase(fetchApplicationById.fulfilled, (state, action) => {
        if (action.payload) {
          applicationsAdapter.addOne(state, action.payload);
        }
      })
      .addCase(fetchApplicationById.rejected, (state, action) => {});
  },
});
