import { makeStyles } from "@material-ui/core";
import {
  Cached,
  CheckCircle,
  Error,
  IndeterminateCheckBox,
} from "@material-ui/icons";
import { DeploymentStatus } from "pipe/pkg/app/web/model/deployment_pb";
import React, { FC } from "react";

const useStyles = makeStyles((theme) => ({
  [DeploymentStatus.DEPLOYMENT_SUCCESS]: {
    color: theme.palette.success.main,
  },
  [DeploymentStatus.DEPLOYMENT_RUNNING]: {
    color: theme.palette.info.main,
    animation: `$running 3s linear infinite`,
  },
  [DeploymentStatus.DEPLOYMENT_ROLLING_BACK]: {
    color: theme.palette.info.main,
  },
  [DeploymentStatus.DEPLOYMENT_FAILURE]: {
    color: theme.palette.error.main,
  },
  [DeploymentStatus.DEPLOYMENT_CANCELLED]: {
    color: theme.palette.error.main,
  },
  [DeploymentStatus.DEPLOYMENT_PENDING]: {
    color: theme.palette.grey[500],
  },
  [DeploymentStatus.DEPLOYMENT_PLANNED]: {
    color: theme.palette.grey[500],
  },
  "@keyframes running": {
    "0%": {
      transform: "rotate(0deg)",
    },
    "100%": {
      transform: "rotate(360deg)",
    },
  },
}));

interface Props {
  status: DeploymentStatus;
}

export const StatusIcon: FC<Props> = ({ status }) => {
  const classes = useStyles();

  switch (status) {
    case DeploymentStatus.DEPLOYMENT_SUCCESS:
      return <CheckCircle className={classes[status]} />;
    case DeploymentStatus.DEPLOYMENT_FAILURE:
    case DeploymentStatus.DEPLOYMENT_CANCELLED:
      return <Error className={classes[status]} />;
    case DeploymentStatus.DEPLOYMENT_RUNNING:
      return <Cached className={classes[status]} />;
    case DeploymentStatus.DEPLOYMENT_ROLLING_BACK:
      return <Cached className={classes[status]} />;
    case DeploymentStatus.DEPLOYMENT_PENDING:
    case DeploymentStatus.DEPLOYMENT_PLANNED:
      return <IndeterminateCheckBox className={classes[status]} />;
  }
};
