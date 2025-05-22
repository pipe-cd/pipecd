import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  DialogProps,
} from "@mui/material";
import { FC } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { SpinnerIcon } from "~/styles/button";

export type DialogConfirmProps = DialogProps & {
  onCancel: () => void;
  onConfirm: () => void;
  title: string;
  description?: string;
  cancelText?: string;
  confirmText?: string;
  loading?: boolean;
};

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
          {loading && <SpinnerIcon />}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DialogConfirm;
