import { dummyApplication } from "../__fixtures__/dummy-application";
import { createStore } from "../../test-utils";
import {
  addApplication,
  Application,
  applicationsSlice,
  ApplicationSyncStatus,
  disableApplication,
  fetchApplication,
  fetchApplications,
  syncApplication,
} from "./applications";
import { CommandModel, CommandStatus, fetchCommand } from "./commands";
import * as applicationsAPI from "../api/applications";

describe("fetchApplications", () => {
  it("should get applications by options", async () => {
    jest
      .spyOn(applicationsAPI, "getApplications")
      .mockImplementation(() =>
        Promise.resolve({ applicationsList: [dummyApplication] })
      );
    const store = createStore({
      applicationFilterOptions: {
        enabled: { value: true },
        envIdsList: ["env-1"],
      },
    });

    await store.dispatch(fetchApplications());

    expect(store.getActions()).toMatchObject([
      { type: fetchApplications.pending.type },
      { type: fetchApplications.fulfilled.type, payload: [dummyApplication] },
    ]);

    expect(applicationsAPI.getApplications).toHaveBeenCalledWith({
      options: {
        enabled: { value: true },
        envIdsList: ["env-1"],
      },
    });
  });
});

describe("applicationsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      applicationsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      adding: false,
      disabling: {},
      entities: {},
      ids: [],
      loading: false,
      syncing: {},
    });
  });

  describe("fetchApplications", () => {
    it(`should handle ${fetchApplications.pending.type}`, () => {
      expect(
        applicationsSlice.reducer(undefined, {
          type: fetchApplications.pending.type,
        })
      ).toEqual({
        adding: false,
        disabling: {},
        entities: {},
        ids: [],
        loading: true,
        syncing: {},
      });
    });

    it(`should handle ${fetchApplications.fulfilled.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            adding: false,
            disabling: {},
            entities: {},
            ids: [],
            loading: false,
            syncing: {},
          },
          {
            type: fetchApplications.fulfilled.type,
            payload: [dummyApplication],
          }
        )
      ).toEqual({
        adding: false,
        disabling: {},
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
        loading: false,
        syncing: {},
      });
    });

    it(`should handle ${fetchApplications.rejected.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            adding: false,
            disabling: {},
            entities: {},
            ids: [],
            loading: true,
            syncing: {},
          },
          {
            type: fetchApplications.rejected.type,
          }
        )
      ).toEqual({
        adding: false,
        disabling: {},
        entities: {},
        ids: [],
        loading: false,
        syncing: {},
      });
    });
  });

  describe("fetchApplication", () => {
    it(`should handle ${fetchApplication.fulfilled.type}`, () => {
      const updatedApplication: Application = {
        ...dummyApplication,
        syncState: {
          ...dummyApplication.syncState,
          status: ApplicationSyncStatus.OUT_OF_SYNC,
        },
      };
      expect(
        applicationsSlice.reducer(
          {
            adding: false,
            disabling: {},
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
            loading: false,
            syncing: {},
          },
          {
            type: fetchApplication.fulfilled.type,
            payload: updatedApplication,
          }
        )
      ).toEqual({
        adding: false,
        disabling: {},
        entities: {
          [dummyApplication.id]: updatedApplication,
        },
        ids: [dummyApplication.id],
        loading: false,
        syncing: {},
      });
    });

    it(`should handle ${fetchApplication.fulfilled.type} without payload`, () => {
      expect(
        applicationsSlice.reducer(
          {
            adding: false,
            disabling: {},
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
            loading: false,
            syncing: {},
          },
          {
            type: fetchApplication.fulfilled.type,
          }
        )
      ).toEqual({
        adding: false,
        disabling: {},
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
        loading: false,
        syncing: {},
      });
    });
  });

  describe("addApplication", () => {
    it(`should handle ${addApplication.pending.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            adding: false,
            disabling: {},
            entities: {},
            ids: [],
            loading: false,
            syncing: {},
          },
          {
            type: addApplication.pending.type,
          }
        )
      ).toEqual({
        adding: true,
        disabling: {},
        entities: {},
        ids: [],
        loading: false,
        syncing: {},
      });
    });

    it(`should handle ${addApplication.fulfilled.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            adding: true,
            disabling: {},
            entities: {},
            ids: [],
            loading: false,
            syncing: {},
          },
          {
            type: addApplication.fulfilled.type,
          }
        )
      ).toEqual({
        adding: false,
        disabling: {},
        entities: {},
        ids: [],
        loading: false,
        syncing: {},
      });
    });

    it(`should handle ${addApplication.rejected.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            adding: true,
            disabling: {},
            entities: {},
            ids: [],
            loading: false,
            syncing: {},
          },
          {
            type: addApplication.rejected.type,
          }
        )
      ).toEqual({
        adding: false,
        disabling: {},
        entities: {},
        ids: [],
        loading: false,
        syncing: {},
      });
    });
  });

  it(`should handle ${fetchCommand.fulfilled.type}`, () => {
    expect(
      applicationsSlice.reducer(
        {
          adding: false,
          disabling: {},
          entities: {},
          ids: [],
          loading: false,
          syncing: {
            "app-1": true,
          },
        },
        {
          type: fetchCommand.fulfilled.type,
          payload: {
            type: CommandModel.Type.SYNC_APPLICATION,
            status: CommandStatus.COMMAND_SUCCEEDED,
            applicationId: "app-1",
          },
        }
      )
    ).toEqual({
      adding: false,
      disabling: {},
      entities: {},
      ids: [],
      loading: false,
      syncing: {
        "app-1": false,
      },
    });

    expect(
      applicationsSlice.reducer(
        {
          adding: false,
          disabling: {},
          entities: {},
          ids: [],
          loading: false,
          syncing: {
            "app-1": true,
          },
        },
        {
          type: fetchCommand.fulfilled.type,
          payload: {
            type: CommandModel.Type.SYNC_APPLICATION,
            status: CommandStatus.COMMAND_NOT_HANDLED_YET,
            applicationId: "app-1",
          },
        }
      )
    ).toEqual({
      adding: false,
      disabling: {},
      entities: {},
      ids: [],
      loading: false,
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
            adding: false,
            disabling: {},
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
            loading: false,
            syncing: {},
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
        adding: false,
        disabling: {
          [dummyApplication.id]: true,
        },
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
        loading: false,
        syncing: {},
      });
    });

    it(`should handle ${disableApplication.fulfilled.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            adding: false,
            disabling: {
              [dummyApplication.id]: true,
            },
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
            loading: false,
            syncing: {},
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
        adding: false,
        disabling: {
          [dummyApplication.id]: false,
        },
        entities: {},
        ids: [],
        loading: false,
        syncing: {},
      });
    });

    it(`should handle ${disableApplication.rejected.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            adding: false,
            disabling: {
              [dummyApplication.id]: true,
            },
            entities: {
              [dummyApplication.id]: dummyApplication,
            },
            ids: [dummyApplication.id],
            loading: false,
            syncing: {},
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
        adding: false,
        disabling: {
          [dummyApplication.id]: false,
        },
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
        loading: false,
        syncing: {},
      });
    });
  });

  describe("syncApplication", () => {
    it(`should handle ${syncApplication.pending.type}`, () => {
      expect(
        applicationsSlice.reducer(
          {
            adding: false,
            disabling: {},
            entities: {},
            ids: [],
            loading: false,
            syncing: {},
          },
          {
            type: syncApplication.pending.type,
            meta: {
              arg: {
                applicationId: dummyApplication.id,
              },
            },
          }
        )
      ).toEqual({
        adding: false,
        disabling: {},
        entities: {},
        ids: [],
        loading: false,
        syncing: {
          [dummyApplication.id]: true,
        },
      });
    });
  });
});
