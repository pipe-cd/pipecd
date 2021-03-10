import { AnyAction, combineReducers } from "redux";
import { ThunkAction } from "redux-thunk";
import { activeStageSlice } from "./active-stage";
import { applicationsSlice } from "./applications";
import { applicationLiveStateSlice } from "./applications-live-state";
import { deploymentsSlice } from "./deployments";
import { environmentsSlice } from "./environments";
import { pipedsSlice } from "./pipeds";
import { stageLogsSlice } from "./stage-logs";
import { toastsSlice } from "./toasts";
import { commandsSlice } from "./commands";
import { meSlice } from "./me";
import { projectSlice } from "./project";
import { deploymentConfigsSlice } from "./deployment-configs";
import { sealedSecretSlice } from "./sealed-secret";
import { apiKeysSlice } from "./api-keys";
import { updateApplicationSlice } from "./update-application";
import { insightSlice } from "./insight";
import { deploymentFrequencySlice } from "./deployment-frequency";
import { deleteApplicationSlice } from "./delete-application";

export const reducers = combineReducers({
  deployments: deploymentsSlice.reducer,
  applicationLiveState: applicationLiveStateSlice.reducer,
  applications: applicationsSlice.reducer,
  updateApplication: updateApplicationSlice.reducer,
  deleteApplication: deleteApplicationSlice.reducer,
  stageLogs: stageLogsSlice.reducer,
  activeStage: activeStageSlice.reducer,
  pipeds: pipedsSlice.reducer,
  environments: environmentsSlice.reducer,
  commands: commandsSlice.reducer,
  toasts: toastsSlice.reducer,
  me: meSlice.reducer,
  project: projectSlice.reducer,
  deploymentConfigs: deploymentConfigsSlice.reducer,
  sealedSecret: sealedSecretSlice.reducer,
  apiKeys: apiKeysSlice.reducer,
  insight: insightSlice.reducer,
  deploymentFrequency: deploymentFrequencySlice.reducer,
});

export type AppState = ReturnType<typeof reducers>;
export type AppThunk = ThunkAction<
  Promise<unknown>,
  AppState,
  unknown,
  AnyAction
>;
