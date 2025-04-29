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
import { useAppSelector } from "~/hooks/redux";
import { APIKey, selectById } from "~/modules/api-keys";

const useStyles = makeStyles((theme) => ({
  disableTargetName: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  description: {
    marginBottom: theme.spacing(2),
  },
}));

export interface DisableAPIKeyConfirmDialogProps {
  apiKeyId: string | null;
  onCancel: () => void;
  onDisable: (id: string) => void;
}

const DIALOG_TITLE = "Disable API Key";
const DESCRIPTION = "Are you sure you want to disable the API key?";

export const DisableAPIKeyConfirmDialog: FC<DisableAPIKeyConfirmDialogProps> = memo(
  function DisableAPIKeyConfirmDialog({ apiKeyId, onDisable, onCancel }) {
    const classes = useStyles();
    const apiKey = useAppSelector<APIKey.AsObject | undefined>((state) =>
      apiKeyId ? selectById(state.apiKeys, apiKeyId) : undefined
    );
    const open = Boolean(apiKey);

    return (
      <Dialog open={open} onClose={onCancel}>
        <DialogTitle>{DIALOG_TITLE}</DialogTitle>
        <DialogContent>
          <Alert severity="warning" className={classes.description}>
            {DESCRIPTION}
          </Alert>
          <Typography variant="caption">NAME</Typography>
          <Typography variant="body1" className={classes.disableTargetName}>
            {apiKey?.name}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={onCancel}>Cancel</Button>
          <Button
            color="primary"
            onClick={() => {
              if (apiKeyId) {
                onDisable(apiKeyId);
              }
            }}
          >
            Disable
          </Button>
        </DialogActions>
      </Dialog>
    );
  }
);
