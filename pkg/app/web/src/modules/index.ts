import { combineReducers } from "redux";
import { activeStageSlice } from "./active-stage";
import { apiKeysSlice } from "./api-keys";
import { applicationCountsSlice } from "./application-counts";
import { applicationsSlice } from "./applications";
import { applicationLiveStateSlice } from "./applications-live-state";
import { commandsSlice } from "./commands";
import { deleteApplicationSlice } from "./delete-application";
import { deletingEnvSlice } from "./deleting-env";
import { deploymentFrequencySlice } from "./deployment-frequency";
import { deploymentsSlice } from "./deployments";
import { environmentsSlice } from "./environments";
import { insightSlice } from "./insight";
import { meSlice } from "./me";
import { pipedsSlice } from "./pipeds";
import { projectSlice } from "./project";
import { sealedSecretSlice } from "./sealed-secret";
import { stageLogsSlice } from "./stage-logs";
import { toastsSlice } from "./toasts";
import { updateApplicationSlice } from "./update-application";
import { unregisteredApplicationsSlice } from "./unregistered-applications";

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
  sealedSecret: sealedSecretSlice.reducer,
  apiKeys: apiKeysSlice.reducer,
  insight: insightSlice.reducer,
  deploymentFrequency: deploymentFrequencySlice.reducer,
  applicationCounts: applicationCountsSlice.reducer,
  deletingEnv: deletingEnvSlice.reducer,
  unregisteredApplications: unregisteredApplicationsSlice.reducer,
});
