import React from "react";
import { Provider } from "react-redux";
import { createStore } from "../../test-utils";
import { dummyApplication } from "../__fixtures__/dummy-application";
import { dummyApplicationLiveState } from "../__fixtures__/dummy-application-live-state";
import { AppLiveState } from "./app-live-state";

export default {
  title: "APPLICATION/AppLiveState",
  component: AppLiveState,
};

export const overview: React.FC = () => (
  <Provider
    store={createStore({
      applicationLiveState: {
        entities: {
          [dummyApplicationLiveState.applicationId]: dummyApplicationLiveState,
        },
        ids: [dummyApplicationLiveState.applicationId],
        hasError: {},
        loading: {},
      },
    })}
  >
    <AppLiveState applicationId={dummyApplication.id} />
  </Provider>
);

export const loading: React.FC = () => (
  <Provider
    store={createStore({
      applicationLiveState: {
        entities: {},
        ids: [],
        hasError: {},
        loading: {
          [dummyApplication.id]: true,
        },
      },
    })}
  >
    <AppLiveState applicationId={dummyApplication.id} />
  </Provider>
);

export const refresh: React.FC = () => (
  <Provider
    store={createStore({
      applicationLiveState: {
        entities: {
          [dummyApplication.id]: dummyApplicationLiveState,
        },
        ids: [dummyApplication.id],
        hasError: {},
        loading: {
          [dummyApplication.id]: true,
        },
      },
    })}
  >
    <AppLiveState applicationId={dummyApplication.id} />
  </Provider>
);

export const notAvailable: React.FC = () => (
  <Provider
    store={createStore({
      applicationLiveState: {
        entities: {},
        ids: [],
        hasError: {},
        loading: {},
      },
    })}
  >
    <AppLiveState applicationId={dummyApplication.id} />
  </Provider>
);
