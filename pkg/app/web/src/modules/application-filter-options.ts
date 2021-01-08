import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { ApplicationKind, ApplicationSyncStatus } from "./applications";

export type ApplicationFilterOptions = {
  enabled?: { value: boolean };
  kindsList: ApplicationKind[];
  envIdsList: string[];
  syncStatusesList: ApplicationSyncStatus[];
};

const initialState: ApplicationFilterOptions = {
  enabled: undefined,
  kindsList: [],
  envIdsList: [],
  syncStatusesList: [],
};

export const applicationFilterOptionsSlice = createSlice({
  name: "applicationFilterOptions",
  initialState,
  reducers: {
    updateApplicationFilter(
      state,
      action: PayloadAction<Partial<ApplicationFilterOptions>>
    ) {
      return { ...state, ...action.payload };
    },
    clearApplicationFilter() {
      return initialState;
    },
  },
});

export const {
  updateApplicationFilter,
  clearApplicationFilter,
} = applicationFilterOptionsSlice.actions;
