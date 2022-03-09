import {
  createAsyncThunk,
  createSlice,
  SerializedError,
  createEntityAdapter,
} from "@reduxjs/toolkit";
import { APIKey } from "pipecd/web/model/apikey_pb";
import * as APIKeysAPI from "~/api/api-keys";

const MODULE_NAME = "apiKeys";

const apiKeysAdapter = createEntityAdapter<APIKey.AsObject>();

export const generateAPIKey = createAsyncThunk<
  string,
  { name: string; role: APIKey.Role }
>(`${MODULE_NAME}/generate`, async ({ name, role }) => {
  const res = await APIKeysAPI.generateAPIKey({ name, role });
  return res.key;
});

export const fetchAPIKeys = createAsyncThunk<
  APIKey.AsObject[],
  { enabled: boolean }
>(`${MODULE_NAME}/getList`, async (options) => {
  const res = await APIKeysAPI.getAPIKeys({ options });
  return res.keysList;
});

export const disableAPIKey = createAsyncThunk<void, { id: string }>(
  `${MODULE_NAME}/disable`,
  async ({ id }) => {
    await APIKeysAPI.disableAPIKey({ id });
  }
);

export const apiKeysSlice = createSlice({
  name: "apiKeys",
  initialState: apiKeysAdapter.getInitialState<{
    generatedKey: string | null;
    loading: boolean;
    generating: boolean;
    disabling: boolean;
    error: null | SerializedError;
  }>({
    generatedKey: null,
    loading: false,
    generating: false,
    disabling: false,
    error: null,
  }),
  reducers: {
    clearGeneratedKey(state) {
      state.generatedKey = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // generateAPIKey
      .addCase(generateAPIKey.pending, (state) => {
        state.generating = true;
        state.generatedKey = null;
        state.error = null;
      })
      .addCase(generateAPIKey.rejected, (state, action) => {
        state.generating = false;
        state.error = action.error;
      })
      .addCase(generateAPIKey.fulfilled, (state, action) => {
        state.generating = false;
        state.generatedKey = action.payload;
      })
      // fetchAPIKeys
      .addCase(fetchAPIKeys.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchAPIKeys.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error;
      })
      .addCase(fetchAPIKeys.fulfilled, (state, action) => {
        state.loading = false;
        apiKeysAdapter.setAll(state, action.payload);
      })
      // disableAPIKey
      .addCase(disableAPIKey.pending, (state) => {
        state.disabling = true;
      })
      .addCase(disableAPIKey.rejected, (state, action) => {
        state.disabling = false;
        state.error = action.error;
      })
      .addCase(disableAPIKey.fulfilled, (state) => {
        state.disabling = false;
      });
  },
});

export const { clearGeneratedKey } = apiKeysSlice.actions;
export const { selectAll, selectById } = apiKeysAdapter.getSelectors();

export { APIKey } from "pipecd/web/model/apikey_pb";
