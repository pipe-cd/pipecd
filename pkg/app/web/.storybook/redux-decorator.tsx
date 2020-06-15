import React from "react";
import { Provider } from "react-redux";
import configureMockStore from "redux-mock-store";
import { AppState } from "../src/modules";
import { reducers } from "../src/modules";
import { configureStore } from "@reduxjs/toolkit";
import { DeepPartial } from "@reduxjs/toolkit";

const mockStore = configureMockStore([]);
const store = configureStore({ reducer: reducers });
const baseState = store.getState();

export const createDecoratorRedux = (initialState: DeepPartial<AppState>) => (
  storyFn: any
) => {
  return (
    <Provider store={mockStore(Object.assign({}, baseState, initialState))}>
      {storyFn()}
    </Provider>
  );
};
