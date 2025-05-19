import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";
import Alert from "@mui/material/Alert";
import { FC, memo } from "react";

export interface DeleteRoleConfirmDialogProps {
  role: string | null;
  onClose: () => void;
  onDelete: (role: string) => void;
}

const DIALOG_TITLE = "Delete Role";
const DESCRIPTION = "Are you sure you want to delete the Role?";

export const DeleteRoleConfirmDialog: FC<DeleteRoleConfirmDialogProps> = memo(
  function DeleteRoleConfirmDialog({ role, onDelete, onClose }) {
    return (
      <Dialog open={Boolean(role)} onClose={onClose}>
        <DialogTitle>{DIALOG_TITLE}</DialogTitle>
        <DialogContent>
          <Alert severity="warning" sx={{ marginBottom: 2 }}>
            {DESCRIPTION}
          </Alert>
          <Typography variant="caption">Role</Typography>
          <Typography
            variant="body1"
            sx={(theme) => ({
              color: theme.palette.text.primary,
              fontWeight: theme.typography.fontWeightMedium,
            })}
          >
            {role}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button
            color="primary"
            onClick={() => {
              if (role) {
                onDelete(role);
              }
            }}
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    );
  }
);
