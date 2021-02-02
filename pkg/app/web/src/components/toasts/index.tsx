import { Button, Snackbar } from "@material-ui/core";
import MuiAlert from "@material-ui/lab/Alert";
import React, { FC } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import { AppState } from "../../modules";
import { IToast, removeToast, selectAll } from "../../modules/toasts";

const AUTO_HIDE_DURATION = 5000;

export const Toasts: FC = () => {
  const dispatch = useDispatch();
  const toasts = useSelector<AppState, IToast[]>((state) =>
    selectAll(state.toasts)
  );

  return (
    <>
      {toasts.map((item) => {
        const handleClose = (): void => {
          dispatch(removeToast({ id: item.id }));
        };
        return (
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
            onClose={handleClose}
            message={item.message}
          >
            {item.severity && (
              <MuiAlert
                onClose={handleClose}
                severity={item.severity}
                elevation={6}
                action={
                  item.to ? (
                    <Button
                      onClick={handleClose}
                      component={RouterLink}
                      to={item.to}
                    >
                      OPEN
                    </Button>
                  ) : null
                }
              >
                {item.message}
              </MuiAlert>
            )}
          </Snackbar>
        );
      })}
    </>
  );
};
