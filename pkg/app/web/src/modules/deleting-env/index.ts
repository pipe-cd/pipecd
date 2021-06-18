import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import type { AppState } from "~/store";
import * as applicationsAPI from "~/api/applications";
import { fetchApplicationsByEnv } from "../applications";
import { deleteEnvironment } from "../environments";

export interface DeletingEnvState {
  env: { id: string; name: string } | null;
  targetApplications: string[];
}

const initialState: DeletingEnvState = {
  env: null,
  targetApplications: [],
};

export const deletingEnv = createAsyncThunk<
  void,
  void,
  {
    state: AppState;
  }
>("deleting-env/delete", async (_, thunkAPI) => {
  const state = thunkAPI.getState();

  if (state.deleteApplication.applicationId) {
    await applicationsAPI.deleteApplication({
      applicationId: state.deleteApplication.applicationId,
    });
  }
});

export const deletingEnvSlice = createSlice({
  name: "deleting-env",
  initialState,
  reducers: {
    setTargetEnv(state, action: PayloadAction<{ id: string; name: string }>) {
      state.env = action.payload;
    },
    clearTargetEnv(state) {
      state.env = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchApplicationsByEnv.fulfilled, (state, action) => {
        state.targetApplications = action.payload.map((app) => app.id);
      })
      .addCase(deleteEnvironment.fulfilled, (state) => {
        state.env = null;
      })
      .addCase(deleteEnvironment.rejected, (state) => {
        state.env = null;
      });
  },
});

export const { setTargetEnv, clearTargetEnv } = deletingEnvSlice.actions;
