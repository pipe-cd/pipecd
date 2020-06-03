import React, { FC } from "react";
import { makeStyles } from "@material-ui/core";
import {
  CheckCircle,
  Error,
  Cached,
  Stop,
  IndeterminateCheckBox
} from "@material-ui/icons";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";

const useStyles = makeStyles(theme => ({
  [StageStatus.STAGE_SUCCESS]: {
    color: theme.palette.success.main
  },
  [StageStatus.STAGE_RUNNING]: {
    color: theme.palette.info.main
  },
  [StageStatus.STAGE_FAILURE]: {
    color: theme.palette.error.main
  },
  [StageStatus.STAGE_CANCELLED]: {
    color: theme.palette.error.main
  },
  [StageStatus.STAGE_NOT_STARTED_YET]: {
    color: theme.palette.grey[500]
  }
}));

interface Props {
  status: StageStatus;
}

export const StageStatusIcon: FC<Props> = ({ status }) => {
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
  }
};
