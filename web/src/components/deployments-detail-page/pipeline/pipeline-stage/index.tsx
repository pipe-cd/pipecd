import { Paper, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import clsx from "clsx";
import { FC, memo } from "react";
import { StageStatus } from "~/modules/deployments";
import { StageStatusIcon } from "./stage-status-icon";

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
    fontFamily: theme.typography.fontFamilyMono,
  },
  main: {
    display: "flex",
    justifyContent: "flex-start",
    alignItems: "center",
  },
  metadata: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(4),
    textAlign: "left",
  },
}));

export interface PipelineStageProps {
  id: string;
  name: string;
  status: StageStatus;
  active: boolean;
  isDeploymentRunning: boolean;
  approver?: string;
  skipper?: string;
  metadata: [string, string][];
  onClick: (stageId: string, stageName: string) => void;
}

const TRAFFIC_PERCENTAGE_META_KEY = {
  PRIMARY: "primary-percentage",
  CANARY: "canary-percentage",
  BASELINE: "baseline-percentage",
  PROMOTE: "promote-percentage",
};

const trafficPercentageMetaKey: Record<string, string> = {
  [TRAFFIC_PERCENTAGE_META_KEY.PRIMARY]: "Primary",
  [TRAFFIC_PERCENTAGE_META_KEY.CANARY]: "Canary",
  [TRAFFIC_PERCENTAGE_META_KEY.BASELINE]: "Baseline",
  [TRAFFIC_PERCENTAGE_META_KEY.PROMOTE]: "Promoted",
};

const createTrafficPercentageText = (meta: [string, string][]): string => {
  const map = meta.reduce<Record<string, string>>((prev, [key, value]) => {
    if (trafficPercentageMetaKey[key]) {
      if (key === TRAFFIC_PERCENTAGE_META_KEY.PROMOTE) {
        prev[key] = `${value}% ${trafficPercentageMetaKey[key]}`;
      } else {
        prev[key] = `${trafficPercentageMetaKey[key]} ${value}%`;
      }
    }
    return prev;
  }, {});

  // Serverless promote stage detail.
  if (map[TRAFFIC_PERCENTAGE_META_KEY.PROMOTE]) {
    return `${map[TRAFFIC_PERCENTAGE_META_KEY.PROMOTE]}`;
  }

  // Traffic routing stage detail.
  let detail = "";
  if (map[TRAFFIC_PERCENTAGE_META_KEY.PRIMARY]) {
    detail += `${map[TRAFFIC_PERCENTAGE_META_KEY.PRIMARY]}`;
  }
  if (map[TRAFFIC_PERCENTAGE_META_KEY.CANARY]) {
    detail += `, ${map[TRAFFIC_PERCENTAGE_META_KEY.CANARY]}`;
  }
  if (map[TRAFFIC_PERCENTAGE_META_KEY.BASELINE]) {
    detail += `, ${map[TRAFFIC_PERCENTAGE_META_KEY.BASELINE]}`;
  }

  return detail;
};

export const PipelineStage: FC<PipelineStageProps> = memo(
  function PipelineStage({
    id,
    name,
    status,
    onClick,
    active,
    approver,
    skipper,
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
        ) : skipper !== undefined ? (
          <div className={classes.metadata}>
            <Typography
              variant="body2"
              color="inherit"
            >{`Skipped by ${skipper}`}</Typography>
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
  }
);
