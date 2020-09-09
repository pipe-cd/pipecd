import { makeStyles } from "@material-ui/core";
import React, { FC, memo, useEffect } from "react";
import { useDispatch } from "react-redux";
import { GithubSSOForm } from "../../components/github-sso-form";
import { RBACForm } from "../../components/rbac-form";
import { StaticAdminForm } from "../../components/static-admin-form";
import { fetchProject } from "../../modules/project";
import { AppDispatch } from "../../store";

const useStyles = makeStyles((theme) => ({
  main: {
    overflow: "auto",
    padding: theme.spacing(3),
    background: theme.palette.background.paper,
    flex: 1,
  },
}));

export const SettingsProjectPage: FC = memo(function SettingsProjectPage() {
  const classes = useStyles();
  const dispatch = useDispatch<AppDispatch>();

  useEffect(() => {
    dispatch(fetchProject());
  }, [dispatch]);

  return (
    <div className={classes.main}>
      <StaticAdminForm />
      <GithubSSOForm />
      <RBACForm />
    </div>
  );
});
