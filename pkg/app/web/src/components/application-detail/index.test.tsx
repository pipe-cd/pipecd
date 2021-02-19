import { DeepPartial } from "@reduxjs/toolkit";
import userEvent from "@testing-library/user-event";
import React from "react";
import { MemoryRouter } from "react-router-dom";
import { createStore, render, screen, waitFor } from "../../../test-utils";
import { server } from "../../mocks/server";
import { AppState } from "../../modules";
import {
  Application,
  syncApplication,
  ApplicationSyncStatus,
} from "../../modules/applications";
import { SyncStrategy } from "../../modules/deployments";
import { dummyApplication } from "../../__fixtures__/dummy-application";
import { dummyApplicationLiveState } from "../../__fixtures__/dummy-application-live-state";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { dummyPiped } from "../../__fixtures__/dummy-piped";
import { ApplicationDetail } from "./";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

const deployingApp: Application.AsObject = {
  ...dummyApplication,
  deploying: true,
  id: "deploying-app",
  syncState: {
    status: ApplicationSyncStatus.DEPLOYING,
    headDeploymentId: "",
    reason: "",
    shortReason: "",
    timestamp: 0,
  },
};

const baseState: DeepPartial<AppState> = {
  applications: {
    ids: [dummyApplication.id, deployingApp.id],
    entities: {
      [dummyApplication.id]: dummyApplication,
      [deployingApp.id]: deployingApp,
    },
    adding: false,
    disabling: {},
    loading: false,
    syncing: {},
  },
  applicationLiveState: {
    ids: [dummyApplicationLiveState.applicationId],
    entities: {
      [dummyApplicationLiveState.applicationId]: dummyApplicationLiveState,
    },
    loading: {},
    hasError: {},
  },
  environments: {
    entities: {
      [dummyEnv.id]: dummyEnv,
    },
    ids: [dummyEnv.id],
  },
  pipeds: {
    entities: {
      [dummyPiped.id]: dummyPiped,
    },
    ids: [dummyPiped.id],
  },
};

describe("ApplicationDetail", () => {
  it("shows application detail and live state", () => {
    const store = createStore(baseState);
    render(
      <MemoryRouter>
        <ApplicationDetail applicationId={dummyApplication.id} />
      </MemoryRouter>,
      {
        store,
      }
    );

    expect(screen.getByText(dummyApplication.name)).toBeInTheDocument();
    expect(screen.getByText(dummyPiped.name)).toBeInTheDocument();
    expect(screen.getByText(dummyEnv.name)).toBeInTheDocument();
    expect(screen.getByText("Healthy")).toBeInTheDocument();
    expect(screen.getByText("Synced")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /sync$/i })).toBeInTheDocument();
  });

  it("shows application sync state as deploying if application is deploying", () => {
    const store = createStore(baseState);
    render(
      <MemoryRouter>
        <ApplicationDetail applicationId={deployingApp.id} />
      </MemoryRouter>,
      {
        store,
      }
    );

    expect(screen.getByText("Deploying")).toBeInTheDocument();
  });

  describe("sync", () => {
    it("dispatch sync action if click sync button", async () => {
      const store = createStore(baseState);
      render(
        <MemoryRouter>
          <ApplicationDetail applicationId={dummyApplication.id} />
        </MemoryRouter>,
        {
          store,
        }
      );

      userEvent.click(screen.getByRole("button", { name: /sync$/i }));

      await waitFor(() =>
        expect(store.getActions()).toMatchObject([
          {
            type: syncApplication.pending.type,
            meta: {
              arg: {
                applicationId: dummyApplication.id,
                syncStrategy: SyncStrategy.AUTO,
              },
            },
          },
        ])
      );
    });

    it("dispatch sync action with selected sync strategy if changed strategy and click the sync button", async () => {
      const store = createStore(baseState);
      render(
        <MemoryRouter>
          <ApplicationDetail applicationId={dummyApplication.id} />
        </MemoryRouter>,
        {
          store,
        }
      );

      userEvent.click(
        screen.getByRole("button", { name: /select sync strategy/i })
      );
      userEvent.click(screen.getByRole("menuitem", { name: /pipeline sync/i }));
      userEvent.click(screen.getByRole("button", { name: /pipeline sync/i }));

      await waitFor(() =>
        expect(store.getActions()).toMatchObject([
          {
            type: syncApplication.pending.type,
            meta: {
              arg: {
                applicationId: dummyApplication.id,
                syncStrategy: SyncStrategy.PIPELINE,
              },
            },
          },
        ])
      );
    });
  });
});
