import {
  Box,
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
import {
  ProjectDescription,
  ProjectTitleWrap,
  ProjectTitle,
  ProjectValues,
  ProjectValuesWrapper,
} from "~/styles/project-setting";
import { ProjectSettingLabeledText } from "../project-setting-labeled-text";
import { useUpdateGithubSso } from "~/queries/project/use-update-github-sso";
import { useToast } from "~/contexts/toast-context";
import { useGetProject } from "~/queries/project/use-get-project";
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
  const [isEdit, setIsEdit] = useState(false);
  const [clientId, setClientID] = useState("");
  const [clientSecret, setClientSecret] = useState("");
  const [baseUrl, setBaseUrl] = useState("");
  const [uploadUrl, setUploadUrl] = useState("");

  const { data: projectDetail } = useGetProject();
  const { github: sso, sharedSSO: sharedSSO } = projectDetail || {};

  const { mutateAsync: updateGithubSso } = useUpdateGithubSso();
  const { addToast } = useToast();
  const handleClose = (): void => {
    setIsEdit(false);
  };

  const handleSave = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    updateGithubSso({ clientId, clientSecret, baseUrl, uploadUrl }).then(() => {
      addToast({
        message: UPDATE_SSO_SUCCESS,
        severity: "success",
      });
    });
    setIsEdit(false);
  };

  const isInvalid = clientId === "" || clientSecret === "";

  return (
    <>
      <ProjectTitleWrap>
        <ProjectTitle variant="h5">{SECTION_TITLE}</ProjectTitle>
      </ProjectTitleWrap>
      <ProjectDescription variant="body1" color="textSecondary">
        {SSO_DESCRIPTION}
      </ProjectDescription>
      <ProjectValuesWrapper>
        {sso ? (
          <>
            <ProjectValues>
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
            </ProjectValues>

            <div>
              <IconButton onClick={() => setIsEdit(true)} size="large">
                <EditIcon />
              </IconButton>
            </div>
          </>
        ) : (
          <ProjectValues>
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
          </ProjectValues>
        )}
      </ProjectValuesWrapper>
      <Dialog
        open={isEdit}
        onClose={handleClose}
        slotProps={{
          transition: {
            onEnter: () => {
              setBaseUrl(sso?.baseUrl ?? "");
              setUploadUrl(sso?.uploadUrl ?? "");
            },
          },
        }}
      >
        <form onSubmit={handleSave}>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <Box
              sx={{
                display: "grid",
                gap: 2,
                py: 2,
                pt: 1,
                minWidth: 300,
                maxWidth: "100%",
                width: "100%",
              }}
            >
              <TextField
                value={clientId}
                variant="outlined"
                size="small"
                label="Client ID"
                fullWidth
                required
                autoFocus
                onChange={(e) => setClientID(e.currentTarget.value)}
              />
              <TextField
                value={clientSecret}
                variant="outlined"
                size="small"
                label="Client Secret"
                fullWidth
                required
                onChange={(e) => setClientSecret(e.currentTarget.value)}
              />
              <TextField
                value={baseUrl}
                variant="outlined"
                size="small"
                label="Base URL"
                fullWidth
                onChange={(e) => setBaseUrl(e.currentTarget.value)}
              />
              <TextField
                value={uploadUrl}
                variant="outlined"
                size="small"
                label="Upload URL"
                fullWidth
                onChange={(e) => setUploadUrl(e.currentTarget.value)}
              />
            </Box>
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
