import { DeepPartial } from "@reduxjs/toolkit";
import { Provider } from "react-redux";
import { AppState } from "../src/modules";
import { createStore } from "../test-utils";

export const createDecoratorRedux = (initialState: DeepPartial<AppState>) => (
  storyFn: any
) => <Provider store={createStore(initialState)}>{storyFn()}</Provider>;
