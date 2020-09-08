import {
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import EditIcon from "@material-ui/icons/Edit";
import React, { FC, memo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { BUTTON_TEXT_CANCEL, BUTTON_TEXT_SAVE } from "../constants/button-text";
import { SSO_DESCRIPTION } from "../constants/text";
import { UPDATE_SSO_FAILED, UPDATE_SSO_SUCCESS } from "../constants/toast-text";
import { AppState } from "../modules";
import { fetchProject, GitHubSSO, updateGitHubSSO } from "../modules/project";
import { addToast } from "../modules/toasts";
import { AppDispatch } from "../store";
import { ProjectSettingLabeledText } from "./project-setting-labeled-text";

const useStyles = makeStyles((theme) => ({
  title: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  titleWithIcon: {
    display: "flex",
    alignItems: "center",
  },
  values: {
    padding: theme.spacing(3),
  },
  indent: {
    padding: theme.spacing(1),
  },
  name: {
    color: theme.palette.text.secondary,
    marginRight: theme.spacing(2),
    minWidth: 120,
  },
  item: {
    display: "flex",
    alignItems: "center",
  },
}));

export interface GitHubSSOFormParams {
  clientId: string;
  clientSecret: string;
  baseUrl: string;
  uploadUrl: string;
  org: string;
  adminTeam: string;
  editorTeam: string;
  viewerTeam: string;
}

const SECTION_TITLE = "Single Sign-On";
const DIALOG_TITLE = `Edit ${SECTION_TITLE}`;

export const GithubSSOForm: FC = memo(function GithubSSOForm() {
  const classes = useStyles();
  const dispatch = useDispatch<AppDispatch>();
  const [isEdit, setIsEdit] = useState(false);
  const [clientId, setClientID] = useState("");
  const [clientSecret, setClientSecret] = useState("");
  const [baseUrl, setBaseUrl] = useState("");
  const [uploadUrl, setUploadUrl] = useState("");
  const sso = useSelector<AppState, GitHubSSO | null>(
    (state) => state.project.github
  );

  const handleClose = (): void => {
    setIsEdit(false);
  };

  const handleSave = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    dispatch(
      updateGitHubSSO({ clientId, clientSecret, baseUrl, uploadUrl })
    ).then((result) => {
      if (updateGitHubSSO.rejected.match(result)) {
        dispatch(
          addToast({
            message: UPDATE_SSO_FAILED,
            severity: "error",
          })
        );
      } else {
        dispatch(fetchProject());
        dispatch(
          addToast({
            message: UPDATE_SSO_SUCCESS,
            severity: "success",
          })
        );
      }
    });
    setIsEdit(false);
  };

  const isInvalid = clientId === "" || clientSecret === "";

  return (
    <>
      <div className={classes.title}>
        <Typography variant="h5" className={classes.titleWithIcon}>
          {SECTION_TITLE}
          <IconButton onClick={() => setIsEdit(true)}>
            <EditIcon />
          </IconButton>
        </Typography>
      </div>

      <Typography variant="body1" color="textSecondary">
        {SSO_DESCRIPTION}
      </Typography>

      <div className={classes.values}>
        {sso ? (
          <>
            <ProjectSettingLabeledText label="Client ID" value="********" />
            <ProjectSettingLabeledText label="Client Secret" value="********" />
            <ProjectSettingLabeledText label="Base URL" value={sso.baseUrl} />
            <ProjectSettingLabeledText
              label="Upload URL"
              value={sso.uploadUrl}
            />
          </>
        ) : (
          <CircularProgress />
        )}
      </div>

      <Dialog
        open={isEdit}
        onEnter={() => {
          setBaseUrl(sso?.baseUrl ?? "");
          setUploadUrl(sso?.uploadUrl ?? "");
        }}
        onClose={handleClose}
      >
        <form onSubmit={handleSave}>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <TextField
              value={clientId}
              variant="outlined"
              margin="dense"
              label="Client ID"
              fullWidth
              required
              autoFocus
              onChange={(e) => setClientID(e.currentTarget.value)}
            />
            <TextField
              value={clientSecret}
              variant="outlined"
              margin="dense"
              label="Client Secret"
              fullWidth
              required
              onChange={(e) => setClientSecret(e.currentTarget.value)}
            />
            <TextField
              value={baseUrl}
              variant="outlined"
              margin="dense"
              label="Base URL"
              fullWidth
              onChange={(e) => setBaseUrl(e.currentTarget.value)}
            />
            <TextField
              value={uploadUrl}
              variant="outlined"
              margin="dense"
              label="Upload URL"
              fullWidth
              onChange={(e) => setUploadUrl(e.currentTarget.value)}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>{BUTTON_TEXT_CANCEL}</Button>
            <Button type="submit" color="primary" disabled={isInvalid}>
              {BUTTON_TEXT_SAVE}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </>
  );
});
