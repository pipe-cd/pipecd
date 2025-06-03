import { FC, PropsWithChildren, useCallback, useState } from "react";
import { IToast, ToastContext } from "./toast-context";
import { ToastMessages } from "~/components/toast-messages";

export const ToastProvider: FC<PropsWithChildren<unknown>> = ({ children }) => {
  const [toasts, setToasts] = useState<IToast[]>([]);

  const addToast = useCallback((payload: Omit<IToast, "id">): void => {
    const id = `${Date.now()}`;
    setToasts((prev) => {
      const lastToast = prev[prev.length - 1];

      if (
        lastToast?.issuer === payload.issuer &&
        lastToast?.message === payload.message
      ) {
        return prev;
      }
      return [...prev, { ...payload, id }];
    });
  }, []);

  const removeToast = useCallback((id: string): void => {
    setToasts((prev) => {
      return prev.filter((toast) => toast.id !== id);
    });
  }, []);

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast }}>
      {children}
      <ToastMessages toasts={toasts} onRemoveToast={removeToast} />
    </ToastContext.Provider>
  );
};
