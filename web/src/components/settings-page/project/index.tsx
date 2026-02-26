import { FC, memo } from "react";
import { Box } from "@mui/material";
import { GithubSSOForm } from "./components/github-sso-form";
import { ProjectStatusForm } from "./components/project-status-form";
import { RBACForm } from "./components/rbac-form";
import { StaticAdminForm } from "./components/static-admin-form";

export const SettingsProjectPage: FC = memo(function SettingsProjectPage() {
  return (
    <Box
      sx={(theme) => ({
        overflow: "auto",
        padding: theme.spacing(3),
        background: theme.palette.background.paper,
        flex: 1,
      })}
    >
      <ProjectStatusForm />
      <Box sx={{ marginTop: 1 }}>
        <StaticAdminForm />
      </Box>
      <GithubSSOForm />
      <RBACForm />
    </Box>
  );
});
