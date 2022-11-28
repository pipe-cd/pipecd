jest.spyOn(Date, "now").mockImplementation(() => 1);

import {
  insightSlice,
  InsightState,
  changeApplication,
  changeRangeFrom,
  changeRangeTo,
} from "./";

const initialState: InsightState = {
  rangeFrom: -2678399999,
  rangeTo: 1,
  applicationId: "",
  labels: [],
};

describe("insightSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      insightSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  it(`should handle ${changeApplication.type}`, () => {
    expect(
      insightSlice.reducer(initialState, {
        type: changeApplication.type,
        payload: "application-1",
      })
    ).toEqual({ ...initialState, applicationId: "application-1" });
  });

  it(`should handle ${changeRangeFrom.type}`, () => {
    expect(
      insightSlice.reducer(initialState, {
        type: changeRangeFrom.type,
        payload: 2,
      })
    ).toEqual({ ...initialState, rangeFrom: 2 });
  });

  it(`should handle ${changeRangeTo.type}`, () => {
    expect(
      insightSlice.reducer(initialState, {
        type: changeRangeTo.type,
        payload: 3,
      })
    ).toEqual({ ...initialState, rangeTo: 3 });
  });
});
