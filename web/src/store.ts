import { configureStore } from "@reduxjs/toolkit";
import { thunkErrorHandler } from "./middlewares/thunk-error-handler";
import { reducers } from "./modules";

export const store = configureStore({
  reducer: reducers,
  devTools: process.env.NODE_ENV === "development",
  // @see https://redux-toolkit.js.org/usage/usage-with-typescript#correct-typings-for-the-dispatch-type
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().prepend(thunkErrorHandler),
});

export type AppState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
