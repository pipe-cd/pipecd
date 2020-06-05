import React, { FC } from "react";
import { makeStyles, Box } from "@material-ui/core";
import { PipelineStage as PipelineStageModel } from "pipe/pkg/app/web/model/deployment_pb";
import { PipelineStage } from "./pipeline-stage";

const useStyles = makeStyles(theme => ({
  requireLine: {
    position: "relative",
    "&::before": {
      content: '""',
      position: "absolute",
      top: "48%",
      left: -theme.spacing(2),
      borderTop: `2px solid ${theme.palette.divider}`,
      width: theme.spacing(4),
      height: 1
    }
  },
  requireCurvedLine: {
    position: "relative",
    "&::before": {
      content: '""',
      position: "absolute",
      bottom: "50%",
      left: 0,
      borderLeft: `2px solid ${theme.palette.divider}`,
      borderBottom: `2px solid ${theme.palette.divider}`,
      width: theme.spacing(2),
      height: 56 + theme.spacing(4)
    }
  }
}));

interface Props {
  stages: PipelineStageModel.AsObject[][];
}

export const Pipeline: FC<Props> = ({ stages }) => {
  const classes = useStyles();
  return (
    <Box display="flex">
      {stages.map((stageColumn, columnIndex) => (
        <Box display="flex" flexDirection="column">
          {stageColumn.map((stage, stageIndex) => (
            <Box
              display="flex"
              p={2}
              className={
                columnIndex > 0
                  ? stageIndex > 0
                    ? classes.requireCurvedLine
                    : classes.requireLine
                  : undefined
              }
            >
              <PipelineStage name={stage.name} status={stage.status} />
            </Box>
          ))}
        </Box>
      ))}
    </Box>
  );
};
