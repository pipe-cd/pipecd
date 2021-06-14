import {
  updateApplicationSlice,
  clearUpdateTarget,
  setUpdateTargetId,
  updateApplication,
  UpdateApplicationState,
} from "./";

const initialState: UpdateApplicationState = {
  updating: false,
  targetId: null,
};

describe("updateApplicationSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      updateApplicationSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  it(`should handle ${setUpdateTargetId.type}`, () => {
    expect(
      updateApplicationSlice.reducer(initialState, {
        type: setUpdateTargetId.type,
        payload: "application-1",
      })
    ).toEqual({ ...initialState, targetId: "application-1" });
  });

  it(`should handle ${clearUpdateTarget.type}`, () => {
    expect(
      updateApplicationSlice.reducer(
        { ...initialState, targetId: "application-1" },
        {
          type: clearUpdateTarget.type,
        }
      )
    ).toEqual(initialState);
  });

  describe("updateApplication", () => {
    it(`should handle ${updateApplication.pending.type}`, () => {
      expect(
        updateApplicationSlice.reducer(initialState, {
          type: updateApplication.pending.type,
        })
      ).toEqual({ ...initialState, updating: true });
    });

    it(`should handle ${updateApplication.rejected.type}`, () => {
      expect(
        updateApplicationSlice.reducer(
          { ...initialState, updating: true },
          {
            type: updateApplication.rejected.type,
          }
        )
      ).toEqual(initialState);
    });

    it(`should handle ${updateApplication.fulfilled.type}`, () => {
      expect(
        updateApplicationSlice.reducer(
          { ...initialState, updating: true },
          {
            type: updateApplication.fulfilled.type,
          }
        )
      ).toEqual(initialState);
    });
  });
});
