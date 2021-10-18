import { setupServer } from "msw/node";
import { createStore } from "~~/test-utils";
import { listApplicationsHandler } from "~/mocks/services/application";
import {
  dummyApplication,
  dummyApplicationSyncState,
} from "~/__fixtures__/dummy-application";
import {
  addApplication,
  Application,
  applicationsSlice,
  ApplicationsState,
  ApplicationSyncStatus,
  disableApplication,
  fetchApplication,
  fetchApplications,
  syncApplication,
} from "./";
import { Command, CommandStatus, fetchCommand } from "../commands";

const server = setupServer(listApplicationsHandler);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

const baseState: ApplicationsState = {
  addedApplicationId: null,
  adding: false,
  disabling: {},
  entities: {},
  ids: [],
  loading: false,
  syncing: {},
  fetchApplicationError: null,
};

describe("fetchApplications", () => {
  it("should get applications by options", async () => {
    const store = createStore({});

    await store.dispatch(
      fetchApplications({ activeStatus: "enabled", envId: "env-1" })
    );

    expect(store.getActions()).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ type: fetchApplications.pending.type }),
        expect.objectContaining({
          type: fetchApplications.fulfilled.type,
          payload: expect.arrayContaining([dummyApplication]),
        }),
      ])
    );
  });
});

describe("applicationsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      applicationsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(baseState);
  });

  describe("fetchApplications", () => {
    it(`should handle ${fetchApplications.pending.type}`, () => {
      expect(
        applicationsSlice.reducer(undefined, {
          type: fetchApplications.pending.type,
        })
      ).toEqual({
        ...baseState,
        loading: true,
      });
    });

    it(`should handle ${fetchApplications.fulfilled.type}`, () => {
      expect(
        applicationsSlice.reducer(baseState, {
          type: fetchApplications.fulfilled.type,
          payload: [dummyApplication],
          loading: true,
        })
      ).toEqual({
        ...baseState,
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      });

      const deletedApp = { ...dummyApplication, id: "app-2", deleted: true };
      const oldApp = { ...dummyApplication, id: "app-3", deleted: false };
      expect(
        applicationsSlice.reducer(baseState, {
          type: fetchApplications.fulfilled.type,
          payload: [dummyApplication, deletedApp, oldApp],
          loading: true,
        })
      ).toEqual({
        ...baseState,
        entities: {
          [dummyApplication.id]: dummyApplication,
          [oldApp.id]: oldApp,
        },
        ids: [dummyApplication.id, oldApp.id],
      });
    });

    it(`should handle ${fetchApplications.rejected.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            ...baseState,
            loading: true,
          },
          {
            type: fetchApplications.rejected.type,
          }
        )
      ).toEqual(baseState);
    });
  });

  describe("fetchApplication", () => {
    it(`should handle ${fetchApplication.fulfilled.type}`, () => {
      const updatedApplication: Application.AsObject = {
        ...dummyApplication,
        syncState: {
          ...dummyApplicationSyncState,
          status: ApplicationSyncStatus.OUT_OF_SYNC,
        },
      };
      expect(
        applicationsSlice.reducer(
          {
            ...baseState,
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
          },
          {
            type: fetchApplication.fulfilled.type,
            payload: updatedApplication,
          }
        )
      ).toEqual({
        ...baseState,
        entities: {
          [dummyApplication.id]: updatedApplication,
        },
        ids: [dummyApplication.id],
      });
    });

    it(`should handle ${fetchApplication.fulfilled.type} without payload`, () => {
      expect(
        applicationsSlice.reducer(
          {
            ...baseState,
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
          },
          {
            type: fetchApplication.fulfilled.type,
          }
        )
      ).toEqual({
        ...baseState,
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      });
    });
  });

  describe("addApplication", () => {
    it(`should handle ${addApplication.pending.type}`, () => {
      expect(
        applicationsSlice.reducer(baseState, {
          type: addApplication.pending.type,
        })
      ).toEqual({
        ...baseState,
        adding: true,
      });
    });

    it(`should handle ${addApplication.fulfilled.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            ...baseState,
            adding: true,
          },
          {
            type: addApplication.fulfilled.type,
            payload: dummyApplication.id,
          }
        )
      ).toEqual({
        ...baseState,
        addedApplicationId: dummyApplication.id,
      })
    });

    it(`should handle ${addApplication.rejected.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            ...baseState,
            adding: true,
          },
          {
            type: addApplication.rejected.type,
          }
        )
      ).toEqual(baseState);
    });
  });

  it(`should handle ${fetchCommand.fulfilled.type}`, () => {
    expect(
      applicationsSlice.reducer(
        {
          ...baseState,
          syncing: {
            "app-1": true,
          },
        },
        {
          type: fetchCommand.fulfilled.type,
          payload: {
            type: Command.Type.SYNC_APPLICATION,
            status: CommandStatus.COMMAND_SUCCEEDED,
            applicationId: "app-1",
          },
        }
      )
    ).toEqual({
      ...baseState,
      syncing: {
        "app-1": false,
      },
    });

    expect(
      applicationsSlice.reducer(
        {
          ...baseState,
          syncing: {
            "app-1": true,
          },
        },
        {
          type: fetchCommand.fulfilled.type,
          payload: {
            type: Command.Type.SYNC_APPLICATION,
            status: CommandStatus.COMMAND_NOT_HANDLED_YET,
            applicationId: "app-1",
          },
        }
      )
    ).toEqual({
      ...baseState,
      syncing: {
        "app-1": true,
      },
    });
  });

  describe("disableApplication", () => {
    it(`should handle ${disableApplication.pending.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            ...baseState,
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
          },
          {
            type: disableApplication.pending.type,
            meta: {
              arg: {
                applicationId: dummyApplication.id,
              },
            },
          }
        )
      ).toEqual({
        ...baseState,
        disabling: {
          [dummyApplication.id]: true,
        },
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      });
    });

    it(`should handle ${disableApplication.fulfilled.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            ...baseState,
            disabling: {
              [dummyApplication.id]: true,
            },
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
          },
          {
            type: disableApplication.fulfilled.type,
            meta: {
              arg: {
                applicationId: dummyApplication.id,
              },
            },
          }
        )
      ).toEqual({
        ...baseState,
        disabling: {
          [dummyApplication.id]: false,
        },
      });
    });

    it(`should handle ${disableApplication.rejected.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            ...baseState,
            disabling: {
              [dummyApplication.id]: true,
            },
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
          },
          {
            type: disableApplication.rejected.type,
            meta: {
              arg: {
                applicationId: dummyApplication.id,
              },
            },
          }
        )
      ).toEqual({
        ...baseState,
        disabling: {
          [dummyApplication.id]: false,
        },
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      });
    });
  });

  describe("syncApplication", () => {
    it(`should handle ${syncApplication.pending.type}`, () => {
      expect(
        applicationsSlice.reducer(baseState, {
          type: syncApplication.pending.type,
          meta: {
            arg: {
              applicationId: dummyApplication.id,
            },
          },
        })
      ).toEqual({
        ...baseState,
        syncing: {
          [dummyApplication.id]: true,
        },
      });
    });
  });
});
