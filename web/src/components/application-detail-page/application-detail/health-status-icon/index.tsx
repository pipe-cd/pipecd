import UnknownIcon from "@mui/icons-material/ErrorOutline";
import FavoriteIcon from "@mui/icons-material/Favorite";
import OtherIcon from "@mui/icons-material/HelpOutline";
import { FC, memo } from "react";
import { ApplicationLiveStateSnapshot } from "~/types/applications-live-state";

export interface ApplicationHealthStatusIconProps {
  health: ApplicationLiveStateSnapshot.Status;
}

export const ApplicationHealthStatusIcon: FC<ApplicationHealthStatusIconProps> = memo(
  function HealthStatusIcon({ health }) {
    switch (health) {
      case ApplicationLiveStateSnapshot.Status.UNKNOWN:
        return <UnknownIcon color={"warning"} />;
      case ApplicationLiveStateSnapshot.Status.HEALTHY:
        return <FavoriteIcon color={"success"} />;
      case ApplicationLiveStateSnapshot.Status.OTHER:
        return <OtherIcon color={"info"} />;
      case ApplicationLiveStateSnapshot.Status.UNHEALTHY:
        return <OtherIcon color={"info"} />;
    }
  }
);
