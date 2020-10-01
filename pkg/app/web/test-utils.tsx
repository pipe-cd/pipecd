import { ThemeProvider } from "@material-ui/core";
import {
  AnyAction,
  configureStore,
  DeepPartial,
  getDefaultMiddleware,
  Store,
  ThunkDispatch,
} from "@reduxjs/toolkit";
import { render, RenderOptions, RenderResult } from "@testing-library/react";
import React from "react";
import { Provider } from "react-redux";
import configureMockStore from "redux-mock-store";
import { AppState, reducers } from "./src/modules";
import { theme } from "./src/theme";

const mockStore = configureMockStore<
  AppState,
  ThunkDispatch<AppState, void, AnyAction>
>(getDefaultMiddleware());
const store = configureStore({ reducer: reducers });
const baseState = store.getState();

export const createStore = (
  initialState: DeepPartial<AppState>
): ReturnType<typeof mockStore> => {
  return mockStore(Object.assign({}, baseState, initialState));
};

/**
 *
 * If you want to test the dispatched action effect, pass the real redux store instead of the initialState.
 */
const customRender = (
  ui: React.ReactElement,
  {
    initialState = {},
    store = createStore(initialState),
    ...renderOptions
  }: {
    initialState?: DeepPartial<AppState>;
    store?: Store<any, AnyAction>;
  } & Omit<RenderOptions, "queries">
): RenderResult => {
  const Wrapper: React.ComponentType = ({ children }) => (
    <Provider store={store}>
      <ThemeProvider theme={theme}>{children}</ThemeProvider>
    </Provider>
  );
  return render(ui, { wrapper: Wrapper, ...renderOptions });
};

// re-export everything
export * from "@testing-library/react";
// override render method
export { customRender as render };
