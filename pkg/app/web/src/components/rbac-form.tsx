import {
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  TextField,
  Typography,
} from "@material-ui/core";
import EditIcon from "@material-ui/icons/Edit";
import React, { FC, memo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { BUTTON_TEXT_CANCEL, BUTTON_TEXT_SAVE } from "../constants/button-text";
import { RBAC_DESCRIPTION } from "../constants/text";
import {
  UPDATE_RBAC_FAILED,
  UPDATE_RBAC_SUCCESS,
} from "../constants/toast-text";
import { AppState } from "../modules";
import { fetchProject, Teams, updateRBAC } from "../modules/project";
import { addToast } from "../modules/toasts";
import { AppDispatch } from "../store";
import { useProjectSettingStyles } from "../styles/project-setting";
import { ProjectSettingLabeledText } from "./project-setting-labeled-text";

const SECTION_TITLE = "Role-Based Access Control";
const DIALOG_TITLE = `Edit ${SECTION_TITLE}`;
const TEAM_LABELS = {
  ADMIN: "Admin Team",
  EDITOR: "Editor Team",
  VIEWER: "Viewer Team",
};

export const RBACForm: FC = memo(function RBACForm() {
  const projectSettingClasses = useProjectSettingStyles();
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
      <div className={projectSettingClasses.title}>
        <Typography
          variant="h5"
          className={projectSettingClasses.titleWithIcon}
        >
          {SECTION_TITLE}
        </Typography>
      </div>

      <Typography
        variant="body1"
        color="textSecondary"
        className={projectSettingClasses.description}
      >
        {RBAC_DESCRIPTION}
      </Typography>

      <div className={projectSettingClasses.valuesWrapper}>
        {teams ? (
          <>
            <div className={projectSettingClasses.values}>
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
