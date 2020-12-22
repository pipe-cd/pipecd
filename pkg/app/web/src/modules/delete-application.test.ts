import {
  deleteApplicationSlice,
  DeleteApplicationState,
} from "./delete-application";

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
});
