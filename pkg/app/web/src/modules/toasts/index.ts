import {
  createEntityAdapter,
  createSlice,
  PayloadAction,
} from "@reduxjs/toolkit";

export type ToastSeverity = "error" | "success" | "warning" | undefined;

export interface IToast {
  id: string;
  message: string;
  severity?: ToastSeverity;
  to?: string;
}

const toastsAdapter = createEntityAdapter<IToast>();

export const { selectAll } = toastsAdapter.getSelectors();

export const toastsSlice = createSlice({
  name: "toasts",
  initialState: toastsAdapter.getInitialState({}),
  reducers: {
    addToast(
      state,
      action: PayloadAction<{
        message: string;
        severity?: ToastSeverity;
        to?: string;
      }>
    ) {
      toastsAdapter.addOne(state, {
        id: `${Date.now()}`,
        message: action.payload.message,
        severity: action.payload.severity,
        to: action.payload.to,
      });
    },
    removeToast(state, action: PayloadAction<{ id: string }>) {
      toastsAdapter.removeOne(state, action.payload.id);
    },
  },
});

export const { addToast, removeToast } = toastsSlice.actions;
