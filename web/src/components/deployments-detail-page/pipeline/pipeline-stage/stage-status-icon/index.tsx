import { makeStyles } from "@material-ui/core";
import {
  Cached,
  CheckCircle,
  Error,
  IndeterminateCheckBox,
  Stop,
  Block,
} from "@material-ui/icons";
import { FC } from "react";
import { StageStatus } from "~/modules/deployments";

const useStyles = makeStyles((theme) => ({
  [StageStatus.STAGE_SUCCESS]: {
    color: theme.palette.success.main,
  },
  [StageStatus.STAGE_RUNNING]: {
    color: theme.palette.info.main,
    animation: `$running 3s linear infinite`,
  },
  [StageStatus.STAGE_FAILURE]: {
    color: theme.palette.error.main,
  },
  [StageStatus.STAGE_CANCELLED]: {
    color: theme.palette.error.main,
  },
  [StageStatus.STAGE_NOT_STARTED_YET]: {
    color: theme.palette.grey[500],
  },
  [StageStatus.STAGE_SKIPPED]: {
    color: theme.palette.grey[500],
  },
  [StageStatus.STAGE_EXITED]: {
    color: theme.palette.success.main,
  },
  "@keyframes running": {
    "0%": {
      transform: "rotate(360deg)",
    },
    "100%": {
      transform: "rotate(0deg)",
    },
  },
}));

export interface StageStatusIconProps {
  status: StageStatus;
}

export const StageStatusIcon: FC<StageStatusIconProps> = ({ status }) => {
  const classes = useStyles();

  switch (status) {
    case StageStatus.STAGE_SUCCESS:
      return <CheckCircle className={classes[status]} />;
    case StageStatus.STAGE_FAILURE:
      return <Error className={classes[status]} />;
    case StageStatus.STAGE_CANCELLED:
      return <Stop className={classes[status]} />;
    case StageStatus.STAGE_NOT_STARTED_YET:
      return <IndeterminateCheckBox className={classes[status]} />;
    case StageStatus.STAGE_RUNNING:
      return <Cached className={classes[status]} />;
    case StageStatus.STAGE_SKIPPED:
      return <Block className={classes[status]} />;
    case StageStatus.STAGE_EXITED:
      return <CheckCircle className={classes[status]} />;
  }
};
