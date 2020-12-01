import { dummyAPIKey } from "../__fixtures__/dummy-api-key";
import {
  APIKeyModel,
  apiKeysSlice,
  disableAPIKey,
  generateAPIKey,
  getAPIKeys,
} from "./api-keys";

const baseState = {
  error: null,
  items: [],
  generatedKey: null,
  loading: false,
  generating: false,
  disabling: false,
};

describe("apiKeysSlice reducer", () => {
  it("should handle initial state", () => {
    expect(
      apiKeysSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(baseState);
  });

  describe("generateAPIKey", () => {
    const arg = {
      name: "new API key",
      role: APIKeyModel.Role.READ_ONLY,
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

  describe("getAPIKeys", () => {
    const arg = {
      enabled: true,
    };
    it(`should handle ${getAPIKeys.pending.type}`, () => {
      expect(
        apiKeysSlice.reducer(baseState, {
          type: getAPIKeys.pending.type,
          meta: {
            arg,
          },
        })
      ).toEqual({ ...baseState, loading: true });
    });

    it(`should handle ${getAPIKeys.rejected.type}`, () => {
      expect(
        apiKeysSlice.reducer(
          { ...baseState, loading: true },
          {
            type: getAPIKeys.rejected.type,
            error: { message: "API_ERROR" },
            meta: {
              arg,
            },
          }
        )
      ).toEqual({ ...baseState, error: { message: "API_ERROR" } });
    });

    it(`should handle ${getAPIKeys.fulfilled.type}`, () => {
      expect(
        apiKeysSlice.reducer(
          { ...baseState, loading: true },
          {
            type: getAPIKeys.fulfilled.type,
            payload: [dummyAPIKey],
            meta: { arg },
          }
        )
      ).toEqual({ ...baseState, items: [dummyAPIKey] });
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
