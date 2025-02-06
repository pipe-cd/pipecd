import {
  Button,
  Dialog,
  CircularProgress,
  DialogActions,
  DialogContent,
  DialogTitle,
  DialogProps,
  makeStyles,
} from "@material-ui/core";
import { FC } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";

export type DialogConfirmProps = DialogProps & {
  onCancel: () => void;
  onConfirm: () => void;
  title: string;
  description?: string;
  cancelText?: string;
  confirmText?: string;
  loading?: boolean;
};

const useStyles = makeStyles((theme) => ({
  progress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
}));

const DialogConfirm: FC<DialogConfirmProps> = ({
  onCancel,
  onClose,
  onConfirm,
  title,
  description,
  cancelText = UI_TEXT_CANCEL,
  confirmText = UI_TEXT_SAVE,
  loading = false,
  ...dialogProps
}) => {
  const classes = useStyles();

  return (
    <Dialog
      onClose={(event, reason) => {
        if (loading) return;
        onClose?.(event, reason);
      }}
      {...dialogProps}
    >
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>{description}</DialogContent>
      <DialogActions>
        <Button onClick={onCancel} disabled={loading}>
          {cancelText}
        </Button>
        <Button color="primary" onClick={onConfirm} disabled={loading}>
          {confirmText}
          {loading && (
            <CircularProgress size={24} className={classes.progress} />
          )}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DialogConfirm;
