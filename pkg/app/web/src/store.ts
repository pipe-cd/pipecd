import { configureStore, getDefaultMiddleware } from "@reduxjs/toolkit";
import { reducers } from "./modules";
import { thunkErrorHandler } from "./middlewares/thunk-error-handler";

export const store = configureStore({
  reducer: reducers,
  devTools: process.env.NODE_ENV === "development",
  middleware: [
    ...getDefaultMiddleware({}),
    thunkErrorHandler,
    // @see https://redux-toolkit.js.org/usage/usage-with-typescript#correct-typings-for-the-dispatch-type
  ] as const,
});

export type AppDispatch = typeof store.dispatch;
