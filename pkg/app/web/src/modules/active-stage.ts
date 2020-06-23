import { createSlice, PayloadAction } from "@reduxjs/toolkit";

export type ActiveStage = {
  id: string;
  name: string;
} | null;

const initialState: ActiveStage = null;

export const activeStageSlice = createSlice({
  name: "activeStage",
  initialState: initialState as ActiveStage,
  reducers: {
    updateActiveStage: (_, action: PayloadAction<ActiveStage>) => {
      return action.payload;
    },
    clearActiveStage: () => {
      return null;
    },
  },
});

export const { updateActiveStage, clearActiveStage } = activeStageSlice.actions;
