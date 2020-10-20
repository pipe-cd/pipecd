import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { generateApplicationSealedSecret } from "../api/piped";

export type SealedSecret = {
  isLoading: boolean;
  data: string | null;
};

const initialState: SealedSecret = {
  isLoading: false,
  data: null,
};

export const generateSealedSecret = createAsyncThunk<
  string,
  { pipedId: string; data: string }
>("sealedSecret/generate", async (params) => {
  const res = await generateApplicationSealedSecret(params);
  return res.data;
});

export const sealedSecretSlice = createSlice({
  name: "sealedSecret",
  initialState,
  reducers: {
    clearSealedSecret(state) {
      state.data = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(generateSealedSecret.pending, (state) => {
        state.isLoading = true;
      })
      .addCase(generateSealedSecret.fulfilled, (state, action) => {
        state.isLoading = false;
        state.data = action.payload;
      })
      .addCase(generateSealedSecret.rejected, (state) => {
        state.isLoading = false;
      });
  },
});

export const { clearSealedSecret } = sealedSecretSlice.actions;
