import UnknownIcon from "@mui/icons-material/ErrorOutline";
import FavoriteIcon from "@mui/icons-material/Favorite";
import OtherIcon from "@mui/icons-material/HelpOutline";
import { FC, memo } from "react";
import { HealthStatus } from "~/modules/applications-live-state";
export interface KubernetesResourceHealthStatusIconProps {
  health: HealthStatus;
}

export const KubernetesResourceHealthStatusIcon: FC<KubernetesResourceHealthStatusIconProps> = memo(
  function HealthStatusIcon({ health }) {
    switch (health) {
      case HealthStatus.UNKNOWN:
        return <UnknownIcon fontSize="small" color="warning" />;
      case HealthStatus.HEALTHY:
        return <FavoriteIcon fontSize="small" color="success" />;
      case HealthStatus.OTHER:
        return <OtherIcon fontSize="small" color="info" />;
    }
  }
);
