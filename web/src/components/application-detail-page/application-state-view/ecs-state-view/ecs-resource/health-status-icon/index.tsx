import { makeStyles } from "@material-ui/core";
import UnknownIcon from "@material-ui/icons/ErrorOutline";
import FavoriteIcon from "@material-ui/icons/Favorite";
import OtherIcon from "@material-ui/icons/HelpOutline";
import { FC, memo } from "react";
import { HealthStatus } from "~/modules/applications-live-state";

const useStyles = makeStyles((theme) => ({
  healthy: {
    color: theme.palette.success.main,
  },
  unknown: {
    color: theme.palette.warning.main,
  },
  other: {
    color: theme.palette.info.main,
  },
}));

export interface ECSResourceHealthStatusIconProps {
  health: HealthStatus;
}

export const ECSResourceHealthStatusIcon: FC<ECSResourceHealthStatusIconProps> = memo(
  function HealthStatusIcon({ health }) {
    const classes = useStyles();
    switch (health) {
      case HealthStatus.UNKNOWN:
        return <UnknownIcon fontSize="small" className={classes.unknown} />;
      case HealthStatus.HEALTHY:
        return <FavoriteIcon fontSize="small" className={classes.healthy} />;
      case HealthStatus.OTHER:
        return <OtherIcon fontSize="small" className={classes.other} />;
    }
  }
);
