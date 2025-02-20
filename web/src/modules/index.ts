import { combineReducers } from "redux";
import { activeStageSlice } from "./active-stage";
import { apiKeysSlice } from "./api-keys";
import { applicationCountsSlice } from "./application-counts";
import { applicationsSlice } from "./applications";
import { applicationLiveStateSlice } from "./applications-live-state";
import { commandsSlice } from "./commands";
import { deleteApplicationSlice } from "./delete-application";
import { deploymentFrequencySlice } from "./deployment-frequency";
import { deploymentChangeFailureRateSlice } from "./deployment-change-failure-rate";
import { deploymentsSlice } from "./deployments";
import { insightSlice } from "./insight";
import { meSlice } from "./me";
import { pipedsSlice } from "./pipeds";
import { projectSlice } from "./project";
import { sealedSecretSlice } from "./sealed-secret";
import { stageLogsSlice } from "./stage-logs";
import { toastsSlice } from "./toasts";
import { updateApplicationSlice } from "./update-application";
import { unregisteredApplicationsSlice } from "./unregistered-applications";
import { eventsSlice } from "./events";
import { deploymentTraceSlice } from "./deploymentTrace";

export const reducers = combineReducers({
  deployments: deploymentsSlice.reducer,
  deploymentTrace: deploymentTraceSlice.reducer,
  applicationLiveState: applicationLiveStateSlice.reducer,
  applications: applicationsSlice.reducer,
  updateApplication: updateApplicationSlice.reducer,
  deleteApplication: deleteApplicationSlice.reducer,
  stageLogs: stageLogsSlice.reducer,
  activeStage: activeStageSlice.reducer,
  pipeds: pipedsSlice.reducer,
  commands: commandsSlice.reducer,
  toasts: toastsSlice.reducer,
  me: meSlice.reducer,
  project: projectSlice.reducer,
  sealedSecret: sealedSecretSlice.reducer,
  apiKeys: apiKeysSlice.reducer,
  insight: insightSlice.reducer,
  deploymentFrequency: deploymentFrequencySlice.reducer,
  deploymentChangeFailureRate: deploymentChangeFailureRateSlice.reducer,
  applicationCounts: applicationCountsSlice.reducer,
  unregisteredApplications: unregisteredApplicationsSlice.reducer,
  events: eventsSlice.reducer,
});
