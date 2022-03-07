import {
  createAsyncThunk,
  createSlice,
  createEntityAdapter,
} from "@reduxjs/toolkit";
import { ApplicationInfo } from "pipecd/pkg/app/web/model/common_pb";
import type { AppState } from "~/store";
import * as applicationsAPI from "~/api/applications";

const MODULE_NAME = "unregistered-applications";

export const unregisteredApplicationsAdapter = createEntityAdapter<
  ApplicationInfo.AsObject
>({});

export const selectAllUnregisteredApplications = (
  state: AppState
): ApplicationInfo.AsObject[] => {
  return state.unregisteredApplications.apps;
};

export const fetchUnregisteredApplications = createAsyncThunk<
  ApplicationInfo.AsObject[]
>(`${MODULE_NAME}/fetchList`, async () => {
  const {
    applicationsList,
  } = await applicationsAPI.getUnregisteredApplications();
  return applicationsList as ApplicationInfo.AsObject[];
});

export { ApplicationInfo } from "pipecd/pkg/app/web/model/common_pb";

export const unregisteredApplicationsSlice = createSlice({
  name: MODULE_NAME,
  initialState: unregisteredApplicationsAdapter.getInitialState({
    apps: [] as ApplicationInfo.AsObject[],
  }),
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(
      fetchUnregisteredApplications.fulfilled,
      (state, action) => {
        unregisteredApplicationsAdapter.removeAll(state);
        state.apps = action.payload;
      }
    );
  },
});
