import React from "react";
import { Provider } from "react-redux";
import { ApplicationSyncStatus } from "pipe/pkg/app/web/model/application_pb";
import { createStore } from "../../test-utils";
import { dummyApplication } from "../__fixtures__/dummy-application";
import { dummyApplicationLiveState } from "../__fixtures__/dummy-application-live-state";
import { dummyEnv } from "../__fixtures__/dummy-environment";
import { dummyPiped } from "../__fixtures__/dummy-piped";
import { ApplicationDetail } from "./application-detail";
import { AppState } from "../modules";

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
    fetchApplicationError: null,
  },
  environments: {
    entities: {
      [dummyEnv.id]: dummyEnv,
    },
    ids: [dummyEnv.id],
  },
  applicationLiveState: {
    entities: {
      [dummyApplicationLiveState.applicationId]: dummyApplicationLiveState,
    },
    ids: [dummyApplicationLiveState.applicationId],
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

export const overview: React.FC = () => (
  <Provider store={createStore(dummyStore)}>
    <ApplicationDetail applicationId={dummyApplication.id} />
  </Provider>
);

export const error: React.FC = () => (
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

export const loadingLiveState: React.FC = () => (
  <Provider
    store={createStore({
      ...dummyStore,
      applicationLiveState: {
        entities: {},
        ids: [],
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
