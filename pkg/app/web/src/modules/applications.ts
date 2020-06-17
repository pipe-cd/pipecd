import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import { Application as ApplicationModel } from "pipe/pkg/app/web/model/application_pb";
import { getApplications } from "../api/applications";

export type Application = Required<ApplicationModel.AsObject>;

export const applicationsAdapter = createEntityAdapter<Application>({
  selectId: (app) => app.id,
});

export const { selectAll, selectById } = applicationsAdapter.getSelectors();

export const fetchApplications = createAsyncThunk<Application[], void>(
  "applications/fetchList",
  async () => {
    const { applicationsList } = await getApplications();
    return applicationsList as Application[];
  }
);

export const applicationsSlice = createSlice({
  name: "applications",
  initialState: applicationsAdapter.getInitialState(),
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchApplications.fulfilled, (state, action) => {
      applicationsAdapter.addMany(state, action.payload);
    });
  },
});
