import { CircularProgress, makeStyles, Typography } from "@material-ui/core";
import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { GithubSSOForm } from "../../components/github-sso-form";
import { StaticAdminForm } from "../../components/static-admin-form";
import {
  UPDATE_STATIC_ADMIN_PASSWORD_SUCCESS,
  UPDATE_STATIC_ADMIN_USERNAME_SUCCESS,
} from "../../constants/toast-text";
import { AppState } from "../../modules";
import {
  fetchProject,
  ProjectState,
  toggleAvailability,
  updatePassword,
  updateUsername,
} from "../../modules/project";
import { addToast } from "../../modules/toasts";
import { AppDispatch } from "../../store";

const useStyles = makeStyles((theme) => ({
  main: {
    overflow: "auto",
    padding: theme.spacing(3),
    background: theme.palette.background.paper,
    flex: 1,
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
      dispatch(
        addToast({
          message: UPDATE_STATIC_ADMIN_USERNAME_SUCCESS,
          severity: "success",
        })
      );
    });
  };
  const handleUpdatePassword = async (password: string): Promise<unknown> => {
    return dispatch(updatePassword({ password })).then(() => {
      dispatch(
        addToast({
          message: UPDATE_STATIC_ADMIN_PASSWORD_SUCCESS,
          severity: "success",
        })
      );
    });
  };
  const handleToggleAvailability = (): void => {
    dispatch(toggleAvailability()).then(() => {
      dispatch(fetchProject());
    });
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
      <GithubSSOForm />
    </div>
  );
});
