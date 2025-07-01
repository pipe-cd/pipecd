import { Box, Paper, Typography } from "@mui/material";
import { FC, memo } from "react";
import { StageStatus } from "~/modules/deployments";
import { StageStatusIcon } from "./stage-status-icon";

export interface PipelineStageProps {
  id: string;
  name: string;
  status: StageStatus;
  active: boolean;
  isDeploymentRunning: boolean;
  metadata: [string, string][];
  displayMetadataText?: string;
  onClick: (stageId: string, stageName: string) => void;
}

// TODO: Use METADATA_STAGE_DISPLAY_KEY instead for all fields in pipedv1.
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
    metadata,
    isDeploymentRunning,
    displayMetadataText,
  }) {
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
        // className={clsx(classes.root, {
        //   [classes.active]: active,
        //   [classes.notStartedYet]: disabled,
        // })}
        onClick={handleOnClick}
        sx={(theme) => ({
          flex: 1,
          display: "inline-flex",
          flexDirection: "column",
          padding: theme.spacing(2),
          cursor: disabled ? "unset" : "pointer",
          backgroundColor: active
            ? theme.palette.primary.main + "1e" // NOTE: 12%
            : undefined,
          "&:hover": {
            backgroundColor: disabled
              ? theme.palette.background.paper
              : theme.palette.action.hover,
          },

          color: disabled ? theme.palette.text.disabled : undefined,
        })}
      >
        <Box
          sx={{
            display: "flex",
            justifyContent: "flex-start",
            alignItems: "center",
          }}
        >
          <StageStatusIcon status={status} />
          <Typography
            variant="subtitle2"
            sx={{
              marginLeft: 1,
              maxWidth: 200,
              whiteSpace: "nowrap",
              textOverflow: "ellipsis",
              overflow: "hidden",
            }}
          >
            <Box
              title={name}
              component={"span"}
              sx={{ fontFamily: "fontFamilyMono" }}
            >
              {name}
            </Box>
          </Typography>
        </Box>
        {displayMetadataText && (
          <Box
            sx={{
              color: "text.secondary",
              marginLeft: 4,
              textAlign: "left",
            }}
          >
            <Typography
              variant="body2"
              color="inherit"
            >{`${displayMetadataText}`}</Typography>
          </Box>
        )}
        {/* TODO: remove trafficPercentage and use only displayMetadataText */}
        {trafficPercentage && (
          <Box
            sx={{
              color: "text.secondary",
              marginLeft: 4,
              textAlign: "left",
            }}
          >
            <Typography variant="body2" color="inherit">
              {trafficPercentage}
            </Typography>
          </Box>
        )}
      </Paper>
    );
  }
);
