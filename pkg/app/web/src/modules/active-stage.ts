import { createSlice, PayloadAction } from "@reduxjs/toolkit";

type ActiveStage = string | null;

const initialState: ActiveStage = null;

export const activeStageSlice = createSlice({
  name: "activeStage",
  initialState: initialState as ActiveStage,
  reducers: {
    updateActiveStage: (_, action: PayloadAction<string>) => {
      return action.payload;
    },
    clearActiveStage: () => {
      return null;
    },
  },
});

export const { updateActiveStage } = activeStageSlice.actions;
