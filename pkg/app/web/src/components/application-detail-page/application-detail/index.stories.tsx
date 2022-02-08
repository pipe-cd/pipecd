import { Story } from "@storybook/react";
import { Provider } from "react-redux";
import { ApplicationSyncStatus } from "~/modules/applications";
import type { AppState } from "~/store";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { dummyApplicationLiveState } from "~/__fixtures__/dummy-application-live-state";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { createStore } from "~~/test-utils";
import { ApplicationDetail } from ".";

const dummyStore: Partial<AppState> = {
  applications: {
    entities: {
      [dummyApplication.id]: dummyApplication,
    },
    ids: [dummyApplication.id],
    syncing: {},
    disabling: {},
    adding: false,
    loading: false,
    addedApplicationId: null,
    fetchApplicationError: null,
  },
  applicationLiveState: {
    entities: {
      [dummyApplicationLiveState.applicationId]: dummyApplicationLiveState,
    },
    ids: [dummyApplicationLiveState.applicationId],
    loading: {},
    hasError: {},
  },
  pipeds: {
    updating: false,
    entities: {
      [dummyPiped.id]: dummyPiped,
    },
    ids: [dummyPiped.id],
    registeredPiped: null,
  },
};

export default {
  title: "APPLICATION/ApplicationDetail",
  component: ApplicationDetail,
};

export const Overview: Story = () => (
  <Provider store={createStore(dummyStore)}>
    <ApplicationDetail applicationId={dummyApplication.id} />
  </Provider>
);

export const HasDescription: Story = () => (
  <Provider
    store={createStore({
      ...dummyStore,
      applications: {
        adding: false,
        disabling: {},
        syncing: {},
        entities: {
          [dummyApplication.id]: {
            ...dummyApplication,
            description: "# Heading\n[document page](https://pipecd.dev)\n",
          },
        },
        ids: [dummyApplication.id],
      },
    })}
  >
    <ApplicationDetail applicationId={dummyApplication.id} />
  </Provider>
);

export const Error: Story = () => (
  <Provider
    store={createStore({
      ...dummyStore,
      applications: {
        adding: false,
        disabling: {},
        syncing: {},
        entities: {
          [dummyApplication.id]: {
            ...dummyApplication,
            syncState: {
              shortReason: "deployment has error",
              reason:
                "long reason of deployment failed. you can see this error if click show detail button.",
              headDeploymentId: "deployment-id",
              timestamp: 0,
              status: ApplicationSyncStatus.OUT_OF_SYNC,
            },
          },
        },
        ids: [dummyApplication.id],
      },
    })}
  >
    <ApplicationDetail applicationId={dummyApplication.id} />
  </Provider>
);

export const LoadingLiveState: Story = () => (
  <Provider
    store={createStore({
      ...dummyStore,
      applicationLiveState: {
        entities: {},
        ids: [],
        loading: {
          [dummyApplication.id]: true,
        },
      },
      applications: {
        adding: false,
        disabling: {},
        syncing: {},
        entities: {
          [dummyApplication.id]: {
            ...dummyApplication,
            syncState: {
              shortReason: "deployment has error",
              reason:
                "long reason of deployment failed. you can see this error if click show detail button.",
              headDeploymentId: "deployment-id",
              timestamp: 0,
              status: ApplicationSyncStatus.OUT_OF_SYNC,
            },
          },
        },
        ids: [dummyApplication.id],
      },
    })}
  >
    <ApplicationDetail applicationId={dummyApplication.id} />
  </Provider>
);

export const NotAvailable: Story = () => (
  <Provider
    store={createStore({
      ...dummyStore,
      applicationLiveState: {
        entities: {},
        ids: [],
        loading: {},
      },
      applications: {
        adding: false,
        disabling: {},
        syncing: {},
        entities: {
          [dummyApplication.id]: {
            ...dummyApplication,
            syncState: undefined,
          },
        },
        ids: [dummyApplication.id],
      },
    })}
  >
    <ApplicationDetail applicationId={dummyApplication.id} />
  </Provider>
);
