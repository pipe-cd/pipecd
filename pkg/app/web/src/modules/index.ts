import { AnyAction, combineReducers } from "redux";
import { ThunkAction, ThunkDispatch } from "redux-thunk";
import { deploymentsSlice } from "./deployments";
import { applicationLiveStateSlice } from "./applications-live-state";
import { stageLogsSlice } from "./stage-logs";
import { activeStageSlice } from "./active-stage";
import { applicationsSlice } from "./applications";
import { pipedsSlice } from "./pipeds";

export const reducers = combineReducers({
  deployments: deploymentsSlice.reducer,
  applicationLiveState: applicationLiveStateSlice.reducer,
  applications: applicationsSlice.reducer,
  stageLogs: stageLogsSlice.reducer,
  activeStage: activeStageSlice.reducer,
  pipeds: pipedsSlice.reducer,
});

export type AppState = ReturnType<typeof reducers>;
export type AppThunk = ThunkAction<
  Promise<unknown>,
  AppState,
  unknown,
  AnyAction
>;
export type AppDispatch = ThunkDispatch<AppState, null, AnyAction>;
