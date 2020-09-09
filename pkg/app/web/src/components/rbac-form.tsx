import React, { FC, useState, memo } from "react";
import EditIcon from "@material-ui/icons/Edit";
import {
  IconButton,
  makeStyles,
  Typography,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Button,
  TextField,
  CircularProgress,
} from "@material-ui/core";
import { BUTTON_TEXT_CANCEL, BUTTON_TEXT_SAVE } from "../constants/button-text";
import { RBAC_DESCRIPTION } from "../constants/text";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import { fetchProject, Teams, updateRBAC } from "../modules/project";
import { ProjectSettingLabeledText } from "./project-setting-labeled-text";
import { AppDispatch } from "../store";
import { addToast } from "../modules/toasts";
import {
  UPDATE_RBAC_FAILED,
  UPDATE_RBAC_SUCCESS,
} from "../constants/toast-text";

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
}));

const SECTION_TITLE = "Role-Based Access Control";
const DIALOG_TITLE = `Edit ${SECTION_TITLE}`;
const TEAM_LABELS = {
  ADMIN: "Admin Team",
  EDITOR: "Editor Team",
  VIEWER: "Viewer Team",
};

export const RBACForm: FC = memo(function RBACForm() {
  const classes = useStyles();
  const teams = useSelector<AppState, Teams | null>(
    (state) => state.project.teams
  );
  const dispatch = useDispatch<AppDispatch>();
  const [isEdit, setIsEdit] = useState(false);
  const [admin, setAdmin] = useState("");
  const [editor, setEditor] = useState("");
  const [viewer, setViewer] = useState("");

  const handleClose = (): void => {
    setIsEdit(false);
  };
  const handleSave = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    dispatch(updateRBAC({ admin, editor, viewer })).then((result) => {
      if (updateRBAC.rejected.match(result)) {
        dispatch(
          addToast({
            message: UPDATE_RBAC_FAILED,
            severity: "success",
          })
        );
      } else {
        dispatch(fetchProject());
        dispatch(
          addToast({
            message: UPDATE_RBAC_SUCCESS,
            severity: "success",
          })
        );
      }
    });
    setIsEdit(false);
  };

  const isNotModified =
    teams !== null &&
    teams.admin === admin &&
    teams.editor === editor &&
    teams.viewer === viewer;

  return (
    <>
      <div className={classes.title}>
        <Typography variant="h5" className={classes.titleWithIcon}>
          {SECTION_TITLE}
        </Typography>
      </div>

      <Typography
        variant="body1"
        color="textSecondary"
        className={classes.description}
      >
        {RBAC_DESCRIPTION}
      </Typography>

      <div className={classes.valuesWrapper}>
        {teams ? (
          <>
            <div className={classes.values}>
              <ProjectSettingLabeledText
                label={TEAM_LABELS.ADMIN}
                value={teams.admin}
              />
              <ProjectSettingLabeledText
                label={TEAM_LABELS.EDITOR}
                value={teams.editor}
              />
              <ProjectSettingLabeledText
                label={TEAM_LABELS.VIEWER}
                value={teams.viewer}
              />
            </div>
            <div>
              <IconButton onClick={() => setIsEdit(true)}>
                <EditIcon />
              </IconButton>
            </div>
          </>
        ) : (
          <CircularProgress />
        )}
      </div>

      <Dialog
        open={isEdit}
        onEnter={() => {
          setAdmin(teams?.admin ?? "");
          setEditor(teams?.editor ?? "");
          setViewer(teams?.viewer ?? "");
        }}
        onClose={handleClose}
      >
        <form onSubmit={handleSave}>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <TextField
              value={admin}
              variant="outlined"
              margin="dense"
              label={TEAM_LABELS.ADMIN}
              fullWidth
              autoFocus
              onChange={(e) => setAdmin(e.currentTarget.value)}
            />
            <TextField
              value={editor}
              variant="outlined"
              margin="dense"
              label={TEAM_LABELS.EDITOR}
              fullWidth
              onChange={(e) => setEditor(e.currentTarget.value)}
            />
            <TextField
              value={viewer}
              variant="outlined"
              margin="dense"
              label={TEAM_LABELS.VIEWER}
              fullWidth
              onChange={(e) => setViewer(e.currentTarget.value)}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>{BUTTON_TEXT_CANCEL}</Button>
            <Button type="submit" color="primary" disabled={isNotModified}>
              {BUTTON_TEXT_SAVE}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </>
  );
});
