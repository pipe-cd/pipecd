import React, { FC, memo } from "react";
import {
  makeStyles,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
} from "@material-ui/core";
import { useSelector, useDispatch } from "react-redux";
import {
  selectById,
  Application,
  disableApplication,
} from "../modules/applications";
import { AppState } from "../modules";
import Alert from "@material-ui/lab/Alert";
import { AppDispatch } from "../store";

const useStyles = makeStyles((theme) => ({
  disableButton: {
    color: theme.palette.error.contrastText,
    backgroundColor: theme.palette.error.main,
    "&:hover": {
      backgroundColor: theme.palette.error.dark,
    },
  },
  description: {
    marginBottom: theme.spacing(2),
  },
}));

interface Props {
  applicationId: string | null;
  onCancel: () => void;
  onDisable: () => void;
}

export const DisableApplicationDialog: FC<Props> = memo(
  function DisableApplicationDialog({ applicationId, onDisable, onCancel }) {
    const classes = useStyles();
    const dispatch = useDispatch<AppDispatch>();

    const application = useSelector<AppState, Application | undefined>(
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
      <Dialog open={Boolean(application)}>
        <DialogTitle>Disable application</DialogTitle>
        <DialogContent>
          <Alert severity="warning" className={classes.description}>
            Are you sure you want to disable the application?
          </Alert>
          <div>
            <Typography variant="caption">NAME</Typography>
            <Typography variant="body1">{application.name}</Typography>
          </div>
        </DialogContent>
        <DialogActions>
          <Button onClick={onCancel}>Cancel</Button>
          <Button className={classes.disableButton} onClick={handleDisable}>
            Disable
          </Button>
        </DialogActions>
      </Dialog>
    );
  }
);
