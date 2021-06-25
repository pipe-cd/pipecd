import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Link,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { Link as RouterLink } from "react-router-dom";
import { Alert } from "@material-ui/lab";
import { FC, memo, useCallback } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_DELETE } from "~/constants/ui-text";
import { red } from "@material-ui/core/colors";
import { shallowEqual } from "react-redux";
import { useEffect } from "react";
import {
  fetchApplicationsByEnv,
  selectApplicationsByEnvId,
} from "~/modules/applications";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { clearTargetEnv } from "~/modules/deleting-env";

const useStyles = makeStyles((theme) => ({
  targetName: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  description: {
    marginBottom: theme.spacing(2),
    whiteSpace: "break-spaces",
  },
  deleteButton: {
    color: theme.palette.getContrastText(red[400]),
    backgroundColor: red[800],
    "&:hover": {
      backgroundColor: red[800],
    },
  },
  appName: {
    marginRight: theme.spacing(1),
  },
}));

const DIALOG_TITLE = "Deleting Environment";
const DESCRIPTION =
  "Once deleted, environment cannot be restored.\nThis action will delete all applications that belong to this environments.";

interface Props {
  onDelete: (envId: string) => void;
}

export const DeleteEnvironmentDialog: FC<Props> = memo(
  function DeleteEnvironmentDialog({ onDelete }) {
    const classes = useStyles();
    const dispatch = useAppDispatch();
    const env = useAppSelector((state) => state.deletingEnv.env);
    const apps = useAppSelector(
      (state) => (env ? selectApplicationsByEnvId(env.id)(state) : []),
      shallowEqual
    );

    useEffect(() => {
      if (env) {
        dispatch(fetchApplicationsByEnv({ envId: env.id }));
      }
    }, [dispatch, env]);

    const handleClose = useCallback(() => {
      dispatch(clearTargetEnv());
    }, [dispatch]);

    const handleDelete = useCallback(() => {
      if (env) {
        onDelete(env.id);
      }
    }, [onDelete, env]);

    return (
      <Dialog open={env !== null} onClose={handleClose} fullWidth>
        <form>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <Alert severity="warning" className={classes.description}>
            {DESCRIPTION}
          </Alert>
          <DialogContent>
            {apps.length > 0 ? (
              <>
                <Typography variant="caption">
                  Applications to be deleted
                </Typography>
                <Typography variant="body1" className={classes.targetName}>
                  <Link
                    onClick={handleClose}
                    component={RouterLink}
                    to={`/applications?envId=${env?.id}`}
                  >
                    View applications
                  </Link>
                </Typography>
              </>
            ) : null}
            <Typography variant="caption">Environment to be deleted</Typography>
            <Typography variant="body1" className={classes.targetName}>
              {env ? env.name : ""}
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>{UI_TEXT_CANCEL}</Button>
            <Button
              variant="contained"
              onClick={handleDelete}
              className={classes.deleteButton}
            >
              {UI_TEXT_DELETE}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    );
  }
);
