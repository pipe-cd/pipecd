import { Typography } from "@mui/material";
import { FC, memo } from "react";
import { RBAC_DESCRIPTION } from "~/constants/text";
import { useProjectSettingStyles } from "~/styles/project-setting";
import { RoleTable } from "./components/role";
import { UserGroupTable } from "./components/user-group";

const SECTION_TITLE = "Role-Based Access Control";

export const RBACForm: FC = memo(function RBACForm() {
  const projectSettingClasses = useProjectSettingStyles();

  return (
    <>
      <div className={projectSettingClasses.title}>
        <Typography
          variant="h5"
          className={projectSettingClasses.titleWithIcon}
        >
          {SECTION_TITLE}
        </Typography>
      </div>

      <Typography
        variant="body1"
        color="textSecondary"
        className={projectSettingClasses.description}
      >
        {RBAC_DESCRIPTION}
      </Typography>

      <RoleTable />

      <UserGroupTable />
    </>
  );
});
