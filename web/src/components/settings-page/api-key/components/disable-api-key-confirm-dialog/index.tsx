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
import { useAppSelector } from "~/hooks/redux";
import { APIKey, selectById } from "~/modules/api-keys";

export interface DisableAPIKeyConfirmDialogProps {
  apiKeyId: string | null;
  onCancel: () => void;
  onDisable: (id: string) => void;
}

const DIALOG_TITLE = "Disable API Key";
const DESCRIPTION = "Are you sure you want to disable the API key?";

export const DisableAPIKeyConfirmDialog: FC<DisableAPIKeyConfirmDialogProps> = memo(
  function DisableAPIKeyConfirmDialog({ apiKeyId, onDisable, onCancel }) {
    const apiKey = useAppSelector<APIKey.AsObject | undefined>((state) =>
      apiKeyId ? selectById(state.apiKeys, apiKeyId) : undefined
    );
    const open = Boolean(apiKey);

    return (
      <Dialog open={open} onClose={onCancel}>
        <DialogTitle>{DIALOG_TITLE}</DialogTitle>
        <DialogContent>
          <Alert severity="warning" sx={{ marginBottom: 2 }}>
            {DESCRIPTION}
          </Alert>
          <Typography variant="caption">NAME</Typography>
          <Typography
            variant="body1"
            sx={(theme) => ({
              color: theme.palette.text.primary,
              fontWeight: theme.typography.fontWeightMedium,
            })}
          >
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
