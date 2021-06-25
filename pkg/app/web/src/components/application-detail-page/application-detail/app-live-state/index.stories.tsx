import { Story } from "@storybook/react";
import { Provider } from "react-redux";
import { createStore } from "~~/test-utils";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { dummyApplicationLiveState } from "~/__fixtures__/dummy-application-live-state";
import { AppLiveState } from ".";

export default {
  title: "APPLICATION/AppLiveState",
  component: AppLiveState,
};

export const Overview: Story = () => (
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

export const loading: Story = () => (
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

export const refresh: Story = () => (
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

export const notAvailable: Story = () => (
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
