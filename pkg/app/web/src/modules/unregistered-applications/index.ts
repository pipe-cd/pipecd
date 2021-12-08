import {
  createAsyncThunk,
  createSlice,
  createEntityAdapter,
} from "@reduxjs/toolkit";
import { ApplicationInfo } from "pipe/pkg/app/web/model/common_pb";
import type { AppState } from "~/store";
import * as applicationsAPI from "~/api/applications";

const MODULE_NAME = "unregistered-applications";

export const unregisteredApplicationsAdapter = createEntityAdapter<
  ApplicationInfo.AsObject
>({});

const { selectAll } = unregisteredApplicationsAdapter.getSelectors();

export const selectAllUnregisteredApplications = (
  state: AppState
): ApplicationInfo.AsObject[] => selectAll(state.unregisteredApplications);

export const fetchUnregisteredApplications = createAsyncThunk<
  ApplicationInfo.AsObject[]
>(`${MODULE_NAME}/fetchList`, async () => {
  const {
    applicationsList,
  } = await applicationsAPI.getUnregisteredApplications();
  return applicationsList as ApplicationInfo.AsObject[];
});

export { ApplicationInfo } from "pipe/pkg/app/web/model/common_pb";

export const unregisteredApplicationsSlice = createSlice({
  name: MODULE_NAME,
  initialState: unregisteredApplicationsAdapter.getInitialState<{
    apps: ApplicationInfo | null;
  }>({
    apps: null,
  }),
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(
      fetchUnregisteredApplications.fulfilled,
      (state, action) => {
        unregisteredApplicationsAdapter.removeAll(state);
        unregisteredApplicationsAdapter.addMany(state, action.payload);
      }
    );
  },
});
