import { CircularProgress, makeStyles, Typography } from "@material-ui/core";
import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { StaticAdminForm } from "../../components/static-admin-form";
import { AppState } from "../../modules";
import {
  fetchProject,
  ProjectState,
  toggleAvailability,
  updatePassword,
  updateUsername,
  ProjectSingleSignOn,
  updateGitHubSSO,
} from "../../modules/project";
import { AppDispatch } from "../../store";
import { GithubSSOForm } from "../../components/github-sso-form";

const useStyles = makeStyles((theme) => ({
  main: {
    overflow: "auto",
    padding: theme.spacing(3),
    background: theme.palette.background.paper,
  },
  group: {
    padding: theme.spacing(1),
  },
  titleMargin: {
    marginTop: theme.spacing(2),
  },
}));

export const SettingsProjectPage: FC = memo(function SettingsProjectPage() {
  const classes = useStyles();
  const dispatch = useDispatch<AppDispatch>();
  const project = useSelector<AppState, ProjectState>((state) => state.project);

  useEffect(() => {
    dispatch(fetchProject());
  }, [dispatch]);

  const handleUpdateUsername = (username: string): void => {
    dispatch(updateUsername({ username })).then(() => {
      dispatch(fetchProject());
    });
  };
  const handleUpdatePassword = (password: string): Promise<unknown> => {
    return dispatch(updatePassword({ password }));
  };
  const handleToggleAvailability = (): void => {
    dispatch(toggleAvailability());
  };

  const handleSaveGitHubSSO = (
    params: ProjectSingleSignOn.GitHub.AsObject
  ): Promise<unknown> => {
    return dispatch(updateGitHubSSO(params));
  };

  if (!project) {
    return <CircularProgress />;
  }

  return (
    <div className={classes.main}>
      <StaticAdminForm
        username={project.username}
        staticAdminDisabled={project.staticAdminDisabled}
        isUpdatingUsername={project.isUpdatingUsername}
        isUpdatingPassword={project.isUpdatingPassword}
        onUpdateUsername={handleUpdateUsername}
        onUpdatePassword={handleUpdatePassword}
        onToggleAvailability={handleToggleAvailability}
      />
      <Typography variant="h5">Single Sign On</Typography>
      <div className={classes.group}>
        <GithubSSOForm
          onSave={handleSaveGitHubSSO}
          isSaving={project.isUpdatingGitHubSSO}
        />
      </div>
    </div>
  );
});
