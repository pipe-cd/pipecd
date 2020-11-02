import { makeStyles, Paper, Typography } from "@material-ui/core";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";
import React, { FC, memo } from "react";
import { StageStatusIcon } from "./stage-status-icon";
import clsx from "clsx";

const useStyles = makeStyles((theme) => ({
  root: {
    flex: 1,
    display: "inline-flex",
    flexDirection: "column",
    padding: theme.spacing(2),
    cursor: "pointer",
    "&:hover": {
      backgroundColor: theme.palette.action.hover,
    },
  },
  active: {
    // NOTE: 12%
    backgroundColor: theme.palette.primary.main + "1e",
  },
  notStartedYet: {
    color: theme.palette.text.disabled,
    cursor: "unset",
    "&:hover": {
      backgroundColor: theme.palette.background.paper,
    },
  },
  name: {
    marginLeft: theme.spacing(1),
    maxWidth: 200,
    whiteSpace: "nowrap",
    textOverflow: "ellipsis",
    overflow: "hidden",
  },
  stageName: {
    fontFamily: "Roboto Mono",
  },
  main: {
    display: "flex",
    justifyContent: "flex-start",
    alignItems: "center",
  },
  metadata: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(4),
  },
}));

interface Props {
  id: string;
  name: string;
  status: StageStatus;
  active: boolean;
  isDeploymentRunning: boolean;
  approver?: string;
  metadata: [string, string][];
  onClick: (stageId: string, stageName: string) => void;
}

const TRAFFIC_PERCENTAGE_META_KEY = {
  PRIMARY: "primary-percentage",
  CANARY: "canary-percentage",
  BASELINE: "baseline-percentage",
};

const trafficPercentageMetaKey: Record<string, string> = {
  [TRAFFIC_PERCENTAGE_META_KEY.PRIMARY]: "Primary",
  [TRAFFIC_PERCENTAGE_META_KEY.CANARY]: "Canary",
  [TRAFFIC_PERCENTAGE_META_KEY.BASELINE]: "Baseline",
};

const createTrafficPercentageText = (meta: [string, string][]): string => {
  const map = meta.reduce<Record<string, string>>((prev, [key, value]) => {
    if (trafficPercentageMetaKey[key]) {
      prev[key] = `${trafficPercentageMetaKey[key]} ${value}%`;
    }
    return prev;
  }, {});

  // If the primary exists, other params also exist.
  if (map[TRAFFIC_PERCENTAGE_META_KEY.PRIMARY]) {
    return `${map[TRAFFIC_PERCENTAGE_META_KEY.PRIMARY]}, ${
      map[TRAFFIC_PERCENTAGE_META_KEY.CANARY]
    }, ${map[TRAFFIC_PERCENTAGE_META_KEY.BASELINE]}`;
  }

  return "";
};

export const PipelineStage: FC<Props> = memo(function PipelineStage({
  id,
  name,
  status,
  onClick,
  active,
  approver,
  metadata,
  isDeploymentRunning,
}) {
  const classes = useStyles();
  const disabled =
    isDeploymentRunning === false &&
    status === StageStatus.STAGE_NOT_STARTED_YET;

  function handleOnClick(): void {
    if (disabled) {
      return;
    }
    onClick(id, name);
  }

  const trafficPercentage = createTrafficPercentageText(metadata);

  return (
    <Paper
      square
      className={clsx(classes.root, {
        [classes.active]: active,
        [classes.notStartedYet]: disabled,
      })}
      onClick={handleOnClick}
    >
      <div className={classes.main}>
        <StageStatusIcon status={status} />
        <Typography variant="subtitle2" className={classes.name}>
          <span title={name} className={classes.stageName}>
            {name}
          </span>
        </Typography>
      </div>
      {approver !== undefined ? (
        <div className={classes.metadata}>
          <Typography
            variant="body2"
            color="inherit"
          >{`Approved by ${approver}`}</Typography>
        </div>
      ) : null}
      {trafficPercentage && (
        <div className={classes.metadata}>
          <Typography variant="body2" color="inherit">
            {trafficPercentage}
          </Typography>
        </div>
      )}
    </Paper>
  );
});
