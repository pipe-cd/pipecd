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
  issuer?: string;
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
        issuer?: string;
        to?: string;
      }>
    ) {
      const toasts = selectAll(state);
      const lastToast = toasts[toasts.length - 1];
      if (
        lastToast?.issuer === action.payload.issuer &&
        lastToast?.message === action.payload.message
      ) {
        return;
      }
      toastsAdapter.addOne(state, {
        id: `${Date.now()}`,
        message: action.payload.message,
        severity: action.payload.severity,
        issuer: action.payload.issuer,
        to: action.payload.to,
      });
    },
    removeToast(state, action: PayloadAction<{ id: string }>) {
      toastsAdapter.removeOne(state, action.payload.id);
    },
  },
});

export const { addToast, removeToast } = toastsSlice.actions;
