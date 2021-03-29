import { ThemeProvider } from "@material-ui/core";
import {
  AnyAction,
  configureStore,
  DeepPartial,
  EnhancedStore,
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

const middlewares = getDefaultMiddleware({
  immutableCheck: false,
  serializableCheck: false,
});

export const createReduxStore = (
  preloadedState?: Partial<AppState>
): EnhancedStore<AppState, AnyAction, typeof middlewares> => {
  return configureStore({
    reducer: reducers,
    middleware: getDefaultMiddleware({
      immutableCheck: false,
      serializableCheck: false,
    }),
    preloadedState,
  });
};

const store = createReduxStore();
const baseState = store.getState();

const mockStore = configureMockStore<
  AppState,
  ThunkDispatch<AppState, void, AnyAction>
>(middlewares);

export const createStore = (
  initialState: DeepPartial<AppState> | undefined = {}
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
    store?: Store<AppState, AnyAction>;
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
