import {
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  makeStyles,
  Switch,
  TextField,
  Typography,
} from "@material-ui/core";
import EditIcon from "@material-ui/icons/Edit";
import clsx from "clsx";
import React, { FC, memo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { BUTTON_TEXT_CANCEL, BUTTON_TEXT_SAVE } from "../constants/button-text";
import { STATIC_ADMIN_DESCRIPTION } from "../constants/text";
import {
  UPDATE_STATIC_ADMIN_INFO_FAILED,
  UPDATE_STATIC_ADMIN_INFO_SUCCESS,
} from "../constants/toast-text";
import { AppState } from "../modules";
import {
  fetchProject,
  toggleAvailability,
  updateStaticAdmin,
} from "../modules/project";
import { addToast } from "../modules/toasts";
import { AppDispatch } from "../store";
import { ProjectSettingLabeledText } from "./project-setting-labeled-text";

const useStyles = makeStyles((theme) => ({
  title: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  description: {
    paddingRight: theme.spacing(6),
  },
  titleWithIcon: {
    display: "flex",
    alignItems: "center",
  },
  valuesWrapper: {
    padding: theme.spacing(1),
    display: "flex",
    justifyContent: "space-between",
  },
  values: {
    padding: theme.spacing(2),
  },
  disabled: {
    opacity: 0.5,
  },
}));

const SECTION_TITLE = "Static Admin";
const DIALOG_TITLE = `Edit ${SECTION_TITLE}`;

export const StaticAdminForm: FC = memo(function StaticAdminForm() {
  const classes = useStyles();
  const dispatch = useDispatch<AppDispatch>();
  const [isEnabled, currentUsername] = useSelector<
    AppState,
    [boolean, string | null]
  >((state) => [
    state.project.staticAdminDisabled === false,
    state.project.username,
  ]);
  const [isEdit, setIsEdit] = useState(false);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleClose = (): void => {
    setIsEdit(false);
  };

  const handleToggleAvailability = (): void => {
    dispatch(toggleAvailability()).then(() => {
      dispatch(fetchProject());
    });
  };

  const handleSave = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    dispatch(updateStaticAdmin({ username, password })).then((result) => {
      if (updateStaticAdmin.rejected.match(result)) {
        dispatch(
          addToast({
            message: UPDATE_STATIC_ADMIN_INFO_FAILED,
            severity: "error",
          })
        );
      } else {
        dispatch(fetchProject());
        dispatch(
          addToast({
            message: UPDATE_STATIC_ADMIN_INFO_SUCCESS,
            severity: "success",
          })
        );
      }
    });
    setIsEdit(false);
  };

  const isInvalidValues = username === "" || password === "";

  return (
    <>
      <div className={classes.title}>
        <Typography variant="h5" className={classes.titleWithIcon}>
          {SECTION_TITLE}
        </Typography>

        <Switch
          checked={isEnabled}
          color="primary"
          onClick={handleToggleAvailability}
          disabled={currentUsername === null}
        />
      </div>

      <Typography
        variant="body1"
        color="textSecondary"
        className={classes.description}
      >
        {STATIC_ADMIN_DESCRIPTION}
      </Typography>

      <div
        className={clsx(classes.valuesWrapper, {
          [classes.disabled]: isEnabled === false,
        })}
      >
        {currentUsername === null ? (
          <CircularProgress />
        ) : (
          <>
            <div className={classes.values}>
              <ProjectSettingLabeledText
                label="Username"
                value={currentUsername}
              />
              <ProjectSettingLabeledText label="Password" value="********" />
            </div>
            <div>
              <IconButton
                onClick={() => setIsEdit(true)}
                disabled={isEnabled === false}
              >
                <EditIcon />
              </IconButton>
            </div>
          </>
        )}
      </div>

      <Dialog
        open={isEdit}
        onEnter={() => {
          setUsername(currentUsername || "");
          setPassword("");
        }}
        onClose={handleClose}
      >
        <form onSubmit={handleSave}>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <TextField
              value={username}
              variant="outlined"
              margin="dense"
              label="Username"
              fullWidth
              required
              autoFocus
              onChange={(e) => setUsername(e.currentTarget.value)}
            />
            <TextField
              value={password}
              autoComplete="new-password"
              variant="outlined"
              margin="dense"
              label="Password"
              type="password"
              fullWidth
              required
              onChange={(e) => setPassword(e.currentTarget.value)}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>{BUTTON_TEXT_CANCEL}</Button>
            <Button type="submit" color="primary" disabled={isInvalidValues}>
              {BUTTON_TEXT_SAVE}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </>
  );
});
