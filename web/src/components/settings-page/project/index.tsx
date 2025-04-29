import makeStyles from "@mui/styles/makeStyles";
import { FC, memo, useEffect } from "react";
import { useAppDispatch } from "~/hooks/redux";
import { fetchProject } from "~/modules/project";
import { GithubSSOForm } from "./components/github-sso-form";
import { RBACForm } from "./components/rbac-form";
import { StaticAdminForm } from "./components/static-admin-form";

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
  const dispatch = useAppDispatch();

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
