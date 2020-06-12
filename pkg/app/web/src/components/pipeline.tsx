import React, { FC, memo } from "react";
import { makeStyles, Box } from "@material-ui/core";
import { PipelineStage } from "./pipeline-stage";
import { useSelector } from "react-redux";
import { AppState } from "../modules";
import { selectById, Deployment, Stage } from "../modules/deployments";

const useConvertedStages = (deploymentId: string) => {
  const stages: Stage[][] = [];
  const deployment = useSelector<AppState, Deployment | undefined>((state) =>
    selectById(state.deployments, deploymentId)
  );

  if (!deployment) {
    return stages;
  }

  stages[0] = deployment.stagesList.filter(
    (stage) => stage.requiresList.length === 0
  );

  let index = 0;
  while (stages[index].length > 0) {
    const previousIds = stages[index].map((stage) => stage.id);
    index++;
    stages[index] = deployment.stagesList.filter((stage) =>
      stage.requiresList.some((id) => previousIds.includes(id))
    );
  }
  return stages;
};

const useStyles = makeStyles((theme) => ({
  requireLine: {
    position: "relative",
    "&::before": {
      content: '""',
      position: "absolute",
      top: "48%",
      left: -theme.spacing(2),
      borderTop: `2px solid ${theme.palette.divider}`,
      width: theme.spacing(4),
      height: 1,
    },
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
      height: 56 + theme.spacing(4),
    },
  },
}));

interface Props {
  deploymentId: string;
}

export const Pipeline: FC<Props> = memo(function Pipeline({ deploymentId }) {
  const classes = useStyles();
  const stages = useConvertedStages(deploymentId);

  return (
    <Box display="flex">
      {stages.map((stageColumn, columnIndex) => (
        <Box
          display="flex"
          flexDirection="column"
          key={`pipeline-${columnIndex}`}
        >
          {stageColumn.map((stage, stageIndex) => (
            <Box
              display="flex"
              p={2}
              key={stage.id}
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
});
