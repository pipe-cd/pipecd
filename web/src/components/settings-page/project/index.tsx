import { FC, memo, useEffect } from "react";
import { useAppDispatch } from "~/hooks/redux";
import { fetchProject } from "~/modules/project";
import { GithubSSOForm } from "./components/github-sso-form";
import { RBACForm } from "./components/rbac-form";
import { StaticAdminForm } from "./components/static-admin-form";
import { Box } from "@mui/material";

export const SettingsProjectPage: FC = memo(function SettingsProjectPage() {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(fetchProject());
  }, [dispatch]);

  return (
    <Box
      sx={(theme) => ({
        overflow: "auto",
        padding: theme.spacing(3),
        background: theme.palette.background.paper,
        flex: 1,
      })}
    >
      <StaticAdminForm />
      <GithubSSOForm />
      <RBACForm />
    </Box>
  );
});
