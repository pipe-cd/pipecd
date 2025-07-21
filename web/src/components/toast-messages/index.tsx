import { Button, Snackbar, SnackbarCloseReason } from "@mui/material";
import MuiAlert from "@mui/material/Alert";
import { FC, memo } from "react";
import { Link as RouterLink } from "react-router-dom";
import { IToast } from "~/contexts/toast-context";

const AUTO_HIDE_DURATION = 5000;

type Props = {
  toasts?: IToast[];
  onRemoveToast?: (id: string) => void;
};

export const ToastMessages: FC<Props> = memo(function Toasts({
  toasts,
  onRemoveToast,
}) {
  const handleRemoveToast = (
    id: string,
    reason?: SnackbarCloseReason
  ): void => {
    if (!onRemoveToast) return;
    if (reason === "clickaway") return;

    onRemoveToast(id);
  };

  return (
    <>
      {toasts?.map((item) => {
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
            onClose={(e, reason) => handleRemoveToast(item.id, reason)}
            message={item.message}
          >
            {item.severity && (
              <MuiAlert
                onClose={() => handleRemoveToast(item.id)}
                severity={item.severity}
                elevation={6}
                action={
                  item.to ? (
                    <Button
                      onClick={() => handleRemoveToast(item.id)}
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
