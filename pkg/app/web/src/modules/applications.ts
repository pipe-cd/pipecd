import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import { Application as ApplicationModel } from "pipe/pkg/app/web/model/application_pb";
import * as applicationsApi from "../api/applications";
import { ApplicationKind as ApplicationKindModel } from "pipe/pkg/app/web/model/common_pb";
import { fetchCommand, CommandStatus, CommandModel } from "./commands";
export { ApplicationSyncStatus } from "pipe/pkg/app/web/model/application_pb";

export type Application = Required<ApplicationModel.AsObject>;
export const ApplicationKind = ApplicationKindModel;
export type ApplicationKind = ApplicationKindModel;

export const applicationsAdapter = createEntityAdapter<Application>({
  selectId: (app) => app.id,
});

export const { selectAll, selectById } = applicationsAdapter.getSelectors();

export const fetchApplications = createAsyncThunk<Application[], void>(
  "applications/fetchList",
  async () => {
    const { applicationsList } = await applicationsApi.getApplications();
    return applicationsList as Application[];
  }
);

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
    configPath?: string;
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
      configPath: props.configPath || "",
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
