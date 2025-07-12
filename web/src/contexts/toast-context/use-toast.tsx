import { useContext } from "react";
import { ToastContext, ToastContextType } from "./toast-context";

export const useToast = (): ToastContextType => useContext(ToastContext);

export default useToast;
