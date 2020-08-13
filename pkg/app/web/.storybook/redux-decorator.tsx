import React from "react";
import { Provider } from "react-redux";
import configureMockStore from "redux-mock-store";
import { AppState } from "../src/modules";
import { reducers } from "../src/modules";
import {
  configureStore,
  DeepPartial,
  getDefaultMiddleware,
} from "@reduxjs/toolkit";

const mockStore = configureMockStore(getDefaultMiddleware());
const store = configureStore({ reducer: reducers });
const baseState = store.getState();

export const createStore = (initialState: DeepPartial<AppState>) => {
  return mockStore(Object.assign({}, baseState, initialState));
};

export const createDecoratorRedux = (initialState: DeepPartial<AppState>) => (
  storyFn: any
) => <Provider store={createStore(initialState)}>{storyFn()}</Provider>;
