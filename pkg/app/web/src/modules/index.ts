import { AnyAction, combineReducers } from "redux";
import { ThunkAction, ThunkDispatch } from "redux-thunk";

export const reducers = combineReducers({});

export type AppState = ReturnType<typeof reducers>;
export type AppThunk = ThunkAction<
  Promise<unknown>,
  AppState,
  unknown,
  AnyAction
>;
export type AppDispatch = ThunkDispatch<AppState, null, AnyAction>;
