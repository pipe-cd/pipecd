import { Middleware, AnyAction, MiddlewareAPI } from "redux";
import { addToast } from "~/modules/toasts";
import { AppDispatch } from "~/store";

function isPlainAction(action: AnyAction): action is AnyAction {
  return typeof action !== "function";
}

export const thunkErrorHandler: Middleware = ({
  dispatch,
}: MiddlewareAPI<AppDispatch>) => (next) => async (action) => {
  let res;
  try {
    res = await next(action);
  } catch (err) {
    if (process.env.NODE_ENV === "development") {
      console.error(err);
    }

    if (err && err.code && err.message) {
      dispatch(addToast({ message: err.message, severity: "error" }));
    } else {
      throw err;
    }
  }

  if (isPlainAction(action)) {
    if (action.type.includes("rejected")) {
      dispatch(
        addToast({
          message: action.error.message,
          severity: "error",
          issuer: action.type,
        })
      );
    }
  }

  return res;
};
