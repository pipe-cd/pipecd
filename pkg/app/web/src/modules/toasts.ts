import {
  createSlice,
  createEntityAdapter,
  PayloadAction,
} from "@reduxjs/toolkit";

export type ToastSeverity = "error" | "success" | "warning" | undefined;

export interface IToast {
  id: string;
  message: string;
  severity?: ToastSeverity;
}

const toastsAdapter = createEntityAdapter<IToast>();

export const { selectAll } = toastsAdapter.getSelectors();

export const toastsSlice = createSlice({
  name: "toasts",
  initialState: toastsAdapter.getInitialState(),
  reducers: {
    addToast(
      state,
      action: PayloadAction<{ message: string; severity?: ToastSeverity }>
    ) {
      toastsAdapter.addOne(state, {
        id: `${Date.now()}`,
        message: action.payload.message,
        severity: action.payload.severity,
      });
    },
    removeToast(state, action: PayloadAction<{ id: string }>) {
      toastsAdapter.removeOne(state, action.payload.id);
    },
  },
});

export const { addToast, removeToast } = toastsSlice.actions;
