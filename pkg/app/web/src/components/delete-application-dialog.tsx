import React, { FC, memo, useCallback } from "react";
import {
  makeStyles,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  CircularProgress,
} from "@material-ui/core";
import { useSelector, useDispatch, shallowEqual } from "react-redux";
import {
  selectById,
  Application,
  fetchApplications,
} from "../modules/applications";
import {
  clearDeletingApp,
  deleteApplication,
} from "../modules/delete-application";
import { AppState } from "../modules";
import Alert from "@material-ui/lab/Alert";
import { AppDispatch } from "../store";
import { red } from "@material-ui/core/colors";
import { UI_TEXT_CANCEL, UI_TEXT_DELETE } from "../constants/ui-text";
import { useStyles as useButtonStyles } from "../styles/button";
import { Skeleton } from "@material-ui/lab";
import { addToast } from "../modules/toasts";
import { DELETE_APPLICATION_SUCCESS } from "../constants/toast-text";

const useStyles = makeStyles((theme) => ({
  applicationName: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  description: {
    marginBottom: theme.spacing(2),
  },
  deleteButton: {
    color: theme.palette.getContrastText(red[400]),
    backgroundColor: red[700],
    "&:hover": {
      backgroundColor: red[700],
    },
  },
}));

const TITLE = "Delete Application";
const ALERT_TEXT = "Are you sure you want to delete the application?";

export const DeleteApplicationDialog: FC = memo(
  function DeleteApplicationDialog() {
    const classes = useStyles();
    const buttonClasses = useButtonStyles();
    const dispatch = useDispatch<AppDispatch>();

    const [application, isDeleting] = useSelector<
      AppState,
      [Application | undefined, boolean]
    >(
      (state) => [
        state.deleteApplication.applicationId
          ? selectById(
              state.applications,
              state.deleteApplication.applicationId
            )
          : undefined,
        state.deleteApplication.deleting,
      ],
      shallowEqual
    );

    const handleDelete = useCallback(() => {
      dispatch(deleteApplication()).then(() => {
        dispatch(fetchApplications());
        dispatch(
          addToast({ severity: "success", message: DELETE_APPLICATION_SUCCESS })
        );
      });
    }, [dispatch]);

    const handleCancel = useCallback(() => {
      dispatch(clearDeletingApp());
    }, [dispatch]);

    return (
      <Dialog open={Boolean(application)} disableBackdropClick={isDeleting}>
        <DialogTitle>{TITLE}</DialogTitle>
        <DialogContent>
          <Alert severity="error" className={classes.description}>
            {ALERT_TEXT}
          </Alert>
          <Typography variant="caption">Name</Typography>
          <Typography variant="body1" className={classes.applicationName}>
            {application ? (
              application.name
            ) : (
              <Skeleton height={24} width={200} />
            )}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCancel} disabled={isDeleting}>
            {UI_TEXT_CANCEL}
          </Button>
          <Button
            variant="contained"
            color="primary"
            onClick={handleDelete}
            className={classes.deleteButton}
            disabled={isDeleting}
          >
            {UI_TEXT_DELETE}
            {isDeleting && (
              <CircularProgress size={24} className={buttonClasses.progress} />
            )}
          </Button>
        </DialogActions>
      </Dialog>
    );
  }
);
