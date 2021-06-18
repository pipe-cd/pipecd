import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import type { AppState } from "~/store";
import * as applicationsAPI from "~/api/applications";

export interface DeleteApplicationState {
  applicationId: string | null;
  deleting: boolean;
}

const initialState: DeleteApplicationState = {
  applicationId: null,
  deleting: false,
};

export const deleteApplication = createAsyncThunk<
  void,
  void,
  {
    state: AppState;
  }
>("applications/delete", async (_, thunkAPI) => {
  const state = thunkAPI.getState();

  if (state.deleteApplication.applicationId) {
    await applicationsAPI.deleteApplication({
      applicationId: state.deleteApplication.applicationId,
    });
  }
});

export const deleteApplicationSlice = createSlice({
  name: "deleteApplication",
  initialState,
  reducers: {
    setDeletingAppId(state, action: PayloadAction<string>) {
      state.applicationId = action.payload;
    },
    clearDeletingApp(state) {
      state.applicationId = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(deleteApplication.pending, (state) => {
        state.deleting = true;
      })
      .addCase(deleteApplication.rejected, (state) => {
        state.deleting = false;
      })
      .addCase(deleteApplication.fulfilled, (state) => {
        state.deleting = false;
        state.applicationId = null;
      });
  },
});

export const {
  clearDeletingApp,
  setDeletingAppId,
} = deleteApplicationSlice.actions;
