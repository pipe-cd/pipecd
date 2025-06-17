import React from "react";

export type ToastSeverity = "error" | "success" | "warning" | undefined;

export type IToast = {
  id: string;
  message: string;
  severity?: ToastSeverity;
  issuer?: string;
  to?: string;
};

export type ToastContextType = {
  toasts: IToast[];
  addToast: (payload: Omit<IToast, "id">) => void;
  removeToast: (id: string) => void;
};

export const ToastContext = React.createContext<ToastContextType>({
  toasts: [],
  addToast: () => Promise.resolve(),
  removeToast: () => Promise.resolve(),
});
