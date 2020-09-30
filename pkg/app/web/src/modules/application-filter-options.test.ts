import { ApplicationKind, ApplicationSyncStatus } from "./applications";
import {
  applicationFilterOptionsSlice,
  clearApplicationFilter,
  updateApplicationFilter,
} from "./application-filter-options";

describe("applicationFilterOptionsSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      applicationFilterOptionsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      enabled: {
        value: true,
      },
      envIdsList: [],
      kindsList: [],
      syncStatusesList: [],
    });
  });

  it(`should handle ${updateApplicationFilter.type}`, () => {
    expect(
      applicationFilterOptionsSlice.reducer(
        {
          enabled: {
            value: true,
          },
          envIdsList: [],
          kindsList: [],
          syncStatusesList: [],
        },
        {
          type: updateApplicationFilter.type,
          payload: {
            envIdsList: ["env1"],
            kindsList: [ApplicationKind.TERRAFORM],
            syncStatusesList: [ApplicationSyncStatus.SYNCED],
          },
        }
      )
    ).toEqual({
      enabled: {
        value: true,
      },
      envIdsList: ["env1"],
      kindsList: [ApplicationKind.TERRAFORM],
      syncStatusesList: [ApplicationSyncStatus.SYNCED],
    });
  });

  it(`should handle ${clearApplicationFilter.type}`, () => {
    expect(
      applicationFilterOptionsSlice.reducer(
        {
          enabled: {
            value: true,
          },
          envIdsList: ["env1"],
          kindsList: [ApplicationKind.KUBERNETES],
          syncStatusesList: [ApplicationSyncStatus.DEPLOYING],
        },
        { type: clearApplicationFilter.type }
      )
    ).toEqual({
      enabled: {
        value: true,
      },
      envIdsList: [],
      kindsList: [],
      syncStatusesList: [],
    });
  });
});
