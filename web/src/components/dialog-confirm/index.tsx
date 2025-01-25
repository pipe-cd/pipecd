import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  DialogProps,
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
};

const DialogConfirm: FC<DialogConfirmProps> = ({
  onCancel,
  onConfirm,
  title,
  description,
  cancelText = UI_TEXT_CANCEL,
  confirmText = UI_TEXT_SAVE,
  ...dialogProps
}) => {
  return (
    <Dialog {...dialogProps}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>{description}</DialogContent>
      <DialogActions>
        <Button onClick={onCancel}>{cancelText}</Button>
        <Button color="primary" onClick={onConfirm}>
          {confirmText}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DialogConfirm;
