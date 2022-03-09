import {
  clearDeletingApp,
  deleteApplication,
  deleteApplicationSlice,
  DeleteApplicationState,
  setDeletingAppId,
} from "./";

const initialState: DeleteApplicationState = {
  applicationId: null,
  deleting: false,
};

describe("deleteApplicationSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      deleteApplicationSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  it("should handle setDeletingAppId", () => {
    expect(
      deleteApplicationSlice.reducer(initialState, {
        type: setDeletingAppId.type,
        payload: "application-1",
      })
    ).toEqual({ ...initialState, applicationId: "application-1" });
  });

  it("should handle clearDeletingApp", () => {
    expect(
      deleteApplicationSlice.reducer(
        { ...initialState, applicationId: "application-1" },
        {
          type: clearDeletingApp.type,
        }
      )
    ).toEqual(initialState);
  });

  describe("deleteApplication", () => {
    it(`should handle ${deleteApplication.pending.type}`, () => {
      expect(
        deleteApplicationSlice.reducer(initialState, {
          type: deleteApplication.pending.type,
        })
      ).toEqual({ ...initialState, deleting: true });
    });

    it(`should handle ${deleteApplication.rejected.type}`, () => {
      expect(
        deleteApplicationSlice.reducer(
          { ...initialState, deleting: true },
          {
            type: deleteApplication.rejected.type,
          }
        )
      ).toEqual(initialState);
    });

    it(`should handle ${deleteApplication.fulfilled.type}`, () => {
      expect(
        deleteApplicationSlice.reducer(
          { ...initialState, deleting: true, applicationId: "application-1" },
          {
            type: deleteApplication.fulfilled.type,
          }
        )
      ).toEqual(initialState);
    });
  });
});
