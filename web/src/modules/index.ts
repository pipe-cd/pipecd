import { combineReducers } from "redux";
import { activeStageSlice } from "./active-stage";
import { applicationCountsSlice } from "./application-counts";
import { applicationsSlice } from "./applications";
import { commandsSlice } from "./commands";
import { pipedsSlice } from "./pipeds";
import { toastsSlice } from "./toasts";
import { updateApplicationSlice } from "./update-application";
import { unregisteredApplicationsSlice } from "./unregistered-applications";

export const reducers = combineReducers({
  applications: applicationsSlice.reducer,
  updateApplication: updateApplicationSlice.reducer,
  activeStage: activeStageSlice.reducer,
  pipeds: pipedsSlice.reducer,
  commands: commandsSlice.reducer,
  toasts: toastsSlice.reducer,
  applicationCounts: applicationCountsSlice.reducer,
  unregisteredApplications: unregisteredApplicationsSlice.reducer,
});
