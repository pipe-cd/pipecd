import {
  sealedSecretSlice,
  SealedSecretState,
  clearSealedSecret,
  generateSealedSecret,
} from "./";

const initialState: SealedSecretState = {
  isLoading: false,
  data: null,
};

describe("sealedSecretSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      sealedSecretSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  it(`should handle ${clearSealedSecret.type}`, () => {
    expect(
      sealedSecretSlice.reducer(
        { ...initialState, data: "secret" },
        {
          type: clearSealedSecret.type,
        }
      )
    ).toEqual(initialState);
  });

  describe("generateSealedSecret", () => {
    it(`should handle ${generateSealedSecret.pending.type}`, () => {
      expect(
        sealedSecretSlice.reducer(initialState, {
          type: generateSealedSecret.pending.type,
        })
      ).toEqual({ ...initialState, isLoading: true });
    });

    it(`should handle ${generateSealedSecret.rejected.type}`, () => {
      expect(
        sealedSecretSlice.reducer(
          { ...initialState, isLoading: true },
          {
            type: generateSealedSecret.rejected.type,
          }
        )
      ).toEqual(initialState);
    });

    it(`should handle ${generateSealedSecret.fulfilled.type}`, () => {
      expect(
        sealedSecretSlice.reducer(
          { ...initialState, isLoading: true },
          {
            type: generateSealedSecret.fulfilled.type,
            payload: "secret",
          }
        )
      ).toEqual({ ...initialState, data: "secret" });
    });
  });
});
