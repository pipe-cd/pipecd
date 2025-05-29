import { Button, Snackbar, SnackbarCloseReason } from "@mui/material";
import MuiAlert from "@mui/material/Alert";
import { FC, memo, SyntheticEvent } from "react";
import { Link as RouterLink } from "react-router-dom";
import { useAppSelector, useAppDispatch } from "~/hooks/redux";
import { IToast, removeToast, selectAll } from "~/modules/toasts";

const AUTO_HIDE_DURATION = 5000;

export const Toasts: FC = memo(function Toasts() {
  const dispatch = useAppDispatch();
  const toasts = useAppSelector<IToast[]>((state) => selectAll(state.toasts));

  return (
    <>
      {toasts.map((item) => {
        const handleClose = (
          _: Event | SyntheticEvent<unknown, Event>,
          reason?: SnackbarCloseReason
        ): void => {
          if (reason === "clickaway") {
            return;
          }
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
});
