import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Skeleton,
  TextField,
  Typography,
} from "@mui/material";
import EditIcon from "@mui/icons-material/Edit";

import * as React from "react";
import { FC, memo, useState } from "react";
import { SSO_DESCRIPTION } from "~/constants/text";
import { UPDATE_SSO_SUCCESS } from "~/constants/toast-text";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { fetchProject, GitHubSSO, updateGitHubSSO } from "~/modules/project";
import { addToast } from "~/modules/toasts";
import { useProjectSettingStyles } from "~/styles/project-setting";
import { ProjectSettingLabeledText } from "../project-setting-labeled-text";
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
  const projectSettingStyles = useProjectSettingStyles();
  const dispatch = useAppDispatch();
  const [isEdit, setIsEdit] = useState(false);
  const [clientId, setClientID] = useState("");
  const [clientSecret, setClientSecret] = useState("");
  const [baseUrl, setBaseUrl] = useState("");
  const [uploadUrl, setUploadUrl] = useState("");
  const [sso, sharedSSO] = useAppSelector<
    [GitHubSSO | null | undefined, string | null | undefined]
  >((state) => [state.project.github, state.project.sharedSSO]);

  const handleClose = (): void => {
    setIsEdit(false);
  };

  const handleSave = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    dispatch(
      updateGitHubSSO({ clientId, clientSecret, baseUrl, uploadUrl })
    ).then((result) => {
      if (updateGitHubSSO.fulfilled.match(result)) {
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
      <div className={projectSettingStyles.title}>
        <Typography variant="h5" className={projectSettingStyles.titleWithIcon}>
          {SECTION_TITLE}
        </Typography>
      </div>
      <Typography
        variant="body1"
        color="textSecondary"
        className={projectSettingStyles.description}
      >
        {SSO_DESCRIPTION}
      </Typography>
      <div className={projectSettingStyles.valuesWrapper}>
        {sso ? (
          <>
            <div className={projectSettingStyles.values}>
              <ProjectSettingLabeledText label="Client ID" value="********" />
              <ProjectSettingLabeledText
                label="Client Secret"
                value="********"
              />
              <ProjectSettingLabeledText label="Base URL" value={sso.baseUrl} />
              <ProjectSettingLabeledText
                label="Upload URL"
                value={sso.uploadUrl}
              />
            </div>

            <div>
              <IconButton onClick={() => setIsEdit(true)} size="large">
                <EditIcon />
              </IconButton>
            </div>
          </>
        ) : (
          <div className={projectSettingStyles.values}>
            {sso === undefined ? (
              <>
                <Skeleton width={200} height={28} />
                <Skeleton width={200} height={28} />
                <Skeleton width={200} height={28} />
                <Skeleton width={200} height={28} />
              </>
            ) : sharedSSO ? (
              <ProjectSettingLabeledText label="Shared SSO" value={sharedSSO} />
            ) : (
              <Typography variant="body1" color="textSecondary">
                Not Configured
              </Typography>
            )}
          </div>
        )}
      </div>
      <Dialog
        open={isEdit}
        onClose={handleClose}
        TransitionProps={{
          onEnter: () => {
            setBaseUrl(sso?.baseUrl ?? "");
            setUploadUrl(sso?.uploadUrl ?? "");
          },
        }}
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
            <Button onClick={handleClose}>{UI_TEXT_CANCEL}</Button>
            <Button type="submit" color="primary" disabled={isInvalid}>
              {UI_TEXT_SAVE}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </>
  );
});
