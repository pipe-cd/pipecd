import { combineReducers } from "redux";
import { activeStageSlice } from "./active-stage";
import { applicationCountsSlice } from "./application-counts";
import { applicationsSlice } from "./applications";
import { applicationLiveStateSlice } from "./applications-live-state";
import { commandsSlice } from "./commands";
import { deleteApplicationSlice } from "./delete-application";
import { deploymentsSlice } from "./deployments";
import { pipedsSlice } from "./pipeds";
import { projectSlice } from "./project";
import { sealedSecretSlice } from "./sealed-secret";
import { stageLogsSlice } from "./stage-logs";
import { toastsSlice } from "./toasts";
import { updateApplicationSlice } from "./update-application";
import { unregisteredApplicationsSlice } from "./unregistered-applications";
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
  project: projectSlice.reducer,
  sealedSecret: sealedSecretSlice.reducer,
  applicationCounts: applicationCountsSlice.reducer,
  unregisteredApplications: unregisteredApplicationsSlice.reducer,
});
