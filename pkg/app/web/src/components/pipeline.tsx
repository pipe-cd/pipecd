import React, { FC, memo, useCallback, useState } from "react";
import {
  makeStyles,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
  Button,
} from "@material-ui/core";
import { PipelineStage } from "./pipeline-stage";
import { useSelector, useDispatch } from "react-redux";
import { AppState } from "../modules";
import {
  selectById,
  Deployment,
  Stage,
  approveStage,
} from "../modules/deployments";
import { fetchStageLog } from "../modules/stage-logs";
import { updateActiveStage, ActiveStage } from "../modules/active-stage";
import { ApprovalStage } from "./approval-stage";
import clsx from "clsx";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";

const WAIT_APPROVAL_NAME = "WAIT_APPROVAL";

const useConvertedStages = (deploymentId: string): Stage[][] => {
  const stages: Stage[][] = [];
  const deployment = useSelector<AppState, Deployment | undefined>((state) =>
    selectById(state.deployments, deploymentId)
  );

  if (!deployment) {
    return stages;
  }

  stages[0] = deployment.stagesList.filter(
    (stage) => stage.requiresList.length === 0 && stage.visible
  );

  let index = 0;
  while (stages[index].length > 0) {
    const previousIds = stages[index].map((stage) => stage.id);
    index++;
    stages[index] = deployment.stagesList.filter(
      (stage) =>
        stage.requiresList.some((id) => previousIds.includes(id)) &&
        stage.visible
    );
  }
  return stages;
};

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
  },
  pipelineColumn: {
    display: "flex",
    flexDirection: "column",
  },
  stage: {
    display: "flex",
    padding: theme.spacing(2),
  },
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
  approveDialog: {
    display: "flex",
  },
}));

interface Props {
  deploymentId: string;
}

export const Pipeline: FC<Props> = memo(function Pipeline({ deploymentId }) {
  const classes = useStyles();
  const dispatch = useDispatch();
  const [approveTargetId, setApproveTargetId] = useState<string | null>(null);
  const isOpenApproveDialog = Boolean(approveTargetId);
  const stages = useConvertedStages(deploymentId);
  const activeStage = useSelector<AppState, ActiveStage>(
    (state) => state.activeStage
  );

  const handleOnClickStage = useCallback(
    (stageId: string, stageName: string) => {
      dispatch(
        fetchStageLog({
          deploymentId,
          stageId,
          offsetIndex: 0,
          retriedCount: 0,
        })
      );
      dispatch(
        updateActiveStage({ id: `${deploymentId}/${stageId}`, name: stageName })
      );
    },
    [dispatch, deploymentId]
  );

  const handleApprove = (): void => {
    if (approveTargetId) {
      dispatch(approveStage({ deploymentId, stageId: approveTargetId }));
      setApproveTargetId(null);
    }
  };

  return (
    <div className={classes.root}>
      {stages.map((stageColumn, columnIndex) => (
        <div className={classes.pipelineColumn} key={`pipeline-${columnIndex}`}>
          {stageColumn.map((stage, stageIndex) => (
            <div
              key={stage.id}
              className={clsx(
                classes.stage,
                columnIndex > 0
                  ? stageIndex > 0
                    ? classes.requireCurvedLine
                    : classes.requireLine
                  : undefined
              )}
            >
              {stage.name === WAIT_APPROVAL_NAME &&
              stage.status === StageStatus.STAGE_RUNNING ? (
                <ApprovalStage
                  id={stage.id}
                  name={stage.name}
                  onClick={() => {
                    setApproveTargetId(stage.id);
                  }}
                  active={activeStage?.id === `${deploymentId}/${stage.id}`}
                />
              ) : (
                <PipelineStage
                  id={stage.id}
                  name={stage.name}
                  status={stage.status}
                  onClick={handleOnClickStage}
                  active={activeStage?.id === `${deploymentId}/${stage.id}`}
                />
              )}
            </div>
          ))}
        </div>
      ))}

      <Dialog
        open={isOpenApproveDialog}
        onClose={() => setApproveTargetId(null)}
      >
        <DialogTitle>Approve stage</DialogTitle>
        <DialogContent>
          <DialogContentText>
            {`To continue deploying, click "APPROVE".`}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setApproveTargetId(null)}>CANCEL</Button>
          <Button color="primary" onClick={handleApprove}>
            APPROVE
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
});
