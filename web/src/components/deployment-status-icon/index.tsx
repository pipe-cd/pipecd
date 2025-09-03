import {
  Cached,
  CheckCircle,
  Error,
  IndeterminateCheckBox,
  Cancel,
} from "@mui/icons-material";
import { DeploymentStatus } from "~/types/deployment";
import { FC } from "react";
import clsx from "clsx";

export interface DeploymentStatusIconProps {
  status: DeploymentStatus;
  className?: string;
}

export const DeploymentStatusIcon: FC<DeploymentStatusIconProps> = ({
  status,
  className,
}) => {
  switch (status) {
    case DeploymentStatus.DEPLOYMENT_SUCCESS:
      return (
        <CheckCircle
          className={clsx(className)}
          sx={{ color: "success.main" }}
          data-testid="deployment-success-icon"
        />
      );
    case DeploymentStatus.DEPLOYMENT_FAILURE:
      return (
        <Error
          className={clsx(className)}
          sx={{ color: "error.main" }}
          data-testid="deployment-error-icon"
        />
      );
    case DeploymentStatus.DEPLOYMENT_CANCELLED:
      return (
        <Cancel
          className={clsx(className)}
          sx={{ color: "grey.500" }}
          data-testid="deployment-cancel-icon"
        />
      );
    case DeploymentStatus.DEPLOYMENT_RUNNING:
      return (
        <Cached
          className={clsx(className)}
          sx={{
            color: "info.main",
            animation: "spin 3s linear infinite",
            "@keyframes spin": {
              "0%": { transform: "rotate(360deg)" },
              "100%": { transform: "rotate(0deg)" },
            },
          }}
          data-testid="deployment-running-icon"
        />
      );
    case DeploymentStatus.DEPLOYMENT_ROLLING_BACK:
      return (
        <Cached
          className={clsx(className)}
          sx={{ color: "info.main" }}
          data-testid="deployment-rollback-icon"
        />
      );
    case DeploymentStatus.DEPLOYMENT_PENDING:
    case DeploymentStatus.DEPLOYMENT_PLANNED:
      return (
        <IndeterminateCheckBox
          sx={{ color: "grey.500" }}
          className={clsx(className)}
          data-testid="deployment-pending-icon"
        />
      );
  }
};
