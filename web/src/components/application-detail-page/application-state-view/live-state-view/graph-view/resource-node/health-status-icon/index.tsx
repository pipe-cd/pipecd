import UnknownIcon from "@mui/icons-material/ErrorOutline";
import FavoriteIcon from "@mui/icons-material/Favorite";
import OtherIcon from "@mui/icons-material/HelpOutline";
import { FC, memo } from "react";
import { ResourceState } from "~~/model/application_live_state_pb";

type HealthStatusIconProps = {
  health: ResourceState.HealthStatus;
};

export const HealthStatusIcon: FC<HealthStatusIconProps> = memo(
  function HealthStatusIcon({ health }) {
    switch (health) {
      case ResourceState.HealthStatus.UNKNOWN:
        return <UnknownIcon fontSize="small" color="warning" />;
      case ResourceState.HealthStatus.HEALTHY:
        return <FavoriteIcon fontSize="small" color="success" />;
      case ResourceState.HealthStatus.UNHEALTHY:
        return <OtherIcon fontSize="small" color="error" />;
      default:
        return null;
    }
  }
);
