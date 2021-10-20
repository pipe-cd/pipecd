import { toastsSlice, addToast, removeToast } from "./";

describe("toastsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      toastsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      entities: {},
      ids: [],
    });
  });

  it(`should handle ${addToast.type}`, () => {
    jest.spyOn(Date, "now").mockImplementation(() => 1);
    expect(
      toastsSlice.reducer(
        {
          entities: {},
          ids: [],
        },
        {
          type: addToast.type,
          payload: {
            message: "toast message",
            severity: "success",
          },
        }
      )
    ).toEqual({
      entities: {
        "1": { id: "1", message: "toast message", severity: "success" },
      },
      ids: ["1"],
    });
  });

  it(`should handle ${removeToast.type}`, () => {
    expect(
      toastsSlice.reducer(
        {
          entities: {
            "1": { id: "1", message: "toast message", severity: "success" },
          },
          ids: ["1"],
        },
        {
          type: removeToast.type,
          payload: {
            id: "1",
          },
        }
      )
    ).toEqual({
      entities: {},
      ids: [],
    });
  });

  it(`should not add the same toast which cause by the same reason with the latest toast`, () => {
    expect(
      toastsSlice.reducer(
        {
          entities: {
            "1": { id: "1", message: "toast message", severity: "error", issuer: "api/rejected" },
          },
          ids: ["1"],
        },
        {
          type: addToast.type,
          payload: {
            message: "toast message",
            severity: "error",
            issuer: "api/rejected",
          },
        }
      )
    ).toEqual({
      entities: {
        "1": { id: "1", message: "toast message", severity: "error", issuer: "api/rejected" },
      },
      ids: ["1"],
    })
  });
});
