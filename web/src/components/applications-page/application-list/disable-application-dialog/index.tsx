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
import { useDisableApplication } from "~/queries/applications/use-disable-application";
import { Application } from "~/types/applications";

export interface DisableApplicationDialogProps {
  open: boolean;
  application?: Application.AsObject | null;
  onCancel: () => void;
  onDisable: () => void;
}

export const DisableApplicationDialog: FC<DisableApplicationDialogProps> = memo(
  function DisableApplicationDialog({
    application,
    open,
    onDisable,
    onCancel,
  }) {
    const { mutate: disableApplication } = useDisableApplication();

    const handleDisable = (): void => {
      if (application) {
        disableApplication(
          { applicationId: application.id },
          { onSuccess: () => onDisable() }
        );
      }
    };

    if (!application) {
      return null;
    }

    return (
      <Dialog open={Boolean(application) && open}>
        <DialogTitle>Disable application</DialogTitle>
        <DialogContent>
          <Alert
            severity="warning"
            sx={{
              marginBottom: 2,
            }}
          >
            Are you sure you want to disable the application?
          </Alert>
          <Typography variant="caption">NAME</Typography>
          <Typography
            variant="body1"
            sx={(theme) => ({
              color: theme.palette.text.primary,
              fontWeight: theme.typography.fontWeightMedium,
            })}
          >
            {application.name}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={onCancel}>Cancel</Button>
          <Button color="primary" onClick={handleDisable}>
            Disable
          </Button>
        </DialogActions>
      </Dialog>
    );
  },
  (prev, next) =>
    prev.open === next.open &&
    prev.application?.id === next.application?.id &&
    prev.application?.name === next.application?.name
);
