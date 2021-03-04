import { ApplicationKind, ApplicationSyncStatus } from "./applications";
import {
  applicationFilterOptionsSlice,
  clearApplicationFilter,
  updateApplicationFilter,
} from "./application-filter-options";

describe("applicationFilterOptionsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      applicationFilterOptionsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      enabled: undefined,
      envIdsList: [],
      kindsList: [],
      syncStatusesList: [],
      name: "",
    });
  });

  it(`should handle ${updateApplicationFilter.type}`, () => {
    expect(
      applicationFilterOptionsSlice.reducer(
        {
          enabled: undefined,
          envIdsList: [],
          kindsList: [],
          syncStatusesList: [],
          name: "",
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
      enabled: undefined,
      envIdsList: ["env1"],
      kindsList: [ApplicationKind.TERRAFORM],
      syncStatusesList: [ApplicationSyncStatus.SYNCED],
      name: "",
    });
  });

  it(`should handle ${clearApplicationFilter.type}`, () => {
    expect(
      applicationFilterOptionsSlice.reducer(
        {
          enabled: undefined,
          envIdsList: ["env1"],
          kindsList: [ApplicationKind.KUBERNETES],
          syncStatusesList: [ApplicationSyncStatus.DEPLOYING],
          name: "",
        },
        { type: clearApplicationFilter.type }
      )
    ).toEqual({
      enabled: undefined,
      envIdsList: [],
      kindsList: [],
      syncStatusesList: [],
      name: "",
    });
  });
});
