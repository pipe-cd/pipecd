import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
  SerializedError,
  isFulfilled,
} from "@reduxjs/toolkit";
import {
  Application,
  ApplicationSyncStatus,
} from "pipecd/web/model/application_pb";
import * as applicationsAPI from "~/api/applications";
import {
  ApplicationGitRepository,
  ApplicationKind,
} from "pipecd/web/model/common_pb";
import { SyncStrategy } from "../deployments";
import { fetchCommand, CommandStatus, Command } from "../commands";
import type { AppState } from "~/store";

const MODULE_NAME = "applications";

export type ApplicationSyncStatusKey = keyof typeof ApplicationSyncStatus;
export type ApplicationKindKey = keyof typeof ApplicationKind;

export const applicationsAdapter = createEntityAdapter<Application.AsObject>({
  selectId: (app) => app.id,
});

export const { selectAll, selectById } = applicationsAdapter.getSelectors();

export interface ApplicationsFilterOptions {
  activeStatus?: string;
  kind?: string;
  syncStatus?: string;
  name?: string;
  pipedId?: string;
  // Suppose to be like ["key-1:value-1"]
  // sindresorhus/query-string doesn't support multidimensional arrays, that's why the format is a bit tricky.
  labels?: Array<string>;
}

export const fetchApplications = createAsyncThunk<
  Application.AsObject[],
  ApplicationsFilterOptions | undefined,
  { state: AppState }
>(`${MODULE_NAME}/fetchList`, async (options = {}) => {
  const labels = new Array<[string, string]>();
  if (options.labels) {
    for (const label of options.labels) {
      const pair = label.split(":");
      pair.length === 2 && labels.push([pair[0], pair[1]]);
    }
  }
  const req = {
    options: {
      envIdsList: [],
      kindsList: options.kind
        ? [parseInt(options.kind, 10) as ApplicationKind]
        : [],
      name: options.name ?? "",
      pipedId: options.pipedId ?? "",
      syncStatusesList: options.syncStatus
        ? [parseInt(options.syncStatus, 10) as ApplicationSyncStatus]
        : [],
      enabled: options.activeStatus
        ? { value: options.activeStatus === "enabled" }
        : undefined,
      labelsMap: labels,
    },
  };
  const { applicationsList } = await applicationsAPI.getApplications(req);
  return applicationsList as Application.AsObject[];
});

export const fetchApplication = createAsyncThunk<
  Application.AsObject | undefined,
  string
>(`${MODULE_NAME}/fetchById`, async (applicationId) => {
  const { application } = await applicationsAPI.getApplication({
    applicationId,
  });
  return application as Application.AsObject;
});

export const syncApplication = createAsyncThunk<
  void,
  { applicationId: string; syncStrategy: SyncStrategy }
>(`${MODULE_NAME}/sync`, async (values, thunkAPI) => {
  const { commandId } = await applicationsAPI.syncApplication(values);

  await thunkAPI.dispatch(fetchCommand(commandId));
});

export const addApplication = createAsyncThunk<
  string,
  {
    name: string;
    pipedId: string;
    repo: ApplicationGitRepository.AsObject;
    repoPath: string;
    configFilename?: string;
    kind?: ApplicationKind;
    platformProvider?: string;
    labels: Array<[string, string]>;
    deployTargets?: Array<{ pluginName: string; deployTarget: string }>;
  }
>(`${MODULE_NAME}/add`, async (props) => {
  const deployTargetsMap =
    props.deployTargets?.reduce((all, { pluginName, deployTarget }) => {
      if (!all[pluginName]) all[pluginName] = [];
      all[pluginName].push(deployTarget);
      return all;
    }, {} as Record<string, string[]>) || {};

  const deployTargetsByPluginMap = Object.entries(deployTargetsMap).map(
    ([pluginName, deployTargetsList]) => {
      return [pluginName, { deployTargetsList }] as [
        string,
        { deployTargetsList: string[] }
      ];
    }
  );

  const { applicationId } = await applicationsAPI.addApplication({
    name: props.name,
    pipedId: props.pipedId,
    gitPath: {
      repo: props.repo,
      path: props.repoPath,
      configFilename: props.configFilename || "",
      url: "",
    },
    platformProvider: props.platformProvider,
    kind: props.kind,
    deployTargetsByPluginMap,
    description: "",
    labelsMap: props.labels,
  });

  return applicationId;
});

export const disableApplication = createAsyncThunk<
  void,
  { applicationId: string }
>(`${MODULE_NAME}/disable`, async (props) => {
  await applicationsAPI.disableApplication(props);
});

export const enableApplication = createAsyncThunk<
  void,
  { applicationId: string }
>(`${MODULE_NAME}/enable`, async (props) => {
  await applicationsAPI.enableApplication(props);
});

const initialState = applicationsAdapter.getInitialState<{
  adding: boolean;
  loading: boolean;
  addedApplicationId: string | null;
  syncing: Record<string, boolean>;
  disabling: Record<string, boolean>;
  fetchApplicationError: SerializedError | null;
}>({
  adding: false,
  loading: false,
  addedApplicationId: null,
  syncing: {},
  disabling: {},
  fetchApplicationError: null,
});

export type ApplicationsState = typeof initialState;

export const applicationsSlice = createSlice({
  name: MODULE_NAME,
  initialState,
  reducers: {
    clearAddedApplicationId(state) {
      state.addedApplicationId = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchApplications.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchApplications.rejected, (state) => {
        state.loading = false;
      })
      .addCase(fetchApplication.pending, (state) => {
        state.fetchApplicationError = null;
      })
      .addCase(fetchApplication.fulfilled, (state, action) => {
        state.fetchApplicationError = null;
        if (action.payload) {
          applicationsAdapter.upsertOne(state, action.payload);
        }
      })
      .addCase(fetchApplication.rejected, (state, action) => {
        state.fetchApplicationError = action.error;
      })
      .addCase(addApplication.pending, (state) => {
        state.adding = true;
      })
      .addCase(addApplication.fulfilled, (state, action) => {
        state.adding = false;
        state.addedApplicationId = action.payload;
      })
      .addCase(addApplication.rejected, (state) => {
        // TODO: Show alert when failed to add an application
        state.adding = false;
      })
      .addCase(syncApplication.pending, (state, action) => {
        state.syncing[action.meta.arg.applicationId] = true;
      })
      .addCase(fetchCommand.fulfilled, (state, action) => {
        if (
          action.payload.type === Command.Type.SYNC_APPLICATION &&
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
      })
      .addMatcher(isFulfilled(fetchApplications), (state, action) => {
        applicationsAdapter.removeAll(state);
        applicationsAdapter.upsertMany(
          state,
          action.payload.filter((app) => app.deleted === false)
        );
        state.loading = false;
      });
  },
});

export const { clearAddedApplicationId } = applicationsSlice.actions;

export {
  Application,
  ApplicationSyncState,
  ApplicationSyncStatus,
  ApplicationDeploymentReference,
} from "pipecd/web/model/application_pb";
export {
  ApplicationKind,
  ApplicationActiveStatus,
  ApplicationGitRepository,
} from "pipecd/web/model/common_pb";
