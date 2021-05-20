// @see https://redux-toolkit.js.org/tutorials/typescript#define-typed-hooks
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import type { AppDispatch, AppState } from "../store";

// eslint-disable-next-line @typescript-eslint/explicit-function-return-type,@typescript-eslint/explicit-module-boundary-types
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<AppState> = useSelector;
