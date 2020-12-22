import { loginSlice, clearProjectName, setProjectName } from "./login";

describe("loginSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      loginSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      projectName: null,
    });
  });

  it(`should handle ${setProjectName.type}`, () => {
    expect(
      loginSlice.reducer(
        {
          projectName: null,
        },
        {
          type: setProjectName.type,
          payload: "pipecd",
        }
      )
    ).toEqual({
      projectName: "pipecd",
    });
  });

  it(`should handle ${clearProjectName.type}`, () => {
    expect(
      loginSlice.reducer(
        {
          projectName: "pipecd",
        },
        {
          type: clearProjectName.type,
        }
      )
    ).toEqual({
      projectName: null,
    });
  });
});
