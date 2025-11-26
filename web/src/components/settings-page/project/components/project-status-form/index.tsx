import { Alert, Switch } from "@mui/material";
import { FC, memo, useCallback } from "react";
import {
  PROJECT_DISABLED_WARNING,
  PROJECT_STATUS_DESCRIPTION,
} from "~/constants/text";
import {
  ProjectDescription,
  ProjectTitle,
  ProjectTitleWrap,
} from "~/styles/project-setting";
import { useGetProject } from "~/queries/project/use-get-project";
import { useToggleProjectAvailability } from "~/queries/project/use-toggle-project-availability";
import { useToast } from "~/contexts/toast-context";
import {
  DISABLE_PROJECT_SUCCESS,
  ENABLE_PROJECT_SUCCESS,
} from "~/constants/toast-text";

const SECTION_TITLE = "Project Status";

export const ProjectStatusForm: FC = memo(function ProjectStatusForm() {
  const { data: projectDetail } = useGetProject();
  const isProjectDisabled = projectDetail?.disabled ?? false;
  const hasProjectId = Boolean(projectDetail?.id);

  const { mutateAsync: toggleAvailability, isLoading } =
    useToggleProjectAvailability();
  const { addToast } = useToast();

  const handleToggle = useCallback(() => {
    if (!hasProjectId || isLoading) {
      return;
    }
    toggleAvailability({ enable: isProjectDisabled }).then(() => {
      addToast({
        message: isProjectDisabled
          ? ENABLE_PROJECT_SUCCESS
          : DISABLE_PROJECT_SUCCESS,
        severity: "success",
      });
    });
  }, [addToast, hasProjectId, isLoading, isProjectDisabled, toggleAvailability]);

  return (
    <>
      <ProjectTitleWrap>
        <ProjectTitle variant="h5">{SECTION_TITLE}</ProjectTitle>
        <Switch
          checked={!isProjectDisabled}
          onChange={handleToggle}
          disabled={!hasProjectId || isLoading}
        />
      </ProjectTitleWrap>
      <ProjectDescription variant="body1" color="textSecondary">
        {PROJECT_STATUS_DESCRIPTION}
      </ProjectDescription>
      {isProjectDisabled && (
        <Alert severity="warning" sx={{ mt: 2 }}>
          {PROJECT_DISABLED_WARNING}
        </Alert>
      )}
    </>
  );
});
