import { IconButton, makeStyles, Typography } from "@material-ui/core";
import EditIcon from "@material-ui/icons/Edit";
import React, { FC, memo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import {
  fetchProject,
  GitHubSSO,
  Teams,
  updateGitHubSSO,
  updateRBAC,
} from "../modules/project";
import { AppDispatch } from "../store";
import { InputForm } from "./input-form";
import { SSOEditDialog } from "./sso-edit-dialog";

const useStyles = makeStyles((theme) => ({
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

export const GithubSSOForm: FC = memo(function GithubSSOForm() {
  const classes = useStyles();
  const dispatch = useDispatch<AppDispatch>();
  const [isEditSSO, setIsEditSSO] = useState(false);
  const teams = useSelector<AppState, Teams | null>(
    (state) => state.project.teams
  );
  const sso = useSelector<AppState, GitHubSSO | null>(
    (state) => state.project.github
  );

  const handleSaveTeams = (params: Partial<Teams>): void => {
    dispatch(updateRBAC(params)).finally(() => {
      dispatch(fetchProject());
    });
  };

  const handleSaveSSO = (
    params: Partial<GitHubSSO> & { clientId: string; clientSecret: string }
  ): void => {
    dispatch(updateGitHubSSO(params)).finally(() => {
      dispatch(fetchProject());
    });
  };

  return (
    <>
      <Typography variant="h6">GitHub</Typography>
      <div className={classes.indent}>
        <Typography variant="subtitle2">Team</Typography>
        <div className={classes.indent}>
          <InputForm
            currentValue={teams?.admin}
            name="Admin Team"
            onSave={(value) => handleSaveTeams({ admin: value })}
          />
          <InputForm
            currentValue={teams?.editor}
            name="Editor Team"
            onSave={(value) => handleSaveTeams({ editor: value })}
          />
          <InputForm
            currentValue={teams?.viewer}
            name="Viewer Team"
            onSave={(value) => handleSaveTeams({ viewer: value })}
          />
        </div>
      </div>

      <div className={classes.indent}>
        <Typography variant="subtitle2">
          SSO
          <IconButton onClick={() => setIsEditSSO(true)}>
            <EditIcon />
          </IconButton>
          <div className={classes.indent}>
            <div className={classes.item}>
              <Typography variant="subtitle1" className={classes.name}>
                Client ID
              </Typography>
              <Typography variant="body1">{sso?.clientId}</Typography>
            </div>

            <div className={classes.item}>
              <Typography variant="subtitle1" className={classes.name}>
                Client Secret
              </Typography>
              <Typography variant="body1">{sso?.clientSecret}</Typography>
            </div>
            <div className={classes.item}>
              <Typography variant="subtitle1" className={classes.name}>
                Base URL
              </Typography>
              <Typography variant="body1">{sso?.baseUrl}</Typography>
            </div>
            <div className={classes.item}>
              <Typography variant="subtitle1" className={classes.name}>
                Upload URL
              </Typography>
              <Typography variant="body1">{sso?.uploadUrl}</Typography>
            </div>
          </div>
        </Typography>
      </div>
      <SSOEditDialog
        currentBaseURL={sso?.baseUrl ?? ""}
        currentUploadURL={sso?.uploadUrl ?? ""}
        onSave={handleSaveSSO}
        open={isEditSSO}
        onClose={() => setIsEditSSO(false)}
      />
    </>
  );
});
