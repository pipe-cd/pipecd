import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  TextField,
  Typography,
} from "@material-ui/core";
import EditIcon from "@material-ui/icons/Edit";
import Skeleton from "@material-ui/lab/Skeleton/Skeleton";
import * as React from "react";
import { FC, memo, useState } from "react";
import { RBAC_DESCRIPTION } from "~/constants/text";
import { UPDATE_RBAC_SUCCESS } from "~/constants/toast-text";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { fetchProject, Teams, updateRBAC } from "~/modules/project";
import { addToast } from "~/modules/toasts";
import { useProjectSettingStyles } from "~/styles/project-setting";
import { ProjectSettingLabeledText } from "../project-setting-labeled-text";

const SECTION_TITLE = "Role-Based Access Control";
const DIALOG_TITLE = `Edit ${SECTION_TITLE}`;
const TEAM_LABELS = {
  ADMIN: "Admin Team",
  EDITOR: "Editor Team",
  VIEWER: "Viewer Team",
};

export const RBACForm: FC = memo(function RBACForm() {
  const projectSettingClasses = useProjectSettingStyles();
  const teams = useAppSelector<Teams | null | undefined>(
    (state) => state.project.teams
  );
  const dispatch = useAppDispatch();
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
      if (updateRBAC.fulfilled.match(result)) {
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
    !!teams &&
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
          </>
        ) : (
          <div className={projectSettingClasses.values}>
            {teams === undefined ? (
              <>
                <Skeleton width={200} height={28} />
                <Skeleton width={200} height={28} />
                <Skeleton width={200} height={28} />
              </>
            ) : (
              <>
                <Typography variant="body1" color="textSecondary">
                  Not Configured
                </Typography>
              </>
            )}
          </div>
        )}
        <div>
          <IconButton onClick={() => setIsEdit(true)}>
            <EditIcon />
          </IconButton>
        </div>
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
            <Button onClick={handleClose}>{UI_TEXT_CANCEL}</Button>
            <Button type="submit" color="primary" disabled={isNotModified}>
              {UI_TEXT_SAVE}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </>
  );
});
