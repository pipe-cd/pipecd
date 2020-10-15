import { Snackbar } from "@material-ui/core";
import MuiAlert from "@material-ui/lab/Alert";
import React, { FC } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import { IToast, removeToast, selectAll } from "../modules/toasts";

const AUTO_HIDE_DURATION = 5000;

export const Toasts: FC = () => {
  const dispatch = useDispatch();
  const toasts = useSelector<AppState, IToast[]>((state) =>
    selectAll(state.toasts)
  );

  return (
    <>
      {toasts.map((item) => (
        <Snackbar
          open
          key={item.id}
          anchorOrigin={{
            vertical: "bottom",
            horizontal: "left",
          }}
          autoHideDuration={
            item.severity === "error" ? null : AUTO_HIDE_DURATION
          }
          onClose={() => dispatch(removeToast({ id: item.id }))}
          message={item.message}
        >
          {item.severity && (
            <MuiAlert
              onClose={() => dispatch(removeToast({ id: item.id }))}
              severity={item.severity}
            >
              {item.message}
            </MuiAlert>
          )}
        </Snackbar>
      ))}
    </>
  );
};
