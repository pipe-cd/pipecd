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

export interface DeleteUserGroupConfirmDialogProps {
  ssoGroup: string | null;
  onCancel: () => void;
  onDelete: (ssoGroup: string) => void;
}

const DIALOG_TITLE = "Delete User Group";
const DESCRIPTION = "Are you sure you want to delete the User Group?";

export const DeleteUserGroupConfirmDialog: FC<DeleteUserGroupConfirmDialogProps> = memo(
  function DeleteUserGroupConfirmDialog({ ssoGroup, onDelete, onCancel }) {
    return (
      <Dialog open={Boolean(ssoGroup)} onClose={onCancel}>
        <DialogTitle>{DIALOG_TITLE}</DialogTitle>
        <DialogContent>
          <Alert severity="warning" sx={{ marginBottom: 2 }}>
            {DESCRIPTION}
          </Alert>
          <Typography variant="caption">Group</Typography>
          <Typography
            variant="body1"
            sx={(theme) => ({
              color: theme.palette.text.primary,
              fontWeight: theme.typography.fontWeightMedium,
            })}
          >
            {ssoGroup}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={onCancel}>Cancel</Button>
          <Button
            color="primary"
            onClick={() => {
              if (ssoGroup) {
                onDelete(ssoGroup);
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
