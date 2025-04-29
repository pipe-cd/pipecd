import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import Alert from "@mui/material/Alert";
import { FC, memo } from "react";

const useStyles = makeStyles((theme) => ({
  deleteTargetName: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  description: {
    marginBottom: theme.spacing(2),
  },
}));

export interface DeleteUserGroupConfirmDialogProps {
  ssoGroup: string | null;
  onCancel: () => void;
  onDelete: (ssoGroup: string) => void;
}

const DIALOG_TITLE = "Delete User Group";
const DESCRIPTION = "Are you sure you want to delete the User Group?";

export const DeleteUserGroupConfirmDialog: FC<DeleteUserGroupConfirmDialogProps> = memo(
  function DeleteUserGroupConfirmDialog({ ssoGroup, onDelete, onCancel }) {
    const classes = useStyles();

    return (
      <Dialog open={Boolean(ssoGroup)} onClose={onCancel}>
        <DialogTitle>{DIALOG_TITLE}</DialogTitle>
        <DialogContent>
          <Alert severity="warning" className={classes.description}>
            {DESCRIPTION}
          </Alert>
          <Typography variant="caption">Group</Typography>
          <Typography variant="body1" className={classes.deleteTargetName}>
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
