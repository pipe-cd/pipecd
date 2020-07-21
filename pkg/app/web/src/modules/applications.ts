import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import {
  Application as ApplicationModel,
  ApplicationSyncStatus,
} from "pipe/pkg/app/web/model/application_pb";
import * as applicationsApi from "../api/applications";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import { fetchCommand, CommandStatus, CommandModel } from "./commands";
import { AppState } from ".";

export type Application = Required<ApplicationModel.AsObject>;
export type ApplicationSyncStatusKey = keyof typeof ApplicationSyncStatus;
export type ApplicationKindKey = keyof typeof ApplicationKind;

export const applicationsAdapter = createEntityAdapter<Application>({
  selectId: (app) => app.id,
});

export const { selectAll, selectById } = applicationsAdapter.getSelectors();

export const fetchApplications = createAsyncThunk<
  Application[],
  void,
  { state: AppState }
>("applications/fetchList", async (_, thunkAPI) => {
  const { applicationFilterOptions } = thunkAPI.getState();
  const { applicationsList } = await applicationsApi.getApplications({
    options: applicationFilterOptions,
  });
  return applicationsList as Application[];
});

export const fetchApplication = createAsyncThunk<
  Application | undefined,
  string
>("applications/fetchById", async (applicationId) => {
  const { application } = await applicationsApi.getApplication({
    applicationId,
  });
  return application as Application;
});

export const syncApplication = createAsyncThunk<
  void,
  { applicationId: string }
>("applications/sync", async ({ applicationId }, thunkAPI) => {
  const { commandId } = await applicationsApi.syncApplication({
    applicationId,
  });

  await thunkAPI.dispatch(fetchCommand(commandId));
});

export const addApplication = createAsyncThunk<
  void,
  {
    name: string;
    env: string;
    pipedId: string;
    repoId: string;
    repoPath: string;
    configFilename?: string;
    kind: ApplicationKind;
    cloudProvider: string;
  }
>("applications/add", async (props) => {
  await applicationsApi.addApplication({
    name: props.name,
    envId: props.env,
    pipedId: props.pipedId,
    gitPath: {
      repoId: props.repoId,
      path: props.repoPath,
      configFilename: props.configFilename || "",
      url: "",
    },
    cloudProvider: props.cloudProvider,
    kind: props.kind,
  });
});

export const disableApplication = createAsyncThunk<
  void,
  { applicationId: string }
>("applications/disable", async (props) => {
  await applicationsApi.disableApplication(props);
});

export const enableApplication = createAsyncThunk<
  void,
  { applicationId: string }
>("applications/enable", async (props) => {
  await applicationsApi.enableApplication(props);
});

export const applicationsSlice = createSlice({
  name: "applications",
  initialState: applicationsAdapter.getInitialState<{
    adding: boolean;
    syncing: Record<string, boolean>;
    disabling: Record<string, boolean>;
  }>({
    adding: false,
    syncing: {},
    disabling: {},
  }),
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchApplications.fulfilled, (state, action) => {
        applicationsAdapter.removeAll(state);
        applicationsAdapter.upsertMany(state, action.payload);
      })
      .addCase(fetchApplication.fulfilled, (state, action) => {
        if (action.payload) {
          applicationsAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(addApplication.pending, (state) => {
        state.adding = true;
      })
      .addCase(addApplication.fulfilled, (state) => {
        state.adding = false;
      })
      .addCase(addApplication.rejected, (state, action) => {
        // TODO: Show alert when failed to add an application
        console.error(action);
        state.adding = false;
      })
      .addCase(syncApplication.pending, (state, action) => {
        state.syncing[action.meta.arg.applicationId] = true;
      })
      .addCase(fetchCommand.fulfilled, (state, action) => {
        if (
          action.payload.type === CommandModel.Type.SYNC_APPLICATION &&
          action.payload.status !== CommandStatus.COMMAND_NOT_HANDLED_YET
        ) {
          // If command type is sync application and that process is finished, change syncing status to false
          state.syncing[action.payload.applicationId] = false;
        }
      })
      .addCase(disableApplication.pending, (state, action) => {
        state.disabling[action.meta.arg.applicationId] = true;
      })
      .addCase(disableApplication.fulfilled, (state, action) => {
        state.disabling[action.meta.arg.applicationId] = false;
        applicationsAdapter.removeOne(state, action.meta.arg.applicationId);
      })
      .addCase(disableApplication.rejected, (state, action) => {
        state.disabling[action.meta.arg.applicationId] = false;
      });
  },
});

export { ApplicationSyncStatus } from "pipe/pkg/app/web/model/application_pb";
export { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
