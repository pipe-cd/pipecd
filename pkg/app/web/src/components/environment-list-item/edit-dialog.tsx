import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from "@material-ui/core";
import { FC, useState } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "../../constants/ui-text";

const DIALOG_TITLE = "Edit Environment description";

interface Props {
  open: boolean;
  description: string;
  onClose: () => void;
  onSave?: () => void;
}

export const EditEnvironmentDialog: FC<Props> = ({
  open,
  description,
  onClose,
  onSave = () => null,
}) => {
  const [desc, setDesc] = useState(description);

  return (
    <Dialog open={open} onClose={onClose} fullWidth>
      <form onSubmit={onSave}>
        <DialogTitle>{DIALOG_TITLE}</DialogTitle>
        <DialogContent>
          <TextField
            value={desc}
            variant="outlined"
            margin="dense"
            label="Description"
            fullWidth
            required
            autoFocus
            onChange={(e) => setDesc(e.currentTarget.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>{UI_TEXT_CANCEL}</Button>
          <Button
            type="submit"
            color="primary"
            disabled={desc === "" || desc === description}
          >
            {UI_TEXT_SAVE}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};
