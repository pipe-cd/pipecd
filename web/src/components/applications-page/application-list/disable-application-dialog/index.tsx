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
import { useAppSelector, useAppDispatch } from "~/hooks/redux";
import {
  Application,
  disableApplication,
  selectById,
} from "~/modules/applications";

const useStyles = makeStyles((theme) => ({
  disableTargetName: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  description: {
    marginBottom: theme.spacing(2),
  },
}));

export interface DisableApplicationDialogProps {
  open: boolean;
  applicationId: string | null;
  onCancel: () => void;
  onDisable: () => void;
}

export const DisableApplicationDialog: FC<DisableApplicationDialogProps> = memo(
  function DisableApplicationDialog({
    applicationId,
    open,
    onDisable,
    onCancel,
  }) {
    const classes = useStyles();
    const dispatch = useAppDispatch();

    const application = useAppSelector<Application.AsObject | undefined>(
      (state) =>
        applicationId
          ? selectById(state.applications, applicationId)
          : undefined
    );

    const handleDisable = (): void => {
      if (applicationId) {
        dispatch(disableApplication({ applicationId })).then(() => {
          onDisable();
        });
      }
    };

    if (!application) {
      return null;
    }

    return (
      <Dialog open={Boolean(application) && open}>
        <DialogTitle>Disable application</DialogTitle>
        <DialogContent>
          <Alert severity="warning" className={classes.description}>
            Are you sure you want to disable the application?
          </Alert>
          <Typography variant="caption">NAME</Typography>
          <Typography variant="body1" className={classes.disableTargetName}>
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
  }
);
