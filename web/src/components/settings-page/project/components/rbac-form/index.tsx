import { FC, memo } from "react";
import { RBAC_DESCRIPTION } from "~/constants/text";
import {
  ProjectDescription,
  ProjectTitle,
  ProjectTitleWrap,
} from "~/styles/project-setting";
import { RoleTable } from "./components/role";
import { UserGroupTable } from "./components/user-group";
import { useGetProject } from "~/queries/project/use-get-project";

const SECTION_TITLE = "Role-Based Access Control";

export const RBACForm: FC = memo(function RBACForm() {
  const { data: projectDetail } = useGetProject();
  const isProjectDisabled = projectDetail?.disabled ?? false;

  return (
    <>
      <ProjectTitleWrap>
        <ProjectTitle variant="h5">{SECTION_TITLE}</ProjectTitle>
      </ProjectTitleWrap>

      <ProjectDescription variant="body1" color="textSecondary">
        {RBAC_DESCRIPTION}
      </ProjectDescription>

      <RoleTable isProjectDisabled={isProjectDisabled} />

      <UserGroupTable isProjectDisabled={isProjectDisabled} />
    </>
  );
});
