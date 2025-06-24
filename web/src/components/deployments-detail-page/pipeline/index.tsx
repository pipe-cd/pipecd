import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from "@mui/material";
import { FC, memo, useCallback, useEffect, useState } from "react";
import {
  METADATA_APPROVED_BY,
  METADATA_SKIPPED_BY,
} from "~/constants/metadata-keys";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { ActiveStage, updateActiveStage } from "~/modules/active-stage";
import {
  approveStage,
  Deployment,
  isDeploymentRunning,
  selectById,
  Stage,
  StageStatus,
} from "~/modules/deployments";
import { fetchStageLog } from "~/modules/stage-logs";
import { ApprovalStage } from "./approval-stage";
import { PipelineStage } from "./pipeline-stage";
import { ManualOperation } from "~~/model/deployment_pb";

const STAGE_HEIGHT = 56;
const APPROVED_STAGE_HEIGHT = 66;

// Find stage that is running or latest
const findDefaultActiveStage = (
  deployment: Deployment.AsObject | undefined
): Stage | null => {
  if (!deployment) {
    return null;
  }

  const stages = deployment.stagesList.filter(
    (stage) =>
      stage.visible && stage.status !== StageStatus.STAGE_NOT_STARTED_YET
  );

  const runningStage = stages.find(
    (stage) => stage.status === StageStatus.STAGE_RUNNING
  );

  if (runningStage) {
    return runningStage;
  }

  return stages[stages.length - 1];
};

const createStagesForRendering = (
  deployment: Deployment.AsObject | undefined
): Stage[][] => {
  if (!deployment) {
    return [];
  }

  const stages: Stage[][] = [];
  const visibleStages = deployment.stagesList.filter((stage) => stage.visible);

  stages[0] = visibleStages.filter((stage) => stage.requiresList.length === 0);

  let index = 0;
  while (stages[index].length > 0) {
    const previousIds = stages[index].map((stage) => stage.id);
    index++;
    stages[index] = visibleStages.filter((stage) =>
      stage.requiresList.some((id) => previousIds.includes(id))
    );
  }
  return stages;
};

const LARGE_STAGE_NAMES = ["WAIT_APPROVAL", "K8S_TRAFFIC_ROUTING"];

export interface PipelineProps {
  deploymentId: string;
}

const findApprover = (
  metadata: Array<[string, string]>
): string | undefined => {
  const res = metadata.find(([key]) => key === METADATA_APPROVED_BY);

  if (res) {
    return res[1];
  }

  return undefined;
};

const findSkipper = (metadata: Array<[string, string]>): string | undefined => {
  const res = metadata.find(([key]) => key === METADATA_SKIPPED_BY);

  if (res) {
    return res[1];
  }

  return undefined;
};

export const Pipeline: FC<PipelineProps> = memo(function Pipeline({
  deploymentId,
}) {
  const dispatch = useAppDispatch();
  const deployment = useAppSelector<Deployment.AsObject | undefined>((state) =>
    selectById(state.deployments, deploymentId)
  );
  const [approveTargetId, setApproveTargetId] = useState<string | null>(null);
  const isOpenApproveDialog = Boolean(approveTargetId);

  const defaultActiveStage = findDefaultActiveStage(deployment);
  const stages = createStagesForRendering(deployment);
  const isRunning = isDeploymentRunning(deployment?.status);

  const activeStage = useAppSelector<ActiveStage>((state) => state.activeStage);

  useEffect(() => {
    if (defaultActiveStage) {
      dispatch(
        updateActiveStage({
          deploymentId,
          stageId: defaultActiveStage.id,
          name: defaultActiveStage.name,
        })
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [dispatch, deploymentId, defaultActiveStage === null]);

  useEffect(() => {
    if (activeStage) {
      dispatch(
        fetchStageLog({
          deploymentId,
          stageId: activeStage.stageId,
          offsetIndex: 0,
          retriedCount: 0,
        })
      );
    }
  }, [dispatch, deploymentId, activeStage]);

  const handleOnClickStage = useCallback(
    (stageId: string, stageName: string) => {
      dispatch(updateActiveStage({ deploymentId, stageId, name: stageName }));
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
    <Box
      sx={{
        textAlign: "center",
        overflow: "scroll",

        "&::-webkit-scrollbar": {
          height: "7px",
        },

        "&::-webkit-scrollbar-thumb": {
          borderRadius: 8,
          backgroundColor: "rgba(0,0,0,0.3)",
        },
      }}
    >
      <Box
        sx={{
          display: "inline-flex",
        }}
      >
        {stages.map((stageColumn, columnIndex) => {
          let isPrevStageLarge = false;
          return (
            <Box
              sx={{
                display: "flex",
                flexDirection: "column",
              }}
              key={`pipeline-${columnIndex}`}
            >
              {stageColumn.map((stage, stageIndex) => {
                const approver = findApprover(stage.metadataMap);
                const skipper = findSkipper(stage.metadataMap);
                const isActive = activeStage
                  ? activeStage.deploymentId === deploymentId &&
                    activeStage.stageId === stage.id
                  : false;

                const showLine = columnIndex > 0;
                const showStraightLine = showLine && stageIndex === 0;
                const showCurvedLine = showLine && stageIndex > 0;
                const isCurvedLineExtend =
                  showCurvedLine && (Boolean(approver) || isPrevStageLarge);

                const stageComp = (
                  <Box
                    key={stage.id}
                    sx={(theme) => ({
                      display: "flex",
                      padding: theme.spacing(2),
                      ...(showStraightLine && {
                        position: "relative",
                        "&::before": {
                          content: '""',
                          position: "absolute",
                          top: "48%",
                          left: theme.spacing(-2),
                          borderTop: `2px solid ${theme.palette.divider}`,
                          width: theme.spacing(4),
                          height: 1,
                        },
                      }),

                      ...(showCurvedLine && {
                        position: "relative",
                        "&::before": {
                          content: '""',
                          position: "absolute",
                          bottom: "50%",
                          left: 0,
                          borderLeft: `2px solid ${theme.palette.divider}`,
                          borderBottom: `2px solid ${theme.palette.divider}`,
                          width: theme.spacing(2),
                          height: `calc(${
                            isCurvedLineExtend
                              ? APPROVED_STAGE_HEIGHT
                              : STAGE_HEIGHT
                          }px + ${theme.spacing(4)})`,
                        },
                      }),
                    })}
                  >
                    {stage.availableOperation ===
                      ManualOperation.MANUAL_OPERATION_APPROVE &&
                    stage.status === StageStatus.STAGE_RUNNING ? (
                      <ApprovalStage
                        id={stage.id}
                        name={stage.name}
                        onClick={() => {
                          setApproveTargetId(stage.id);
                        }}
                        active={isActive}
                      />
                    ) : (
                      <PipelineStage
                        id={stage.id}
                        name={stage.name}
                        status={stage.status}
                        metadata={stage.metadataMap}
                        onClick={handleOnClickStage}
                        active={isActive}
                        approver={approver}
                        skipper={skipper}
                        isDeploymentRunning={isRunning}
                      />
                    )}
                  </Box>
                );
                isPrevStageLarge = LARGE_STAGE_NAMES.includes(stage.name);
                return stageComp;
              })}
            </Box>
          );
        })}

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
      </Box>
    </Box>
  );
});
