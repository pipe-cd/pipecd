import { makeStyles } from "@material-ui/core";
import UnknownIcon from "@material-ui/icons/ErrorOutline";
import FavoriteIcon from "@material-ui/icons/Favorite";
import OtherIcon from "@material-ui/icons/HelpOutline";
import { FC, memo } from "react";
import { ResourceState } from "~~/model/application_live_state_pb";

const useStyles = makeStyles((theme) => ({
  healthy: {
    color: theme.palette.success.main,
  },
  unknown: {
    color: theme.palette.warning.main,
  },
  unhealthy: {
    color: theme.palette.error.main,
  },
}));

type HealthStatusIconProps = {
  health: ResourceState.HealthStatus;
};

export const HealthStatusIcon: FC<HealthStatusIconProps> = memo(
  function HealthStatusIcon({ health }) {
    const classes = useStyles();
    switch (health) {
      case ResourceState.HealthStatus.UNKNOWN:
        return <UnknownIcon fontSize="small" className={classes.unknown} />;
      case ResourceState.HealthStatus.HEALTHY:
        return <FavoriteIcon fontSize="small" className={classes.healthy} />;
      case ResourceState.HealthStatus.UNHEALTHY:
        return <OtherIcon fontSize="small" className={classes.unhealthy} />;
      default:
        return null;
    }
  }
);
