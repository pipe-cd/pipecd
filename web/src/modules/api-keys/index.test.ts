import { dummyAPIKey } from "~/__fixtures__/dummy-api-key";
import {
  APIKey,
  apiKeysSlice,
  disableAPIKey,
  generateAPIKey,
  fetchAPIKeys,
  clearGeneratedKey,
} from ".";

const baseState = {
  error: null,
  generatedKey: null,
  loading: false,
  generating: false,
  disabling: false,
  entities: {},
  ids: [],
};

describe("apiKeysSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      apiKeysSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(baseState);
  });

  it("should handle clearGeneratedKey", () => {
    expect(
      apiKeysSlice.reducer(
        {
          ...baseState,
          generatedKey:
            "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.bspmf2xvt74area19iaxl0yh33jzwelq493vzil0orgzylrdb1",
        },
        {
          type: clearGeneratedKey.type,
        }
      )
    ).toEqual(baseState);
  });

  describe("generateAPIKey", () => {
    const arg = {
      name: "new API key",
      role: APIKey.Role.READ_ONLY,
    };
    it(`should handle ${generateAPIKey.pending.type}`, () => {
      expect(
        apiKeysSlice.reducer(baseState, {
          type: generateAPIKey.pending.type,
          meta: {
            arg,
          },
        })
      ).toEqual({ ...baseState, generating: true });
    });

    it(`should handle ${generateAPIKey.rejected.type}`, () => {
      expect(
        apiKeysSlice.reducer(
          { ...baseState, generating: true },
          {
            type: generateAPIKey.rejected.type,
            error: { message: "API_ERROR" },
            meta: {
              arg,
            },
          }
        )
      ).toEqual({ ...baseState, error: { message: "API_ERROR" } });
    });

    it(`should handle ${generateAPIKey.fulfilled.type}`, () => {
      expect(
        apiKeysSlice.reducer(
          { ...baseState, generating: true },
          {
            type: generateAPIKey.fulfilled.type,
            payload: "API_KEY",
            meta: { arg },
          }
        )
      ).toEqual({ ...baseState, generatedKey: "API_KEY", generating: false });
    });
  });

  describe("fetchAPIKeys", () => {
    const arg = {
      enabled: true,
    };
    it(`should handle ${fetchAPIKeys.pending.type}`, () => {
      expect(
        apiKeysSlice.reducer(baseState, {
          type: fetchAPIKeys.pending.type,
          meta: {
            arg,
          },
        })
      ).toEqual({ ...baseState, loading: true });
    });

    it(`should handle ${fetchAPIKeys.rejected.type}`, () => {
      expect(
        apiKeysSlice.reducer(
          { ...baseState, loading: true },
          {
            type: fetchAPIKeys.rejected.type,
            error: { message: "API_ERROR" },
            meta: {
              arg,
            },
          }
        )
      ).toEqual({ ...baseState, error: { message: "API_ERROR" } });
    });

    it(`should handle ${fetchAPIKeys.fulfilled.type}`, () => {
      expect(
        apiKeysSlice.reducer(
          { ...baseState, loading: true },
          {
            type: fetchAPIKeys.fulfilled.type,
            payload: [dummyAPIKey],
            meta: { arg },
          }
        )
      ).toEqual({
        ...baseState,
        entities: { [dummyAPIKey.id]: dummyAPIKey },
        ids: [dummyAPIKey.id],
      });

      expect(
        apiKeysSlice.reducer(
          {
            ...baseState,
            entities: { [dummyAPIKey.id]: dummyAPIKey },
            ids: [dummyAPIKey.id],
            loading: true,
          },
          {
            type: fetchAPIKeys.fulfilled.type,
            payload: [],
            meta: { arg },
          }
        )
      ).toEqual({
        ...baseState,
        entities: {},
        ids: [],
      });
    });
  });

  describe("disableAPIKey", () => {
    const arg = {
      id: "api-key-1",
    };

    it(`should handle ${disableAPIKey.pending.type}`, () => {
      expect(
        apiKeysSlice.reducer(baseState, {
          type: disableAPIKey.pending.type,
          meta: {
            arg,
          },
        })
      ).toEqual({ ...baseState, disabling: true });
    });

    it(`should handle ${disableAPIKey.rejected.type}`, () => {
      expect(
        apiKeysSlice.reducer(
          { ...baseState, disabling: true },
          {
            type: disableAPIKey.rejected.type,
            error: { message: "API_ERROR" },
            meta: {
              arg,
            },
          }
        )
      ).toEqual({ ...baseState, error: { message: "API_ERROR" } });
    });

    it(`should handle ${disableAPIKey.fulfilled.type}`, () => {
      expect(
        apiKeysSlice.reducer(
          { ...baseState, disabling: true },
          {
            type: disableAPIKey.fulfilled.type,
            payload: [],
            meta: { arg },
          }
        )
      ).toEqual(baseState);
    });
  });
});
